package service

import (
	"context"

	"github.com/MK-BK/tk-colly/common"
	"github.com/MK-BK/tk-colly/models"
)

type CategoryManager struct{}

func NewCategoryManager() *CategoryManager {
	return &CategoryManager{}
}

func (*CategoryManager) List(ctx context.Context) ([]*models.Category, error) {
	categories := make([]*models.Category, 0)
	if err := common.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (*CategoryManager) Create(ctx context.Context, categories ...*models.Category) error {
	return common.DB.Create(categories).Error
}
