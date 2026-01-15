package pkg

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/russross/blackfriday"
)

func MarkdownToHTML(content string) string {
	htmlContent := blackfriday.MarkdownCommon([]byte(content))
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(htmlContent)))
	doc.Find("script").Remove()
	doc.Find("iframe").Remove()
	doc.Find("object").Remove()
	safeHTML, _ := doc.Html()
	return safeHTML
}
