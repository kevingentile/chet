package user

import (
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
)

const (
	CreateUserCommand        eh.CommandType = "CreateUser"
	MarkUserVerifiedCommand  eh.CommandType = "VerifyUser"
	DenyCreateUserCommand    eh.CommandType = "DenyCreateUser"
	ConfirmUserCreateCommand eh.CommandType = "ConfirmUserCreate"
)

type CreateUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}

func (c CreateUser) AggregateID() uuid.UUID          { return c.ID }
func (c CreateUser) AggregateType() eh.AggregateType { return UserAggregateType }
func (c CreateUser) CommandType() eh.CommandType     { return CreateUserCommand }

type DenyCreateUser struct {
	ID uuid.UUID
}

func (c DenyCreateUser) AggregateID() uuid.UUID          { return c.ID }
func (c DenyCreateUser) AggregateType() eh.AggregateType { return UserAggregateType }
func (c DenyCreateUser) CommandType() eh.CommandType     { return DenyCreateUserCommand }

type ConfirmUserCreate struct {
	ID    uuid.UUID
	Email string
}

func (c ConfirmUserCreate) AggregateID() uuid.UUID          { return c.ID }
func (c ConfirmUserCreate) AggregateType() eh.AggregateType { return UserAggregateType }
func (c ConfirmUserCreate) CommandType() eh.CommandType     { return ConfirmUserCreateCommand }

type DenyUserCreate struct {
	ID    uuid.UUID
	Email string
}

func (c DenyUserCreate) AggregateID() uuid.UUID          { return c.ID }
func (c DenyUserCreate) AggregateType() eh.AggregateType { return UserAggregateType }
func (c DenyUserCreate) CommandType() eh.CommandType     { return DenyCreateUserCommand }
