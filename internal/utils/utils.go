package utils

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/daviduzondu/website/internal/structs"
	figure "github.com/mangoumbrella/goldmark-figure"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
)

func ValidateDate(dateStr string) bool {
	_, err := time.Parse("2 January 2006", dateStr)
	CheckErr(err)
	return true
}

func copyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	CheckErr(err)
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	CheckErr(err)
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	CheckErr(err)

	return nil
}

// Recursively copies the contents of the source directory onto the destination directory.
func CopyDir(srcDir string, dstDir string) error {
	EnsureDirExists(dstDir)

	entries, err := os.ReadDir(srcDir)
	CheckErr(err)

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			err := CopyDir(srcPath, dstPath)
			CheckErr(err)
		} else {
			err := copyFile(srcPath, dstPath)
			CheckErr(err)
		}
	}

	return nil
}

func EnsureDirExists(dir string) {
	_, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(dir, os.ModePerm)
		CheckErr(err)
	}
}

var ErrMessage = ""

func CheckErr(err error) {
	if err != nil {
		if ErrMessage != "" {
			fmt.Println(ErrMessage)
		}
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

func extractProps(propString string) map[string]string {
	props := make(map[string]string)
	// Regex to extract key="value" pairs
	re := regexp.MustCompile(`(\w+)="([^"]*)"`)

	matches := re.FindAllStringSubmatch(propString, -1)
	for _, match := range matches {
		props[match[1]] = match[2]
	}

	return props
}

func ConvertToHtml(content []byte) []byte {
	var buf bytes.Buffer
	err := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()), goldmark.WithRendererOptions(html.WithUnsafe()), goldmark.WithExtensions(
		extension.GFM,
		figure.Figure,
		&anchor.Extender{Position: anchor.After},
		extension.Footnote,
		highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
		),
	)).Convert(content, &buf)

	if err != nil {
		log.Fatal("Something went wrong while trying to parse Markdown.", err)
	}

	// Load the snippets file

	htmlContent := buf.String()

	snippetRegex := regexp.MustCompile(`<!--\s*([\w-]+)\s*(.*?)-->`)

	tmpl, err := template.ParseFiles(filepath.Join(First(os.Getwd()), "templates", "snippets", "snippets.html"))
	CheckErr(err)

	finalHtml := snippetRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		// Extract snippet name and props
		snippetMatch := snippetRegex.FindStringSubmatch(match)
		snippetName := snippetMatch[1]
		propString := snippetMatch[2]
		props := extractProps(propString)

		// Use template to render the snippet's HTML
		var snippetBuf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&snippetBuf, snippetName, props); err != nil {
			log.Printf("Error rendering snippet %s: %v", snippetName, err)
			return match // Return the original comment if rendering fails
		}

		return snippetBuf.String() // Return the rendered HTML
	})

	return []byte(finalHtml)
}
