package notion

import (
	"github.com/dstotijn/go-notion"
	"github.com/soundreaper/g-notion-sync/pkg/google"
)

type Service interface {
	CheckExistingEvents(*notion.DatabaseQueryResponse, []google.EventItem) ([]string, []google.EventItem, error)
	GetEvents() (*notion.DatabaseQueryResponse, error)
	UpdateEvents([]google.EventItem, []string) error
	Create([]google.EventItem) error
	CheckNewEvents() (*notion.DatabaseQueryResponse, error)
}

// Client ...
type Client struct {
	Client *notion.Client
	DB string
}

// NewNotionService ...
func NewNotionService(client *notion.Client, db string) Service {
	return &Client{
		Client: client,
		DB: db,
	}
}




