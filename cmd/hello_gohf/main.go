package main

import (
	"errors"
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

		return response.JSON(http.StatusOK, map[string]string{
			"Hello": name,
		})
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
