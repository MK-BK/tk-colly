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

type MovieView struct {
	gorm.Model
	Name        string
	DisplayName string
	Href        string `gorm:"primarykey"`
	ImagePath   string
	Description string
	Rating      string
	Director    string
	Actors      string
	Category    string
	SubCategory string
	Region      string
	Language    string
	Released    string
}

type MovieListOption struct {
	Name     string
	Categoty string
	Limit    int
	Offset   int
}

type MoviePlayer struct {
	gorm.Model
	MovieID int       `gorm:"primarykey"`
	Players []*Player `gorm:"json"`
}

type Player struct {
	Resource string
	Name     string
	URL      string
}

type MoviesManager interface {
	List(ctx context.Context, options *MovieListOption) ([]*MovieView, error)
	Get(ctx context.Context, id string) (*MovieView, error)
	GetPlayer(ctx context.Context, id string) (*MoviePlayer, error)
	CreateMovieView(tx *gorm.DB, movieView *MovieView) error
	CreateMoviePlayer(tx *gorm.DB, moviePlayer *MoviePlayer) error
}
