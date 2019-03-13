package main

import (
	"github.com/labstack/echo/v4"
	"text/template"
	"io"
	"path/filepath"
	"net/http"
    "fmt"
	log "github.com/sirupsen/logrus"
	"bytes"
	"errors"
	"github.com/labstack/gommon/color"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

const debug = true

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	rendertype string
	templates map[string]*template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(base_path string, w io.Writer, name string, data interface{}, c echo.Context) error {

	// RELOAD
	cmsTemplateGlob := filepath.Join(base_path, "cms/templates/*.html")
	cmsPartialsGlob := filepath.Join(base_path, "cms/templates/_*.html")
	cms_templates := make_templates(cmsPartialsGlob,cmsTemplateGlob)
	cms_renderer := &TemplateRenderer{rendertype: "cms", templates: cms_templates}
	t = cms_renderer

	
	template, exists := t.templates[name]
	if exists {
		err := template.ExecuteTemplate(w, name, data)
		if err!=nil {
			log.Printf(err.Error())
		}
		return err
	} else {
		log.Printf("Template "+name+" not found")
		return errors.New("Template "+name+" not found")
	}
}

func UserSession(email string, c echo.Context) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60*10,  // 86400 * 7,
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
			//log.Println("Make template:",f)
			t:=template.New(f)
			t.ParseFiles(x)
			t.ParseFiles(cms_partials...)
			cms_templates[f]=t
		}
	}
	return cms_templates
}

func render_site_page(site_templates map[string]*template.Template, p *Page, cms CMS, c echo.Context) ([]byte, error) {
	site_renderer := &TemplateRenderer{rendertype: "site", templates: site_templates}
	data:=map[string]interface{}{"config":cms.Config, "pages":cms.Root, "current": p}
	buf := new(bytes.Buffer)
	err := site_renderer.Render(cms.Path, buf, "page.html", data, c)
	return buf.Bytes(), err
}

func versionInfo(e *echo.Echo) {
	colorer := color.New()
	fmt.Println("")
	fmt.Println(`   _____       .__               `)
	fmt.Println(`  /  _  \___  _|__|__  _______   `)
	fmt.Println(` /  /_\  \  \/ /  \  \/ /\__  \  `)
	fmt.Println(`/    |    \   /|  |\   /  / __ \_`)
	fmt.Println(`\____|__  /\_/ |__| \_/  (____  /`)
	fmt.Println(`        \/                    \/ `,colorer.Red("v1.0.0"))
	fmt.Println("High performance minimal CMS")
	//fmt.Printf("%v\n", "♥" == "\u2665")
	fmt.Println(colorer.Blue("https://aviva.dk"))
}

func RunServer(base_path string, config_path string) {
	//	db := initDB("storage.db")
	//	migrate(db)

	e := echo.New()
	versionInfo(e)
	e.Debug=true
	e.HideBanner=true
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

	//for _, t := range cms_renderer.templates {
	//	fmt.Println(t.Name(),t)
	//}

	config, err := LoadConfig(config_path)
	if err!=nil {
		log.Println("Config file "+config_path+" could not be loaded")
		return
	}

	cms := RandomCMS()
	root_page := cms.Root
	cms.Path = base_path

	// Attach all user defined pages to Echos router
	root_page.Apply(func(p *Page) {
		e.GET(p.AbsSlug(), func(c echo.Context) error {
			//log.Printf("%+v\n", p)
			buf,err:=render_site_page(site_templates,p,cms, c)
			if err==nil {
				return c.HTMLBlob(http.StatusOK, buf)
			} else {
				return err
			}
		})
	})


	e.GET("/aviva/login", func(c echo.Context) error {
		return ViewLogin(cms,c,cms_renderer);
	})

	e.GET("/aviva", func(c echo.Context) error {
		return ViewDashboard(cms,c,cms_renderer);
	})

	e.GET("/aviva/page/:id", func(c echo.Context) error {
		return ViewPage(cms,c,cms_renderer);
	})

	e.GET("/aviva/page/json/:id", func(c echo.Context) error {
		return JsonGetPage(cms,c,cms_renderer);
	})
	e.POST("/aviva/page/json/:id", func(c echo.Context) error {
		return JsonPostPage(cms,c,cms_renderer);
	})

	e.POST("/aviva/login/authenticate", func(c echo.Context) error {
		return PostAuthenticate(cms,c,cms_renderer);
	})

	e.GET("/aviva/logout", func(c echo.Context) error {
		return GetLogout(cms,c,cms_renderer)
	})


	e.File("/favicon.ico", filepath.Join(base_path, "themes/alfa/assets/img/favicon.png"))

	e.Static("/cms_assets", filepath.Join(base_path, "cms/assets"))
	e.Static("/theme_assets", filepath.Join(base_path, "themes/alfa/assets"))

	//	e.File("/", "public/index.html")
	//	e.GET("/tasks", handlers.GetTasks(db))
	//	e.PUT("/tasks", handlers.PutTask(db))
	//	e.DELETE("/tasks/:id", handlers.DeleteTask(db))

//	go func() {
//		if err := e.Start(config.Port); err != nil {
//			e.Logger.Info("shutting down the server")
//		}
		e.Logger.Fatal(e.Start(config.Port))
//	}()

/*	<-quit
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
*/	
}

func ConvertStruct(s []ContentLink) string {
	fmt.Printf("%+v\n", s)
	fmt.Printf("(%v, %T)\n", s, s)
	
//  map[string]interface{}
//	var result map[string]interface{}
	return "ok"
}

func TestServer() {
    //run()
    //gdb.TestRnd()
	cms:=RandomCMS()
    root_page := cms.Root
  //  root_page.Print(0,0)
	//fmt.Printf("%+v\n", root_page)
	fmt.Printf("%+v\n\n", root_page.ContentLinkJsons())
//	fmt.Printf("%+v\n\n", ConvertStruct(root_page.ContentLinks))
	
}
