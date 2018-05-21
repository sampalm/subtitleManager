package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sampalm/subtitleManager/osb"
)

type Flag struct {
	Get     []string
	Options []bool
}

var fg Flag

// flags
const (
	path = iota
	ext
	version
	move
	lang
)

// options
const (
	del = iota
	help
	only
	org
	dl
)

func init() {
	// OSB Package
	dl := flag.Bool("dl", false, "Download subtitles to all selected files.")
	lang := flag.String("lang", "eng", "Set language to download subtitles. Default: 'eng'.")

	// Manager
	p := flag.String("p", "", "Set the root path.")
	e := flag.String("e", ".srt", "Set the extension of the file.")
	v := flag.String("v", "", "Set the version of the subtitle.")
	m := flag.String("m", "", "Only move files to this selected directory.")

	d := flag.Bool("d", false, "Only delete files in selected directory.")
	h := flag.Bool("h", false, "Returns basic instructions to use Subtitle Manager.")
	o := flag.Bool("only", false, "Runs search only into the main path.")
	org := flag.Bool("o", false, "Organize all files in selected directy.")

	flag.Parse()
	fg.Get = []string{*p, *e, *v, *m, *lang}
	fg.Options = []bool{*d, *h, *o, *org, *dl}
}

func main() {
	switch {
	case fg.Options[help]:
		PrintInfo()
		PrintHelp()
		break
	case fg.Get[path] == "":
		log.Println("Path must be defined.")
		break
	case fg.Options[dl]:
		var cErr = make(chan error)
		var cHash = make(chan []string)
		var cSub = make(chan []osb.Subtitle)
		files, err := fg.FetchAll()
		CheckErr("Request", 1, err)
		for _, file := range files {
			go func(file *os.File) {
				hash, size, err := osb.HashFile(file)
				if err != nil {
					cErr <- err
				}
				cHash <- []string{fmt.Sprintf("%x", hash), fmt.Sprint(size)}
			}(file)
		}
		for range files {
			select {
			case hashSize := <-cHash:
				cDone := make(chan string)
				subs, err := osb.SearchHashSub(hashSize[0], hashSize[1], fg.Get[lang])
				if err != nil {
					fmt.Fprintf(os.Stdout, "Request: searchHashSub: %s", err.Error())
				}
				for _, sub := range subs {
					go func(sub osb.Subtitle) {
						//fmt.Fprintf(os.Stdout, "SUB FOUND: %s - Language: %s\n", sub.FileName, sub.LanguageName)
						err := osb.DownloadSub(&sub)
						if err != nil {
							cErr <- err
						}
						// CREATE SUB
						path := fg.Get[path]
						if fg.Get[move] != "" {
							path = fg.Get[move]
						}
						create(sub.FileName, sub.Body.Bytes(), path)
						cDone <- sub.FileName
					}(sub)
				}
				for range subs {
					select {
					case err := <-cErr:
						fmt.Fprintf(os.Stdout, "Request: %s\n", err.Error())
					case file := <-cDone:
						fmt.Fprintf(os.Stdout, "Request: File %s downloaded with success!\n", file)
					}
				}
			case err := <-cErr:
				fmt.Fprintf(os.Stdout, "Request: hashFile: %s\n", err.Error())
			case subs := <-cSub:
				for _, sub := range subs {
					fmt.Fprintf(os.Stdout, "Sub name: %s\nSub link: %s\n", sub.FileName, sub.DownloadLink)
				}
			}
		}
		break
	case fg.Options[org]:
		if !ConfirmAction("This task will move all files that will be found") {
			log.Println("Task canceled.")
			os.Exit(1)
		}
		subs, err := fg.Getall()
		CheckErr("path GetAll", 1, err)
		fg.OrganizeAll(subs)
		break
	case fg.Get[move] != "":
		if fg.Get[version] != "" {
			if !ConfirmAction("This task will move any files thats match with given version") {
				log.Println("Task canceled.")
				os.Exit(1)
			}
		} else if !ConfirmAction("This task will move all files that will be found") {
			log.Println("Task canceled.")
			os.Exit(1)
		}
		subs, err := fg.Getall()
		CheckErr("path GetAll", 1, err)
		fg.Moveall(subs)
		break
	case fg.Options[del]:
		if fg.Get[version] != "" {
			if !ConfirmAction("This task will delete any files thats match with given version") {
				log.Println("Task canceled.")
				os.Exit(1)
			}
		} else if !ConfirmAction("This task will delete all files that will be found") {
			log.Println("Task canceled.")
			os.Exit(1)
		}
		fg.Deleteall(nil)
		break
	case fg.Get[path] != "":
		if fg.Get[version] == "" {
			log.Println("Must define a version to use only path flag.")
			return
		}
		if !ConfirmAction("This task will delete any files that doesn't match with given version") {
			log.Println("Task canceled.")
			os.Exit(1)
		}
		// Execute path
		subs, err := fg.Getall()
		CheckErr("path GetAll", 1, err)
		// delete files
		fg.Deleteall(subs)
		// copy files
		cDone := make(chan bool)
		for _, sub := range subs {
			go create(sub.Name, sub.Body.Bytes(), fg.Get[path])
			cDone <- true
		}
		for range subs {
			<-cDone
		}
		break
	default:
		log.Println("Must define a flag before use, for help use -h.")
	}
	CloseFilter()
}
