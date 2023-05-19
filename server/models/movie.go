package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

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
	Categoty string
	Name     string
}

type Play struct {
	Resource string
	Name     string
	URL      string
}

type PlayLists []*Play

func (c PlayLists) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *PlayLists) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if v, ok := value.([]byte); ok {
		return json.Unmarshal(v, c)
	}

	return fmt.Errorf("failed to scan CustomType")
}

type MovieInterface interface {
	List(ctx context.Context, options *MovieListOption) ([]*Movie, error)
	View(ctx context.Context, id string) (*MovieView, error)
	Save(ctx context.Context, movies ...*Movie) error
	SaveView(ctx context.Context, movieViews ...*MovieView) error
}

type MovieView struct {
	gorm.Model
	MovieID     int
	Description string
	Name        string
	DisplayName string
	Rating      float32
	Director    string
	Actors      string
	Category    string
	Region      string
	Language    string
	Released    string
	PlayLists
}
