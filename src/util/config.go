package util

import (
	"encoding/json"
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
	Google      Google
}

func NewConfigReader() *ConfigReader {
	return &ConfigReader{
		Server:      Server{},
		MenuService: MenuService{},
		AI:          AI{},
		Google:      Google{},
	}
}

func (c *ConfigReader) Setup() {
	config := os.Getenv("CONFIG_YAML_AI")

	configPath := filepath.Join("configs", "config.yaml")

	fmt.Println("Making the directory for config")
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		log.Fatal(err)
	}

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
	viper.AddConfigPath("configs")

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

	// Parse Google Console config
	c.Google.SpreadsheetID = viper.GetString("google_console.spreadsheet_id")
	c.Google.Credentials.Type = viper.GetString("google_console.credentials.type")
	c.Google.Credentials.ProjectID = viper.GetString("google_console.credentials.project_id")
	c.Google.Credentials.PrivateKeyID = viper.GetString("google_console.credentials.private_key_id")
	c.Google.Credentials.PrivateKey = viper.GetString("google_console.credentials.private_key")
	c.Google.Credentials.ClientEmail = viper.GetString("google_console.credentials.client_email")
	c.Google.Credentials.ClientID = viper.GetString("google_console.credentials.client_id")
	c.Google.Credentials.AuthURI = viper.GetString("google_console.credentials.auth_uri")
	c.Google.Credentials.TokenURI = viper.GetString("google_console.credentials.token_uri")
	c.Google.Credentials.AuthProviderX509CertURL = viper.GetString("google_console.credentials.auth_provider_x509_cert_url")
	c.Google.Credentials.ClientX509CertURL = viper.GetString("google_console.credentials.client_x509_cert_url")
	c.Google.Credentials.UniverseDomain = viper.GetString("google_console.credentials.universe_domain")
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

type Google struct {
	SpreadsheetID string
	Credentials   GoogleCredentials
}

type GoogleCredentials struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

// ToJSON converts credentials to JSON bytes (for API usage)
func (gc *GoogleCredentials) ToJSON() ([]byte, error) {
	return json.Marshal(gc)
}

// Validate ensures all required fields are set
func (g *Google) Validate() error {
	if g.SpreadsheetID == "" {
		return fmt.Errorf("google_console.spreadsheet_id is not set in config")
	}
	if g.Credentials.ClientEmail == "" {
		return fmt.Errorf("google_console.credentials.client_email is not set in config")
	}
	if g.Credentials.PrivateKey == "" {
		return fmt.Errorf("google_console.credentials.private_key is not set in config")
	}
	if g.Credentials.ProjectID == "" {
		return fmt.Errorf("google_console.credentials.project_id is not set in config")
	}
	return nil
}
