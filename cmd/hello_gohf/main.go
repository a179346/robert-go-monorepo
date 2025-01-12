package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gohf-http/gohf/v4"
	"github.com/gohf-http/gohf/v4/gohf_responses"
)

func main() {
	router := gohf.New()

	router.Handle("GET /greeting", func(c *gohf.Context) gohf.Response {
		name := c.Req.GetQuery("name")
		if name == "" {
			return gohf_responses.NewErrorResponse(
				http.StatusBadRequest,
				errors.New("Name is required"),
			)
		}

		greeting := fmt.Sprintf("Hello, %s!", name)
		return gohf_responses.NewTextResponse(http.StatusOK, greeting)
	})

	router.Use(func(c *gohf.Context) gohf.Response {
		return gohf_responses.NewErrorResponse(
			http.StatusNotFound,
			errors.New("Page not found"),
		)
	})

	mux := router.CreateServeMux()
	log.Fatal(http.ListenAndServe(":8080", mux))
}
