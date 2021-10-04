package notion

import (
	"github.com/dstotijn/go-notion"
	"github.com/soundreaper/g-notion-sync/pkg/google"
)

// CheckExistingEvents ....
func (c *Client) CheckExistingEvents(notionToday *notion.DatabaseQueryResponse, calToday []google.EventItem) ([]string, []google.EventItem, error) {
	var gcalIdList []string
	var checkNeedsUpdate []string
	var notInNotion []google.EventItem

	// We check here if there are any already existing events in Notion
	if len(notionToday.Results) != 0 {
		for _, u := range notionToday.Results {
			// The database query has properties if Notion is not COMPLETELY empty
			if val, ok := u.Properties.(notion.DatabasePageProperties); ok {
				gcalIdList = append(gcalIdList, val["GCal_ID"].RichText[0].Text.Content)
			}
		}

		for _, gCal := range calToday {
			_, found := google.Find(gcalIdList, gCal.ID)
			if !found {
				// Check which events aren't in Notion and prepare a list for creating them
				notInNotion = append(notInNotion, gCal)
			} else {
				// Check which events are in Notion and prepare a list of ID's for potential updating
				checkNeedsUpdate = append(checkNeedsUpdate, gCal.ID)
			}
		}
		// If the returned result list was empty, Notion is missing events from Google Calendar
	} else {
		notInNotion = append(notInNotion, calToday...)
	}

	return checkNeedsUpdate, notInNotion, nil

}
