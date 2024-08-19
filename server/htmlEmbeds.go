package server

import (
	_ "embed"
	"strings"
)

//go:embed pagewrap.html
var pageTemplate string

//go:embed style.css
var pageCSS string

//go:embed addForm.html
var addForm string

func htmlWrap(html string) string {
	output := strings.ReplaceAll(pageTemplate, "%%CSS%%", "<style>"+pageCSS+"</style>")
	return strings.ReplaceAll(output, "%%PAGECONTENT%%", html)
}
