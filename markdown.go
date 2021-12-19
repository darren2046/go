package golanglibs

import (
	"bytes"

	"github.com/alecthomas/chroma"
	htmlChromaFormatter "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday"
)

func md2html(md string) string {
	flags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_TOC

	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	htmlContent := string(blackfriday.MarkdownOptions([]byte(md), blackfriday.HtmlRenderer(flags, "", ""), blackfriday.Options{Extensions: extensions}))

	insideCode := false
	htmlCodeContent := ""
	lang := ""
	code := ""
	for _, line := range String(htmlContent).Splitlines() {
		if line.StartsWith("<pre><code") {
			insideCode = true
			if !String("<pre><code>").In(line.S) {
				res := Re.FindAll("<pre><code class=\"language-(.+?)\">", line.S)
				if len(res) != 0 {
					line = line.Replace(res[0][0], "")
					lang = res[0][1]
				}
			} else {
				line = line.Replace("<pre><code>", "")
			}
		}
		if line.StartsWith("</code></pre>") {
			insideCode = false
			codeHTML := getHightLightHTML(Html.Decode(code), lang)
			htmlCodeContent += codeHTML
			code = ""
			lang = ""
		}
		if !insideCode && !line.StartsWith("</code></pre>") {
			// res := reFindAll("\\[(.+?)\\]\\((.+?)\\)", line)
			// if len(res) != 0 && !strIn("!"+res[0][0], line) {
			// 	links = append(links, res[0])
			// 	line = strReplace(line, res[0][0], "<a href=\""+res[0][2]+"\">"+res[0][1]+"</a>")
			// }
			htmlCodeContent += line.S + "\n"
		}
		if insideCode {
			code += line.S + "\n"
		}
	}

	return htmlCodeContent
}

func getHightLightHTML(code string, codeType ...string) string {
	style := "emacs"

	var l chroma.Lexer
	var displayLineNumber bool
	if len(codeType) != 0 && codeType[0] != "" {
		l = lexers.Get(codeType[0])
		if l == nil {
			l = lexers.Fallback
		} else {
			displayLineNumber = true
		}
	} else {
		l = lexers.Fallback
	}

	l = chroma.Coalesce(l)

	f := htmlChromaFormatter.New(htmlChromaFormatter.Standalone(false), htmlChromaFormatter.WithLineNumbers(displayLineNumber), htmlChromaFormatter.TabWidth(2))
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}
	it, err := l.Tokenise(nil, code)
	Panicerr(err)

	var buf bytes.Buffer
	err = f.Format(&buf, s, it)
	Panicerr(err)
	return buf.String()
}
