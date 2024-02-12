package graph

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	accpb "github.com/infinimesh/proto/node/accounts"
	devpb "github.com/infinimesh/proto/node/devices"
	nspb "github.com/infinimesh/proto/node/namespaces"
	"github.com/infinimesh/proto/plugins"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InfinimeshProtobufEntity interface {
	*devpb.Device | *accpb.Account | *nspb.Namespace | *plugins.Plugin
}

type ListQueryResult[T InfinimeshProtobufEntity] struct {
	Result []T `json:"result"`
	Count  int `json:"count"`
}

type InfinimeshGenericActionsRepo[T InfinimeshProtobufEntity] interface {
	ListQuery(ctx context.Context, log *zap.Logger, from InfinimeshGraphNode) (*ListQueryResult[T], error)
}

type infinimeshGenericActionsRepo[T InfinimeshProtobufEntity] struct {
	db driver.Database
}

func NewGenericRepo[T InfinimeshProtobufEntity](db driver.Database) InfinimeshGenericActionsRepo[T] {
	return &infinimeshGenericActionsRepo[T]{
		db: db,
	}
}

const ListObjectsOfKind = `
LET result = (
	FOR node, edge, path IN 0..@depth OUTBOUND @from
	GRAPH @permissions_graph
	OPTIONS {order: "bfs", uniqueVertices: "global"}
	FILTER IS_SAME_COLLECTION(@@kind, node)
	FILTER edge.level > 0
	%s
	LET perm = path.edges[0]
	LET last = path.edges[-1]
	RETURN MERGE(node, {
		uuid: node._key,
		access: {
		level: last.role == 2 ? last.level : perm.level,
		role:  last.role == 2 ? last.role : perm.role,
		namespace: last.role == 2 ? null : path.vertices[-2]._key
		}
	})
)
RETURN { 
	result: (@offset > 0 && @limit > 0) ? SLICE(result, @offset, @offset + @limit) : result,
	count: LENGTH(result)
}
`

// List children nodes
// ctx - context
// log - logger
// db - Database connection
// from - Graph node to start traversal from
// children - children type(collection name)
// depth
func (r *infinimeshGenericActionsRepo[T]) ListQuery(ctx context.Context, log *zap.Logger, from InfinimeshGraphNode) (*ListQueryResult[T], error) {
	offset := OffsetValue(ctx)
	limit := LimitValue(ctx)

	var kind string
	switch fmt.Sprintf("%T", *new(T)) {
	case "*devices.Device":
		kind = schema.DEVICES_COL
	case "*accounts.Account":
		kind = schema.ACCOUNTS_COL
	case "*namespaces.Namespace":
		kind = schema.NAMESPACES_COL
	case "*plugins.Plugin":
		kind = schema.PLUGINS_COL
	default:
		return nil, fmt.Errorf("unknown type %T", *new(T))
	}

	bindVars := map[string]interface{}{
		"depth":             DepthValue(ctx),
		"from":              from.ID(),
		"permissions_graph": schema.PERMISSIONS_GRAPH.Name,
		"@kind":             kind,
		"offset":            offset,
		"limit":             limit,
	}
	log.Debug("Ready to build query", zap.Any("bindVars", bindVars))

	filters := ""
	if ns := NSFilterValue(ctx); ns != "" {
		filters += fmt.Sprintf("FILTER path.vertices[-2]._key == \"%s\"\n", ns)
	}

	cr, err := r.db.Query(ctx, fmt.Sprintf(ListObjectsOfKind, filters), bindVars)

	if err != nil {
		log.Debug("Error while executing query", zap.Error(err))
		return nil, err
	}

	defer cr.Close()

	var resp ListQueryResult[T]
	_, err = cr.ReadDocument(ctx, &resp)
	if err != nil {
		log.Warn("Error unmarshalling Document", zap.Error(err))
		return nil, status.Error(codes.Internal, "Couldn't execute query")
	}

	return &resp, nil
}
