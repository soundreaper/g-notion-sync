package google

import (
	"context"
	"github.com/dstotijn/go-notion"
	"google.golang.org/api/calendar/v3"
	"time"
)

// CreateEvents ...
func (g *CalendarService) CreateEvent(calendarID string, notInGCal *notion.DatabaseQueryResponse, client *notion.Client) error {
	srv := g.Client
	ctx := context.Background()

	for _, u := range notInGCal.Results {
		if val, ok := u.Properties.(notion.DatabasePageProperties); ok {
			summary := val["Name"].Title[0].Text.Content
			startTime := val["Date"].Date.Start.Time.Format(time.RFC3339)
			endTime := val["Date"].Date.End.Time.Format(time.RFC3339)

			event := &calendar.Event{
				Summary: summary,
				Start: &calendar.EventDateTime{
					DateTime: startTime,
					TimeZone: startTime,
				},
				End: &calendar.EventDateTime{
					DateTime: endTime,
					TimeZone: endTime,
				},
			}

			event, err := srv.Events.Insert(calendarID, event).Do()
			if err != nil {
				return err
			}


			_, err = client.UpdatePageProps(ctx, u.ID, notion.UpdatePageParams{
				DatabasePageProperties: &notion.DatabasePageProperties{
					"GCal_ID": notion.DatabasePageProperty{
						RichText: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: event.Id,
								},
							},
						},
					},
				},
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}