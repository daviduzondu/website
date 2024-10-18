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
		} else if filepath.Ext(entry.Name()) == ".md" && !strings.HasPrefix(entry.Name(), "~") && !strings.HasPrefix(entry.Name(), "-") {
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

	if !*matter.Draft {
		page.Frontmatter = matter
		page.Html = template.HTML(string(utils.ConvertToHtml(content)))
		page.Src = mdPath
		page.Dest = htmlPath
		page.Href = filepath.Join(string(filepath.Separator), filepath.Base(filepath.Dir(htmlPath)), filepath.Base(htmlPath))
		page.PageUrlOnGitHub = utils.GetPageUrlOnGitHub(mdPath, siteConfig)
		page.SrcName = strings.Replace(fileName, ".md", "", 1)
		page.Active = filepath.Dir(page.Href)

		if strings.HasPrefix(filepath.Base(filepath.Dir(mdPath)), "_") || page.Series != "" {
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

		if len(matter.Tags) > 0 {
			for _, m := range matter.Tags {
				var tag structs.Tag

				for _, t := range siteConfig.AllTags {
					if t.Name == m {
						tag = t
					}
				}
				dest := filepath.Join(utils.First(os.Getwd()), "dist", "tags", m+".html")
				utils.EnsureDirExists(filepath.Dir(dest))
				utils.CheckErr(err)

				tag.Name = m
				tag.Pages = append(tag.Pages, page)
				tag.Dest = dest
				tag.Title = string("Posts tagged ") + string('"') + tag.Name + string('"')
				tag.Href = filepath.Join(string(filepath.Separator), filepath.Base(filepath.Dir(dest)), filepath.Base(dest))
				page.Tags = append(page.Tags, tag)
				siteConfig.AllTags = append(siteConfig.AllTags, tag)
			}
		}

		siteConfig.AllPages = append(siteConfig.AllPages, page)

		utils.EnsureDirExists(filepath.Dir(htmlPath))
		utils.CheckErr(err)
	}
}

func getFrontmatter(data io.Reader) ([]byte, structs.Matter) {
	var matter structs.Matter
	content, err := frontmatter.Parse(data, &matter)
	utils.CheckErr(err)

	if matter.Date != "" {
		utils.ValidateDate(matter.Date)
	}

	if matter.Draft == nil {
		defaultDraftValue := true
		matter.Draft = &defaultDraftValue
	}

	return content, matter
}
