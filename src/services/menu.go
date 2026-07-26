package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/yazu-codes/scanme-ai.git/src/model"
)

type Menu struct {
	httpClient *http.Client
	menuApiURL string
	logger     *slog.Logger
}

func NewMenu(m string, l *slog.Logger) *Menu {
	return &Menu{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		menuApiURL: m,
		logger:     l,
	}
}

func (c *Menu) GetMenu(ctx context.Context, menuCode string) (string, error) {
	url := fmt.Sprintf("%s/c/%s", c.menuApiURL, menuCode)

	menuName := ""

	err := c.doGet(ctx, url, menuName)
	if err != nil {
		return "", err
	}

	return menuName, nil
}

func (c *Menu) GetMenuDataByMenuName(ctx context.Context, menuName string) (model.MenuDTO, error) {
	url := fmt.Sprintf("%s/%s", c.menuApiURL, menuName)

	menuData := model.MenuDTO{}

	fmt.Println("DOING GET", menuName)

	err := c.doGet(ctx, url, &menuData)
	if err != nil {
		return model.MenuDTO{}, err
	}

	fmt.Println(menuData)

	fmt.Println("DONE GET")

	return menuData, nil
}

func (c *Menu) doGet(ctx context.Context, url string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("RAW RESPONSE:", string(body))

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
