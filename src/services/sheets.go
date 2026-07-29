package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Sheets struct {
	spreadsheetID string
	sheetsAPI     *sheets.Service
}

// Initialize the Sheets service using credentials from config
func NewSheetsServiceFromConfig(spreadsheetID string, credentialsJSON []byte) (*Sheets, error) {
	ctx := context.Background()

	// Create service from credentials bytes
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, fmt.Errorf("unable to create sheets service: %v", err)
	}

	return &Sheets{
		spreadsheetID: spreadsheetID,
		sheetsAPI:     srv,
	}, nil
}

// SaveLead appends a lead to the spreadsheet
func (s *Sheets) SaveLead(name, email, phone, message string) error {
	// ctx := context.Background()

	// Prepare the row data
	values := []interface{}{
		name,
		email,
		phone,
		message,
		time.Now().Format("2006-01-02 15:04:05"),
	}

	fmt.Println("VALUES:")
	fmt.Println(values...)

	// Append to Sheet1
	_, err := s.sheetsAPI.Spreadsheets.Values.Append(
		s.spreadsheetID,
		"Sheet1!A2",
		&sheets.ValueRange{
			Values: [][]interface{}{values},
		},
	).Do()

	if err != nil {
		log.Printf("Unable to append data to sheet: %v", err)
		return err
	}

	log.Printf("Lead saved: %s (%s)", name, email)
	return nil
}

// SaveLeadFromMap saves a lead from a map
func (s *Sheets) SaveLeadFromMap(leadData map[string]interface{}) error {
	name := leadData["name"].(string)
	email := leadData["email"].(string)
	phone := leadData["phone"].(string)
	message := leadData["message"].(string)

	return s.SaveLead(name, email, phone, message)
}
