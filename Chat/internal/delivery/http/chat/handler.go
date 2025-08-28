package chat

import (
	"ChatRoom/chat/internal/domain"
	"log"
	"net/http"

	cfg "ChatRoom/chat/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	usecase  domain.ChatUsecase
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

func RegisterRoutes(rg *gin.RouterGroup, uc domain.ChatUsecase) {
	h := &ChatHandler{
		usecase: uc,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[*websocket.Conn]bool),
	}

	rg.GET("/ws", h.HandleWS)
	rg.GET("/messages", h.GetMessages)
}

func (h *ChatHandler) HandleWS(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer conn.Close()
	h.clients[conn] = true
	defer delete(h.clients, conn)

	for {
		var msg domain.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("read error: %v", err)
			break
		}

		// 呼叫 Usecase
		err = h.usecase.PostMessage(c.Request.Context(), msg.Username, msg.Content)
		if err != nil {
			cfg.Logger.Errorln("PostMessage error:", err.Error())
		}

		// 廣播給所有 client
		for c := range h.clients {
			_ = c.WriteJSON(msg)
		}
	}
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	msgs, _ := h.usecase.GetMessages(c.Request.Context())
	c.JSON(http.StatusOK, msgs)
}
