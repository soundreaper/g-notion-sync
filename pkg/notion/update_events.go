package notion

import (
	"context"
	"github.com/dstotijn/go-notion"
	"github.com/soundreaper/g-notion-sync/pkg/google"
	"time"
)


// UpdateEvents ...
func (c *Client) UpdateEvents(notionNeedsUpdate []google.EventItem, notionNeedsUpdateId []string) error {
	client := c.Client
	ctx := context.Background()

	for idx, update := range notionNeedsUpdate {
		start, _ := time.Parse(time.RFC3339, update.StartDate)
		end, _ := time.Parse(time.RFC3339, update.EndDate)
		last := notion.NewDateTime(end, true)

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