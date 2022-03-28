package handler

import (
	"awesomeProjectRentaTeam/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.Default()

	landingPage := r.Group("/")
	{
		landingPage.GET("/get_posts/", h.GetPostsWithLimit)
	}
	authorized := r.Group("/")
	authorized.Use(gin.BasicAuth(gin.Accounts{"admin": "admin"}))
	{
		authorized.POST("/insert_post", h.InsertPost)
		authorized.GET("/generate_random", h.GenerateRandomPosts)
	}
	return r
}
