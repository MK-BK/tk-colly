package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

var Routers = []Route{
	{Method: http.MethodGet, Path: "/list", Handler: ListMovie},
	{Method: http.MethodGet, Path: "/view/:id", Handler: ViewMovie},
}
