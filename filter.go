// Subtitles Manager - Suman ??
// TODO:
// 1) Implement a new function that only moves files instead delete them.
// 2) Implement a new function that organize files by folders that have the name of the serie or movie

package main

import (
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

type File struct {
	Path   string
	Name   string
	Ext    string
	Size   int64
	Folder string
}

var sub File
var options = make([]string, 4)
var flags = make([]string, 3)
var wg sync.WaitGroup

const (
	path    = 0
	ext     = 1
	version = 2
)

func init() {
	p := flag.String("p", "", "Set the root path.")
	e := flag.String("e", ".srt", "Set the extension of the file.")
	v := flag.String("v", "", "Set the version of the subtitle.")
	flag.Parse()
	flags = []string{*p, *e, *v}
}

func main() {
	var dst string
	start := time.Now()
	PrintInfo()
	if flags[path] == "" {
		fmt.Fprintf(os.Stdout, "Task stopped b/c path isn't declared.")
		os.Exit(1)
	}
	sf, err := getall()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Task stopped because: %v\n", err)
		os.Exit(1)
	}
	// Copy files to temp folder
	dst = os.TempDir() + "\\subfilter\\"
	copyall(sf, dst)

	// Delete all files from source folder
	deleteall(flags[path], flags[ext])

	// Copy files to source folder
	copyall(sf, flags[path])

	// Delete all files from temp folder
	deleteall(dst, flags[ext])

	end := time.Since(start).Seconds()
	fmt.Fprintf(os.Stdout, "*** TASK COMPLETED in %.2fs without any errors. ***", end)
}

func getall() (files []File, err error) {
	var named = false
	if flags[version] != "" {
		named = !named
	}
	err = filepath.Walk(flags[path], func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == flags[ext] {
			if named {
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

func copyall(files []File, dst string) {
	// Create or check if the output dst exists
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.Mkdir(dst, 0666)
		fmt.Fprintf(os.Stdout, "Temp dir: %s\n", dst)
	}
	fmt.Fprintf(os.Stdout, "Starting to copy all files...\n")

	wg.Add(len(files))
	for _, file := range files {
		go copy(file, dst)
	}
	wg.Wait()
	fmt.Fprintf(os.Stdout, "All files have been copied!\n")
}

func copy(file File, dst string) {
	defer wg.Done()
	start := time.Now()
	dsti := fmt.Sprintf("%s\\%s", dst, file.Name)
	// First verify if file already exists inside dst folder
	dstfile, err := os.OpenFile(dsti, syscall.O_CREAT|syscall.O_EXCL, 0666)
	if err != nil {
		fmt.Fprintf(os.Stdout, "While checking an error occurs: %v\n", err)
		return
	}
	defer dstfile.Close()

	srcfile, err := os.Open(file.Path)
	if err != nil {
		fmt.Fprintf(os.Stdout, "While opening an error occurs: %v\n", err)
		return
	}
	defer srcfile.Close()

	if _, err := io.Copy(dstfile, srcfile); err != nil {
		fmt.Fprintf(os.Stdout, "Task stopped, was not possible to copy the files.\n File: %s\n", file.Path)
		os.Exit(2)
	}

	// Set path to new folder
	file.Path = dsti

	end := time.Since(start).Seconds()
	fmt.Fprintf(os.Stdout, "File copied from %s to %s in %.2fs.\n", file.Path, dsti, end)
	return
}

func deleteall(dst string, ext string) {
	if err := filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(os.Stdout, "While deleting an error occurs: %v\n", err)
			}
		}
		fmt.Fprintf(os.Stdout, "File %s has been deleted.\n", info.Name())
		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "while deleting %v\n", err)
		os.Exit(3)
	}
}
