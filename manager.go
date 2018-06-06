package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

const bufSize = 4 * 1024 * 1024

type File struct {
	Name string
	Path string
	Ext  string
}

func crawler(fl *[]File) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".srt" {
			return nil
		}
		f := File{
			Name: info.Name(),
			Path: path,
			Ext:  filepath.Ext(path),
		}
		*fl = append(*fl, f)
		return nil
	}
}

func PullTreeDir(root string, fl *[]File) error {
	err := filepath.Walk(root, crawler(fl))
	if err != nil {
		return err
	}
	return nil
}

func PullDir(root string, fl *[]File) error {
	m, _ := filepath.Glob(filepath.Join(root, "*"))
	for i := range m {
		d, err := os.Lstat(m[i])
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(m[i]) != ".srt" {
			continue
		}
		f := File{
			Name: d.Name(),
			Path: m[i],
			Ext:  filepath.Ext(m[i]),
		}
		*fl = append(*fl, f)
	}
	return nil
}

func PullOut(dst, src string) (int64, error) {
	buf := make([]byte, bufSize)
	bw := int64(0)

	// Open source file
	fs, err := os.OpenFile(src, syscall.O_RDONLY, os.ModePerm)
	if err != nil {
		return bw, err // Handle Error
	}
	// try to close fs
	defer func() {
		if err := fs.Close(); err != nil {
			panic(err)
		}
	}()
	stt, err := fs.Stat()
	if err != nil {
		return bw, err
	}
	// Get file size
	chunckSize := stt.Size() / 1024

	// Open destine file
	fd, err := os.OpenFile(dst, syscall.O_CREAT|syscall.O_WRONLY, os.ModePerm)
	if err != nil {
		return bw, err // Handle Error
	}
	defer func() {
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}()

	// write fs to fd
	for i := 0; i < int(chunckSize/4); i++ {
		r, err := fs.Read(buf)
		if err == io.EOF {
			return bw, nil
		}
		if err != nil {
			return bw, err
		}
		bw = int64(r)
		_, err = fd.Write(buf[:r])
		if err != nil && err != io.EOF {
			panic(err)
		}
	}
	return bw, nil
}

func main() {
	var fl []File
	err := PullDir(os.Args[1], &fl)
	if err != nil {
		panic(err)
	}
	for _, f := range fl {
		// Check source dir
		if _, err := os.Stat(os.Args[2]); os.IsNotExist(err) {
			if err = os.MkdirAll(os.Args[2], os.ModePerm); err != nil {
				panic(err)
			}
		}
		// PullOut files
		dst := filepath.Join(os.Args[2], f.Name)
		fmt.Printf("Creating File: %s\n", f.Name)
		bw, err := PullOut(dst, f.Path)
		if err != nil {
			fmt.Println("Error occurs: ", err)
			return
		}
		fmt.Println("File created, bytes written: ", bw)
	}
}
