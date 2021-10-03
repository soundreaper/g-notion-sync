package config

import (
	"os"
)

type NotionConfig struct {
	NotionKey string
	NotionDB  string
}

func GetNotionConfig() *NotionConfig {
	return &NotionConfig{
		NotionKey: os.Getenv("NOTION_SECRET"),
		NotionDB:  os.Getenv("NOTION_DB"),
	}
}
