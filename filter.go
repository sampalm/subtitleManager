// TODO: Implement a new function that organize files by folders that have the name of the serie or movie.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
	sub = Sub{
		Name: name,
		Path: path,
		Body: b,
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
		wg.Add(2)
		go delete(sub.Path)
		go create(sub, fg.Get[move])
	}
	wg.Wait()
}

func delete(path string) {
	defer wg.Done()
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

func create(sub Sub, path string) {
	defer wg.Done()
	file := filepath.Join(path, sub.Name)
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
	fmt.Fprintf(os.Stdout, "createFile: file created: %s\n", file)
}
