package main

import (
    "fmt"
    "strings"
	"encoding/json"
)


type Reply struct {
	Status     string
	Data       interface{} // almost everything
}

type Config struct {
	Title      string
	Theme      string
	Port       string
	Debug      bool
}

type Content struct {
	// Content
	Title      string
	Teaser     string  // text
	Content    string  // text / html / markdown
	LinkUrl    string
	LinkText   string
	LinkTarget string
	ImageUrl   string  // https://picsum.photos/200/300
	ImageText  string
	// Meta
	Name       string  // unique internal name
	Id         string  // unique internal id (permanent)
	Class      string  // a class designator
	Visible    bool
	Type       string  // Types : Article, List, Link, Image
	// Hierarchy
	Parent     *Content
	Children   []*Content
}


/*
func (c *Content) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Title     string
	}{
		Title:       c.Title,
	})
}
*/


type ContentLink struct {
	Content    *Content
	Position   string  // name of position in template
	Index      int     // index in position
	Visible    bool
}


type ContentLinkJson struct {
	ContentId    string
	ContentTitle string
	Index        int
	Position     string
	Visible      bool
}

func (c *ContentLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Content string
		Position string
		Index int
		Visible bool
	}{
		Content: c.Content.Id,
		Position: c.Position,
		Index: c.Index,
		Visible: c.Visible,
	})
}


func (c Page) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Title        string
		Description  string
		Keywords     string
		//Content  string
		ContentLinks []ContentLink
		// Meta
		Name     string  // unique internal name
		Id       string  // unique internal id (permanent)
		Domain   string
		Slug     string
		Template string // template file to use
		Class    string
		Redirect string // takes precedence if defined
		Index    int    // position in respect to other pages
		Visible  bool   // visible in menu
		MenuOnly bool   // visible in menu
		// Hierarchy
	//	Parent   *Page
		Children []*Page
	}{
		Title:        c.Title,
		Description:  c.Description,
		Keywords:     c.Keywords,
		//Content  string
		ContentLinks: c.ContentLinks,
		// Meta
		Name:     c.Name,  // unique internal name
		Id:       c.Id,  // unique internal id (permanent)
		Domain:   c.Domain,
		Slug:     c.Slug,
		Template: c.Template, // template file to use
		Class:    c.Class,
		Redirect: c.Redirect, // takes precedence if defined
		Index:    c.Index,    // position in respect to other pages
		Visible:  c.Visible,   // visible in menu
		MenuOnly: c.MenuOnly,   // visible in menu
		// Hierarchy
	//	Parent   *Page
		Children: c.Children,
	})
}

type Page struct {
	// Content
	Title        string
	Description  string
	Keywords     string
	//Content  string
	ContentLinks []ContentLink
	// Meta
	Name     string  // unique internal name
	Id       string  // unique internal id (permanent)
	Domain   string
	Slug     string
	Template string // template file to use
	Class    string
	Redirect string // takes precedence if defined
	Index    int    // position in respect to other pages
	Visible  bool   // visible in menu
	MenuOnly bool   // visible in menu
	// Hierarchy
	Parent   *Page
	Children []*Page
}

type PageWrap struct {
	Page     *Page
	Level    int
	Index    int
	GlobalIndex int
}

func (c *CMS) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Config     Config
		Path     string
		Content  []*Content
		Root     *Page
	}{
		Config:       c.Config,
		Path:         c.Path,
		Content:      c.Content,
		Root:         c.Root,
	})
}


type CMS struct {
	Config   Config
	Root     *Page
	Content  []*Content
	Path     string
}




func (p *Page) ContentLinkJsons() []ContentLinkJson {
	cl := make([]ContentLinkJson,0)
	for _, c := range p.ContentLinks {
		clj:=ContentLinkJson{
			ContentId: c.Content.Id,
			ContentTitle: c.Content.Title,
			Index:     c.Index,
			Position:  c.Position,
			Visible:   c.Visible,
		}
		cl=append(cl,clj)
	}
	return cl
}


// Return true if Page with argument id is a child or decendant
// of Page p.
func (p *Page) ChildOf(id string) bool {
	if p.Id==id {
		return true
	}
	parent := p.Parent
	if parent != nil {
		for parent != nil {
			if parent.Id==id {
				return true
			}
			parent = parent.Parent
		}
	}
	return false
}

func (p *Page) RemoveChild(child *Page) {
	found:=-1
	for index,c := range p.Children {
		if c==child {
			found=index
		}
	}
	if found>-1 {
		p.Children=append(p.Children[:found], p.Children[found+1:]...)
	}
}

func (p *Page) AddChild(child *Page) {
	p.Children=append(p.Children,child)
	child.Parent=p
}



func (p *Page) ContentForPosition(pos string) []ContentLink {
	cl := make([]ContentLink,0)
	for _, c := range p.ContentLinks {
		if c.Position==pos {
			cl=append(cl,c)
		}
	}
	return cl
}




func (page *Page) AppendPages(page_list *[]PageWrap,level int, index int, gindex *int) {
	pw := PageWrap{
		Page: page,
		Level: level,
		Index: index,
	    GlobalIndex: *gindex}
	*gindex+=1
	*page_list = append(*page_list,pw)
	for idx, c := range page.Children {
		c.AppendPages(page_list,level+1,idx,gindex)
	}
}

func (cms CMS) PageList() []PageWrap {
	page_list := make([]PageWrap,0)
	gi := 0
	cms.Root.AppendPages(&page_list,0,0,&gi)
	return page_list
}

func (p Page) Print(index int, level int) {
    fmt.Printf("%s %d Title: %s\n", strings.Repeat("  ",level), index, p.Title)
    fmt.Printf("%s %d Name:  %s\n", strings.Repeat("  ",level), index, p.Name)
    for i,c := range p.Children {
        c.Print(i,level+1)
    }
    return
}


func (p *Page) Apply(fn func(*Page)) {
	fn(p)
    for _,c := range p.Children {
//		(&c).Apply(fn)
		c.Apply(fn)
    }
    return
}

func (p Page) BCT() string {
	bct:="<li class=\"is-active\"><a href=\"#\" aria-current=\"page\">"+p.Title+"</a></li>\n"
	for p.Parent != nil {
		p=*p.Parent
		bct=fmt.Sprintf("<li><a href=\"%s\">%s</a></li>\n%s", p.AbsSlug(), p.Title, bct)
	}
    return "<nav class=\"breadcrumb has-arrow-separator\" aria-label=\"breadcrumbs\"><ul>"+bct+"</ul></nav>"
}

func (p Page) CMSBCT() string {
	bct:=fmt.Sprintf("<li class=\"is-active\"><a class=\"%s\" href=\"#\" aria-current=\"page\">%s</a></li>\n", p.Class, p.Name)
	for p.Parent != nil {
		p=*p.Parent
		bct=fmt.Sprintf("<li><a class=\"%s\" href=\"/aviva/page/%s\">%s</a></li>\n%s", p.Class, p.Id, p.Name, bct)
	}
    return "<nav class=\"breadcrumb has-arrow-separator\" aria-label=\"breadcrumbs\"><ul>"+bct+"</ul></nav>"
}


func (p Page) AbsSlug() string {
	slug := p.Slug
	parent := p.Parent
	if parent != nil {
		for parent != nil {
			slug = strings.Join( []string{parent.Slug , slug} , "/")
			parent = parent.Parent
		}
	} else {
		slug = "/"
	}
	return slug
}

func (p *Page) Find(field string, value string) *Page {
	if p.Id==value {
		return p
	}
	for _,c := range p.Children {
		r := c.Find(field,value)
		if r!=nil {
			return r
		}
	}
	return nil
}



func (c Content) HasChildren() bool {
	return len(c.Children)>0
}



/*
func (p *Page) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Title     string
		Id        string
		Domain    string
		Slug      string
		Template  string // template file to use
		Class     string
		Redirect  string // takes precedence if defined
		
	}{
		Id:       p.Id,
		Title:    p.Title,
		Domain:    p.Domain,
		Slug:       p.Slug,
		Template:       p.Template,
		Class:       p.Class,
		Redirect:       p.Redirect,
	})
}
*/

