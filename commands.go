package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var in io.Reader = os.Stdin
var c = Controller{}

func SliceString(value string) []string {
	return strings.Split(strings.TrimSpace(value), ",")
}

func main() {
	// Set subcommands
	copyCommand := flag.NewFlagSet("copy", flag.ExitOnError)
	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	categorizeCommand := flag.NewFlagSet("categorize", flag.ExitOnError)

	//queryCommand := flag.NewFlagSet("query", flag.ExitOnError)
	downloadCommand := flag.NewFlagSet("download", flag.ExitOnError)

	// Copy subcommand flags
	copySrc := copyCommand.String("src", "", "Defines source folder to copy.")
	copyDst := copyCommand.String("dst", "", "Defines destine folder to copy.")
	copyOnly := copyCommand.Bool("only", false, "Ignores all subfolders inside source folder.")
	// Delete subcommand flags
	deletePath := deleteCommand.String("path", "", "Defines root folder to delete.")
	deleteIgn := deleteCommand.String("ignore", "", "Ignores this folder.")
	deleteOnly := deleteCommand.Bool("only", false, "Ignores all subfolders inside source folder.")
	// Categorize subcommand flags
	categorizeSrc := categorizeCommand.String("src", "", "Defines source folder to organize.")
	categorizeDst := categorizeCommand.String("dst", "", "Defines destine folder to organize.")
	// Query subcommand flags

	// Download subcommand flags
	downloadPath := downloadCommand.String("path", "", "Defines root file's folder.")
	downloadLang := downloadCommand.String("lang", "eng", "Defines default language to download subtitles.")
	downloadMLang := downloadCommand.String("multi", "", "Defines multiples languages. (Sep by comma)")
	downloadScore := downloadCommand.Int("score", 0, "Defines rating score for subtitles.")
	// Check that subcommand has been provided
	if len(os.Args) < 2 {
		fmt.Println("You must provide a command task.")
		os.Exit(1)
	}

	// Select the subcommand
	switch os.Args[1] {
	case "copy":
		fmt.Println("Copy task choosed.")
		copyCommand.Parse(os.Args[2:])
	case "delete":
		fmt.Println("Delete task choosed.")
		deleteCommand.Parse(os.Args[2:])
	case "categorize":
		fmt.Println("Categorize task choosed.")
		categorizeCommand.Parse(os.Args[2:])
	case "download":
		fmt.Println("Download task choosed.")
		downloadCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Handle CopyCommand flags
	if copyCommand.Parsed() {
		if *copySrc == "" || *copyDst == "" {
			copyCommand.PrintDefaults()
			os.Exit(1)
		}
		if *copyOnly {
			MoveFiles(*copyDst, *copySrc, PullDir)
		} else {
			MoveFiles(*copyDst, *copySrc, PullTreeDir)
		}
	}

	// Handle DeleteCommand flags
	if deleteCommand.Parsed() {
		if *deletePath == "" {
			deleteCommand.PrintDefaults()
			os.Exit(1)
		}
		if *deleteOnly {
			DeleteFolder(*deletePath, *deleteIgn, PullDir)
		} else {
			DeleteFolder(*deletePath, *deleteIgn, PullTreeDir)
		}
	}

	// Handle CategorizeCommand flags
	if categorizeCommand.Parsed() {
		if *categorizeSrc == "" || *categorizeDst == "" {
			categorizeCommand.PrintDefaults()
			os.Exit(1)
		}
		Categorize(*categorizeDst, *categorizeSrc, PullCategorized)
	}

	// Handle DownloadCommand flags
	if downloadCommand.Parsed() {
		if *downloadPath == "" {
			downloadCommand.PrintDefaults()
			os.Exit(1)
		}
		c := Controller{
			RootFolder:      *downloadPath,
			DefaultLanguage: *downloadLang,
			MultiLanguage:   SliceString(*downloadMLang),
			RatingScore:     *downloadScore,
		}
		GetHashFiles(&c, *downloadPath, PullDir)
	}
}
