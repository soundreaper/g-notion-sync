package notion

import (
	"context"
	"github.com/dstotijn/go-notion"
)

// CheckNewEvents ...
func (c *Client )CheckNewEvents() (*notion.DatabaseQueryResponse, error) {
	client := c.Client
	ctx := context.Background()
	notionDB := c.DB

	notInGCal, err := client.QueryDatabase(ctx, notionDB, &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "GCal_ID",
			Text: &notion.TextDatabaseQueryFilter{
				IsEmpty: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &notInGCal, nil
}
