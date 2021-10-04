package google


// CalendarType ....
type CalendarType struct {
	CID     string
	Summary string
	Events  []EventItem
}

// EventItem ...
type EventItem struct {
	ID        string
	Summary   string
	StartDate  string
	EndDate    string
	EditedDate string
}
