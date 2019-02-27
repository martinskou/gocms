package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"errors"
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"encoding/json"
)

func bundle_files(abs_dest string, files_path string, files []string) {

	out, err := os.OpenFile(abs_dest, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("failed to open first file for writing:", err)
	}
	defer out.Close()

	for _, f := range files {

		out.WriteString(fmt.Sprintf("\n\n/* ---- [ %s ]----- */\n\n", f))

		in, err := os.Open(filepath.Join(files_path, f))
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


// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func build(path string, b Build) {
	dest_path:=filepath.Join(path,b.Destination)
	bundle_files(dest_path, path, b.Files)
}

type Build struct {
		Task         string
		Type         string
		Files        []string
		Destination  string
}

func LoadBuild(path string) ([]Build, error) {
	var builds []Build
	bundle_file:=filepath.Join(path,"build.json")
	if Exists(bundle_file) {
		dat, _ := ioutil.ReadFile(bundle_file)
		err := json.Unmarshal(dat  , &builds)
		if err != nil {
			fmt.Println("error:", err)
			return builds, errors.New("Error reading json file")
		} else {
			return builds, nil
		}
	} else {
		return builds, errors.New("Error, missing json file")
	}

}

func Bundle(path string) {
	fmt.Println(path)
	bundle_file:=filepath.Join(path,"build.json")
	if Exists(bundle_file) {
		dat, _ := ioutil.ReadFile(bundle_file)
		var builds []Build
		err := json.Unmarshal(dat  , &builds)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			for _,b := range builds {
				fmt.Printf("%+v\n", b)
				build(path,b)
			}
		}
	} else {
		fmt.Println("build.json file not found in path")
	}
	
	
}

func NewWatcher(files []string, path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					Bundle(path)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for _,f := range files {
		err = watcher.Add(f)
		if err != nil {
			log.Fatal(err)
		}}
	<-done
}

func Watch(path string) {
	builds,err:=LoadBuild(path)
	Bundle(path)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		var files []string
		for _,b := range builds {
//			fmt.Printf("%+v\n", b.Files)
			for _,f := range b.Files {
				files=append(files,filepath.Join(path,f))
			}
		}
		fmt.Printf("%+v\n", files)
		NewWatcher(files,path)
	}
}
