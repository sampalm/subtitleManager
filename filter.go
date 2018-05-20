// TODO: Implement a new function that organize files by folders that have the name of the serie or movie.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
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
	Media
}

type Media struct {
	Title    string
	Season   string
	Size     int64
	Language string
}

// Log used to save errors log
type Log struct {
	Func string
	Err  error
}

var wg sync.WaitGroup
var mtx sync.Mutex
var errFound bool
var logs []Log

func (fg Flag) Getall() (files []Sub, err error) {
	err = filepath.Walk(fg.Get[path], func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != fg.Get[ext] {
			return nil
		}
		if fg.Options[only] {
			onlypath := filepath.Join(fg.Get[0], info.Name())
			if onlypath != path {
				return nil
			}
		}
		if fg.Get[version] != "" {
			if match, _ := regexp.MatchString("([a-zA-Z0-9]+)."+fg.Get[version], path); !match {
				return nil
			}
		}
		files = append(files, buffering(info.Name(), path))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, err
}

func buffering(name, path string) (sub Sub) {
	var b = &bytes.Buffer{}
	// Open file
	fl, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stdout, "copyAllFiles: cannot copy file %s err: %s", name, err.Error())
		return
	}
	defer fl.Close()
	// Copy file to buffer
	if _, err := io.Copy(b, fl); err != nil {
		fmt.Fprintf(os.Stdout, "copyAllFiles: cannot copy file %s err: %s", name, err.Error())
		return
	}
	st, _ := fl.Stat()
	sub = Sub{
		Name: name,
		Path: path,
		Body: b,
		Media: Media{
			Size: st.Size(),
		},
	}
	return
}

func (fg Flag) Deleteall(subs []Sub) {
	fmt.Println("SUBS LENGTH: ", len(subs), "DELETE STATE: ", fg.Options[del])
	if len(subs) == 0 && !fg.Options[del] {
		return // None file was copy neither flag -d was set
	}
	if err := filepath.Walk(fg.Get[path], func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != fg.Get[ext] {
			return nil
		}
		if fg.Options[only] {
			onlypath := filepath.Join(fg.Get[0], info.Name())
			if onlypath != path {
				return nil
			}
		}
		if fg.Get[version] != "" && fg.Options[del] {
			if match, _ := regexp.MatchString("([a-zA-Z0-9]+)."+fg.Get[version], path); !match {
				return nil
			}
		}
		if err := os.Remove(path); err != nil {
			fmt.Fprintf(os.Stdout, "deleteAllFiles: could not detele file: %s, Error: %v\n", info.Name(), err)
			log := Log{
				Func: "deleteAll",
				Err:  err,
			}
			logs = append(logs, log)
			errFound = true
		}
		return nil
	}); err != nil {
		fmt.Fprintf(os.Stderr, "deleteAllFiles: unable to delete any file. %v\n", err)
		os.Exit(3)
	}
}

func (fg Flag) Moveall(subs []Sub) {
	cDone := make(chan bool)
	// Create or check if the output dst exists
	if _, err := os.Stat(fg.Get[move]); os.IsNotExist(err) {
		if err := os.MkdirAll(fg.Get[move], 0642); err != nil {
			fmt.Fprintf(os.Stdout, "moveFile: could not create directory: %s, Error: %v\n", fg.Get[move], err)
			os.Exit(3)
		}
		fmt.Fprintf(os.Stdout, "moveFile: directory created: %s\n", fg.Get[move])
	}

	// Delete and move files
	for _, sub := range subs {
		go delete(sub.Path)
		go create(sub.Name, sub.Body.Bytes(), fg.Get[move])
		cDone <- true
	}
	// Wait every gorountine finish
	for range subs {
		<-cDone
	}
}

func delete(path string) {
	if err := os.RemoveAll(path); err != nil {
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

func create(name string, body []byte, path string) {
	file := filepath.Join(path, name)
	f, err := os.OpenFile(file, syscall.O_RDWR|syscall.O_CREAT, 0624)
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

	if _, err := f.Write(body); err != nil {
		fmt.Fprintf(os.Stdout, "createFile: could not save file: %s, Error: %v\n", file, err)
		errFound = true
		log := Log{
			Func: "create",
			Err:  err,
		}
		logs = append(logs, log)
		return
	}
}

func (fg Flag) OrganizeAll(subs []Sub) {
	var ch = make(chan Sub)
	var cn = make(chan string)
	if fg.Options[del] {
		fg.Options[only] = true
		fg.Deleteall(subs)
	}
	for _, sub := range subs {
		go getParams(sub, ch)
	}
	for c := 0; c < len(subs); {
		select {
		case cs := <-ch:
			go organize(cs, cn)
		case <-cn:
			c++
		}
	}
}

func organize(sub Sub, cn chan string) {
	defer func() {
		// Notify that this function is done
		cn <- "done"
	}()
	newPath := filepath.Join(fg.Get[path], sub.Title, sub.Season)
	if fg.Get[move] != "" {
		newPath = filepath.Join(fg.Get[move], sub.Title, sub.Season)
	}
	// Create new folder to sub
	if err := os.MkdirAll(newPath, 0642); err != nil {
		fmt.Fprintf(os.Stdout, "organize: could not create directory: %s, Error: %v\n", fg.Get[move], err)
		os.Exit(3)
	}
	fmt.Fprintf(os.Stdout, "organize: directory created: %s\n", newPath)
	// Copy sub to new folder
	file := filepath.Join(newPath, sub.Name)
	f, err := os.OpenFile(file, syscall.O_RDWR|syscall.O_CREAT, 0624)
	if err != nil {
		fmt.Fprintf(os.Stdout, "organize: could not create file: %s, Error: %v\n", file, err)
		errFound = true
		log := Log{
			Func: "organize",
			Err:  err,
		}
		logs = append(logs, log)
		return
	}
	defer f.Close()
	if _, err := f.Write(sub.Body.Bytes()); err != nil {
		fmt.Fprintf(os.Stdout, "organize: could not save file: %s, Error: %v\n", file, err)
		errFound = true
		log := Log{
			Func: "organize",
			Err:  err,
		}
		logs = append(logs, log)
		return
	}
}

func getParams(sub Sub, ch chan Sub) {
	rawS := regexp.MustCompile("([a-zA-Z]([0-9])+[a-zA-Z]([0-9])+)").FindString(sub.Name)
	rawN := regexp.MustCompile("([a-zA-Z]([0-9])+[a-zA-Z]([0-9])+)").Split(sub.Name, -1)[0]
	sub.Title = strings.TrimSpace(strings.Replace(rawN, ".", " ", -1))
	if rawN != "" {
		sub.Season = strings.TrimSpace(regexp.MustCompile("([a-zA-Z]([0-9])+)").FindString(rawS))
	}
	if rawS == "" {
		sub.Title = strings.TrimSpace(strings.Replace(strings.Split(sub.Name, fg.Get[ext])[0], ".", " ", -1))
	}
	ch <- sub
}

func (fg Flag) FetchAll() (files []*os.File, err error) {
	var bucket []*os.File
	if err := filepath.Walk(fg.Get[path], func(path string, info os.FileInfo, err error) error {
		// Accepted files type .avi, .mp4, .mkv,
		ext := filepath.Ext(path)
		if ext != ".avi" && ext != ".mp4" && ext != ".mkv" {
			return nil
		}
		if fg.Options[only] {
			onlypath := filepath.Join(fg.Get[0], info.Name())
			if onlypath != path {
				return nil
			}
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("fetchAll: openfile: cannot open file %s err: %s", file.Name(), err.Error())
		}
		bucket = append(bucket, file)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("fetchAll: filepath: %s", err.Error())
	}
	if len(bucket) == 0 {
		return nil, fmt.Errorf("fetchAll: bucket: none file found")
	}

	return bucket, nil
}
