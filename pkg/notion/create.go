package notion

import (
	"context"
	"github.com/dstotijn/go-notion"
	"github.com/soundreaper/g-notion-sync/pkg/google"
	"time"
)

// Create ...
func (c *Client) Create(notInNotion []google.EventItem) error {
	client := c.Client
	notionDB := c.DB
	ctx := context.Background()

	for _, create := range notInNotion {
		start, _ := time.Parse(time.RFC3339, create.StartDate)
		end, _ := time.Parse(time.RFC3339, create.EndDate)
		ender := notion.NewDateTime(end, true)

		_ , _ = client.CreatePage(ctx, notion.CreatePageParams{
			ParentType: notion.ParentTypeDatabase,
			ParentID:   notionDB,

			DatabasePageProperties: &notion.DatabasePageProperties{
				"Name": notion.DatabasePageProperty{
					Title: []notion.RichText{
						{
							Type: notion.RichTextTypeText,
							Text: &notion.Text{
								Content: create.Summary,
							},
						},
					},
				},
				"Date": notion.DatabasePageProperty{
					Date: &notion.Date{
						Start: notion.NewDateTime(start, true),
						End:   &ender,
					},
				},
				"GCal_ID": notion.DatabasePageProperty{
					RichText: []notion.RichText{
						{
							Type: notion.RichTextTypeText,
							Text: &notion.Text{
								Content: create.ID,
							},
						},
					},
				},
			},
		})
	}
	return nil
}
