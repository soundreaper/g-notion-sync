package google

import "google.golang.org/api/calendar/v3"




// GetCalendars returns the Auth User's calendars
func (g *CalendarService) GetCalendars() (*calendar.CalendarList, error) {
	// service
	service := g.Client
	cList, err := getAuthUserCalendars(service)
	if err != nil {
		return nil, err
	}
	return cList, nil
}


