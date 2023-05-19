package api

import (
	"errors"
	"net/http"

	"github.com/MK-BK/tk-colly/models"

	"github.com/gin-gonic/gin"
)

func ListMovie(c *gin.Context) {
	var options models.MovieListOption
	if err := c.ShouldBind(&options); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	movies, err := models.GlobalEnvironment.MovieInterface.List(c.Request.Context(), &options)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, movies)
}

func ViewMovie(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("Params id empty"))
		return
	}

	movieView, err := models.GlobalEnvironment.MovieInterface.View(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, movieView)
}
