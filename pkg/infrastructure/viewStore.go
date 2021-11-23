package infrastructure

import "github.com/google/uuid"

type ViewStorer interface {
	Update(id uuid.UUID, view interface{}) error
	Create(view interface{}) error
}
