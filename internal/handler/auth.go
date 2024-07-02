package handler

import (
	"net/http"

	"github.com/File-Sharer/user-service/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) authSignUp(c *gin.Context) {
	var userReq model.User
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	user, token, err := h.services.Auth.SignUp(c.Request.Context(), &userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil, "token": token, "user": user})
}

func (h *Handler) authSignIn(c *gin.Context) {
	var userReq model.User
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	user, token, err := h.services.Auth.SignIn(c.Request.Context(), &userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil, "token": token, "user": user})
}
