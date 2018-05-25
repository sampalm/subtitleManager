package main

import (
	"flag"
	"log"
	"os"
)

type Flag struct {
	Get     []string
	Options []bool
	Const   []int
}

var fg Flag

// flags
const (
	path = iota
	ext
	version
	move
	lang
	mlang
)

// options
const (
	del = iota
	help
	only
	org
	dl
)

// const
const (
	rate = iota
)

func init() {
	// OSB Package
	dl := flag.Bool("dd", false, "Download subtitles to all selected files.")
	lang := flag.String("lang", "eng", "Set language to download subtitles. Default: 'eng'.")
	mlang := flag.String("mlang", "", "Set multiples languages to download subtitles.")
	rate := flag.Int("rate", 0, "Set a minimum rating to download subtitles.")

	// Manager
	p := flag.String("p", "", "Set the root path.")
	e := flag.String("e", ".srt", "Set the extension of the file.")
	v := flag.String("v", "", "Set the version of the subtitle.")
	m := flag.String("move", "", "Only move files to this selected directory.")

	d := flag.Bool("d", false, "Only delete files in selected directory.")
	h := flag.Bool("help", false, "Returns basic instructions to use Subtitle Manager.")
	o := flag.Bool("only", false, "Runs search only into the main path.")
	org := flag.Bool("org", false, "Organize all files in selected directy.")

	flag.Parse()
	fg.Get = []string{*p, *e, *v, *m, *lang, *mlang}
	fg.Options = []bool{*d, *h, *o, *org, *dl}
	fg.Const = []int{*rate}
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
		files, err := fg.FetchAll()
		CheckErr("Request", 1, err)
		fg.SaveAll(files)
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
