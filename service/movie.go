package service

import (
	"context"

	"github.com/MK-BK/tk-colly/common"
	"github.com/MK-BK/tk-colly/models"
)

type MovieManager struct{}

func NewMovieManger() *MovieManager {
	return &MovieManager{}
}

func (*MovieManager) List(ctx context.Context, options *models.MovieListOption) ([]*models.Movie, error) {
	movies := make([]*models.Movie, 0)
	if err := common.DB.Limit(options.Limit).Offset(options.Offset).Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

func (*MovieManager) Get(ctx context.Context, id string) (*models.MovieView, error) {
	var movieView models.MovieView

	if err := common.DB.Model(&models.MovieView{}).Where("movie_id = ?", id).Find(&movieView).Error; err != nil {
		return nil, err
	}

	return &movieView, nil
}

func (*MovieManager) Save(ctx context.Context, movies ...*models.Movie) error {
	return common.DB.Save(movies).Error
}

func (*MovieManager) SaveView(ctx context.Context, movieViews ...*models.MovieView) error {
	return common.DB.Save(movieViews).Error
}
