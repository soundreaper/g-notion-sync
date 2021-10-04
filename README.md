# G-Notion-Sync

## What is it?
G-Notion-Sync is a console based app that uses:
 - [Golang](https://golang.org/)
 - [Notion API](https://developers.notion.com/)
 - [Go-Notion](https://github.com/dstotijn/go-notion)
 - [Google O-Auth2](https://developers.google.com/identity)
 - [Google Calendar API](https://developers.google.com/calendar)

Specifically, the application checks an authorized user's Google Calendar for events within the current day and automatically syncs them to a Notion Workspace. It also checks the Notion Workspace for any events in the current day and syncs them to Google Calendar. It can be seen as a functional 2-Way Sync. Any updates made to events in Google Calendar and/or Notion will sync over to the other platform as well.

## How to Setup?
### Step 1:
Clone this repository.
```git
git clone https://github.com/soundreaper/g-notion-sync
```
### Step 2:
- Update ```config/credentials.json``` with [Create Credentials](https://developers.google.com/workspace/guides/create-credentials)
- Update ```.env``` with [Notion Quickstart](https://developers.notion.com/docs)
### Step 3:
- From the Notion Quickstart guide, open the database table that you created and remove all rows and delete the "Tags" column. Add a Date column titled "Date", a Text column titled "GCal_ID", and an advanced Last edited time column titled "Edited". It should look like the following:
![Example Database](https://raw.githubusercontent.com/soundreaper/g-notion-sync/main/assets/Database.PNG)

- Proceed to the Google Calendar for whichever Google Account will be authorized with the application. Add a few test events if there are none in the current day. 
### Step 4:
Go into the folder and run the app.
```bash
cd g-notion-sync

go run main.go
```

## Result
- The Notion Database should be auto-populated with any events from the authorized user's Google Calendar. Try making changes to the events like the start and finish times or the name of an event and run the app again.

- To add an event to Notion, create a new row but only populate the Name and Date fields and run the app. It should automatically appear on Google Calendar and the GCal_ID field will auto-populate. 

## To-Do
- Currently, there is no functionality for deleting events with this app. If an item is removed from one service it will simply be auto-populated from the other service. To remove an event, it must be manually deleted from both services. Adding functionality for this going forward would be useful!

- The app should theoretically handle timezones properly since all dates that are read-in and comparisons between dates are done in UTC. They are then converted to the user's local timezone before being pushed to either service. If there are any issues with this, please open an issue!