package providers

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/yazu-codes/scanme-ai.git/src/model"
)

type AnthropicProvider struct {
	client       *anthropic.Client
	systemPrompt string
}

func NewAnthropicProvider(apiKey, systemPrompt string) *AnthropicProvider {
	ac := anthropic.NewClient(
		option.WithAPIKey(apiKey), // or set ANTHROPIC_API_KEY env var
	)

	return &AnthropicProvider{
		client:       &ac,
		systemPrompt: systemPrompt,
	}
}

func (p *AnthropicProvider) FinalizeSystemPrompt(data model.MenuDTO) string {
	return fmt.Sprintf(p.systemPrompt, data.ToString())
}

func (p *AnthropicProvider) Chat(request model.ChatRequest, sp string) (*model.ChatResponse, error) {
	msgs := make([]anthropic.MessageParam, 0, len(request.Messages))

	for _, m := range request.Messages {
		switch m.Type {
		case model.RoleUser:
			msgs = append(msgs, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content)))
		case model.RoleAssistant:
			msgs = append(msgs, anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Content)))
		}
	}

	resp, err := p.client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:        anthropic.ModelClaudeHaiku4_5_20251001,
		MaxTokens:    1024,
		CacheControl: anthropic.NewCacheControlEphemeralParam(),
		System: []anthropic.TextBlockParam{
			{Text: sp},
		},
		Temperature: anthropic.Float(float64(request.Temperature)),
		Messages:    msgs,
	})
	if err != nil {
		return nil, err
	}

	var content string
	for _, block := range resp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return model.NewChatResponse(content), nil
}
