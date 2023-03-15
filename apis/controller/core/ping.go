package core

import (
	"github.com/gin-gonic/gin"
)

type resp struct {
	Status string `json:"status" example:"pong"`
}

// Ping
//
//	@Summary		health check
//	@Description	do ping
//	@Tags			core
//	@Produce		json
//	@Success		200	{object}	resp
//	@Router			/ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, resp{Status: "pong"})
}
