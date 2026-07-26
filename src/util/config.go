package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type ConfigReader struct {
	Server      Server
	MenuService MenuService
	AI          AI
}

func NewConfigReader() *ConfigReader {
	return &ConfigReader{Server: Server{}, MenuService: MenuService{}}
}

func (c *ConfigReader) Setup() {
	config := os.Getenv("CONFIG_YAML_AI")
	// config = "a"

	configPath := filepath.Join("configs", "config.yaml")

	fmt.Println("Making the directory for config")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		log.Fatal(err)
	}

	// err = os.WriteFile(configPath, []byte(config), 0600)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if config != "" {
		fmt.Println("CONFIG_YAML environment variable is set. Writing to config.yaml.")
		err := os.WriteFile(configPath, []byte(config), 0600)
		if err != nil {
			fmt.Println("Error writing file.")
			panic(err)
		}
		fmt.Println("Wrote to config")
	} else {
		fmt.Println("CONFIG_YAML environment variable is not set. Using existing config.yaml.")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs") // current directory

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	fmt.Println("config read successfully")

	// -----------------------
	// Extract values
	// -----------------------
	c.Server.Address = viper.GetString("server.address")
	c.Server.Port = viper.GetString("server.port")

	c.MenuService.Address = viper.GetString("menu_service.address")
	c.MenuService.Port = viper.GetString("menu_service.port")

	c.AI.OpenAIKey = viper.GetString("ai.openai_api_key")
	c.AI.AnthropicKey = viper.GetString("ai.anthropic_api_key")
	c.AI.SystemPrompt = viper.GetString("ai.system_prompt")
	// c.MenuService.Port = viper.GetString("menu_service.port")
}

type Server struct {
	Address string
	Port    string
}

func (s *Server) ConstructUrl() string {
	if s.Port == "" {
		return s.Address
	}
	return s.Address + ":" + s.Port
}

type MenuService struct {
	Address string
	Port    string
}

type AI struct {
	OpenAIKey    string
	AnthropicKey string
	SystemPrompt string
}
