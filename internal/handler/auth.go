package handler

import (
	"net/http"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (h *Handler) authSignUp(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	token, err := h.services.Auth.SignUp(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	if err := auth.SendToken(c, token, viper.GetString("app.host")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil})
}

func (h *Handler) authSignIn(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	token, err := h.services.Auth.SignIn(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	if err := auth.SendToken(c, token, viper.GetString("app.host")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil})
}
