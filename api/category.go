package api

import (
	"context"
	"net/http"

	"github.com/MK-BK/tk-colly/models"
	"github.com/gin-gonic/gin"
)

func createCategory(c *gin.Context) {
	var category models.Category
	if err := c.Bind(&category); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := GE.CategoryManager.Create(c.Request.Context(), &category); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
}

func listCategory(c *gin.Context) {
	categories, err := GE.CategoryManager.List(context.Background())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, categories)
}
