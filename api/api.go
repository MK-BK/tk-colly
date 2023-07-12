package api

import (
	"net/http"

	"github.com/MK-BK/tk-colly/models"
	"github.com/gin-gonic/gin"
)

var GE = &models.GlobalEnvironment

type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

var Routers = []Route{
	{Method: http.MethodGet, Path: "/category", Handler: listCategory},
	{Method: http.MethodGet, Path: "/movies", Handler: listMovie},
	{Method: http.MethodGet, Path: "/movies/:id", Handler: getMovie},
}
