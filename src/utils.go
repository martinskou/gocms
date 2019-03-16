package main

import (
	"fmt"
	"strings"
    "math/rand"
    "time"
	"errors"
	"io/ioutil"
	"encoding/json"
	log "github.com/sirupsen/logrus"
)




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

func pseudo_uuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X", b[0:2], b[2:4], b[4:6], b[6:8])
	return
}


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
//var konsonanter = []rune("bcdfghjklmnprstvwxz")
var konsonanter = []rune("bcdddffghhjkllmmnnnprrrsstttvwxz")
//var vokaler = []rune("aeioqyuæøå")
var vokaler = []rune("aaaeeeeeiiioooqyuuæøå")

func randString(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func randWord(n int) string {
    b := make([]rune,0)
    for i:=0; i<n; i++ {
		b=append(b, konsonanter[rand.Intn(len(konsonanter))])
		b=append(b, vokaler[rand.Intn(len(vokaler))])
		if rand.Intn(3)==1 {
			b=append(b, konsonanter[rand.Intn(len(konsonanter))])
		}
    }
    return string(b)
}

func randSentence(n int) string {
	b := make([]string, n)
    for i := range b {
        b[i] = randWord(rand.Intn(3)+1)
    }
    return strings.Join(b," ")
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

func RandomContent(allow_children bool) Content {
	t:=strings.Title(randWord(rand.Intn(3)+2))
	var children []*Content
	if allow_children {
		children=RandomContents(rand.Intn(5),false)
	}
	c:=Content{Title: t,
			   Name: t,
			   Id: pseudo_uuid(),
	           Teaser: randSentence(30),
			   Content: randSentence(150),
			   ImageUrl: "https://source.unsplash.com/random/640x480",
			   ImageText: randSentence(10),
			   Visible: true,
			   Children: children,
		   	   Type: "Article"}
	return c
}

func RandomContents (items int, allow_children bool) []*Content {
    var ap []*Content
    for i := 0; i < items; i++ {
		c:=RandomContent(allow_children)
	    ap=append(ap,&c)
    }
    return ap
}


func RandomPages (max_pages int, parent *Page) []*Page {
    //ap := make([]Page,0)
    var ap []*Page
	if max_pages>0 {
	    pages := rand.Intn(max_pages)
	    for i := 0; i < pages; i++ {
	        t:=strings.Title(randWord(rand.Intn(2)+2))
	        p:= Page{Title: t,
	                 Name: t,
				     Id: pseudo_uuid(),
					 Slug: strings.ToLower(t),
					 Visible: true,
					 Template: "standard_page.html",
					 ContentLinks: make([]ContentLink,0), // randSentence(150),
					 Parent: parent} // &[]Page{}}
			p.Children = RandomPages(max_pages-2,&p)
	        ap=append(ap,&p)
	    }
	}
    return ap //[]Page{ p }
}

func RandomPageHierarchy () *Page {
    p := Page{Title: "Frontpage",
		      Name: "Frontpage",
		      Id: pseudo_uuid(),
			  Slug: "",
			  Template: "standard_page.html",
			  ContentLinks: make([]ContentLink,0), //randSentence(50),
			  Parent: nil}
	p.Children = RandomPages(10,&p)
    return &p
}

func FillPagesWithContent(cms *CMS) {
	cms.Root.Apply(func(p *Page) {
	//	fmt.Println("Adding content to",p.AbsSlug())
	//	p.Title="DEMO!"
		cll := len(cms.Content)
		mx := rand.Intn(5)+3
		id_a, id_b := 0, 0
		var cp string
		var idx int
		for i:=0; i<mx; i++ {
			c := cms.Content[rand.Intn(cll)]
			if rand.Intn(2)==0 {
				cp="b"
				idx=id_a
				id_a+=1
			} else {
				cp="a"
				idx=id_b				
				id_b+=1
			}
			cl := ContentLink{
				Content:    c,    // &(*cms.Content)[0],
				Position:   cp,   // name of position in template
				Index:      idx,  // index in position
				Visible:    true}
			p.ContentLinks=append(p.ContentLinks,cl)
		}

/*		if (cms.Root.Children[0].Id==p.Id) {
			fmt.Printf("RESULT 1 : %s %s %p\n",p.Id, p.Title, p)
		}
		if (cms.Root.Id==p.Id) {
			fmt.Printf("RESULT A : %s %s %p\n",p.Id, p.Title, p)
		}*/
	})
	//cms.Root.Title="DEMO!!!"
//	fmt.Printf("%p \n", cms.Root)
//	fmt.Printf("RESULT 2 : %s %s %p\n",cms.Root.Children[0].Id, cms.Root.Children[0].Title, &cms.Root.Children[0])
//	fmt.Printf("RESULT B : %s %s %p\n",cms.Root.Id, cms.Root.Title, cms.Root)
//	fmt.Printf("ROOT : %+v\n", cms.Root)
}




func Save(cms *CMS) {
	bytes, err := json.MarshalIndent(cms,"","    ")
    if err != nil {
        log.Fatal(err)
    }	
 //   log.Println(string(bytes))
	err = ioutil.WriteFile("output.json", bytes, 0644)
	if err != nil {
        log.Fatal("Error writing JSON", err.Error())
    }	
}

func RandomCMS () CMS {
    cms := CMS{
		Config: Config{Title: randString(10),
			           Theme: "alfa"},
		Root: RandomPageHierarchy(),
		Content: RandomContents(100,true)}
	FillPagesWithContent(&cms)

	Save(&cms)
    return cms
}


func LoadConfig(config_path string) (Config,error) {
	var config Config

	if Exists(config_path) {
		dat, _ := ioutil.ReadFile(config_path)
		err := json.Unmarshal(dat, &config)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			return config, nil
		}
	} else {
		fmt.Println("config.json file not found")
	}
	return config, errors.New("Config not loaded")
}

func init() {
    rand.Seed(time.Now().UnixNano())
}
