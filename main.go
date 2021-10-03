package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/soundreaper/g-notion-sync/config"
	"github.com/soundreaper/g-notion-sync/util"
	"google.golang.org/api/calendar/v3"
)

type calendarType struct {
	c_id    string
	summary string
	events  []eventItem
}

type eventItem struct {
	id         string
	summary    string
	startDate  string
	endDate    string
	editedDate string
}

func main() {
	ctx := context.Background()
	notion_key := config.GetNotionConfig().NotionKey
	notion_db := config.GetNotionConfig().NotionDB
	srv := util.GetCalendarFromGoogle()

	c_list, err := srv.CalendarList.List().Fields("items/id").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of calendars: %v", err)
	}

	calendars := []calendarType{}
	for _, v := range c_list.Items {
		calendars = append(calendars, calendarType{v.Id, v.Summary, []eventItem{}})
	}

	for ind, calendar := range calendars {
		t1 := time.Now().Format(time.RFC3339)
		t2 := time.Now().AddDate(0, 0, 5).Format(time.RFC3339)

		events, _ := srv.Events.List(calendar.c_id).ShowDeleted(false).SingleEvents(true).TimeMin(t1).TimeMax(t2).MaxResults(15).OrderBy("startTime").Do()

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
				calendars[ind].events = append(calendars[ind].events, eventItem{item.Id, item.Summary, startDate, endDate, editedDate})
			}
		}
	}

	client := notion.NewClient(notion_key)

	calToday := []eventItem{}
	for _, cal := range calendars {
		if len(cal.events) != 0 {
			for _, item := range cal.events {
				t, _ := time.Parse(time.RFC3339, item.startDate)
				// adding 1 day to the current time so we can sync the events for the next day, change as needed
				if t.Format("2006-01-02") == time.Now().AddDate(0, 0, 0).Format("2006-01-02") {
					calToday = append(calToday, item)
				}
			}
		}
	}

	t, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	prop := []notion.DatabaseQueryFilter{
		{
			Property: "Date",
			Date: &notion.DateDatabaseQueryFilter{
				OnOrAfter: &t,
			},
		},
		{
			Property: "GCal_ID",
			Text: &notion.TextDatabaseQueryFilter{
				IsNotEmpty: true,
			},
		},
	}
	notionToday, err := client.QueryDatabase(ctx, notion_db, &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			And: prop,
		},
	})
	if err != nil {
		log.Fatalf("Unable to query database: %v", err)
	}

	var gcal_id_list []string
	var check_needs_update []string
	not_in_notion := []eventItem{}
	// We check here if there are any already existing events in Notion
	if len(notionToday.Results) != 0 {
		for _, u := range notionToday.Results {
			// The database query has properties if Notion is not COMPLETELY empty
			if val, ok := u.Properties.(notion.DatabasePageProperties); ok {
				gcal_id_list = append(gcal_id_list, val["GCal_ID"].RichText[0].Text.Content)
			}
		}

		for _, g_cal := range calToday {
			_, found := util.Find(gcal_id_list, g_cal.id)
			if !found {
				// Check which events aren't in Notion and prepare a list for creating them
				not_in_notion = append(not_in_notion, g_cal)
			} else {
				// Check which events are in Notion and prepare a list of ID's for potential updating
				check_needs_update = append(check_needs_update, g_cal.id)
			}
		}
		// If the returned result list was empty, Notion is missing events from Google Calendar
	} else {
		not_in_notion = append(not_in_notion, calToday...)
	}

	notion_needs_update := []eventItem{}
	var notion_needs_update_ID []string

	gcal_needs_update := []notion.DatabaseQueryResponse{}
	var gcal_needs_update_ID []string
	// need to get event from calToday with given gcal_id and compare to event in notion with same gcal_id
	for _, check := range check_needs_update {
		var idx int
		for i, item := range calToday {
			if item.id == check {
				idx = i
			}
		}
		currEventGCal := calToday[idx]
		currEventGCalT, _ := time.Parse(time.RFC3339, currEventGCal.editedDate)
		d := (60 * time.Second)
		currEventGCalTime := currEventGCalT.Truncate(d)

		notionUpdate, _ := client.QueryDatabase(ctx, notion_db, &notion.DatabaseQuery{
			Filter: &notion.DatabaseQueryFilter{
				Property: "GCal_ID",
				Text: &notion.TextDatabaseQueryFilter{
					Equals: check,
				},
			},
		})
		currEventNotion := notionUpdate.Results[0]
		currEventNotionTime, _ := time.Parse(time.RFC3339, currEventNotion.Properties.(notion.DatabasePageProperties)["Edited"].LastEditedTime.Format(time.RFC3339))

		compare_case_1 := currEventGCalTime.Before(currEventNotionTime)
		compare_case_2 := currEventGCalTime.After(currEventNotionTime)
		if compare_case_1 {
			gcal_needs_update = append(gcal_needs_update, notionUpdate)
			gcal_needs_update_ID = append(gcal_needs_update_ID, calToday[idx].id)
		} else if compare_case_2 {
			notion_needs_update = append(notion_needs_update, calToday[idx])
			notion_needs_update_ID = append(notion_needs_update_ID, notionUpdate.Results[0].ID)
		} else {
			continue
		}
	}

	calendarID := "primary"
	for idx, update := range gcal_needs_update {
		summary := update.Results[0].Properties.(notion.DatabasePageProperties)["Name"].Title[0].Text.Content
		startTime := update.Results[0].Properties.(notion.DatabasePageProperties)["Date"].Date.Start.Time.Format(time.RFC3339)
		endTime := update.Results[0].Properties.(notion.DatabasePageProperties)["Date"].Date.End.Time.Format(time.RFC3339)

		event := &calendar.Event{
			Summary: summary,
			Start: &calendar.EventDateTime{
				DateTime: startTime,
				TimeZone: startTime,
			},
			End: &calendar.EventDateTime{
				DateTime: endTime,
				TimeZone: endTime,
			},
		}

		_, err := srv.Events.Update(calendarID, gcal_needs_update_ID[idx], event).Do()
		if err != nil {
			log.Fatalf("Unable to update event. %v\n", err)
		}
	}

	for idx, update := range notion_needs_update {
		start, _ := time.Parse(time.RFC3339, update.startDate)
		end, _ := time.Parse(time.RFC3339, update.endDate)
		ender := notion.NewDateTime(end, true)

		client.UpdatePageProps(ctx, notion_needs_update_ID[idx], notion.UpdatePageParams{
			DatabasePageProperties: &notion.DatabasePageProperties{
				"Name": notion.DatabasePageProperty{
					Title: []notion.RichText{
						{
							Type: notion.RichTextTypeText,
							Text: &notion.Text{
								Content: update.summary,
							},
						},
					},
				},
				"Date": notion.DatabasePageProperty{
					Date: &notion.Date{
						Start: notion.NewDateTime(start, true),
						End:   &ender,
					},
				},
			},
		})
	}

	for _, create := range not_in_notion {
		start, _ := time.Parse(time.RFC3339, create.startDate)
		end, _ := time.Parse(time.RFC3339, create.endDate)
		ender := notion.NewDateTime(end, true)

		client.CreatePage(ctx, notion.CreatePageParams{
			ParentType: notion.ParentTypeDatabase,
			ParentID:   notion_db,

			DatabasePageProperties: &notion.DatabasePageProperties{
				"Name": notion.DatabasePageProperty{
					Title: []notion.RichText{
						{
							Type: notion.RichTextTypeText,
							Text: &notion.Text{
								Content: create.summary,
							},
						},
					},
				},
				"Date": notion.DatabasePageProperty{
					Date: &notion.Date{
						Start: notion.NewDateTime(start, true),
						End:   &ender,
					},
				},
				"GCal_ID": notion.DatabasePageProperty{
					RichText: []notion.RichText{
						{
							Type: notion.RichTextTypeText,
							Text: &notion.Text{
								Content: create.id,
							},
						},
					},
				},
			},
		})
	}

	notInGCal, err := client.QueryDatabase(ctx, notion_db, &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "GCal_ID",
			Text: &notion.TextDatabaseQueryFilter{
				IsEmpty: true,
			},
		},
	})
	if err != nil {
		log.Fatalf("Unable to query database: %v", err)
	}

	for _, u := range notInGCal.Results {
		if val, ok := u.Properties.(notion.DatabasePageProperties); ok {
			summary := val["Name"].Title[0].Text.Content
			startTime := val["Date"].Date.Start.Time.Format(time.RFC3339)
			endTime := val["Date"].Date.End.Time.Format(time.RFC3339)

			event := &calendar.Event{
				Summary: summary,
				Start: &calendar.EventDateTime{
					DateTime: startTime,
					TimeZone: startTime,
				},
				End: &calendar.EventDateTime{
					DateTime: endTime,
					TimeZone: endTime,
				},
			}

			event, err := srv.Events.Insert(calendarID, event).Do()
			if err != nil {
				log.Fatalf("Unable to create event. %v\n", err)
			}

			notionUpdateID, _ := client.QueryDatabase(ctx, notion_db, &notion.DatabaseQuery{
				Filter: &notion.DatabaseQueryFilter{
					Property: "Name",
					Text: &notion.TextDatabaseQueryFilter{
						Equals: event.Summary,
					},
				},
			})

			client.UpdatePageProps(ctx, notionUpdateID.Results[0].ID, notion.UpdatePageParams{
				DatabasePageProperties: &notion.DatabasePageProperties{
					"GCal_ID": notion.DatabasePageProperty{
						RichText: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: event.Id,
								},
							},
						},
					},
				},
			})
		}
	}
	fmt.Println("Sync Complete.")
}
