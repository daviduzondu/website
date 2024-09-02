package structs

import (
	"html/template"
	"time"
)

type SiteData struct {
	Title             string `json:"title"`
	BaseUrl           string `json:"base_url"`
	GitHubRepository  string `json:"gh_repository"`
	GitHubContentPath string `json:"gh_content_path"`
	LastBuild         time.Time
	BuildTime         string
	Year              int
	AllPages          []Page
	AllLists          []List
}

type Page struct {
	Frontmatter     Matter
	Html            template.HTML
	Series          any
	Src             string
	Dest            string
	Href            string
	PageUrlOnGitHub string
	SrcName         string
}

type List struct {
	Pages       []Page
	Name        string
	Description string
	Dest        string
}

type Matter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Date        string   `yaml:"date"`
}
