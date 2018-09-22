package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"syscall"
)

const bufSize = 4 * 1024 * 1024

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

type Puller func(root, ignore string, fl *[]File) error

func crawler(ignore string, fl *[]File) filepath.WalkFunc {
	exts := map[string]bool{".srt": true, ".sub": true, ".sbv": true}
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
	exts := map[string]bool{".srt": true, ".sub": true, ".sbv": true}
	ignore = filepath.Join(ignore)
	f, err := os.Open(root)
	if err != nil {
		log.Fatalln(err)
	}
	// read all the files from dir
	m, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	sort.Strings(m)
	for i := range m {
		p := filepath.Join(root, m[i])
		d, err := os.Lstat(p)
		if err != nil {
			return err
		}
		if filepath.Dir(p) == ignore {
			continue
		}
		if _, ok := exts[filepath.Ext(p)]; !ok {
			continue
		}
		f := File{
			Name: d.Name(),
			Path: p,
			Ext:  filepath.Ext(p),
		}
		*fl = append(*fl, f)
	}
	if len(*fl) == 0 {
		log.Fatalln("None files found!")
	}
	return nil
}

func pullFiles(folder []string, fl *[]File) error {
	exts := map[string]bool{".srt": true, ".sub": true, ".sbv": true}
	for i := range folder {
		d, err := os.Lstat(folder[i])
		if err != nil {
			return err
		}
		if d.IsDir() {
			path, _ := globDir(folder[i])
			pullFiles(path, fl)
			continue
		}
		if _, ok := exts[filepath.Ext(d.Name())]; !ok {
			continue
		}
		f := File{
			Name: d.Name(),
			Path: folder[i],
			Ext:  filepath.Ext(d.Name()),
		}

		s := regexp.MustCompile(`(.[0-9]{4,}.([a-z-A-Z]|).*|(.[a-zA-Z]([0-9]+)){2,})`).FindString(f.Name)
		f.Title = strings.TrimSpace(strings.Replace(strings.Split(f.Name, s)[0], ".", " ", -1))
		ss := regexp.MustCompile(`([.][a-zA-Z]([0-9]{2}))`).FindStringSubmatch(s)
		if ss != nil {
			f.Season = "Season " + ss[2]
		}

		*fl = append(*fl, f)
	}
	return nil
}

func globDir(address string) ([]string, error) {
	f, err := os.Open(address)
	if err != nil {
		return nil, fmt.Errorf("Open returns -> %s", err)
	}
	// read all the files from dir.
	m, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, fmt.Errorf("Readdirnames returns -> %s", err)
	}
	// get full directory of the file file.
	for i := 0; i < len(m); i++ {
		m[i] = filepath.Join(address, m[i])
	}
	sort.Strings(m)
	return m, nil
}

func PullCategorized(root, ignore string, fl *[]File) error {
	m, _ := globDir(root)
	if err := pullFiles(m, fl); err != nil {
		return err
	}
	return nil
}

func PullOut(dst, src string) (int64, error) {
	buf := make([]byte, bufSize)
	bw := int64(0)

	fs, err := os.OpenFile(src, syscall.O_RDONLY, os.ModePerm)
	if err != nil {
		return bw, err
	}
	defer func() {
		if err := fs.Close(); err != nil {
			panic(err)
		}
	}()
	stt, err := fs.Stat()
	if err != nil {
		return bw, err
	}
	chunckSize := stt.Size() / 1024

	fd, err := os.OpenFile(dst, syscall.O_CREAT|syscall.O_WRONLY, os.ModePerm)
	if err != nil {
		return bw, err
	}
	defer func() {
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}()

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

func MoveFiles(dst, src string, p Puller) {
	var fl []File
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err = os.MkdirAll(dst, os.ModePerm); err != nil {
			panic(err)
		}
	}
	if err := p(src, dst, &fl); err != nil {
		panic(err)
	}
	for _, f := range fl {
		dst := filepath.Join(dst, f.Name)
		log.Printf("Creating File: %s\n", f.Name)
		bw, err := PullOut(dst, f.Path)
		if err != nil {
			log.Println("MovileFiles => ", err)
			return
		}
		log.Println("File created, bytes written: ", bw)
	}
}

func remove(fl []File) error {
	var dl map[string]bool
	for i := range fl {
		if _, ok := dl[filepath.Dir(fl[i].Path)]; !ok {
			dl[filepath.Dir(fl[i].Path)] = true
			log.Printf("Folder %s has been cleaned.\n", filepath.Dir(fl[i].Path))
		}
		if err := os.RemoveAll(fl[i].Path); err != nil {
			return err
		}
	}
	return nil
}

func deleteFile(path string, stack map[string]bool) error {
	fmt.Println("deleteFile reached....")
	if _, ok := stack[filepath.Dir(path)]; !ok {
		stack[filepath.Dir(path)] = true
		log.Printf("Cleaning %s folder.\n", filepath.Dir(path))
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

func DeleteFolder(path, ignore string, p Puller) {
	var fl []File
	if err := p(path, ignore, &fl); err != nil {
		panic(err)
	}
	if len(fl) == 0 {
		return
	}
	if cerr := remove(fl); cerr != nil {
		panic(cerr) // raising a flag
	}
}

func Categorize(dst, src string, p Puller) {
	var fl []File
	if dst == "" {
		dst = src
	}
	if err := p(src, dst, &fl); err != nil {
		panic(err)
	}
	for _, f := range fl {
		fp := filepath.Join(dst, f.Title, f.Season, f.Name)
		if f.Season == "" {
			fp = filepath.Join(dst, f.Title, f.Name)
		}
		if _, err := os.Stat(filepath.Dir(fp)); os.IsNotExist(err) {
			if err = os.MkdirAll(filepath.Dir(fp), os.ModePerm); err != nil {
				panic(err)
			}
		}
		log.Printf("Creating File: %s\n", fp)
		bw, err := PullOut(fp, f.Path)
		if err != nil {
			log.Println("PullOut => ", err)
			return
		}
		log.Println("File created, bytes written: ", bw)
		dl := map[string]bool{}
		if err := deleteFile(f.Path, dl); err != nil {
			log.Println("deleteFile => ", err)
			return
		}
	}
	fmt.Println("Everything worked out... cya ;)")
}
