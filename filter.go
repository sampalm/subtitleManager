// TODO: Implement a new function that organize files by folders that have the name of the serie or movie

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"time"
)

// File used to save files info
type File struct {
	Path   string
	Name   string
	Ext    string
	Size   int64
	Folder string
}

// Sub used to save files content
type Sub struct {
	Name string
	Path string
	Body *bytes.Buffer
}

// Log used to save errors log
type Log struct {
	Func string
	Err  error
}

var options = make([]string, 4)
var flags = make([]string, 3)
var wg sync.WaitGroup
var mtx sync.Mutex
var dst string
var errFound bool
var logs []Log

const (
	path    = 0
	ext     = 1
	version = 2
	move    = 3
)

func init() {
	p := flag.String("p", "", "Set the root path.")
	e := flag.String("e", ".srt", "Set the extension of the file.")
	v := flag.String("v", "", "Set the version of the subtitle.")
	h := flag.Bool("h", false, "Returns basic instructions to use Subtitle Manager.")
	m := flag.String("m", "", "Only move files to this selected directory.")
	flag.Parse()
	flags = []string{*p, *e, *v, *m}
	if *h {
		PrintInfo()
		PrintHelp()
		os.Exit(1)
	}
}

func main() {
	start := time.Now()

	PrintInfo()
	if flags[path] == "" {
		fmt.Fprintf(os.Stdout, "flags: you must define flag -p.")
		os.Exit(1)
	}

	sf, err := getall()
	if err != nil {
		fmt.Fprintf(os.Stdout, "system: unable to complete task because: %v\n", err)
		os.Exit(1)
	}

	subs := copyall(sf)

	// Main function
	core(subs)

	end := time.Since(start).Seconds()
	if errFound {
		fmt.Fprintf(os.Stdout, "*** TASK COMPLETED in %.2fs with some errors. Open log file to see all program execution errors ***", end)
		if err := CreateLogFile(logs); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}
	fmt.Fprintf(os.Stdout, "*** TASK COMPLETED in %.2fs without any errors. ***", end)
}

func getall() (files []File, err error) {

	err = filepath.Walk(flags[path], func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == flags[ext] {
			if flags[version] != "" {
				if match, _ := regexp.MatchString("([a-zA-Z0-9]+)."+flags[version], path); !match {
					return nil
				}
			}

			file := File{
				Path:   path,
				Name:   info.Name(),
				Ext:    filepath.Ext(path),
				Size:   info.Size(),
				Folder: filepath.Dir(path),
			}
			files = append(files, file)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, err
}

func copyall(files []File) (subs []Sub) {
	fmt.Fprintf(os.Stdout, "copyAllFiles: starting to copy all files...\n")

	wg.Add(len(files))
	for _, file := range files {
		go copy(file, &subs)
	}
	wg.Wait()
	fmt.Fprintf(os.Stdout, "copyAllFiles: all files have been copied!\n")
	return subs
}

func copy(file File, subs *[]Sub) {
	defer wg.Done()
	start := time.Now()
	buf := &bytes.Buffer{}

	src, err := os.Open(file.Path)
	if err != nil {
		fmt.Fprintf(os.Stdout, "copyFile: could not open file: %s.\n Error: %s\n", file.Path, err)
		errFound = true
		log := Log{
			Func: "copy",
			Err:  err,
		}
		logs = append(logs, log)
		return
	}
	defer src.Close()
	io.Copy(buf, src)

	// Save subtitle body
	sub := Sub{
		Name: file.Name,
		Path: file.Path,
		Body: buf,
	}

	mtx.Lock()
	*subs = append(*subs, sub)
	mtx.Unlock()
	end := time.Since(start).Seconds()
	fmt.Fprintf(os.Stdout, "copyFile: file %s copied in %.2fs.\n", file.Name, end)
	return
}

func deleteall(dst string, ext string) {
	if err := filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(os.Stdout, "deleteAllFiles: could not detele file: %s, Error: %v\n", path, err)
				log := Log{
					Func: "deleteAll",
					Err:  err,
				}
				logs = append(logs, log)
				errFound = true
			}
		}
		fmt.Fprintf(os.Stdout, "deleteAllFiles: file deleted without any errors: %s\n", info.Name())
		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "deleteAllFiles: unable to delete any file. %v\n", err)
		os.Exit(3)
	}
}

func moveall(dst string, subs []Sub) {
	// Create or check if the output dst exists
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		err := os.Mkdir(dst, 0777)
		if err != nil {
			fmt.Fprintf(os.Stdout, "moveFile: could not create directory: %s, Error: %v\n", dst, err)
			os.Exit(3)
		}
		fmt.Fprintf(os.Stdout, "moveFile: Output directory created: %s\n", dst)
	}

	// Delete and move files
	for _, sub := range subs {
		wg.Add(2)
		go delete(sub.Path)
		go create(sub, dst)
	}
	wg.Wait()
}

func delete(path string) {
	defer wg.Done()
	if err := os.Remove(path); err != nil {
		fmt.Fprintf(os.Stdout, "deteleFile: could not delete file: %s, Error: %v\n", path, err)
		errFound = true
		log := Log{
			Func: "delete",
			Err:  err,
		}
		logs = append(logs, log)
		return
	}
}

func create(sub Sub, path string) {
	defer wg.Done()
	file := filepath.Join(path, "/", sub.Name)
	f, err := os.OpenFile(file, syscall.O_RDWR|syscall.O_CREAT, 0777)
	if err != nil {
		fmt.Fprintf(os.Stdout, "createFile: could not create file: %s, Error: %v\n", file, err)
		errFound = true
		log := Log{
			Func: "create",
			Err:  err,
		}
		logs = append(logs, log)
		return
	}
	defer f.Close()

	if _, err := f.Write(sub.Body.Bytes()); err != nil {
		fmt.Fprintf(os.Stdout, "createFile: could not save file: %s, Error: %v\n", file, err)
		errFound = true
		log := Log{
			Func: "create",
			Err:  err,
		}
		logs = append(logs, log)
		return
	}

	fmt.Fprintf(os.Stdout, "createFile: file created without any error: %s\n", file)
}

func core(subs []Sub) {
	// Check if flag MOVE is SET
	if flags[move] != "" {
		moveall(flags[move], subs)
		return
	}

	// Default routine
	// Delete all files inside default directory
	deleteall(flags[path], flags[ext])
	for _, sub := range subs {
		wg.Add(1)
		go create(sub, flags[path])
	}
	wg.Wait()
}
