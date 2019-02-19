package gdb

import (
    "fmt"
)

type Page struct {
	Title    string
	Name     string
	Children *[]Page
}

func ExampleSite () Page {
    p := Page{Title: "Frontpage",
		Name: "Frontpage",
		Children: &[]Page{
			{Title: "Products", Children: &[]Page{{Title: "Prod1", Children: &[]Page{}},
				{Title: "Prod2", Children: &[]Page{}}}},
			{Title: "About us", Children: &[]Page{}}}}
    return p
}

func init() {
    fmt.Printf("data.go init\n")
}
