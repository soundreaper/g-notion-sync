package google

import (
	"github.com/dstotijn/go-notion"
	"google.golang.org/api/calendar/v3"
)

type Service interface {
	GetCalendars() (*calendar.CalendarList, error)
	UpdateCalendar(string, []notion.DatabaseQueryResponse, []string) error
	FilterEvents(cList *calendar.CalendarList) ([]EventItem, error)
	CreateEvent(string, *notion.DatabaseQueryResponse, *notion.Client) error
}

// CalendarService ...
type CalendarService struct {
	Client *calendar.Service
}

func NewGCalService(client *calendar.Service) Service {
	return &CalendarService{
		Client: client,
	}
}
