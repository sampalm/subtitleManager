package main

import (
	"flag"
	"log"
	"os"
)

type Flag struct {
	Get     []string
	Options []bool
}

var fg Flag

const (
	// flags
	path    = 0
	ext     = 1
	version = 2
	move    = 3

	// options
	del  = 0
	help = 1
	only = 2
)

func init() {
	p := flag.String("p", "", "Set the root path.")
	e := flag.String("e", ".srt", "Set the extension of the file.")
	v := flag.String("v", "", "Set the version of the subtitle.")
	m := flag.String("m", "", "Only move files to this selected directory.")

	d := flag.Bool("d", false, "Only delete files in selected directory.")
	h := flag.Bool("h", false, "Returns basic instructions to use Subtitle Manager.")
	o := flag.Bool("only", false, "Runs search only into the main path.")

	flag.Parse()
	fg.Get = []string{*p, *e, *v, *m}
	fg.Options = []bool{*d, *h, *o}
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
		fg.Deleteall()
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
		fg.Deleteall()
		// copy files
		for _, sub := range subs {
			wg.Add(1)
			go create(sub, fg.Get[path])
		}
		wg.Wait()
		break
	default:
		log.Println("Must define a flag before use, for help use -h.")
	}
	CloseFilter()
}
