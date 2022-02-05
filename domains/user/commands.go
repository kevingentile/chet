package user

import (
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
)

const (
	CreateUserCommand        eh.CommandType = "CreateUser"
	VerifyUserCommand        eh.CommandType = "VerifyUser"
	ConfirmUserCreateCommand eh.CommandType = "ConfirmCreateUser"
	DenyUserCreateCommand    eh.CommandType = "DenyCreateUser"
)

type CreateUser struct {
	ID       uuid.UUID
	Username string
	Email    string
}

type DenyCreateUser struct {
	ID uuid.UUID
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

type ConfirmUserCreate struct {
	ID uuid.UUID
}

func (c ConfirmUserCreate) AggregateID() uuid.UUID          { return c.ID }
func (c ConfirmUserCreate) AggregateType() eh.AggregateType { return UserAggregateType }
func (c ConfirmUserCreate) CommandType() eh.CommandType     { return ConfirmUserCreateCommand }

type DenyUserCreate struct {
	ID uuid.UUID
}

func (c DenyUserCreate) AggregateID() uuid.UUID          { return c.ID }
func (c DenyUserCreate) AggregateType() eh.AggregateType { return UserAggregateType }
func (c DenyUserCreate) CommandType() eh.CommandType     { return DenyUserCreateCommand }
