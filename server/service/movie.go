package service

import (
	"context"

	"github.com/MK-BK/tk-colly/database"
	"github.com/MK-BK/tk-colly/models"
)

type Movie struct{}

func NewMovieManger() *Movie {
	return &Movie{}
}

func (*Movie) List(ctx context.Context, options *models.MovieListOption) ([]*models.Movie, error) {
	movies := make([]*models.Movie, 0)
	if err := database.DB.Where("category = ? OR name LIKE %%s%", options.Categoty, options.Name).Find(movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (*Movie) View(ctx context.Context, id string) (*models.MovieView, error) {
	movieView := &models.MovieView{}
	if err := database.DB.Where("id = ?", id).Find(movieView).Error; err != nil {
		return nil, err
	}
	return movieView, nil
}

// TODO: SaveOrUpdate
func (*Movie) Save(ctx context.Context, movies ...*models.Movie) error {
	if err := database.DB.Save(movies).Error; err != nil {
		return err
	}

	return nil
}

// TODO: SaveOrUpdate
func (*Movie) SaveView(ctx context.Context, movieViews ...*models.MovieView) error {
	if err := database.DB.Save(movieViews).Error; err != nil {
		return err
	}

	return nil
}
