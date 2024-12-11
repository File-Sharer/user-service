package handler

import (
	"net/http"
	"os"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) authSignUp(c *gin.Context) {
	var userReq model.User
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	user, jwtPair, err := h.services.Auth.SignUp(c.Request.Context(), &userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", jwtPair.RefreshToken, 3600 * 24 * 7, "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil, "accessToken": jwtPair.AccessToken, "user": user})
}

func (h *Handler) authSignIn(c *gin.Context) {
	var userReq model.User
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
		return
	}

	user, jwtPair, err := h.services.Auth.SignIn(c.Request.Context(), &userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", jwtPair.RefreshToken, 3600 * 24 * 7, "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil, "accessToken": jwtPair.AccessToken, "user": user})
}

func (h *Handler) authRefresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": err.Error()})
		return
	}

	decodedJwt, err := h.hasherClient.DecodeJWT(c.Request.Context(), &pb.DecodeJWTReq{
		Secret: os.Getenv("HASHER_SECRET"),
		Jwt: refreshToken,
	})
	if err != nil || !decodedJwt.Ok {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": err.Error()})
		return
	}

	jwtPair, err := h.hasherClient.GenerateJWTPair(c.Request.Context(), &pb.GenerateJWTPairReq{
		Secret: os.Getenv("HASHER_SECRET"),
		UserId: decodedJwt.GetUserId(),
		Role: decodedJwt.GetRole(),
	})
	if err != nil || !jwtPair.Ok {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": err.Error()})
		return
	}

	user, err := h.services.User.FindByID(c.Request.Context(), decodedJwt.GetUserId())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("refreshToken", jwtPair.GetRefreshToken(), 3600 * 24 * 7, "/", "localhost", true, true)

	logrus.Info("refreshed")

	c.JSON(http.StatusOK, gin.H{"ok": true, "error": nil, "accessToken": jwtPair.GetAccessToken(), "user": user})
}

func (h *Handler) authSignout(c *gin.Context) {
	c.SetCookie("refreshToken", "", 0, "", "", true, true)
}
