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
	Description       string `json:"description"`
	Favicon           string `json:"favicon"`
	LastBuild         time.Time
	BuildTime         string
	Year              int
	AllPages          []Page
	AllLists          []List
	AllTags           []Tag
}

type Tag struct {
	Pages []Page
	Name  string
	Dest  string
	Title string
	Href  string
}

type Page struct {
	Frontmatter     Matter
	Html            template.HTML
	Series          string
	Tags            []Tag
	Active          string
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
	Draft       *bool    `yaml:"draft"`
	Date        string   `yaml:"date"`
}
