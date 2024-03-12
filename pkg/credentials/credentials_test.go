package credentials_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/arangodb/go-driver"
	"github.com/google/uuid"
	driver_mocks "github.com/infinimesh/infinimesh/mocks/github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/credentials"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	"github.com/infinimesh/proto/node/accounts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type credentialsControllerFixture struct {
	ctrl credentials.CredentialsController

	mocks struct {
		db   *driver_mocks.MockDatabase
		col  *driver_mocks.MockCollection
		edge *driver_mocks.MockCollection
	}
	data struct {
		ctx context.Context

		acc driver.DocumentID
	}
}

func newCredentialsControllerFixture(t *testing.T) *credentialsControllerFixture {
	f := &credentialsControllerFixture{}
	f.data.ctx = context.TODO()

	f.mocks.db = driver_mocks.NewMockDatabase(t)
	f.mocks.col = driver_mocks.NewMockCollection(t)
	f.mocks.edge = driver_mocks.NewMockCollection(t)

	g := driver_mocks.NewMockGraph(t)
	f.mocks.db.EXPECT().Graph(f.data.ctx, schema.CREDENTIALS_GRAPH.Name).
		Return(g, nil)

	g.EXPECT().VertexCollection(f.data.ctx, schema.CREDENTIALS_COL).
		Return(f.mocks.col, nil)
	g.EXPECT().EdgeCollection(f.data.ctx, schema.CREDENTIALS_EDGE_COL).
		Return(f.mocks.edge, driver.VertexConstraints{}, nil)

	f.ctrl = credentials.NewCredentialsController(f.data.ctx, zap.NewExample(), f.mocks.db)

	f.data.acc = driver.NewDocumentID(schema.ACCOUNTS_COL, uuid.New().String())

	return f
}

// Find
//

func TestFind_FailsOn_NotFound(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	_, err := f.ctrl.Find(context.TODO(), "mock", "Find")
	assert.Error(t, err)
	assert.EqualError(t, err, "couldn't find credentials")
}

func TestFind_FailsOn_Authorize(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	_, err := f.ctrl.Find(context.TODO(), "mock", "Authorize")
	assert.Error(t, err)
	assert.EqualError(t, err, "couldn't authorize")
}

func TestFind_Success(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	_, err := f.ctrl.Find(context.TODO(), "mock", "valid")
	assert.NoError(t, err)
}

// MakeCredentials
//

func TestMakeCredentials_FailsOn_NotGiven(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	_, err := f.ctrl.MakeCredentials(nil)
	assert.Error(t, err)
	assert.EqualError(t, err, "credentials aren't given")
}

func TestMakeCredentials_FailsOn_InvalidType(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	_, err := f.ctrl.MakeCredentials(&accounts.Credentials{Type: "invalid"})
	assert.Error(t, err)
	assert.EqualError(t, err, "unknown auth type")
}

func TestMakeCredentials_FailsOn_Fabric(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	_, err := f.ctrl.MakeCredentials(&accounts.Credentials{Type: "mock", Data: []string{"NewMockCredentials"}})
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid")
}

func TestMakeCredentials_Success(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	cred, err := f.ctrl.MakeCredentials(&accounts.Credentials{Type: "standard", Data: []string{"valid", "valid"}})
	assert.NoError(t, err)

	assert.Equal(t, "standard", cred.Type())
}

// ListCredentialsAndEdges
//

func TestListCredentialsAndEdges_FailsOn_Query(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	acc := driver.NewDocumentID("Accounts", uuid.New().String())

	f.mocks.db.EXPECT().Query(context.TODO(),
		credentials.ListCredentialsAndEdgesQuery,
		map[string]interface{}{
			"account":     acc,
			"credentials": schema.CREDENTIALS_COL,
		}).Return(nil, assert.AnError)

	_, err := f.ctrl.ListCredentialsAndEdges(context.TODO(), acc)
	assert.Error(t, err)
}

func TestListCredentialsAndEdges_FailsOn_ReadDocument(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	ctx := context.TODO()
	acc := driver.NewDocumentID("Accounts", uuid.New().String())

	c := &driver_mocks.MockCursor{}
	f.mocks.db.EXPECT().Query(ctx,
		credentials.ListCredentialsAndEdgesQuery,
		map[string]interface{}{
			"account":     acc,
			"credentials": schema.CREDENTIALS_COL,
		}).Return(c, nil)

	c.EXPECT().Close().Return(nil)
	c.EXPECT().ReadDocument(ctx, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)

	_, err := f.ctrl.ListCredentialsAndEdges(ctx, acc)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestListCredentialsAndEdges_Success(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	ctx := context.TODO()
	acc := driver.NewDocumentID("Accounts", uuid.New().String())

	c := &driver_mocks.MockCursor{}
	f.mocks.db.EXPECT().Query(ctx,
		credentials.ListCredentialsAndEdgesQuery,
		map[string]interface{}{
			"account":     acc,
			"credentials": schema.CREDENTIALS_COL,
		}).Return(c, nil)

	c.EXPECT().Close().Return(nil)
	c.EXPECT().ReadDocument(ctx, mock.MatchedBy(func(*[]string) bool {
		return true
	})).Return(driver.DocumentMeta{}, nil)

	_, err := f.ctrl.ListCredentialsAndEdges(ctx, acc)
	assert.NoError(t, err)
}

// ListCredentials
//

func TestListCredentials_FailsOn_Query(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	acc := driver.NewDocumentID("Accounts", uuid.New().String())

	f.mocks.db.EXPECT().Query(context.TODO(),
		credentials.ListCredentialsQuery,
		map[string]interface{}{
			"account":           acc.String(),
			"credentials_graph": schema.CREDENTIALS_GRAPH.Name,
		}).Return(nil, assert.AnError)

	_, err := f.ctrl.ListCredentials(context.TODO(), acc)
	assert.Error(t, err)
}

func TestListCredentials_FailsOn_Unmarshal(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	ctx := context.TODO()
	acc := driver.NewDocumentID("Accounts", uuid.New().String())

	c := &driver_mocks.MockCursor{}
	f.mocks.db.EXPECT().Query(ctx,
		credentials.ListCredentialsQuery,
		map[string]interface{}{
			"account":           acc.String(),
			"credentials_graph": schema.CREDENTIALS_GRAPH.Name,
		}).Return(c, nil)

	c.EXPECT().Close().Return(nil)
	c.EXPECT().ReadDocument(ctx, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)

	_, err := f.ctrl.ListCredentials(ctx, acc)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestListCredentials_FailsOn_Success(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	ctx := context.TODO()
	acc := driver.NewDocumentID("Accounts", uuid.New().String())

	c := &driver_mocks.MockCursor{}
	f.mocks.db.EXPECT().Query(ctx,
		credentials.ListCredentialsQuery,
		map[string]interface{}{
			"account":           acc.String(),
			"credentials_graph": schema.CREDENTIALS_GRAPH.Name,
		}).Return(c, nil)

	c.EXPECT().Close().Return(nil)

	called := false
	c.EXPECT().ReadDocument(ctx, mock.Anything).RunAndReturn(func(ctx context.Context, i interface{}) (driver.DocumentMeta, error) {
		if called {
			return driver.DocumentMeta{}, driver.NoMoreDocumentsError{}
		}
		called = true
		return driver.DocumentMeta{}, nil
	})

	_, err := f.ctrl.ListCredentials(ctx, acc)
	assert.NoError(t, err)
}

// MakeListable
//

func TestMakeListable_FailsOn_Unlistable(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	_, err := f.ctrl.MakeListable(credentials.ListCredentialsResponse{
		Type: "invalid",
		D:    map[string]interface{}{},
	})
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("Credentials of type %s aren't Listable", "invalid"))
}

func TestMakeListable_Success(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	_, err := f.ctrl.MakeListable(credentials.ListCredentialsResponse{
		Type: "standard",
		D: map[string]interface{}{
			"username": "valid",
			"password": "valid",
		},
	})
	assert.NoError(t, err)
}

// SetCredentials
//

func TestSetCredentials_Update_FailsOn_UpdateDocument(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	key := "standard-" + f.data.acc.Key()
	f.mocks.edge.EXPECT().ReadDocument(f.data.ctx, key, mock.Anything).Return(driver.DocumentMeta{}, nil)
	f.mocks.col.EXPECT().UpdateDocument(f.data.ctx, mock.Anything, mock.Anything).
		Return(driver.DocumentMeta{}, assert.AnError)

	err := f.ctrl.SetCredentials(f.data.ctx, f.data.acc, &credentials.StandardCredentials{})
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid_argument: Error updating Credentials of type")
}

func TestSetCredentials_Update_Success(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	key := "standard-" + f.data.acc.Key()
	f.mocks.edge.EXPECT().ReadDocument(f.data.ctx, key, mock.Anything).Return(driver.DocumentMeta{}, nil)
	f.mocks.col.EXPECT().UpdateDocument(f.data.ctx, mock.Anything, mock.Anything).
		Return(driver.DocumentMeta{}, nil)

	err := f.ctrl.SetCredentials(f.data.ctx, f.data.acc, &credentials.StandardCredentials{})
	assert.NoError(t, err)
}

func TestSetCredentials_FailsOn_CreateDocument(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	key := "standard-" + f.data.acc.Key()
	f.mocks.edge.EXPECT().ReadDocument(f.data.ctx, key, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)

	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(driver.DocumentMeta{}, assert.AnError)

	err := f.ctrl.SetCredentials(f.data.ctx, f.data.acc, &credentials.StandardCredentials{})
	assert.Error(t, err)
	assert.EqualError(t, err, "internal: Couldn't create credentials")
}

func TestSetCredentials_FailsOn_CreateEdge(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	key := "standard-" + f.data.acc.Key()
	f.mocks.edge.EXPECT().ReadDocument(f.data.ctx, key, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)

	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(driver.DocumentMeta{}, nil)
	f.mocks.edge.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(driver.DocumentMeta{}, assert.AnError)
	f.mocks.col.EXPECT().RemoveDocument(f.data.ctx, mock.Anything).
		Return(driver.DocumentMeta{}, nil)

	err := f.ctrl.SetCredentials(f.data.ctx, f.data.acc, &credentials.StandardCredentials{})
	assert.Error(t, err)
	assert.EqualError(t, err, "internal: Couldn't assign credentials")
}

func TestSetCredentials_Success(t *testing.T) {
	f := newCredentialsControllerFixture(t)

	key := "standard-" + f.data.acc.Key()
	f.mocks.edge.EXPECT().ReadDocument(f.data.ctx, key, mock.Anything).Return(driver.DocumentMeta{}, assert.AnError)

	f.mocks.col.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(driver.DocumentMeta{}, nil)
	f.mocks.edge.EXPECT().CreateDocument(f.data.ctx, mock.Anything).
		Return(driver.DocumentMeta{}, nil)

	err := f.ctrl.SetCredentials(f.data.ctx, f.data.acc, &credentials.StandardCredentials{})
	assert.NoError(t, err)
}
