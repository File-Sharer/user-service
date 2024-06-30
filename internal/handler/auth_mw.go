package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) mwAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "user is not authorized"})
		c.Abort()
		return
	}

	token := strings.Split(header, " ")[1]
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "user is not authorized"})
		c.Abort()
		return
	}

	user, err := h.getUserDataFromTokenClaims(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": err.Error()})
		c.Abort()
		return
	}

	c.Set("user", *user.DTO())

	c.Next()
}
