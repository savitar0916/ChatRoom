package cmd

import (
	cfg "ChatRoom/chat/config"
	chatApp "ChatRoom/chat/internal/application/chat"
	chatDelivery "ChatRoom/chat/internal/delivery/http"
	chatInfra "ChatRoom/chat/internal/infra/mysql"
	mysqlClient "ChatRoom/chat/pkg/mysql"

	toolsServer "github.com/Wuli-Giao-Giao/tools/server"
	toolsHttp "github.com/Wuli-Giao-Giao/tools/server/http"

	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// serveCmd 代表 `serve` 指令
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the chat service",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func startServer() {
	// 載入 config
	cfg.LoadConfig()

	// 建立 MySQL client
	mysqlConfig := &mysqlClient.Config{
		User:     cfg.AppConfig.MySQL.User,
		Password: cfg.AppConfig.MySQL.Password,
		Host:     cfg.AppConfig.MySQL.Host,
		Port:     cfg.AppConfig.MySQL.Port,
		Database: cfg.AppConfig.MySQL.Database,
	}

	cfg.Logger.Infof("Connecting to MySQL at %s:%d", mysqlConfig.Host, mysqlConfig.Port)
	// 連接 MySQL 資料庫
	mysqlClient, err := mysqlClient.NewMysqlClient(mysqlConfig)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	cfg.Logger.Infof("Connected to MySQL successfully")

	// 初始化 repository
	chatRepo := chatInfra.NewMySQLRepository(mysqlClient)

	// 初始化 usecase
	chatUsecase := chatApp.NewChatUsecase(chatRepo)

	// 建立 Router
	router := chatDelivery.NewRouter(chatDelivery.RouteOption{
		ChatUsecase: chatUsecase,
	})

	addr := fmt.Sprintf(":%d", cfg.AppConfig.Server.Port)

	httpServer := toolsHttp.NewHTTPServer(addr, router)

	// 使用 runner 啟動並管理 server
	runner := toolsServer.NewRunner(httpServer)
	runner.Run()
}
