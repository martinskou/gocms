package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	log "github.com/sirupsen/logrus"
)


func DNS_register_extension(e *echo.Echo, cms CMS, renderer *TemplateRenderer) {

	log.Printf("Register extension for DNS")

	e.GET("/ext/dns", func(c echo.Context) error {
		
		return c.HTMLBlob(http.StatusOK, []byte("dns extension"))
	})

	
}
