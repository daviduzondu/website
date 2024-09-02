package content

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/daviduzondu/website/internal/structs"
	"github.com/daviduzondu/website/internal/utils"
)

func ApplyTemplate(siteData *structs.SiteData, basePath string) {
	for _, page := range siteData.AllPages {
		var buf bytes.Buffer
		var data struct {
			Page structs.Page
			Site structs.SiteData
		}
		data.Page = page
		data.Site = *siteData

		tmpl, err := template.ParseFiles(filepath.Join(basePath, "templates", "base.html"), filepath.Join(basePath, "templates", "partials", "nav.html"), filepath.Join(basePath, "templates", "article.html"), filepath.Join(basePath, "templates", "partials", "github-page.html"))
		utils.CheckErr(err)
		err = tmpl.Execute(&buf, data)
		utils.CheckErr(err)
		err = os.WriteFile(page.Dest, buf.Bytes(), os.ModePerm)
		utils.CheckErr(err)
	}

	for _, list := range siteData.AllLists {
		var buf bytes.Buffer
		var data struct {
			List structs.List
			Site structs.SiteData
		}
		data.List = list
		data.Site = *siteData

		tmpl, err := template.ParseFiles(filepath.Join(basePath, "templates", "base.html"), filepath.Join(basePath, "templates", "partials", "nav.html"), filepath.Join(basePath, "templates", "list.html"))
		utils.CheckErr(err)
		err = tmpl.Execute(&buf, data)
		utils.CheckErr(err)
		err = os.WriteFile(list.Dest, buf.Bytes(), os.ModePerm)
		utils.CheckErr(err)
	}
}
