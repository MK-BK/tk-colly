package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/MK-BK/tk-colly/models"

	"github.com/gin-gonic/gin"
)

func listCategory(c *gin.Context) {
	categories, err := GE.CategoryManager.List(context.Background())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, categories)
}

func listMovie(c *gin.Context) {
	var options models.MovieListOption

	options.Offset, _ = strconv.Atoi(c.Query("offset"))
	options.Limit, _ = strconv.Atoi(c.Query("limit"))

	movies, err := GE.MoviesManager.List(context.Background(), &options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, movies)
}

func getMovie(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("Params id empty"))
		return
	}

	movieView, err := GE.MoviesManager.Get(context.Background(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, movieView)
}
