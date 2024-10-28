package handler

import (
	"context"
	"os"

	pb "github.com/File-Sharer/user-service/hasher_pbs"
	"github.com/File-Sharer/user-service/internal/model"
	"github.com/File-Sharer/user-service/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Handler struct {
	services *service.Service
	hasherClient pb.HasherClient
}

func New(services *service.Service, hasherClient pb.HasherClient) *Handler {
	return &Handler{
		services: services,
		hasherClient: hasherClient,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{viper.GetString("fileService.origin"), viper.GetString("frontend.origin")},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", h.authSignUp)
			auth.POST("/signin", h.authSignIn)
		}

		user := api.Group("/user")
		{
			user.GET("", h.mwAuth, h.userGet)
			user.GET("/:id", h.mwAuth, h.userFindByID)
		}
	}

	return router
}

func (h *Handler) getUserDataFromTokenClaims(ctx context.Context, token string) (*model.User, error) {
	userRes, err := h.hasherClient.DecodeJWT(ctx, &pb.DecodeJWTReq{Secret: os.Getenv("HASHER_SECRET"), Jwt: token})
	if err != nil {
		return nil, err
	}

	user, err := h.services.User.FindByID(ctx, userRes.UserId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) getUser(c *gin.Context) *model.User {
	userReq, _ := c.Get("user")

	user, ok := userReq.(model.User)
	if !ok {
		return nil
	}

	return &user
}
