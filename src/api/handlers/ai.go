package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yazu-codes/scanme-ai.git/src/model"
	"github.com/yazu-codes/scanme-ai.git/src/services"
)

type AI struct {
	logger        *slog.Logger
	aiService     *services.AI
	menuService   *services.Menu
	sheetsService *services.Sheets
}

func NewAI(ai *services.AI, menu *services.Menu, sheets *services.Sheets, logger *slog.Logger) *AI {
	return &AI{aiService: ai, menuService: menu, logger: logger, sheetsService: sheets}
}

func (i *AI) HealthCheck(c *gin.Context) {
	fmt.Println("what")
	c.JSON(200, "up-and-running")
}

func (i *AI) Leads(c *gin.Context) {
	// Read and log body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "failed to read body"})
		return
	}

	// Restore body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Log the received body
	log.Printf("Received: %s", string(body))

	// Bind JSON
	var lead model.Lead
	if err := c.BindJSON(&lead); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save to Google Sheets
	err = i.sheetsService.SaveLead("", lead.Email, lead.Phone, "")
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save lead"})
		return
	}

	// Success response
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Lead saved successfully",
		"lead":    lead,
	})

}

func (i *AI) Chat(c *gin.Context) {
	i.logger.Info("Chat attempt")

	var chatRequest model.ChatRequest
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Fetch menu data via 2 calls:
	menuData, err := i.menuService.GetMenuDataByMenuName(c, chatRequest.MenuName)
	if err != nil {
		i.logger.Error(err.Error())
		c.JSON(400, gin.H{"error": "menu fetching failed"})
		return
	}

	chatResponse, err := i.aiService.Chat(chatRequest, menuData)
	if err != nil {
		c.JSON(400, gin.H{"error": "chat response failed"})
		return
	}
	// menuID := c.PostForm("menu_id")
	// if menuID == "" {
	// 	c.JSON(400, gin.H{"error": "missing menu id"})
	// 	return
	// }

	// file, header, err := c.Request.FormFile("image")
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": "missing image"})
	// 	return
	// }
	// defer file.Close()

	// created, err := i.service.Create(
	// 	c.Request.Context(),
	// 	file,
	// 	header.Filename,
	// 	menuID,
	// )

	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	fmt.Println(chatResponse)

	c.JSON(201, chatResponse)
}
