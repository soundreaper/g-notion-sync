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

	// Notion
	notionKey := config.GetNotionConfig().NotionKey
	notionDB := config.GetNotionConfig().NotionDB
	client := notion.NewClient(notionKey)
	notionClient := no.NewNotionService(client, notionDB)
	srv := google.GetCalendarFromGoogle()

	// Google
	// cal id
	calendarID := "primary"
	// create calendar client
	googleCalClient := google.NewGCalService(srv)

	// get auth user google calendar
	calList, err := googleCalClient.GetCalendars()
	if err != nil {
		log.Fatal(err)
	}

	// filter and sort google calendar events
	calToday, err := googleCalClient.FilterEvents(calList)
	if err != nil {
		log.Fatal(err)
	}

	// get notion events
	notionToday, err := notionClient.GetEvents()
	if err != nil {
		log.Fatal(err)
	}

	// Check existing events
	checkNeedsUpdate, notInNotion, err := notionClient.CheckExistingEvents(notionToday, calToday)
	if err != nil {
		log.Fatal(err)
	}

	var notionNeedsUpdate []google.EventItem
	var notionNeedsUpdateId []string

	var gcalNeedsUpdate []notion.DatabaseQueryResponse
	var gcalNeedsUpdateId []string
	// need to get event from calToday with given gcal_id and compare to event in notion with same gcal_id
	for _, check := range checkNeedsUpdate {
		var idx int
		for i, item := range calToday {
			if item.ID == check {
				idx = i
			}
		}
		currEventGCal := calToday[idx]
		currEventGCalT, _ := time.Parse(time.RFC3339, currEventGCal.EditedDate)
		d := 60 * time.Second
		currEventGCalTime := currEventGCalT.Truncate(d)

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

		compareCase1 := currEventGCalTime.Before(currEventNotionTime)
		compareCase2 := currEventGCalTime.After(currEventNotionTime)
		if compareCase1 {
			gcalNeedsUpdate = append(gcalNeedsUpdate, notionUpdate)
			gcalNeedsUpdateId = append(gcalNeedsUpdateId, calToday[idx].ID)
		} else if compareCase2 {
			notionNeedsUpdate = append(notionNeedsUpdate, calToday[idx])
			notionNeedsUpdateId = append(notionNeedsUpdateId, notionUpdate.Results[0].ID)
		} else {
			continue
		}
	}

	// Update Google Calendar
	err = googleCalClient.UpdateCalendar(calendarID, gcalNeedsUpdate, gcalNeedsUpdateId)
	if err != nil {
		log.Fatal(err)
	}

	// Update Notion Database
	err = notionClient.UpdateEvents(notionNeedsUpdate, notionNeedsUpdateId)
	if err != nil {
		log.Fatal(err)
	}

	err = notionClient.Create(notInNotion)
	if err != nil {
		log.Fatal(err)
	}
	// Check for new events in Notion that are not in GCal
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
