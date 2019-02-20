package main

import (
	"fmt"
	"os"
	"io"
	"log"
	"github.com/urfave/cli"
	"github.com/martinskou/gocms/gserver"
//	"path/filepath"
  )

func bundle_files(dest string, files []string ) {

	out, err := os.OpenFile(dest, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("failed to open first file for writing:", err)
	}
	defer out.Close()

	for _,f := range files {

		out.WriteString(fmt.Sprintf("\n\n/* ---- [ %s ]----- */\n\n", f)  )

		in, err := os.Open(f)
		if err != nil {
			log.Fatalln("failed to open second file for reading:", err)
		}
		defer in.Close()

		n, err := io.Copy(out, in)
		if err != nil {
			log.Fatalln("failed to append second file to first:", err)
		} else {
			fmt.Printf("%s %d\n", f, n)
		}
		in.Close()
	}

	out.Close()

}

func bundle() {
	fmt.Printf("!\n")

	bundle_files("assets/bundle.js", []string{"dependencies/vue/vue.js"})
	bundle_files("assets/bundle.css", []string{"dependencies/bulma/css/bulma.css", "dependencies/gocms/style.css" })

	bundle_files("assets/bundle.min.js", []string{"dependencies/vue/vue.min.js"})
	bundle_files("assets/bundle.min.css", []string{"dependencies/bulma/css/bulma.min.css", "dependencies/gocms/style.css" })

}

func main() {
	app := cli.NewApp()
	  app.Name = "gocms"
	  app.Usage = "GoCMS utilities"

	  app.Commands = []cli.Command{
		  { Name:    "bundle",
	        Aliases: []string{"b"},
	        Usage:   "bundle dependencies (js and css)",
	        Action:  func(c *cli.Context) error {
	          bundle()
	          return nil
	        }},
			{
			  Name:    "run",
			  Aliases: []string{"r"},
			  Usage:   "start gocms server",
			  Action:  func(c *cli.Context) error {
				gserver.Run()
				return nil
			  }},
			  {
  			  Name:    "test",
  			  Aliases: []string{"t"},
  			  Usage:   "test something",
  			  Action:  func(c *cli.Context) error {
  				gserver.Test()
  				return nil
  			  }}}

/*	  app.Action = func(c *cli.Context) error {
	    fmt.Println("Hello friend!")
	    return nil
	  }
*/
	  err := app.Run(os.Args)
	  if err != nil {
	    log.Fatal(err)
	  }
}
