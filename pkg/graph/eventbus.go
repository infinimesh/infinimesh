package graph

import (
	"context"
	"errors"
	"fmt"
	accpb "github.com/infinimesh/proto/node/accounts"
	"time"

	infinimesh "github.com/infinimesh/infinimesh/pkg/shared"
	proto_eventbus "github.com/infinimesh/proto/eventbus"
	"github.com/infinimesh/proto/node"

	"connectrpc.com/connect"
	"github.com/arangodb/go-driver"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type EventsService struct {
	log *zap.Logger

	bus *EventBus
}

func NewEventsService(log *zap.Logger, bus *EventBus) *EventsService {
	return &EventsService{
		log: log.Named("EventsService"),
		bus: bus,
	}
}

func (e *EventsService) Subscribe(ctx context.Context, req *connect.Request[node.EmptyMessage], stream *connect.ServerStream[proto_eventbus.Event]) error {
	log := e.log.Named("Subscribe")

	outgoingContext, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated subscribe"))
	}

	accountClaim := outgoingContext.Get(infinimesh.INFINIMESH_ACCOUNT_CLAIM)
	if len(accountClaim) != 1 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("invalid auth"))
	}

	uuid := accountClaim[0]

	log.Debug("Requestor", zap.String("uuid", uuid))

	subscribe, err := e.bus.Subscribe(ctx, uuid)
	if err != nil {
		log.Error("Failed to subscribe", zap.Error(err))
		return err
	}

	for event := range subscribe {
		log.Info("Received event", zap.Any("Event", event))
		err := stream.Send(&event)
		if err != nil {
			log.Error("Failed to send event", zap.Error(err))
		}
	}
	return nil
}

type EventBusService interface {
	Subscribe(context.Context, string) (<-chan proto_eventbus.Event, error)
	Notify(context.Context, *proto_eventbus.Event) error
}

type EventBus struct {
	log *zap.Logger

	channel *amqp.Channel

	repo InfinimeshGenericActionsRepo[*accpb.Account]
}

func NewEventBus(log *zap.Logger, db driver.Database, amqp *amqp.Connection) (*EventBus, error) {
	channel, err := amqp.Channel()
	if err != nil {
		log.Error("Failed to create channel", zap.Error(err))
		return nil, err
	}

	err = channel.ExchangeDeclare("events", "direct", true, false, false, false, nil)
	if err != nil {
		log.Error("Failed to create exchange", zap.Error(err))
		return nil, err
	}

	repo := NewGenericRepo[*accpb.Account](db)

	return &EventBus{
		log:     log.Named("Eventbus"),
		channel: channel,
		repo:    repo,
	}, nil
}

func (e *EventBus) Subscribe(ctx context.Context, uuid string) (<-chan proto_eventbus.Event, error) {
	log := e.log.Named("Subscribe").Named(uuid)
	now := time.Now().Unix()

	queue, err := e.channel.QueueDeclare(fmt.Sprintf("event-%s-%d", uuid, now), false, false, true, false, nil)
	if err != nil {
		log.Error("Failed to create queue", zap.Error(err))
		return nil, err
	}

	err = e.channel.QueueBind(queue.Name, uuid, "events", false, nil)
	if err != nil {
		log.Error("Failed to bind queue", zap.Error(err))
		return nil, err
	}

	consume, err := e.channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Error("Failed to consume queue", zap.Error(err))
		return nil, err
	}
	events := make(chan proto_eventbus.Event)

	go func(log *zap.Logger, msgs <-chan amqp.Delivery, events chan<- proto_eventbus.Event) {
		for msg := range msgs {
			log.Debug("Received msg")
			var event proto_eventbus.Event
			err := proto.Unmarshal(msg.Body, &event)
			if err != nil {
				log.Error("Failed to unmarshal msg body", zap.Error(err))
				continue
			}

			log.Debug("Send event", zap.Any("Event", event))
			events <- event
		}
	}(log, consume, events)

	return events, nil
}

func (e *EventBus) Notify(ctx context.Context, event *proto_eventbus.Event) error {
	log := e.log.Named("Notify")
	log.Debug("Invoke")

	marshal, err := proto.Marshal(event)
	if err != nil {
		log.Error("Failed to marshal event", zap.Error(err))
		return err
	}

	var accounts []*accpb.Account
	switch event.GetEntity().(type) {
	case *proto_eventbus.Event_Account:
		result, err := e.repo.ListQuery(ctx, log, NewBlankAccountDocument(event.GetAccount().GetUuid()), "INBOUND")
		if err != nil {
			log.Error("Failed to list accounts", zap.Error(err))
			return err
		}
		accounts = result.Result
	case *proto_eventbus.Event_Device:
		result, err := e.repo.ListQuery(ctx, log, NewBlankAccountDocument(event.GetAccount().GetUuid()), "INBOUND")
		if err != nil {
			log.Error("Failed to list accounts", zap.Error(err))
			return err
		}
		accounts = result.Result
	case *proto_eventbus.Event_Namespace:
		result, err := e.repo.ListQuery(ctx, log, NewBlankAccountDocument(event.GetAccount().GetUuid()), "INBOUND")
		if err != nil {
			log.Error("Failed to list accounts", zap.Error(err))
			return err
		}
		accounts = result.Result
	default:
		return errors.New("failed to define entity type")
	}

	var errs []string

	for _, val := range accounts {
		err = e.channel.PublishWithContext(ctx, "events", val.GetUuid(), false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        marshal,
		})
		if err != nil {
			log.Error("Failed to publish event", zap.Error(err))
			errs = append(errs, err.Error())
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to publish some events. %v", errs)
	}

	return nil
}
