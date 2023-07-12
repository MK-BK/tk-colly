package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type MovieView struct {
	gorm.Model
	MovieID     int
	ImagePath   string
	Name        string
	Description string
	DisplayName string
	Rating      string
	Director    string
	Actors      string
	Category    string
	Region      string
	Language    string
	Released    string
	PlayLists
}

type PlayLists []*Play

type Play struct {
	Resource string
	Name     string
	URL      string
}

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
