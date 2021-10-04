package google

import "google.golang.org/api/calendar/v3"





func getAuthUserCalendars(srv *calendar.Service) (*calendar.CalendarList, error) {
	cList, err := srv.CalendarList.List().Fields("items/id").Do()
	if err != nil {
		return nil, err
	}

	return cList, nil
}

