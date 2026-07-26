package model

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

type Message struct {
	Type    string
	Content string
}

func NewMessage(msgType, content string) Message {
	return Message{
		Type:    msgType,
		Content: content,
	}
}

func NewUserMessage(content string) Message {
	return Message{Type: RoleUser, Content: content}
}

func NewAssistantMessage(content string) Message {
	return Message{Type: RoleAssistant, Content: content}
}

type ChatRequest struct {
	MenuName    string `json:"3fc0469d66b2e72d3a7185687df9d459"`
	Messages    []Message
	Temperature float32
}

func NewChatRequest(menuId string, messages []Message, temperature float32) ChatRequest {
	return ChatRequest{
		MenuName:    menuId,
		Messages:    messages,
		Temperature: temperature,
	}
}

func NewSimpleChatRequest(menuId, userMessage string) ChatRequest {
	return ChatRequest{
		MenuName: menuId,
		Messages: []Message{
			NewUserMessage(userMessage),
		},
		Temperature: 1.0,
	}
}

type ChatResponse struct {
	Content string
}

func NewChatResponse(content string) *ChatResponse {
	return &ChatResponse{Content: content}
}
