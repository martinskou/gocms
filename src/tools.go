package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
