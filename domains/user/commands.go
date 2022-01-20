package user

import (
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
)

const (
	CreateUserCommand eh.CommandType = "CreateUser"
	VerifyUserCommand eh.CommandType = "VerifyUser"
)

type CreateUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}

func (c CreateUser) AggregateID() uuid.UUID          { return c.ID }
func (c CreateUser) AggregateType() eh.AggregateType { return UserAggregateType }
func (c CreateUser) CommandType() eh.CommandType     { return CreateUserCommand }

type VerifyUser struct {
	ID uuid.UUID
}

func (c VerifyUser) AggregateID() uuid.UUID          { return c.ID }
func (c VerifyUser) AggregateType() eh.AggregateType { return UserAggregateType }
func (c VerifyUser) CommandType() eh.CommandType     { return VerifyUserCommand }
