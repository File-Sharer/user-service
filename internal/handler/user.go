package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) userGet(c *gin.Context) {
	user := h.getUser(c)

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil, "data": user})
}

func (h *Handler) userFindByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.services.User.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil, "data": user.DTO()})
}
