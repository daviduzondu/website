package utils

import (
	"bytes"
	"errors"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/daviduzondu/website/internal/structs"
	figure "github.com/mangoumbrella/goldmark-figure"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/anchor"
)

func EnsureDirExists(dir string) {
	_, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(dir, os.ModePerm)
		CheckErr(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func First[T, U any](val T, _ U) T {
	return val
}

func FormatMdOutputPath(dir string, fileName string) string {
	pathComponents := strings.Split(dir, string(filepath.Separator))

	for i, component := range pathComponents {
		if strings.HasPrefix(component, "_") {
			pathComponents[i] = strings.TrimPrefix(component, "_")
		}
	}

	trimmedDir := filepath.Join(pathComponents...)
	trimmedDir = strings.Replace(trimmedDir, filepath.Join("www", "content"), "dist", 1)
	htmlFileName := strings.Replace(fileName, ".md", ".html", 1)

	return filepath.Join(string(filepath.Separator), trimmedDir, htmlFileName)
}

func GetPageUrlOnGitHub(dir string, siteData *structs.SiteData) string {
	_, contentBase, _ := strings.Cut(dir, filepath.Join("www", string(filepath.Separator)))
	url, err := url.JoinPath(siteData.GitHubContentPath, contentBase)
	CheckErr(err)
	return url
}

func ConvertToHtml(content []byte) []byte {
	var buf bytes.Buffer
	err := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()), goldmark.WithExtensions(
		extension.GFM,
		figure.Figure,
		&anchor.Extender{Position: anchor.Before},
		extension.Footnote,
		highlighting.NewHighlighting(
			highlighting.WithStyle("gruvbox"),
		),
	)).Convert(content, &buf)

	if err != nil {
		log.Fatal("Something went wrong while trying to parse Markdown.", err)
	}
	return buf.Bytes()
}
