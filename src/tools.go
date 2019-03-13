package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	log "github.com/sirupsen/logrus"
	"os"
	//	"os/signal"
	"os/exec"
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
			log.Println("error:", err)
			return builds, errors.New("Error reading json file")
		} else {
			return builds, nil
		}
	} else {
		return builds, errors.New("Error, missing json file")
	}

}

func Bundle(path string) {
	log.Println("Bundling",path)
	bundle_file:=filepath.Join(path,"build.json")
	if Exists(bundle_file) {
		dat, _ := ioutil.ReadFile(bundle_file)
		var builds []Build
		err := json.Unmarshal(dat  , &builds)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			for _,b := range builds {
				log.Printf("%+v\n", b)
				build(path,b)
			}
		}
	} else {
		log.Println("build.json file not found in path")
	}
	
	
}

func NewWatcher(files []string, path string, callback func(string), arg string) {
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
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("File change", event.Name)
					callback(arg)
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

func BundleWatch(path string) {
	builds,err:=LoadBuild(path)
	Bundle(path)
	if err != nil {
		log.Println("error:", err)
	} else {
		var files []string
		for _,b := range builds {
			for _,f := range b.Files {
				files=append(files,filepath.Join(path,f))
			}
		}
		// log.Printf("%+v\n", files)
		NewWatcher(files,path,Bundle,path)
	}
}


func WatcherChan(path string, res chan string) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	
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
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("File change", event.Name)
					res <- event.Name
					
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
		if strings.HasSuffix(f.Name(),".go") {
			//log.Println(f.Name())
			err = watcher.Add(f.Name())
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	<-done
}


func build_exec(src_path string) error {
	log.Println("Building gocms")
	args := append([]string{"go", "build", "-o", "../bin/gocms", ".", })
	var command *exec.Cmd
	
	command = exec.Command(args[0], args[1:]...)
	command.Dir=src_path
	
	output, _ := command.CombinedOutput()

	if command.ProcessState.Success() {
		return nil
	} else {
		return errors.New(string(output))
	}

	return nil
}



func runner(base_path string, bin_path string) (*exec.Cmd, error) {
	log.Println("Starting gocms")
	args := append([]string{filepath.Join(bin_path,"gocms"), "startserver", "config/config.json"})
	command := exec.Command(args[0], args[1:]...)
	command.Dir=base_path
	command.Stdout = os.Stdout
    command.Stderr = os.Stderr

/*	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil,err
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return nil,err
	}
*/
	err := command.Start()
	if err != nil {
		return nil,err
	}

//	go io.Copy(r.writer, stdout)
//	go io.Copy(r.writer, stderr)
	go command.Wait()

	return command,nil
}



func ReloaderWatch(path string, config string) {
	log.Println("Starting server with reload")
	src_path:=filepath.Join(path,"src")
	bin_path:=filepath.Join(path,"bin")
	

	var command *exec.Cmd
	command=nil
/*	
	command,err := runner(path,bin_path)
	if err!=nil {
		log.Fatal(err)
	}
	log.Println("PID:",command.Process.Pid)
*/	
	c := make(chan string)
	go WatcherChan(src_path,c)

	err := build_exec(src_path)
	if err!=nil {
		log.Fatal(err)
	}

	for {
		if command==nil {
			command,err = runner(path,bin_path)
			if err!=nil {
				log.Fatal(err)
			}
		}

		// r := <- c
		<- c
		
		err = build_exec(src_path)
		if err!=nil {
			log.Println(err)
		} else {

			log.Println("Build success")
			
			err = command.Process.Kill()
			if err!=nil {
				log.Fatal(err)
			}
			command=nil

			
		}

	}	
} 
