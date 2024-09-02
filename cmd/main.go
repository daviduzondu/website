package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/daviduzondu/website/internal/config"
	"github.com/daviduzondu/website/internal/content"
	"github.com/daviduzondu/website/internal/utils"
)

func main() {

	basePath := utils.First(os.Getwd())
	outputPath := filepath.Join(basePath, "dist")

	siteData := config.LoadConfig(filepath.Join(basePath, "config.json"))
	utils.EnsureDirExists(outputPath)
	os.RemoveAll(outputPath)

	contentPath := filepath.Join(basePath, "/www/content/")
	content.Traverse(contentPath, &siteData)
	siteData.Year = time.Now().Year()
	content.ApplyTemplate(&siteData, basePath)
}
