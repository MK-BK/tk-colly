package models

import "context"

type Category struct {
	Name string
}

type CategoryManager interface {
	Save(ctx context.Context, categories ...*Category) error
	List(ctx context.Context) ([]*Category, error)
}
