package content

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/daviduzondu/website/internal/structs"
	"github.com/daviduzondu/website/internal/utils"
)

func Traverse(dir string, siteData *structs.SiteData) {
	parent := filepath.Dir(dir)
	entries, err := os.ReadDir(dir)
	utils.CheckErr(err)

	for _, entry := range entries {
		if entry.IsDir() {
			Traverse(filepath.Join(dir, entry.Name()), siteData)
		} else if filepath.Ext(entry.Name()) == ".md" && !strings.HasPrefix(entry.Name(), "~") {
			processMarkdown(dir, parent, entry.Name(), siteData)
		}
	}
}

func processMarkdown(dir, parent, fileName string, siteConfig *structs.SiteData) {
	var page structs.Page
	var list structs.List

	mdPath := filepath.Join(dir, fileName)
	htmlPath := utils.FormatMdOutputPath(dir, fileName)

	data, err := os.ReadFile(mdPath)
	utils.CheckErr(err)

	if strings.HasPrefix(filepath.Base(filepath.Dir(mdPath)), "(") && strings.HasSuffix(filepath.Base(filepath.Dir(mdPath)), ")") {
		page.Series = strings.TrimSuffix(strings.TrimPrefix(filepath.Base(filepath.Dir(mdPath)), "("), ")")
		htmlPath = utils.FormatMdOutputPath(parent, fileName)
	}

	content, matter := getFrontmatter(strings.NewReader(string(data)))
	page.Frontmatter = matter
	page.Html = template.HTML(string(utils.ConvertToHtml(content)))
	page.Src = mdPath
	page.Dest = htmlPath
	page.Href = filepath.Join("/", filepath.Base(filepath.Dir(htmlPath)), filepath.Base(htmlPath))
	page.PageUrlOnGitHub = utils.GetPageUrlOnGitHub(mdPath, siteConfig)
	page.SrcName = fileName

	if strings.HasPrefix(filepath.Base(filepath.Dir(mdPath)), "_") {
		var dest string = filepath.Join(strings.Replace(filepath.Dir(htmlPath), fileName, "", 1), "index.html")
		for _, l := range siteConfig.AllLists {
			if l.Dest == dest {
				list = l
			}
		}

		list.Pages = append(list.Pages, page)
		list.Name = strings.ToUpper(string(filepath.Base(filepath.Dir(htmlPath))[0])) + filepath.Base(filepath.Dir(htmlPath))[1:]
		list.Dest = dest
		siteConfig.AllLists = append(siteConfig.AllLists, list)
	}

	siteConfig.AllPages = append(siteConfig.AllPages, page)

	utils.EnsureDirExists(filepath.Dir(htmlPath))
	utils.CheckErr(err)
}

func getFrontmatter(data io.Reader) ([]byte, structs.Matter) {
	var matter structs.Matter
	content, err := frontmatter.Parse(data, &matter)
	utils.CheckErr(err)
	return content, matter
}
