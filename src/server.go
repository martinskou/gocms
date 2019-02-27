package main

import (
	//    "database/sql"
	//    "go-echo-vue/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"path/filepath"
	"net/http"
//    "time"
    "fmt"
//	"time"
	"bytes"
	"errors"
    //"./src/data"
    //"github.com/martinskou/gocms/gdb"

//	"github.com/gorilla/sessions"
	
  "github.com/gorilla/sessions"
  "github.com/labstack/echo-contrib/session"
	//    "github.com/labstack/echo/engine/standard"
	//    _ "github.com/mattn/go-sqlite3"
)

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
//	if debug {
//		t := template.Must(template.ParseGlob(templateGlob))
//		return t.ExecuteTemplate(w, name, data)
//	} else {
		return t.templates.ExecuteTemplate(w, name, data)
//	}
}

func UserSession(email string, c echo.Context) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   180,  // 86400 * 7,
		HttpOnly: true,
	}
	sess.Values["email"] = c.FormValue("email")
	sess.Save(c.Request(), c.Response())
}

/*
type session struct {
	Timestamp time.Time
	Data       interface{} // almost everything
}

var sessions map[string]session
*/

func RunServer(base_path string, config_path string) {

	
	//	db := initDB("storage.db")
	//	migrate(db)

	//	t:=template.Must(template.ParseGlob("../cms/templates/*.html"))
/*	files1,_:=filepath.Glob("../cms/templates/*.html")
	files2,_:=filepath.Glob("../themes/alfa/templates/*.html")
	fmt.Printf("%+v\n", files1)
	fmt.Printf("%+v\n", files2)
	fmt.Printf("%+v\n", append(files1, files2...))*/
//	t,_:=template.ParseGlob("../cms/templates/*.html")
//	fmt.Printf("%+v\n", t)
//	return

	fmt.Printf("%+v\n", base_path)
	
	e := echo.New()
	e.Debug=true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	

	cmsTemplateGlob := filepath.Join(base_path, "cms/templates/*.html")
	siteTemplateGlob := filepath.Join(base_path, "themes/alfa/templates/*.html")

	cms_renderer := &TemplateRenderer{templates: template.Must(template.ParseGlob(cmsTemplateGlob))}
	site_renderer := &TemplateRenderer{templates: template.Must(template.ParseGlob(siteTemplateGlob))}
	
	e.Renderer = cms_renderer

	/*
    if debug { // Print templates
        for _, t := range renderer.templates.Templates() {
            fmt.Printf("%+v\n", t.Name())
        }
    }
	*/

//	abs_config_path := filepath.Join(work_path,config_path)
//	fmt.Println("Config file", abs_config_path)

	config, err := LoadConfig(config_path)
	if err!=nil {
		fmt.Println("Config file could not be loaded")
		return
	}
//	fmt.Printf("%+v\n", config)

	
	root_page := RandomSite()

	// Attach all pages to Echos router
	//fmt.Println("Attaching page routes")
	root_page.Apply(func(p Page) {
		//fmt.Printf("%30s %10s %30s\n", p.Name, p.Slug, p.AbsSlug())
		e.GET(p.AbsSlug(), func(c echo.Context) error {

//	cms_renderer = &TemplateRenderer{templates: template.Must(template.ParseGlob(cmsTemplateGlob))}
	site_renderer = &TemplateRenderer{templates: template.Must(template.ParseGlob(siteTemplateGlob))}
			
			data:=map[string]interface{}{"config":config, "pages":root_page, "current": p}
			
			buf := new(bytes.Buffer)
			if err = site_renderer.Render(buf, "page.html", data, c); err != nil {
				return errors.New("Render of template failed")
			}
			return c.HTMLBlob(http.StatusOK, buf.Bytes())
			
//			c.Render(http.StatusOK, "standard_page.html", data)
//			return nil
		})
	})

	e.GET("/aviva", func(c echo.Context) error {

	cms_renderer = &TemplateRenderer{templates: template.Must(template.ParseGlob(cmsTemplateGlob))}
//	site_renderer = &TemplateRenderer{templates: template.Must(template.ParseGlob(siteTemplateGlob))}

		sess, _ := session.Get("session", c)
		fmt.Println(sess.Values["email"])
		
		data:=map[string]interface{}{"config":config, "pages":root_page, "user": sess.Values["email"]}
			
			buf := new(bytes.Buffer)
			if err = cms_renderer.Render(buf, "login.html", data, c); err != nil {
				return errors.New("Render of template failed")
			}
			return c.HTMLBlob(http.StatusOK, buf.Bytes())

		
//		c.Render(http.StatusOK, "login.html", map[string]interface{}{"config":config, "pages":root_page})
		// renderer.render ? 
//		return nil
	})

	e.POST("/aviva/login/authenticate", func(c echo.Context) error {
		var r Reply
		fmt.Printf("%#v\n", c.FormValue("email"))
		if (c.FormValue("email")=="msd@infoserv.dk") {
			fmt.Println("OK\n")
			UserSession(c.FormValue("email"), c)
			r = Reply{Status: "OK", Data: ""}
		} else {
			r = Reply{Status: "FAIL", Data: ""}
		}
		return c.JSON(http.StatusOK,r)
	})

	e.GET("/aviva/logout", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		fmt.Println(sess.Values["email"])
		sess.Values["email"] = nil
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound,"/aviva")
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
	
	e.File("/favicon.ico", filepath.Join(base_path, "themes/alfa/assets/img/favicon.png"))
	
	e.Static("/cms_assets", filepath.Join(base_path, "cms/assets"))
	e.Static("/theme_assets", filepath.Join(base_path, "themes/alfa/assets"))
	
	//	e.File("/", "public/index.html")
	//	e.GET("/tasks", handlers.GetTasks(db))
	//	e.PUT("/tasks", handlers.PutTask(db))
	//	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

	e.Logger.Fatal(e.Start(config.Port))
}

func TestServer() {
    //run()
    //gdb.TestRnd()
	cms:=RandomCMS()
    root_page := cms.Root
    root_page.Print(0,0)
	//fmt.Printf("%+v\n", root_page)
}

