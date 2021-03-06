package main

import (
	"fmt"
	"log"
	"os"
//	"os/signal"
	"path/filepath"
	"github.com/urfave/cli"
)


func main() {
	app := cli.NewApp()
	app.Name = "gocms"
	app.Usage = "GoCMS utilities"

//	d, _ := filepath.Abs(filepath.Dir(os.Args[0]))
//	fmt.Println("PATH", d)
	
//	fmt.Println("WORK PATH", work_path)
//	cms_path := filepath.Dir(src_path)
//	fmt.Println("CMS", cms_path)

	app.Commands = []cli.Command{
		{Name: "bundle",
			Aliases: []string{"b"},
			Usage:   "bundle dependencies (js and css)",
			Action: func(c *cli.Context) error {
				bundle_path := c.Args().First()
				if bundle_path=="" {
					fmt.Println("You must specify a path containing a build.json file")
				} else {
					work_path, _ := os.Getwd()
					Bundle(filepath.Join(work_path, bundle_path))
				}
				return nil
			}},
		{Name: "watch",
			Aliases: []string{"w"},
			Usage:   "watch for changes and bundle dependencies (js and css)",
			Action: func(c *cli.Context) error {
				bundle_path := c.Args().First()
				if bundle_path=="" {
					fmt.Println("You must specify a path containing a build.json file")
				} else {
					work_path, _ := os.Getwd()
					BundleWatch(filepath.Join(work_path, bundle_path))
				}
				return nil
			}},		
		{
			Name:    "startserver",
			Aliases: []string{"r"},
			Usage:   "start gocms server",
			Action: func(c *cli.Context) error {
				config_path := c.Args().First()
				if config_path=="" {
					fmt.Println("You must specify a config.json file")
				} else {
					work_path, _ := os.Getwd()
					abs_config_path := filepath.Join(work_path,config_path)
					base_path := filepath.Dir(filepath.Dir(abs_config_path))
					//quit := make(chan os.Signal)
					RunServer(base_path,config_path)
				}
				return nil
			}},
		{
			Name:    "reloadserver",
			Aliases: []string{"R"},
			Usage:   "start gocms server and reload on change",
			Action: func(c *cli.Context) error {
				config_path := c.Args().First()
				if config_path=="" {
					fmt.Println("You must specify a config.json file")
				} else {
					work_path, _ := os.Getwd()
					abs_config_path := filepath.Join(work_path,config_path)
					base_path := filepath.Dir(filepath.Dir(abs_config_path))
					ReloaderWatch(base_path,config_path)
				}
				return nil
			}},
		
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test something",
			Action: func(c *cli.Context) error {
				TestServer()
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
