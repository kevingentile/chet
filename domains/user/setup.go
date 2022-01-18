package user

import (
	"context"
	"fmt"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/looplab/eventhorizon/commandhandler/aggregate"
	"github.com/looplab/eventhorizon/commandhandler/bus"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/looplab/eventhorizon/uuid"
)

type HandlerAdder interface {
	AddHandler(context.Context, eh.EventMatcher, eh.EventHandler) error
}

// Setup configures the guestlist.
func Setup(
	ctx context.Context,
	eventStore eh.EventStore,
	local,
	global HandlerAdder,
	commandBus *bus.CommandHandler,
	invitationRepo eh.ReadWriteRepo,
	eventID uuid.UUID) error {

	// Create the aggregate repository.
	aggregateStore, err := events.NewAggregateStore(eventStore)
	if err != nil {
		return fmt.Errorf("could not create aggregate store: %w", err)
	}

	userHandler, err := aggregate.NewCommandHandler(UserAggregateType, aggregateStore)
	if err != nil {
		return fmt.Errorf("could not create command handler: %w", err)
	}
	commandHandler := eh.UseCommandHandlerMiddleware(userHandler)
	for _, cmd := range []eh.CommandType{
		CreateUserCommand,
	} {
		if err := commandBus.SetHandler(commandHandler, cmd); err != nil {
			return fmt.Errorf("could not add command handler for '%s': %w", cmd, err)
		}
	}

	userProjector := projector.NewEventHandler(
		NewUserProjector(), invitationRepo)
	userProjector.SetEntityFactory(func() eh.Entity { return &User{} })
	if err := local.AddHandler(ctx, eh.MatchEvents{
		UserCreatedEvent,
	}, userProjector); err != nil {
		return fmt.Errorf("could not add user projector: %w", err)
	}

	// // Add a logger as an observer.
	// if err := global.AddHandler(ctx, eh.MatchAll{},
	// 	eh.UseEventHandlerMiddleware(&Logger{}, observer.Middleware)); err != nil {
	// 	return fmt.Errorf("could not add logger to event bus: %w", err)
	// }

	return nil
}
