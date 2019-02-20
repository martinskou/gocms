package gserver

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

//	root_page := gdb.ExampleSite()

	// Check if data is a map with string keys (type switch)
	if viewContext, isMap := data.(map[string]interface{}); isMap {
//        viewContext["current"] = root_page
//        viewContext["pages"] = root_page
        fmt.Printf("%#v\n", viewContext)
	}
	if debug {
		t := template.Must(template.ParseGlob(templateGlob))
		return t.ExecuteTemplate(w, name, data)
	} else {
		return t.templates.ExecuteTemplate(w, name, data)
	}
}




func Run() {

	//	db := initDB("storage.db")
	//	migrate(db)

	e := echo.New()
	//e.Debug=true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob(templateGlob)),
	}
	e.Renderer = renderer

	/*
    if debug { // Print templates
        for _, t := range renderer.templates.Templates() {
            fmt.Printf("%+v\n", t.Name())
        }
    }
	*/

	root_page := gdb.RandomSite()

	// Attach all pages to Echos router
	fmt.Println("Attaching page routes")
	root_page.Apply(func(p gdb.Page) {
		fmt.Printf("%30s %10s %30s\n", p.Name, p.Slug, p.AbsSlug())
		e.GET(p.AbsSlug(), func(c echo.Context) error {
			c.Render(http.StatusOK, "standard_page.html", map[string]interface{}{"pages":root_page, "current": p})
			return nil
		})
	})

	/*
	e.GET("/something", func(c echo.Context) error {
		c.Render(http.StatusOK, "standard_page.html", map[string]interface{}{
			"pages":root_page})
        return nil
	}).Name = "foobar"

	e.GET("/z*", func(c echo.Context) error {
		fmt.Printf("%+v\n", c)
		return c.String(http.StatusOK, "/users/1/files/*")
	})
	*/
	e.File("/favicon.ico", "assets/img/favicon.ico")
	e.Static("/assets", "assets")
	//	e.File("/", "public/index.html")
	//	e.GET("/tasks", handlers.GetTasks(db))
	//	e.PUT("/tasks", handlers.PutTask(db))
	//	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

	e.Logger.Fatal(e.Start(":8000"))
}

func Test() {
    //run()
    //gdb.TestRnd()
    root_page := gdb.RandomSite()
    root_page.Print(0,0)
	//fmt.Printf("%+v\n", root_page)
}

func main() {
    Run()
    //for i := 0; i < 10; i++ {
    //    test()
    //}
}
