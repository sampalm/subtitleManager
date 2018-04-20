package main

import "fmt"

// PrintInfo print out subtitle manager logo and some infos.
func PrintInfo() {
	messageInfo := []byte("+--------------------------------------------------------------------------+\n|█▀▀ █░░█ █▀▀▄ ▀▀█▀▀ ░▀░ ▀▀█▀▀ █░░ █▀▀ █▀▀   █▀▀ ░▀░ █░░ ▀▀█▀▀ █▀▀ █▀▀█ █▀▀|\n|▀▀█ █░░█ █▀▀▄ ░░█░░ ▀█▀ ░░█░░ █░░ █▀▀ ▀▀█   █▀▀ ▀█▀ █░░ ░░█░░ █▀▀ █▄▄▀ ▀▀█|\n|▀▀▀ ░▀▀▀ ▀▀▀░ ░░▀░░ ▀▀▀ ░░▀░░ ▀▀▀ ▀▀▀ ▀▀▀   ▀░░ ▀▀▀ ▀▀▀ ░░▀░░ ▀▀▀ ▀░▀▀ ▀▀▀|\n+--------------------------------------------------------------------------+\n")

	fmt.Printf("%s", messageInfo)
}

// PrintHelp print out basic instructions of how to use subtitle manager.
func PrintHelp() {
	fmt.Printf(`
Flags: -p="Define o caminho a ser executado a busca" -e="Define a extensão do arquivo" -v="Define a versão do arquivo"
=== HOW TO USE ===
* Execute o programa com a flag [-p] para definir o local de execução. eg: "./subs/" (Obrigatorio)
* Você pode definir a extensão do arquivo a ser exportado. eg: ".sub" (Opcional)
* Você pode definir a versão do arquivo. eg: "720p-WEB". (Opcional)
* Você pode definiar a pasta para onde os arquivos serão movidos. eg: "./meus-subs/" (Opcional)
	
Caso ocorra algum problema, ou tenha alguma sugestão por favor me envie um email: samuelpalmeira@outlook.com
	`)
}
