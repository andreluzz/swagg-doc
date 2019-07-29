package main

// swagg-doc:api
// openapi: 3.0.0
// info:
//   title: Sandbox API
//   description: Optional multiline or single-line description in [CommonMark](http://commonmark.org/help/) or HTML.
//   version: 0.0.1
// servers:
//   - url: http://localhost:8080/api/v1
//     description: Optional server description, e.g. Development server for testing and development
//
// parameters:
//   language-code:
//     name: language-code
//     in: header
//     description: Translation fields language code
//     required: true
//     schema:
//       type: string
//       example: pt-br

import (
	"net/http"

	"github.com/andreluzz/swagg-doc/mock/api/controller"
	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()

	r.Get("/users", controller.GetUsers)
	http.ListenAndServe(":1234", r)
}
