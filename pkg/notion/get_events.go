package notion

import (
	"context"
	"github.com/dstotijn/go-notion"
	"time"
)


// GetEvents ...
func (c *Client) GetEvents() (*notion.DatabaseQueryResponse, error){
	ctx := context.Background()
	client := c.Client
	notionDB := c.DB
	
	t, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

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
