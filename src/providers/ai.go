package providers

import (
	"github.com/yazu-codes/scanme-ai.git/src/model"
)

type AIProvider interface {
	Chat(request model.ChatRequest, sp string) (*model.ChatResponse, error)
	FinalizeSystemPrompt(data model.MenuDTO) string
}
