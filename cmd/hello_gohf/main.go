package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

func main() {
	router := gohf.New()

	router.GET("/greeting", func(c *gohf.Context) gohf.Response {
		name := c.Req.GetQuery("name")
		if name == "" {
			return response.Error(
				http.StatusBadRequest,
				errors.New("Name is required"),
			)
		}

		greeting := fmt.Sprintf("Hello, %s!", name)
		return response.Text(http.StatusOK, greeting)
	})

	router.Use(func(c *gohf.Context) gohf.Response {
		return response.Error(
			http.StatusNotFound,
			errors.New("Page not found"),
		)
	})

	mux := router.CreateServeMux()
	log.Fatal(http.ListenAndServe(":8080", mux))
}
