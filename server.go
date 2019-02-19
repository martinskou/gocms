package main

import (
	//    "database/sql"
	//    "go-echo-vue/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io"
	"net/http"
//    "time"
    "fmt"

    //"./src/data"
    "github.com/martinskou/gocms/gdb"

	//    "github.com/labstack/echo/engine/standard"
	//    _ "github.com/mattn/go-sqlite3"
)

const templateGlob = "templates/*.html"
const debug = true

/*type Page struct {
	Title    string
	Name     string
	Children *[]Page
}*/

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	root_page := gdb.ExampleSite()

	// Check if data is a map with string keys (type switch)
	if viewContext, isMap := data.(map[string]interface{}); isMap {
        viewContext["current"] = root_page
        viewContext["pages"] = root_page
        fmt.Printf("%#v\n", data)
	}
	if debug {
		t := template.Must(template.ParseGlob(templateGlob))
		return t.ExecuteTemplate(w, name, data)
	} else {
		return t.templates.ExecuteTemplate(w, name, data)
	}
}




func main() {

	//	db := initDB("storage.db")
	//	migrate(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob(templateGlob)),
	}
	e.Renderer = renderer

    if debug {
        fmt.Printf("Loaded templates:\n")
        for _, t := range renderer.templates.Templates() {
            fmt.Printf("%+v\n", t.Name())
        }
    }

	// Named route "foobar"
	e.GET("/something", func(c echo.Context) error {
		c.Render(http.StatusOK, "demo1.html", map[string]interface{}{
			"name": "Dolly!", "menu": [2]string{"test", "test2"},
		})
        return nil
	}).Name = "foobar"

	e.Static("/", "assets")
	e.File("/favicon.ico", "assets/favicon.ico")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.File("/d/", "templates/demo1.html")
	//	e.File("/", "public/index.html")
	//	e.GET("/tasks", handlers.GetTasks(db))
	//	e.PUT("/tasks", handlers.PutTask(db))
	//	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

	e.Logger.Fatal(e.Start(":8000"))
}
