package user

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/kevingentile/chet/domains/user"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/commandhandler/bus"
	localEventBus "github.com/looplab/eventhorizon/eventbus/local"
	memoryEventStore "github.com/looplab/eventhorizon/eventstore/memory"
	"github.com/looplab/eventhorizon/repo/memory"
	"github.com/looplab/eventhorizon/uuid"
)

func TestUser(t *testing.T) {
	eventBus := localEventBus.NewEventBus()
	go func() {
		for e := range eventBus.Errors() {
			log.Printf("eventbus: %s", e.Error())
		}
	}()

	// Create the event store.
	eventStore, err := memoryEventStore.NewEventStore(
		memoryEventStore.WithEventHandler(eventBus), // Add the event bus as a handler after save.
	)
	if err != nil {
		log.Fatalf("could not create event store: %s", err)
	}

	// Create the command bus.
	commandBus := bus.NewCommandHandler()

	// Create the read repositories.
	userRepo := memory.NewRepo()
	userRepo.SetEntityFactory(func() eh.Entity { return &user.User{} })

	ctx := context.Background()

	// Setup a test utility waiter that waits for all 11 events to occur before
	// evaluating results.
	var wg sync.WaitGroup
	wg.Add(1)
	eventBus.AddHandler(ctx, eh.MatchAll{}, eh.EventHandlerFunc(
		func(ctx context.Context, e eh.Event) error {
			wg.Done()
			return nil
		},
	))
	eventID := uuid.New()
	testUserID := uuid.New()
	user.Setup(ctx, eventStore, eventBus, eventBus, commandBus, userRepo, eventID)

	if err := commandBus.HandleCommand(ctx, &user.CreateUser{ID: testUserID, Username: "username", Email: "test@gmail.com"}); err != nil {
		t.Error(err)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	users, err := userRepo.FindAll(ctx)
	if err != nil {
		t.Error(err)
	}

	userList := []string{}
	for _, u := range users {
		if u, ok := u.(*user.User); ok {
			userList = append(userList, u.Username)
		}
	}

	if err := eventBus.Close(); err != nil {
		t.Error(err)
	}

	if err := userRepo.Close(); err != nil {
		t.Error(err)
	}

	if err := eventStore.Close(); err != nil {
		t.Error(err)
	}

	expectedList := []string{"username"}

	if userList[0] != expectedList[0] {
		t.Error("userlist not equal to expected")
	}

}
