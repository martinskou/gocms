package main

import (
	//    "database/sql"
	//    "go-echo-vue/handlers"
	"github.com/labstack/echo/v4"
//	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"path/filepath"
//	"io/ioutil"
	"net/http"
//    "time"
    "fmt"
	"log"
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

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	rendertype string
	templates map[string]*template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	log.Println("Render", name,t.rendertype)
	template, exists := t.templates[name]
	if exists {
		return template.ExecuteTemplate(w, name, data)
	} else {
		return errors.New("Template "+name+" not found")
	}
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

func make_templates(cmsPartialsGlob string ,cmsTemplateGlob string) map[string]*template.Template {
	cms_templates := make(map[string]*template.Template)
	cms_partials,_:=filepath.Glob(cmsPartialsGlob)
	files,_:=filepath.Glob(cmsTemplateGlob)
	for _,x := range files {
		f:=filepath.Base(x)
		if f[0]!='_' {
			log.Println("Make template:",f)
			t:=template.New(f)
			t.ParseFiles(x)
			t.ParseFiles(cms_partials...)			
			cms_templates[f]=t
		}
	}
	return cms_templates
}

func render_site_page(site_templates map[string]*template.Template, p *Page, root_page *Page,config Config, c echo.Context) ([]byte, error) {
	site_renderer := &TemplateRenderer{rendertype: "site", templates: site_templates}			
	data:=map[string]interface{}{"config":config, "pages":root_page, "current": p}
	buf := new(bytes.Buffer)
	err := site_renderer.Render(buf, "page.html", data, c)
	return buf.Bytes(), err
}

func RunServer(base_path string, config_path string) {	
	//	db := initDB("storage.db")
	//	migrate(db)

	log.Println("Base path", base_path)
	
	e := echo.New()
	e.Debug=true
//	e.Use(middleware.Logger())
//	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	cmsTemplateGlob := filepath.Join(base_path, "cms/templates/*.html")
	cmsPartialsGlob := filepath.Join(base_path, "cms/templates/_*.html")
	siteTemplateGlob := filepath.Join(base_path, "themes/alfa/templates/*.html")
	sitePartialsGlob := filepath.Join(base_path, "themes/alfa/templates/_*.html")
	
	cms_templates := make_templates(cmsPartialsGlob,cmsTemplateGlob)
	site_templates := make_templates(sitePartialsGlob,siteTemplateGlob)

	cms_renderer := &TemplateRenderer{rendertype: "cms", templates: cms_templates}

	for _, t := range cms_renderer.templates {
		fmt.Println(t.Name(),t)
	}

	
	config, err := LoadConfig(config_path)
	if err!=nil {
		fmt.Println("Config file could not be loaded")
		return
	}
	
	root_page := RandomSite()

	// Attach all pages to Echos router
	root_page.Apply(func(p Page) {
		e.GET(p.AbsSlug(), func(c echo.Context) error {
			buf,err:=render_site_page(site_templates,&p,root_page,config, c)
			if err==nil {
				return c.HTMLBlob(http.StatusOK, buf)
			} else {
				return err
			}
		})
	})

	

	e.GET("/aviva/login", func(c echo.Context) error {
		log.Println("/AVIVA/LOGIN -----------")

		sess, _ := session.Get("session", c)
		fmt.Println(sess.Values["email"])
		
		data:=map[string]interface{}{"config":config, "pages":root_page, "user": sess.Values["email"]}
			
		buf := new(bytes.Buffer)
		if err = cms_renderer.Render(buf, "login.html", data, c); err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, buf.Bytes())
	})
	
	e.GET("/aviva", func(c echo.Context) error {
		log.Println("/AVIVA -----------")		

		sess, _ := session.Get("session", c)
		fmt.Println(sess.Values["email"])
		
		if sess.Values["email"]==nil {
			return c.Redirect(http.StatusFound,"/aviva/login")
		} else {
			data:=map[string]interface{}{"config":config, "pages":root_page, "user": sess.Values["email"]}
			
			buf := new(bytes.Buffer)
			if err = cms_renderer.Render(buf, "dashboard.html", data, c); err != nil {
				return err
			}
			return c.HTMLBlob(http.StatusOK, buf.Bytes())
		}		
	})

	
	e.POST("/aviva/login/authenticate", func(c echo.Context) error {
		log.Println("/AVIVA/LOGIN/AUTHENTICATE -----------")
		var r Reply
		if (c.FormValue("email")=="msd@infoserv.dk") {
			UserSession(c.FormValue("email"), c)
			data := make(map[string]interface{})
			data["Goto"] = "/aviva"
			r = Reply{Status: "OK", Data: data}
		} else {
			r = Reply{Status: "FAIL", Data: ""}
		}
		return c.JSON(http.StatusOK,r)
	})

	e.GET("/aviva/logout", func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		sess.Values["email"] = nil
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound,"/aviva/login")
	})


	
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

