package google

import "google.golang.org/api/calendar/v3"

// GetCalendars returns the Authorized User's Calendars
func (g *CalendarService) GetCalendars() (*calendar.CalendarList, error) {
	// G-Cal Service
	service := g.Client
	cList, err := getAuthUserCalendars(service)
	if err != nil {
		return nil, err
	}
	return cList, nil
}
