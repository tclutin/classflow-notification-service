package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	APIBaseURL = "https://api.telegram.org/bot"
)

type Client interface {
	SendMessage(chatID int64, message string) error
}

type TGClient struct {
	Token string
}

func NewTGClient(token string) *TGClient {
	return &TGClient{
		Token: token,
	}
}

func (t *TGClient) SendMessage(chatID int64, message string) error {
	apiURL := fmt.Sprintf("%s%s/sendMessage", APIBaseURL, t.Token)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message payload: %w", err)
	}
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status: %s", resp.Status)
	}

	time.Sleep(1 * time.Second)

	return nil
}
