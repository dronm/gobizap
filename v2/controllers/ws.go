package controllers

import (
	"github.com/dronm/gobizap/v2/ws"

	"github.com/gin-gonic/gin"
	"net/http"
)

func WebSocket(c *gin.Context) {
	funcName := "WebSocket"
	sess := GetSession(c, funcName)
	if sess == nil {
		return
	}
	if err := ws.Upgrade(c.Writer, c.Request, sess); err != nil {
		ServeError(c, http.StatusInternalServerError, funcName+" ws.Upgrade()", err)
	}
}
