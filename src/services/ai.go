package services

import (
	"log/slog"

	"github.com/yazu-codes/scanme-ai.git/src/model"
	"github.com/yazu-codes/scanme-ai.git/src/providers"
)

type AI struct {
	logger     *slog.Logger
	aiProvider providers.AIProvider
}

func NewAI(a providers.AIProvider, l *slog.Logger) *AI {
	return &AI{
		logger:     l,
		aiProvider: a,
	}
}

func (s *AI) Chat(
	chatRequest model.ChatRequest,
	menuData model.MenuDTO,
) (*model.ChatResponse, error) {
	sp := s.aiProvider.FinalizeSystemPrompt(menuData)

	chatResponse, err := s.aiProvider.Chat(chatRequest, sp)
	if err != nil {
		s.logger.Error("Failed to get response to a chat request: {%s}", err.Error())
		return nil, err
	}

	return chatResponse, nil
}
