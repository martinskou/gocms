package gdb

import (
    "fmt"
    "strings"
    "math/rand"
    "time"
)

type Content struct {
	// Content
	Title      string
	Teaser     string
	Content    string
	LinkUrl    string
	LinkText   string
	LinkTarget string
	ImageUrl   string
	ImageText  string
	// Meta
	Index      int     // position if part of list
	Position   string  // position in template
	Class      string  // a class designator
	Visible    bool
}

type Page struct {
	// Content
	Title    string
	Name     string
	Content  string
	// Meta
	Slug     string
	Template string // template file to use
	Class    string
	Index    int    // position in respect to other pages
	Visible  bool   // visible in menu
	Parent   *Page
	Children *[]Page
}



func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func joinRunes(runes ...rune) string {
	var sb strings.Builder
	for _, r := range runes {
		sb.WriteRune(r)
	}
	return sb.String()
}


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func randSentence(n int) string {
	b := make([]string, n)
    for i := range b {
        b[i] = randString(6)
    }
    return strings.Join(b," ")
}

func (p Page) Print(index int, level int) {
    fmt.Printf("%s %d Title: %s\n", strings.Repeat("  ",level), index, p.Title)
    fmt.Printf("%s %d Name:  %s\n", strings.Repeat("  ",level), index, p.Name)
    for i,c := range *p.Children {
        c.Print(i,level+1)
    }
    return
}


func (p Page) Apply(fn func(Page)) {
	fn(p)
    for _,c := range *p.Children {
        c.Apply(fn)
    }
    return
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


func TestRnd() {
    var m map[int]int
    m = make(map[int]int,0)
    for i := 0; i < 100000; i++ {
        x:=rand.Intn(10)
        //println(x)
        m[x]+=1
    }
    for k := range m {
        println(k,m[k])
    }
}

func RandomContent() Content {
	t:=strings.Title(randString(rand.Intn(5)+2))
	c:=Content{Title: t,
	           Teaser: randSentence(30),
			   Content: randSentence(150),
			   Visible: true}
	return c
}


func RandomPages (max_pages int, parent *Page) *[]Page {
    //ap := make([]Page,0)
    var ap []Page
	if max_pages>0 {
	    pages := rand.Intn(max_pages)
	    for i := 0; i < pages; i++ {
	        t:=strings.Title(randString(rand.Intn(5)+2))
	        p:= Page{Title: t,
	                 Name: t,
					 Slug: strings.ToLower(t),
					 Visible: true,
					 Template: "standard_page.html",
					 Content: randSentence(150),
					 Parent: parent} // &[]Page{}}
			p.Children = RandomPages(max_pages-2,&p)
	        ap=append(ap,p)
	    }
	}
    return &ap //[]Page{ p }
}


func RandomSite () Page {
    p := Page{Title: "Frontpage",
			  Name: "Frontpage",
			  Slug: "",
			  Template: "standard_page.html",
			  Content: randSentence(50),
			  Parent: nil}
	p.Children = RandomPages(10,&p)
    return p
}

func init() {
    rand.Seed(time.Now().UnixNano())
}
