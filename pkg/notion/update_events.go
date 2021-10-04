package notion

import (
	"context"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/soundreaper/g-notion-sync/pkg/google"
)

// UpdateEvents ...
func (c *Client) UpdateEvents(notionNeedsUpdate []google.EventItem, notionNeedsUpdateId []string) error {
	client := c.Client
	ctx := context.Background()

	// Update Event(s) in Notion
	for idx, update := range notionNeedsUpdate {
		// Retrieve Event info from G-Cal
		start, _ := time.Parse(time.RFC3339, update.StartDate)
		end, _ := time.Parse(time.RFC3339, update.EndDate)
		last := notion.NewDateTime(end, true)

		// Update Event in Notion
		_, _ = client.UpdatePageProps(ctx, notionNeedsUpdateId[idx], notion.UpdatePageParams{
			DatabasePageProperties: &notion.DatabasePageProperties{
				"Name": notion.DatabasePageProperty{
					Title: []notion.RichText{
						{
							Type: notion.RichTextTypeText,
							Text: &notion.Text{
								Content: update.Summary,
							},
						},
					},
				},
				"Date": notion.DatabasePageProperty{
					Date: &notion.Date{
						Start: notion.NewDateTime(start, true),
						End:   &last,
					},
				},
			},
		})
	}
	return nil
}
