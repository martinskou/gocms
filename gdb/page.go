package gdb

import (
    "fmt"
    "strings"
    "math/rand"
    "time"
)

type Page struct {
	Title    string
	Name     string
	Slug     string
	Content  string
	Parent   *Page
	Children *[]Page
}


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
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

func ExampleSite () Page {
    p := Page{Title: "Frontpage",
			  Name: "Frontpage",
			  Slug: "/",
			  Children: &[]Page{
			{Title: "Products", Children: &[]Page{{Title: "Prod1", Children: &[]Page{}},
				{Title: "Prod2", Children: &[]Page{}}}},
			{Title: "About us", Children: &[]Page{}}}}
    return p
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
					 Content: randString(80),
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
			  Content: "Blah blak...",
			  Parent: nil}
	p.Children = RandomPages(10,&p)
    return p
}

func init() {
    rand.Seed(time.Now().UnixNano())
}
