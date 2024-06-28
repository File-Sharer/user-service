package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) mwAuth(c *gin.Context) {
	tokenCookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not authorized"})
		c.Abort()
		return
	}

	user, err := h.getUserDataFromTokenClaims(c.Request.Context(), tokenCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("user", *user.DTO())

	c.Next()
}
