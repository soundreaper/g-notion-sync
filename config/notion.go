package config

import (
	"os"
)

// Notion Secret Info
type NotionConfig struct {
	NotionKey string
	NotionDB  string
}

// Return Notion Secret Info
func GetNotionConfig() *NotionConfig {
	return &NotionConfig{
		NotionKey: os.Getenv("NOTION_SECRET"),
		NotionDB:  os.Getenv("NOTION_DB"),
	}
}
