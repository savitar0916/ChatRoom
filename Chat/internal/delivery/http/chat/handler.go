package chat

import (
	"ChatRoom/chat/internal/domain"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
}

type ChatHandler struct {
	usecase  domain.ChatUsecase
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	mu       sync.RWMutex
	logger   Logger
}

func RegisterRoutes(rg *gin.RouterGroup, uc domain.ChatUsecase, logger Logger) {
	if logger == nil {
		panic("chat handler logger is required")
	}
	h := &ChatHandler{
		usecase: uc,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		clients: make(map[*websocket.Conn]bool),
		logger:  logger,
	}

	rg.GET("/ws", h.HandleWS)
	rg.GET("/messages", h.GetMessages)
}

func (h *ChatHandler) HandleWS(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Errorf("upgrade error: %v", err)
		return
	}
	defer conn.Close()
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
	}()

	for {
		var msg domain.Message
		if err := conn.ReadJSON(&msg); err != nil {
			h.logger.Errorf("read error: %v", err)
			break
		}

		// 呼叫 Usecase
		saved, err := h.usecase.PostMessage(c.Request.Context(), msg.Username, msg.Content)
		if err != nil {
			h.logger.Errorln("PostMessage error:", err.Error())
			continue
		}

		clients := h.snapshotClients()

		// 廣播給所有 client
		for _, c := range clients {
			_ = c.WriteJSON(saved)
		}
	}
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	msgs, err := h.usecase.GetMessages(c.Request.Context())
	if err != nil {
		h.logger.Errorln("GetMessages error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	c.JSON(http.StatusOK, msgs)
}

func (h *ChatHandler) snapshotClients() []*websocket.Conn {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conns := make([]*websocket.Conn, 0, len(h.clients))
	for c := range h.clients {
		conns = append(conns, c)
	}
	return conns
}
