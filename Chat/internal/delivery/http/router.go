package http

import (
	"ChatRoom/chat/internal/delivery/http/chat"
	"ChatRoom/chat/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouteOption struct {
	ChatUsecase domain.ChatUsecase
	Logger      chat.Logger
}

func NewRouter(opt RouteOption) *gin.Engine {
	router := gin.Default()
	chatGroup := router.Group("/chat")
	{
		chat.RegisterRoutes(chatGroup, opt.ChatUsecase, opt.Logger)
	}

	// 直接 serve index.html
	router.LoadHTMLFiles("./public/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return router
}
