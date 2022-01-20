package user

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/looplab/eventhorizon/mocks"
	"github.com/looplab/eventhorizon/uuid"
)

func TestAggregateHandleCommand(t *testing.T) {
	TimeNow = func() time.Time {
		return time.Date(2017, time.July, 10, 23, 0, 0, 0, time.Local)
	}

	id := uuid.New()

	cases := map[string]struct {
		agg            *UserAggregate
		cmd            eh.Command
		expectedEvents []eh.Event
		expectedErr    error
	}{
		"unknown command": {
			agg: &UserAggregate{
				AggregateBase: events.NewAggregateBase(UserAggregateType, id),
			},
			cmd:            &mocks.Command{ID: id},
			expectedEvents: nil,
			expectedErr:    errors.New("could not handle command: Command"),
		},
		"create": {
			agg: &UserAggregate{
				AggregateBase: events.NewAggregateBase(UserAggregateType, id),
				Email:         "test@gmail.com",
				Username:      "username",
			},
			cmd:            &CreateUser{ID: id, Username: "username", Email: "test@gmail.com"},
			expectedEvents: nil,
			expectedErr:    errors.New("user already created"),
		},
		"user already created": {
			agg: &UserAggregate{
				AggregateBase: events.NewAggregateBase(UserAggregateType, id),
			},
			cmd: &CreateUser{ID: id, Username: "username", Email: "test@gmail.com"},
			expectedEvents: []eh.Event{
				eh.NewEvent(UserCreatedEvent, &UserCreatedData{Username: "username", Email: "test@gmail.com"}, TimeNow(), eh.ForAggregate(UserAggregateType, id, 1)),
			},
			expectedErr: nil,
		},
		"veify": {
			agg: &UserAggregate{
				AggregateBase: events.NewAggregateBase(UserAggregateType, id),
			},
			cmd: &VerifyUser{},
			expectedEvents: []eh.Event{
				eh.NewEvent(UserVerifiedEvent, nil, TimeNow(), eh.ForAggregate(UserAggregateType, id, 1)),
			},
			expectedErr: nil,
		},
		"user already verified": {
			agg: &UserAggregate{
				AggregateBase: events.NewAggregateBase(UserAggregateType, id),
				Verified:      true,
			},
			cmd:            &VerifyUser{},
			expectedEvents: nil,
			expectedErr:    errors.New("user already verified"),
		},
	}

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.agg.HandleCommand(context.Background(), tc.cmd)
			if (err != nil && tc.expectedErr == nil) ||
				(err == nil && tc.expectedErr != nil) ||
				(err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
				t.Errorf("test case '%s': incorrect error", name)
				t.Log("exp:", tc.expectedErr)
				t.Log("got:", err)
			}
			events := tc.agg.UncommittedEvents()
			if !reflect.DeepEqual(events, tc.expectedEvents) {
				t.Errorf("test case '%s': incorrect events", name)
				t.Log("exp:\n", pretty.Sprint(tc.expectedEvents))
				t.Log("got:\n", pretty.Sprint(events))
			}
		})
	}
}
