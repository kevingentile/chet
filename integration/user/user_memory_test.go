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

func findUser(ctx context.Context, userRepo *memory.Repo, id uuid.UUID) *user.User {
	users, err := userRepo.FindAll(ctx)
	if err != nil {
		return nil
	}
	for _, u := range users {
		if u, ok := u.(*user.User); ok {
			if u.ID == id {
				return u
			}
		}
	}
	return nil
}

func validateUserCreated(t *testing.T, cmd *user.CreateUser, actualUser *user.User) {
	if actualUser == nil {
		t.Error("missing user")
		t.FailNow()
	}

	if actualUser.Email != cmd.Email {
		t.Error("emails not equal", actualUser.Email, cmd.Email)
	}

	if actualUser.Username != cmd.Username {
		t.Error("usernames not equal", actualUser.Username, cmd.Username)
	}

	if actualUser.ID != cmd.ID {
		t.Error("IDs not equal", actualUser.ID, cmd.ID)
	}

	if actualUser.CreatedAt == 0 {
		t.Error("createdAt not set", actualUser.CreatedAt)
	}
}

func validateUserVerified(t *testing.T, cmd *user.VerifyUser, actualUser *user.User) {
	if !actualUser.Verified {
		t.Error("user not verified")
	}
}

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

	var wg sync.WaitGroup
	wg.Add(2)
	eventBus.AddHandler(ctx, eh.MatchAll{}, eh.EventHandlerFunc(
		func(ctx context.Context, e eh.Event) error {
			wg.Done()
			return nil
		},
	))
	eventID := uuid.New()
	testUserID := uuid.New()
	user.Setup(ctx, eventStore, eventBus, eventBus, commandBus, userRepo, eventID)
	createUserCmd := &user.CreateUser{ID: testUserID, Username: "username", Email: "test@gmail.com"}
	if err := commandBus.HandleCommand(ctx, createUserCmd); err != nil {
		t.Error(err)
		t.FailNow()
	}

	verifyUserCmd := &user.VerifyUser{ID: testUserID}
	if err := commandBus.HandleCommand(ctx, verifyUserCmd); err != nil {
		t.Error(err)
		t.FailNow()
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	actualUser := findUser(ctx, userRepo, testUserID)
	validateUserCreated(t, createUserCmd, actualUser)
	validateUserVerified(t, verifyUserCmd, actualUser)

	if err := eventBus.Close(); err != nil {
		t.Error(err)
	}

	if err := userRepo.Close(); err != nil {
		t.Error(err)
	}

	if err := eventStore.Close(); err != nil {
		t.Error(err)
	}
}
