package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yazu-codes/scanme-ai.git/src/api"
	"github.com/yazu-codes/scanme-ai.git/src/api/handlers"
	"github.com/yazu-codes/scanme-ai.git/src/providers"
	"github.com/yazu-codes/scanme-ai.git/src/services"
	"github.com/yazu-codes/scanme-ai.git/src/util"
)

func main() {
	fmt.Println("Hi")

	var config *util.ConfigReader = util.NewConfigReader()
	config.Setup()

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	logger = logger.With(slog.String("component", "ai_service"))

	// var anthropicProvider providers.AnthropicProvider
	anthropicProvider := providers.NewAnthropicProvider(config.AI.AnthropicKey, config.AI.SystemPrompt)

	aiService := services.NewAI(anthropicProvider, logger)
	menuService := services.NewMenu(config.MenuService.Address, logger)

	aiHandler := handlers.NewAI(aiService, menuService, logger)

	server := api.NewServer(config.Server.ConstructUrl(), logger)

	server.SetupDefaultConfig()

	server.Router.GET("/", aiHandler.HealthCheck)
	server.Router.POST("/leads", aiHandler.Leads)
	server.Router.POST("/chat", aiHandler.Chat)

	server.Run()
}
