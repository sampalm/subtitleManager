package main

import "fmt"

func PrintInfo() {
	messageInfo := []byte("+--------------------------------------------------------------------------+\n|█▀▀ █░░█ █▀▀▄ ▀▀█▀▀ ░▀░ ▀▀█▀▀ █░░ █▀▀ █▀▀   █▀▀ ░▀░ █░░ ▀▀█▀▀ █▀▀ █▀▀█ █▀▀|\n|▀▀█ █░░█ █▀▀▄ ░░█░░ ▀█▀ ░░█░░ █░░ █▀▀ ▀▀█   █▀▀ ▀█▀ █░░ ░░█░░ █▀▀ █▄▄▀ ▀▀█|\n|▀▀▀ ░▀▀▀ ▀▀▀░ ░░▀░░ ▀▀▀ ░░▀░░ ▀▀▀ ▀▀▀ ▀▀▀   ▀░░ ▀▀▀ ▀▀▀ ░░▀░░ ▀▀▀ ▀░▀▀ ▀▀▀|\n+--------------------------------------------------------------------------+\n")

	fmt.Printf("%s", messageInfo)
}
