package providers

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/yazu-codes/scanme-ai.git/src/model"
)

type OpenAIProvider struct {
	client *openai.Client
}

func (p *OpenAIProvider) Chat(
	ctx context.Context,
	req *model.ChatRequest,
) (*model.ChatResponse, error) {

	// translate your internal format
	// into OpenAI API format
	return nil, nil

}
