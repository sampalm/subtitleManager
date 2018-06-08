package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

const bufSize = 4 * 1024 * 1024

var exts = map[string]bool{".srt": true, ".sub": true, ".sbv": true}

type File struct {
	Name string
	Path string
	Ext  string
	Info
}

type Info struct {
	Title  string
	Season string
}

type PullFiles func(root, ignore string, fl *[]File) error

func crawler(ignore string, fl *[]File) filepath.WalkFunc {
	ignore = filepath.Join(ignore)
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path == ignore {
				return filepath.SkipDir
			}
			return nil
		}
		// Check for valid extensions
		if _, ok := exts[filepath.Ext(path)]; !ok {
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

func PullTreeDir(root, ignore string, fl *[]File) error {
	err := filepath.Walk(root, crawler(ignore, fl))
	if err != nil {
		return err
	}
	return nil
}

func PullDir(root, ignore string, fl *[]File) error {
	ignore = filepath.Join(ignore)
	m, _ := filepath.Glob(filepath.Join(root, "*"))
	for i := range m {
		d, err := os.Lstat(m[i])
		if err != nil {
			return err
		}
		// Trying to copy all files to nowhere ?
		if filepath.Dir(m[i]) == ignore {
			continue
		}
		// Check for valid extensions
		if _, ok := exts[filepath.Ext(m[i])]; !ok {
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

func pullPath(folder []string, fl *[]File) error {
	for i := range folder {
		d, err := os.Lstat(folder[i])
		if err != nil {
			return err
		}
		if d.IsDir() {
			path, _ := filepath.Glob(filepath.Join(folder[i], "*"))
			pullPath(path, fl)
			continue
		}
		// Check for valid extensions
		if _, ok := exts[filepath.Ext(d.Name())]; !ok {
			continue
		}
		f := File{
			Name: d.Name(),
			Path: folder[i],
			Ext:  filepath.Ext(d.Name()),
		}
		regexSplit(f.Name, &f.Info)
		*fl = append(*fl, f)
	}
	return nil
}

func regexSplit(filename string, info *Info) { // Raising a flag
	s := regexp.MustCompile("([a-zA-Z]([0-9])+[a-zA-Z]([0-9])+)").FindString(filename)
	n := regexp.MustCompile("([a-zA-Z]([0-9])+[a-zA-Z]([0-9])+)").Split(filename, -1)[0]
	info.Title = strings.TrimSpace(strings.Replace(n, ".", " ", -1))
	if n != "" {
		info.Season = strings.TrimSpace(regexp.MustCompile("([a-zA-Z]([0-9])+)").FindString(s))
	}
	if s != "" {
		info.Title = strings.TrimSpace(strings.Replace(strings.Split(n, filepath.Ext(filename))[0], ".", " ", -1))
	}
}

func PullCategorized(root, ignore string, fl *[]File) error {
	m, _ := filepath.Glob(filepath.Join(root, "*"))
	if err := pullPath(m, fl); err != nil {
		return err
	}
	return nil
}

// PullOut will pull a file from source path and write it to destine path. Returns a written bytes and an error if any occurs. If a critical error occurs this function will execute a panic.
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

func MoveFiles(dst, src string, p PullFiles) {
	var fl []File
	// Check source dir
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err = os.MkdirAll(dst, os.ModePerm); err != nil {
			panic(err)
		}
	}
	if err := p(src, dst, &fl); err != nil {
		panic(err)
	}
	for _, f := range fl {
		// PullOut files
		dst := filepath.Join(dst, f.Name)
		fmt.Printf("Creating File: %s\n", f.Name)
		bw, err := PullOut(dst, f.Path)
		if err != nil {
			fmt.Println("Error occurs: ", err)
			return
		}
		fmt.Println("File created, bytes written: ", bw)
	}
}

func DeleteFiles(path, ignore string, p PullFiles) {
	var fl []File
	if err := p(path, ignore, &fl); err != nil {
		panic(err)
	}
	if len(fl) == 0 {
		return
	}
	for i := range fl {
		os.RemoveAll(fl[i].Path)
	}
	fmt.Println("Root folder cleaned.")
}

func Categorize(dst, src string, p PullFiles) {
	var fl []File
	if dst == "" {
		dst = src
	}
	if err := p(src, dst, &fl); err != nil {
		panic(err)
	}
	for _, f := range fl {
		fp := filepath.Join(dst, f.Title, f.Season, f.Name)
		if f.Season == "" { // Raising a flag
			fp = filepath.Join(dst, f.Name)
		}
		// Check source dir
		if _, err := os.Stat(filepath.Dir(fp)); os.IsNotExist(err) {
			if err = os.MkdirAll(filepath.Dir(fp), os.ModePerm); err != nil {
				panic(err)
			}
		}
		fmt.Println("Source: ", f.Path, "Destine: ", fp)
		fmt.Printf("Creating File: %s\n", fp)
		bw, err := PullOut(fp, f.Path)
		if err != nil {
			fmt.Println("Error occurs: ", err)
			return
		}
		fmt.Println("File created, bytes written: ", bw)
	}
}

func main() {
	//MoveFiles(os.Args[2], os.Args[1], PullDir)
	//DeleteFiles(os.Args[1], os.Args[2], PullDir)
	Categorize("", os.Args[1], PullCategorized)
}
