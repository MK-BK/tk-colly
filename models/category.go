package models

import (
	"context"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string
}

type CategoryManager interface {
	Create(ctx context.Context, categories ...*Category) error
	List(ctx context.Context) ([]*Category, error)
}
