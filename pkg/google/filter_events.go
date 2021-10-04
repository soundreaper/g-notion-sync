package google

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

// FilterEvents ...
func (g *CalendarService) FilterEvents(cList *calendar.CalendarList) ([]EventItem, error) {
	var calendars []CalendarType
	srv := g.Client
	for _, v := range cList.Items {
		calendars = append(calendars, CalendarType{CID: v.Id, Summary: v.Summary, Events: []EventItem{}})
	}

	for ind, cal := range calendars {
		t1 := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
		t2 := time.Now().AddDate(0, 0, 5).Format(time.RFC3339)

		events, err := srv.Events.List(cal.CID).ShowDeleted(false).SingleEvents(true).TimeMin(t1).TimeMax(t2).MaxResults(15).OrderBy("startTime").Do()
		if err != nil {
			return nil, err
		}

		if len(events.Items) != 0 {
			for _, item := range events.Items {
				startDate := item.Start.DateTime
				if startDate == "" {
					startDate = item.Start.Date
				}

				endDate := item.End.DateTime
				if endDate == "" {
					endDate = item.End.Date
				}

				editedDate := item.Updated
				if editedDate == "" {
					editedDate = item.Updated
				}
				calendars[ind].Events = append(calendars[ind].Events, EventItem{ID: item.Id, Summary: item.Summary, StartDate: startDate, EndDate: endDate, EditedDate: editedDate})
			}
		}
	}

	var calToday []EventItem
	for _, cal := range calendars {
		if len(cal.Events) != 0 {
			for _, item := range cal.Events {
				t, _ := time.Parse(time.RFC3339, item.StartDate)
				// adding 1 day to the current time so we can sync the Events for the next day, change as needed
				if t.Format("2006-01-02") == time.Now().Format("2006-01-02") {
					calToday = append(calToday, item)
				}
			}
		}
	}

	return calToday, nil
}
