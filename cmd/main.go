package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/daviduzondu/website/internal/config"
	"github.com/daviduzondu/website/internal/content"
	"github.com/daviduzondu/website/internal/utils"
)

var BasePath = utils.First(os.Getwd())
var OutputPath = filepath.Join(BasePath, "dist")

func main() {
	siteData := config.LoadConfig(filepath.Join(BasePath, "config.json"))
	utils.EnsureDirExists(OutputPath)
	os.RemoveAll(OutputPath)

	contentPath := filepath.Join(BasePath, "www", string(filepath.Separator), "content")
	content.Traverse(contentPath, &siteData)
	siteData.Year = time.Now().Year()
	content.ApplyTemplate(&siteData, BasePath, OutputPath)

}
