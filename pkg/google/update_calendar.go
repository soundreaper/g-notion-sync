package google

import (
	"github.com/dstotijn/go-notion"
	"google.golang.org/api/calendar/v3"
	"time"
)



// UpdateCalendar ...
func (g *CalendarService )UpdateCalendar(calendarID string, gcalNeedsUpdate []notion.DatabaseQueryResponse, gcalNeedsUpdateId []string) error {
	for idx, update := range gcalNeedsUpdate {
		summary := update.Results[0].Properties.(notion.DatabasePageProperties)["Name"].Title[0].Text.Content
		startTime := update.Results[0].Properties.(notion.DatabasePageProperties)["Date"].Date.Start.Time.Format(time.RFC3339)
		endTime := update.Results[0].Properties.(notion.DatabasePageProperties)["Date"].Date.End.Time.Format(time.RFC3339)

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

		_, err := g.Client.Events.Update(calendarID, gcalNeedsUpdateId[idx], event).Do()
		if err != nil {
			return err
		}
	}
	return nil
}
