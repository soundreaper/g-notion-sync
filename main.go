package main

import (
	"context"
	"log"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/soundreaper/g-notion-sync/config"
	"github.com/soundreaper/g-notion-sync/pkg/google"
	no "github.com/soundreaper/g-notion-sync/pkg/notion"
)

func main() {
	ctx := context.Background()
	// Configuration

	// Google
	// Google Calendar Service
	srv := google.GetCalendarFromGoogle()
	// G-Cal ID; Primary refers to the authorized user
	calendarID := "primary"
	// Create G-Cal Client
	googleCalClient := google.NewGCalService(srv)

	// Notion
	// Getting Notion Secret Key and Database ID from .env
	notionKey := config.GetNotionConfig().NotionKey
	notionDB := config.GetNotionConfig().NotionDB
	// Create Notion Client with Specificed Database ID
	client := notion.NewClient(notionKey)
	notionClient := no.NewNotionService(client, notionDB)

	// Get Authorized User for G-Cal
	calList, err := googleCalClient.GetCalendars()
	if err != nil {
		log.Fatal(err)
	}

	// Filter and Sort G-Cal Events
	calToday, err := googleCalClient.FilterEvents(calList)
	if err != nil {
		log.Fatal(err)
	}

	// Get Notion Events
	notionToday, err := notionClient.GetEvents()
	if err != nil {
		log.Fatal(err)
	}

	// Check Existing Events in Both Services
	checkNeedsUpdate, notInNotion, err := notionClient.CheckExistingEvents(notionToday, calToday)
	if err != nil {
		log.Fatal(err)
	}

	// Instantiate Two Lists, Each for Updating Events in Respective Services
	var notionNeedsUpdate []google.EventItem
	var notionNeedsUpdateId []string

	var gcalNeedsUpdate []notion.DatabaseQueryResponse
	var gcalNeedsUpdateId []string
	// Get matching event IDs
	for _, check := range checkNeedsUpdate {
		var idx int
		for i, item := range calToday {
			if item.ID == check {
				idx = i
			}
		}

		// Retrieve Last Edited Time for Given Event in G-Cal
		currEventGCal := calToday[idx]
		currEventGCalT, _ := time.Parse(time.RFC3339, currEventGCal.EditedDate)
		d := 60 * time.Second
		currEventGCalTime := currEventGCalT.Truncate(d)

		// Retrieve Last Edited Time for Given Event in Notion
		notionUpdate, _ := client.QueryDatabase(ctx, notionDB, &notion.DatabaseQuery{
			Filter: &notion.DatabaseQueryFilter{
				Property: "GCal_ID",
				Text: &notion.TextDatabaseQueryFilter{
					Equals: check,
				},
			},
		})
		currEventNotion := notionUpdate.Results[0]
		currEventNotionTime, _ := time.Parse(time.RFC3339, currEventNotion.Properties.(notion.DatabasePageProperties)["Edited"].LastEditedTime.Format(time.RFC3339))

		// Compare Last Edited Times and Append Event to list depending on which Service Needs to Update
		compareCase1 := currEventGCalTime.Before(currEventNotionTime)
		compareCase2 := currEventGCalTime.After(currEventNotionTime)
		if compareCase1 {
			// If Notion's Last Edited Time is more Recent, Update in G-Cal
			gcalNeedsUpdate = append(gcalNeedsUpdate, notionUpdate)
			gcalNeedsUpdateId = append(gcalNeedsUpdateId, calToday[idx].ID)
		} else if compareCase2 {
			// If G-Cal's Last Edited Time is more Recent, Update in Notion
			notionNeedsUpdate = append(notionNeedsUpdate, calToday[idx])
			notionNeedsUpdateId = append(notionNeedsUpdateId, notionUpdate.Results[0].ID)
		} else {
			// If Last Edited Time is same, no need to Update Given Event
			continue
		}
	}

	// Update Event(s) in G-Cal
	err = googleCalClient.UpdateCalendar(calendarID, gcalNeedsUpdate, gcalNeedsUpdateId)
	if err != nil {
		log.Fatal(err)
	}

	// Update Event(s) in Notion
	err = notionClient.UpdateEvents(notionNeedsUpdate, notionNeedsUpdateId)
	if err != nil {
		log.Fatal(err)
	}

	// Create Event(s) in Notion
	err = notionClient.Create(notInNotion)
	if err != nil {
		log.Fatal(err)
	}

	// Create Event(s) in G-Cal
	notInGCal, err := notionClient.CheckNewEvents()
	if err != nil {
		log.Fatal(err)
	}

	err = googleCalClient.CreateEvent(calendarID, notInGCal, client)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sync Complete.")
}
