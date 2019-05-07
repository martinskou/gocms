package main

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)


func register_extensions(e *echo.Echo, cms CMS, renderer *TemplateRenderer) {

	log.Printf("Register extensions")

	DNS_register_extension(e,cms,renderer)
	
}
