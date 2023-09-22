package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/MK-BK/tk-colly/models"
	"github.com/gin-gonic/gin"
)

func listMovie(c *gin.Context) {
	options := &models.MovieListOption{
		Offset: 0,
		Limit:  20,
	}

	if offset := c.Query("offset"); offset != "" {
		options.Offset, _ = strconv.Atoi(offset)
	}

	if limit := c.Query("limit"); limit != "" {
		options.Offset, _ = strconv.Atoi(limit)
	}

	if category := c.Query("category"); category != "" {
		options.Categoty = category
	}

	movies, err := GE.MoviesManager.List(context.Background(), options)
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

func getMoviePlayer(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("Params id empty"))
		return
	}

	player, err := GE.MoviesManager.GetPlayer(context.Background(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, player)
}
