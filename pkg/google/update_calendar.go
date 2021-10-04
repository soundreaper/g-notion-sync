package google

import (
	"time"

	"github.com/dstotijn/go-notion"
	"google.golang.org/api/calendar/v3"
)

// UpdateCalendar ...
func (g *CalendarService) UpdateCalendar(calendarID string, gcalNeedsUpdate []notion.DatabaseQueryResponse, gcalNeedsUpdateId []string) error {
	// Update Event(s) in G-Cal
	for idx, update := range gcalNeedsUpdate {
		// Retrieve Event info from Notion
		summary := update.Results[0].Properties.(notion.DatabasePageProperties)["Name"].Title[0].Text.Content
		startTime := update.Results[0].Properties.(notion.DatabasePageProperties)["Date"].Date.Start.Time.Format(time.RFC3339)
		endTime := update.Results[0].Properties.(notion.DatabasePageProperties)["Date"].Date.End.Time.Format(time.RFC3339)

		// Create Event Object to be Updated
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

		// Update Event in G-Cal
		_, err := g.Client.Events.Update(calendarID, gcalNeedsUpdateId[idx], event).Do()
		if err != nil {
			return err
		}
	}
	return nil
}
