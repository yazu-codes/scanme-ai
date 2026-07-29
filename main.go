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

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	logger = logger.With(slog.String("component", "ai_service"))

	var config *util.ConfigReader = util.NewConfigReader()
	config.Setup()

	// Validate Google config
	if err := config.Google.Validate(); err != nil {
		logger.Error("Google config validation failed: %v", err)
	}

	// Convert credentials to JSON
	credentialsJSON, err := config.Google.Credentials.ToJSON()
	if err != nil {
		logger.Error("Failed to convert credentials to JSON: %v", err)
	}

	// Initialize Sheets service with credentials from config
	sheetsService, err := services.NewSheetsServiceFromConfig(
		config.Google.SpreadsheetID,
		credentialsJSON,
	)
	if err != nil {
		logger.Error("Failed to initialize Sheets service: %v", err)
	}

	logger.Info("Sheets service initialized successfully")

	// var anthropicProvider providers.AnthropicProvider
	anthropicProvider := providers.NewAnthropicProvider(config.AI.AnthropicKey, config.AI.SystemPrompt)

	aiService := services.NewAI(anthropicProvider, logger)
	menuService := services.NewMenu(config.MenuService.Address, logger)

	aiHandler := handlers.NewAI(aiService, menuService, sheetsService, logger)

	server := api.NewServer(config.Server.ConstructUrl(), logger)

	server.SetupDefaultConfig()

	server.Router.GET("/", aiHandler.HealthCheck)
	server.Router.POST("/leads", aiHandler.Leads)
	server.Router.POST("/chat", aiHandler.Chat)

	server.Run()
}
