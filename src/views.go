package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	//"log"
	log "github.com/sirupsen/logrus"
	"bytes"
//	"io/ioutil"
	"github.com/labstack/echo-contrib/session"
)



func ViewDashboard(cms CMS, c echo.Context, renderer *TemplateRenderer) error {
	log.WithFields(log.Fields{"path": c.Path()}).Info("ViewDashboard")
	sess, _ := session.Get("session", c)
	if sess.Values["email"]==nil {
		return c.Redirect(http.StatusFound,"/aviva/login")
	} else {
		data:=map[string]interface{}{"cms":cms, "user": sess.Values["email"]}

		buf := new(bytes.Buffer)
		if err := renderer.Render(cms.Path, buf, "dashboard.html", data, c); err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, buf.Bytes())
	}
}

func ViewPage(cms CMS, c echo.Context, renderer *TemplateRenderer) error {
	id:=c.Param("id")
	log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("ViewPage")
	sess, _ := session.Get("session", c)
	if sess.Values["email"]==nil {
		return c.Redirect(http.StatusFound,"/aviva/login")
	} else {
		// find page in cms root page tree
		current:=cms.Root.Find("Id",id)
		if current==nil {
			return c.Redirect(http.StatusFound,"/aviva")
		}
		data:=map[string]interface{}{"cms":cms, "current":current, "user": sess.Values["email"]}
		buf := new(bytes.Buffer)
		if err := renderer.Render(cms.Path, buf, "page_editor.html", data, c); err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, buf.Bytes())
	}
}



func JsonGetPage(cms CMS, c echo.Context, renderer *TemplateRenderer) error {
	id:=c.Param("id")
	log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("JsonGetPage")
	sess, _ := session.Get("session", c)
	if sess.Values["email"]==nil {
		log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("Not logged in")
		return c.NoContent(http.StatusNotFound)
	} 
	current:=cms.Root.Find("Id",id)
	if current==nil {
		log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("Page not found")
		return c.NoContent(http.StatusNotFound)
	}
	log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("Returning page as Json")
	return c.JSON(http.StatusOK, current.ContentLinkJsons())
}

func find_content(all_content []*Content, id string) *Content {
	for _, c := range all_content {
		if c.Id==id {
			return c
		}
	}
	return nil
}

func JsonPostPage(cms CMS, c echo.Context, renderer *TemplateRenderer) error {
	id:=c.Param("id")
	log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("JsonPostPage")
	sess, _ := session.Get("session", c)
	if sess.Values["email"]==nil {
		log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("Not logged in")
		return c.NoContent(http.StatusNotFound)
	} 


	current:=cms.Root.Find("Id",id)
	if current==nil {
		log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("Page not found")
		return c.NoContent(http.StatusNotFound)
	}
	log.WithFields(log.Fields{"path": c.Path(),"page": id}).Info("Returning page as Json")
//	return c.JSON(http.StatusOK, current.ContentLinkJsons())

	
	new_content_links := make([]ContentLinkJson,0)
	
	if err := c.Bind(&new_content_links); err != nil {
		log.Println(err)
		return err
	}
	log.Println("bind success")
	log.Println(new_content_links)

	// Update current pages contentlinks!
	current.ContentLinks=make([]ContentLink,0)
	for _, ncl := range new_content_links {
		log.Println(ncl)
		rc:=find_content(cms.Content, ncl.ContentId)
		nrcl:=ContentLink{
			Content:    rc,
			Position:   ncl.Position,
			Index:      ncl.Index,
			Visible:    true}
		current.ContentLinks=append(current.ContentLinks,nrcl)
		
	}
	

	//body, _ := ioutil.ReadAll(c.Request().Body)
	//log.Println(body)

	//m := echo.Map{}
	return c.JSON(http.StatusOK, current.ContentLinkJsons())
		
	//return c.NoContent(http.StatusOK)
	//return c.JSON(http.StatusOK, current.ContentLinks)
}

// Authentication

func ViewLogin(cms CMS, c echo.Context, renderer *TemplateRenderer) error {
	log.WithFields(log.Fields{"path": c.Path()}).Info("ViewLogin")
	sess, _ := session.Get("session", c)
	data:=map[string]interface{}{"cms":cms, "user": sess.Values["email"]}
	buf := new(bytes.Buffer)
	if err := renderer.Render(cms.Path, buf, "login.html", data, c); err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, buf.Bytes())
}

func PostAuthenticate(cms CMS, c echo.Context, renderer *TemplateRenderer) error {
	var r Reply
	log.WithFields(log.Fields{"path": c.Path(),"email": c.FormValue("email")}).Info("PostAuthenticate")
	if (c.FormValue("email")=="msd@infoserv.dk") {
		UserSession(c.FormValue("email"), c)
		data := make(map[string]interface{})
		data["Goto"] = "/aviva"
		r = Reply{Status: "OK", Data: data}
	} else {
		r = Reply{Status: "FAIL", Data: ""}
	}
	return c.JSON(http.StatusOK,r)
}

func GetLogout(cms CMS, c echo.Context, renderer *TemplateRenderer) error {
	sess, _ := session.Get("session", c)
	log.WithFields(log.Fields{"path": c.Path(),"email": sess.Values["email"]}).Info("GetLogout")
	sess.Values["email"] = nil
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound,"/aviva/login")
}
