package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/daviduzondu/website/internal/structs"
)

func LoadConfig(filePath string) structs.SiteData {
	var siteConfig structs.SiteData
	f, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Failed to retrieve configuration file")
	}
	json.Unmarshal(f, &siteConfig)
	return siteConfig
}
