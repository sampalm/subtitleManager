# Subtitle Manager

Subtitle Manager is an application to manage movie captions, series, anime, etc more efficiently. You can also download the captions directly to the application and organize them by categories.

The application is still in the testing phase and also does not have all the functions. This is a project for learning purposes only, I don't have the goal of turning it into something beyond that.

I appreciate any suggestions or tips to improve application performance or improve the code itself.

## Setup
If you want to change something in the code you can use: 
```
go get -u github.com/sampalm/subtitleManager.
```

Or you also can download the executable version **[Outdated]**.

## How to use

To use the application you must use flags to define directories and different functions of the application. For more information Use `-h` tag to see a detailed set of instructions explaining the main functions of the application.

### Flags: 
* p "Set main path of the files" 
* e "Set extension" 
* v "Set version" 
* m "Set path to move files"

### Usage:
`subtitlemanager [flags]`

### Flags info:
* **-p**	
    * Execute the program with "-p" flag to set the main path that contains the subtitles. eg: "./subs/" *(Required)*
* **-e**	
    * Use it to set the file extension. eg: ".sub" *(Optional)*
* **-v**	
    * Use it to set the subtitle version. eg: "720p-WEB". *(Optional)*
* **-m**	
    * Use it to set the folder to where the files will be moved. eg: "./my-subs/". *(Optional)*
	    * If you set this flag no files will be deleted, only the files that matched will be moved. 
	
### Additional:
If any problems occur, or you have any suggestions please send me an email: samuelpalmeira@outlook.com
	
