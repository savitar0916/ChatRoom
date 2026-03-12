package http

import (
	"ChatRoom/chat/internal/delivery/http/chat"
	"ChatRoom/chat/internal/domain"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

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

	// 讓 root 執行目錄改變時仍能正確找到前端頁面
	router.LoadHTMLFiles(filepath.Join(resolvePublicDir(), "index.html"))
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return router
}

func resolvePublicDir() string {
	candidates := []string{
		"./public",
		"./chat/public",
	}

	if _, filename, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates, filepath.Join(filepath.Dir(filename), "..", "..", "..", "public"))
	}

	for _, dir := range candidates {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}

		if _, err := os.Stat(filepath.Join(dir, "index.html")); err == nil {
			return dir
		}
	}

	return candidates[len(candidates)-1]
}
