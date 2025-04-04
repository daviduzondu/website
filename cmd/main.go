package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/daviduzondu/website/internal/config"
	"github.com/daviduzondu/website/internal/content"
	"github.com/daviduzondu/website/internal/utils"
)

var BasePath = utils.First(os.Getwd())
var OutputPath = filepath.Join(BasePath, "dist")
var Version = "1.0"

func main() {
	var configPath = os.Getenv("SITE_CONFIG")
	fmt.Println("StaticBuilder Version: ", Version)
	fmt.Println("Loading config from", configPath)
	path := configPath
	if path == "" {
		path = filepath.Join(BasePath, "config.json")
	}

	siteData := config.LoadConfig(path)

	utils.EnsureDirExists(OutputPath)
	os.RemoveAll(OutputPath)

	contentPath := filepath.Join(BasePath, "www", string(filepath.Separator), "content")
	content.Traverse(contentPath, &siteData)
	siteData.Year = time.Now().Year()
	content.ApplyTemplate(&siteData, BasePath, OutputPath)
	utils.CopyDir(filepath.Join(BasePath, "www", "assets"), filepath.Join(OutputPath, "assets"))
	utils.CopyDir(filepath.Join(BasePath, "www", "public"), filepath.Join(OutputPath, "public"))
}
