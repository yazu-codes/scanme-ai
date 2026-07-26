package model

import (
	"fmt"
	"sort"
	"strings"
)

type MenuDTO struct {
	Menu PublicMenu `json:"menu"`
}

// Public-facing shapes — no DB identifiers exposed
type PublicMenu struct {
	MenuOwner         PublicMenuOwner         `json:"menu_owner"`
	MenuConfiguration PublicMenuConfiguration `json:"menu_configuration"`
	MenuItems         []PublicMenuItem        `json:"menu_items"`
}

type PublicMenuOwner struct {
	Name               string `json:"menu_owner_name"`
	Phone              string `json:"menu_owner_phone"`
	LogoURL            string `json:"menu_owner_logo_url"`
	Slogan             string `json:"menu_owner_slogan"`
	PlaceBackgroundURL string `json:"menu_owner_place_background_url"`
}

type PublicMenuConfiguration struct {
	CategoryOrder   string `json:"category_order"`
	BackgroundColor string `json:"background_color"`
	FontColor       string `json:"font_color"`
	FontFamily      string `json:"font_family"`
	FontSize        int    `json:"font_size"`
}

type PublicMenuItem struct {
	Name                 string  `json:"name"`
	Price                float64 `json:"price"`
	Description          string  `json:"description"`
	PictureURL           string  `json:"picture_url"`
	Category             string  `json:"category"`
	Allergens            string  `json:"allergens"`
	DisplayOrderPosition int     `json:"display_order_position"`
}

func (dto MenuDTO) ToString() string {
	var sb strings.Builder

	owner := dto.Menu.MenuOwner
	sb.WriteString("RESTAURANT INFO:\n")
	sb.WriteString(fmt.Sprintf("  Name: %s\n", valueOrNA(owner.Name)))
	if owner.Slogan != "" {
		sb.WriteString(fmt.Sprintf("  Slogan: %s\n", owner.Slogan))
	}
	if owner.Phone != "" {
		sb.WriteString(fmt.Sprintf("  Phone: %s\n", owner.Phone))
	}
	sb.WriteString("\n")

	items := dto.Menu.MenuItems

	byCategory := make(map[string][]PublicMenuItem)
	for _, item := range items {
		byCategory[item.Category] = append(byCategory[item.Category], item)
	}

	categories := make([]string, 0, len(byCategory))
	for cat := range byCategory {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	sb.WriteString("MENU ITEMS:\n")
	for _, cat := range categories {
		catItems := byCategory[cat]
		sort.Slice(catItems, func(i, j int) bool {
			return catItems[i].DisplayOrderPosition < catItems[j].DisplayOrderPosition
		})

		sb.WriteString(fmt.Sprintf("\nCategory: %s\n", valueOrNA(cat)))
		for _, item := range catItems {
			sb.WriteString(fmt.Sprintf("  - %s — %.2f\n", valueOrNA(item.Name), item.Price))
			if item.Description != "" {
				sb.WriteString(fmt.Sprintf("    Description: %s\n", item.Description))
			}
			sb.WriteString(fmt.Sprintf("    Allergens: %s\n", allergensOrNone(item.Allergens)))
		}
	}

	return sb.String()
}

func valueOrNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

func allergensOrNone(s string) string {
	if s == "" {
		return "None listed"
	}
	return s
}
