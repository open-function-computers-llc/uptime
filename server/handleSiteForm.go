package server

import (
	"io"
	"net/http"
)

func (s *Server) handleSiteForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form := `
<h1>Add a URL</h1>
<form action="/store" method="POST">
	<input type="text" name="url" placeholder="https://openfunctioncomputers.com" />
	<input type="submit" value="Add" />
</form>
<a class='button' href="/">Cancel</a>
		`
		io.WriteString(w, htmlWrap(form))
	}
}
