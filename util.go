package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// PrintInfo print out subtitle manager logo and some infos.
func PrintInfo() {
	messageInfo := []byte("+--------------------------------------------------------------------------+\n|█▀▀ █░░█ █▀▀▄ ▀▀█▀▀ ░▀░ ▀▀█▀▀ █░░ █▀▀ █▀▀   █▀▀ ░▀░ █░░ ▀▀█▀▀ █▀▀ █▀▀█ █▀▀|\n|▀▀█ █░░█ █▀▀▄ ░░█░░ ▀█▀ ░░█░░ █░░ █▀▀ ▀▀█   █▀▀ ▀█▀ █░░ ░░█░░ █▀▀ █▄▄▀ ▀▀█|\n|▀▀▀ ░▀▀▀ ▀▀▀░ ░░▀░░ ▀▀▀ ░░▀░░ ▀▀▀ ▀▀▀ ▀▀▀   ▀░░ ▀▀▀ ▀▀▀ ░░▀░░ ▀▀▀ ▀░▀▀ ▀▀▀|\n+--------------------------------------------------------------------------+\n")

	fmt.Fprintf(os.Stdout, "%s", messageInfo)
}

// PrintHelp print out basic instructions of how to use subtitle manager.
func PrintHelp() {
	fmt.Fprintf(os.Stdout, `
	
Usage:
	subtitlemanager [flags] [options...]

Flags:
	p		Set the main path. eg: "./subs/" (Required)
	e		Set the file extension. eg: ".sub" (Optional)
	v		Set the subtitle version. eg: "720p-WEB". (Required)
	m		Set the folder to where the files will be moved. eg: "./my-subs/" (Optional)

Options:
	org		Organize all files in selected directy. (Optional)
	d		Delete files in selected directory and sub-directories. (Optional)
	h		Returns basic instructions to use Subtitle Manager. (Optional)
	only	Restrict all set options to be executed only in the main path. It will ignore all sub-directories. (Optional)
	
Additional:
	If any problems occur, or you have any suggestions please send me an email: samuelpalmeira@outlook.com
	`)
}

// CreateLogFile will create or re-write a log with all error that occurs while the program has been executed.
func CreateLogFile(logs []Log) error {
	time := time.Now().Format("01/_2/2006 - 15:04:05")
	file := filepath.Join("./log.txt")
	f, err := os.OpenFile(file, syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT, 0777)
	if err != nil {
		return err
	}

	for _, log := range logs {
		if _, err := fmt.Fprintf(f, "**************** LOG ERROR %v ****************\r\n**** %s: %s\r\n", time, log.Func, log.Err); err != nil {
			return err
		}
	}
	return nil

}

// ConfirmAction fulfill a menssage and print it out to user confirm a task action and return a boolean, if users input is anything then 'y' this function will return a false boolean otherwise it will be true.
func ConfirmAction(message string) bool {
	var confirm string
	fmt.Fprintf(os.Stdout, "%s. Anything besides 'y' will cancel this task Confirm?(y/N): ", message)
	fmt.Scan(&confirm)
	if confirm == "y" {
		return true
	}
	return false
}

// CheckErr check if exits an error and print it out with the function that called CheckErr. This function will forces the program to close with os.Exit function.
func CheckErr(caller string, gravity int, err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s: %s\n", caller, err.Error())
		switch gravity {
		case 0:
			break
		case 1:
			os.Exit(1)
			break
		case 2:
			panic("_BAD_FUNCTION")
		}

	}
}

// CloseFilter will close the program and report if any error occurred in its execution.
func CloseFilter() {
	if errFound {
		fmt.Fprintln(os.Stdout, "*** TASK COMPLETED with some errors. Open log file to see all program execution errors ***")
		if err := CreateLogFile(logs); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}
	fmt.Fprintln(os.Stdout, "*** TASK COMPLETED without any errors. ***")
}
