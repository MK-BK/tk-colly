package models

import (
	"context"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Name        string
	Description string
	Href        string
	Categoty    string
}

type MovieListOption struct {
	Name     string
	Categoty string
	Limit    int
	Offset   int
}

type MoviesManager interface {
	List(ctx context.Context, options *MovieListOption) ([]*Movie, error)
	Get(ctx context.Context, id string) (*MovieView, error)
	Save(ctx context.Context, movies ...*Movie) error
	SaveView(ctx context.Context, movieViews ...*MovieView) error
}
