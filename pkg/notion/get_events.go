package notion

import (
	"context"
	"time"

	"github.com/dstotijn/go-notion"
)

// GetEvents ...
func (c *Client) GetEvents() (*notion.DatabaseQueryResponse, error) {
	// Retrieve Notion Client
	ctx := context.Background()
	client := c.Client
	notionDB := c.DB

	// Set Time for Query Parameter
	t, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	// Create Query Parameter with Multiple Search Criterias
	prop := []notion.DatabaseQueryFilter{
		{
			Property: "Date",
			Date: &notion.DateDatabaseQueryFilter{
				OnOrAfter: &t,
			},
		},
		{
			Property: "GCal_ID",
			Text: &notion.TextDatabaseQueryFilter{
				IsNotEmpty: true,
			},
		},
	}

	// Query Notion Database to Find Events that Occur Today and DO Have a GCal_ID
	notionToday, err := client.QueryDatabase(ctx, notionDB, &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			And: prop,
		},
	})
	if err != nil {
		return nil, err
	}

	return &notionToday, nil
}
