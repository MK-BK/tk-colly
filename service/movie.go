package service

import (
	"context"

	"github.com/MK-BK/tk-colly/common"
	"github.com/MK-BK/tk-colly/models"
	"gorm.io/gorm"
)

type MovieManager struct{}

func NewMovieManger() *MovieManager {
	return &MovieManager{}
}

func (*MovieManager) List(ctx context.Context, options *models.MovieListOption) ([]*models.MovieView, error) {
	movieViews := make([]*models.MovieView, 0)
	if err := common.DB.Limit(options.Limit).
		Offset(options.Offset).
		Find(&movieViews).Error; err != nil {
		return nil, err
	}

	return movieViews, nil
}

func (*MovieManager) Get(ctx context.Context, id string) (*models.MovieView, error) {
	var movieView models.MovieView

	if err := common.DB.Where("id = ?", id).Find(&movieView).Error; err != nil {
		return nil, err
	}

	return &movieView, nil
}

func (*MovieManager) CreateMovieView(tx *gorm.DB, movieViews *models.MovieView) error {
	if err := tx.Take(&models.MovieView{}, "href = ?", movieViews.Href).Error; err != nil {
		if err := tx.Create(movieViews).Error; err != nil {
			return err
		}
		return nil
	}

	return tx.Where("id = ?", movieViews.ID).Updates(movieViews).Error
}

func (*MovieManager) CreateMoviePlayer(tx *gorm.DB, moviePlayer *models.MoviePlayer) error {
	if err := tx.Take(&models.MoviePlayer{}, "movie_id = ?", moviePlayer.MovieID).Error; err != nil {
		if err := tx.Create(moviePlayer).Error; err != nil {
			return err
		}
		return nil
	}

	return tx.Where("movie_id = ?", moviePlayer.MovieID).Updates(moviePlayer).Error
}

func (*MovieManager) GetPlayer(ctx context.Context, id string) (*models.MoviePlayer, error) {
	var player models.MoviePlayer

	if err := common.DB.Where("movie_id = ?", id).Find(&player).Error; err != nil {
		return nil, err
	}

	return &player, nil
}
