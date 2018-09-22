package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("You must provide a command task.")
	}

	switch os.Args[1] {
	case "copy", "c":
		if err := copy(os.Args[2:]); err != nil {
			log.Fatalln(err)
		}

	case "delete", "dl":
		if err := delete(os.Args[2:]); err != nil {
			log.Fatalln(err)
		}

	case "categorize", "ct":
		if err := categorize(os.Args[2:]); err != nil {
			log.Fatalln(err)
		}

	case "download", "d":
		if err := download(os.Args[2:]); err != nil {
			log.Fatalln(err)
		}

	case "query", "q":
		if err := query(os.Args[2:]); err != nil {
			log.Fatalln(err)
		}

	case "help", "h":
		fmt.Println("Help command choosed.")
		fmt.Printf("\n[subtitleManager command subcommands...]\n\n*** Commands:\n\tcopy, c\n\tdelete, dl\n\tcategorize, ct\n\tdownload, d\n\tquery, q\n*** Subcomands:\nUse [subtitleManager command] to show all subcommands from that command.")
		os.Exit(0)

	default:
		flag.PrintDefaults()
		os.Exit(0)
	}

}

func copy(arg []string) error {
	log.Println("Copy task choosed.")

	copyCommand := flag.NewFlagSet("copy", flag.ExitOnError)
	copySrc := copyCommand.String("src", "", "Defines source folder to copy.")
	copyDst := copyCommand.String("dst", "", "Defines destine folder to copy.")
	copyOnly := copyCommand.Bool("only", false, "Ignores all subfolders inside source folder.")
	if err := copyCommand.Parse(arg); err != nil {
		return err
	}

	if *copySrc == "" || *copyDst == "" {
		copyCommand.PrintDefaults()
		return nil
	}

	if *copyOnly {
		MoveFiles(*copyDst, *copySrc, PullDir)
		return nil
	}
	MoveFiles(*copyDst, *copySrc, PullTreeDir)

	return nil
}

func delete(arg []string) error {
	log.Println("Delete task choosed.")

	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	deletePath := deleteCommand.String("path", "", "Defines root folder to delete.")
	deleteIgn := deleteCommand.String("ignore", "", "Ignores this folder.")
	deleteOnly := deleteCommand.Bool("only", false, "Ignores all subfolders inside source folder.")
	if err := deleteCommand.Parse(arg); err != nil {
		return err
	}

	if *deletePath == "" {
		deleteCommand.PrintDefaults()
		return nil
	}

	if *deleteOnly {
		DeleteFolder(*deletePath, *deleteIgn, PullDir)
		return nil
	}
	DeleteFolder(*deletePath, *deleteIgn, PullTreeDir)

	return nil
}

func categorize(arg []string) error {
	log.Println("Categorize task choosed.")

	categorizeCommand := flag.NewFlagSet("categorize", flag.ExitOnError)
	categorizeSrc := categorizeCommand.String("src", "", "Defines source folder to organize.")
	categorizeDst := categorizeCommand.String("dst", "", "Defines destine folder to organize.")
	if err := categorizeCommand.Parse(arg); err != nil {
		return err
	}

	if *categorizeSrc == "" || *categorizeDst == "" {
		categorizeCommand.PrintDefaults()
		return nil
	}
	Categorize(*categorizeDst, *categorizeSrc, PullCategorized)

	return nil
}

func download(arg []string) error {
	log.Println("Download task choosed.")

	downloadCommand := flag.NewFlagSet("download", flag.ExitOnError)
	downloadPath := downloadCommand.String("path", "", "Defines root file's folder.")
	downloadLang := downloadCommand.String("lang", "eng", "Defines default language to download subtitles.")
	downloadMLang := downloadCommand.String("multi", "", "Defines multiples languages. (Sep by comma)")
	downloadScore := downloadCommand.Int("score", 0, "Defines rating score for subtitles.")
	downloadScan := downloadCommand.Bool("scan", false, "Scans and downloads all subtitles.")
	downloadMaxQueue := downloadCommand.Int("max", 4, "Sets the maximum download limit.")
	if err := downloadCommand.Parse(arg); err != nil {
		return err
	}

	if *downloadPath == "" {
		downloadCommand.PrintDefaults()
		return nil
	}
	c := &Controller{
		RootFolder:      *downloadPath,
		DefaultLanguage: *downloadLang,
		RatingScore:     *downloadScore,
	}
	if *downloadMLang != "" {
		// Splits the user-defined values and returns a map. (eg -multi en,pt,fr)
		c.MultiLanguage = langParse(*downloadMLang)
	}
	if *downloadScan {
		c.ScanFolder = true
		c.QueueMax = *downloadMaxQueue
		GetHashFiles(c, *downloadPath, PullTreeDir)
		return nil
	}
	GetHashFiles(c, *downloadPath, PullDir)

	return nil
}

func query(arg []string) error {
	log.Println("Query task choosed.")

	queryCommand := flag.NewFlagSet("query", flag.ExitOnError)
	queryPath := queryCommand.String("path", "", "Defines root to save the downloaded subtitles.")
	queryName := queryCommand.String("name", "", "Defines query's name.")
	querySeason := queryCommand.String("season", "", "Defines query's season.")
	queryEpisode := queryCommand.String("episode", "", "Defines query's episode.")
	queryLang := queryCommand.String("lang", "eng", "Defines default language.")
	queryMLang := queryCommand.String("multi", "", "Defines multiples languages. (Sep by comma)")
	queryScore := queryCommand.Int("score", 0, "Defines rating score for subtitles.")
	if err := queryCommand.Parse(arg); err != nil {
		return err
	}

	if *queryPath == "" || *queryName == "" {
		queryCommand.PrintDefaults()
		return nil
	}
	c := &Controller{
		RootFolder:      *queryPath,
		DefaultLanguage: *queryLang,
		RatingScore:     *queryScore,
	}
	if *queryMLang != "" {
		// Splits the user-defined values and returns a map. (eg -multi en,pt,fr)
		c.MultiLanguage = langParse(*queryMLang)
	}
	params := url.Values{}
	params.Add("query", *queryName)
	if *querySeason != "" {
		params.Add("season", *querySeason)
	}
	if *queryEpisode != "" {
		params.Add("episode", *queryEpisode)
	}
	DownloadQuery(c, params)

	return nil
}

func langParse(value string) map[string]bool {
	ml := map[string]bool{}
	for _, str := range strings.Split(strings.TrimSpace(value), ",") {
		ml[str] = true
	}
	return ml
}
