package infrastructure

type EventStorer interface {
	Save(event interface{}) error
}
