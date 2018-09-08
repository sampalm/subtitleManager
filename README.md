# Subtitle Manager

Subtitle Manager is an application to manage movie captions, series, animes, etc more efficiently. You can also download the captions directly to the application and organize them by categories.

The application is still in the testing version and also does not have all of your functions. This is a project for learning purposes only, I don't have the goal of turning it into something beyond that.

I appreciate any suggestions or tips to improve application performance or improve the code itself.

## Getting Started
If you want to change something in the code you can use: 
```
go get github.com/sampalm/subtitleManager
```

## How to use

To use the application you must use flags to define directories and different functions of the application. For more information use `-help` tag to see a set of instructions explaining the main functions of the application and of each flag or option. You can also combine different flags and options to make multiples task at once, but there are an rank of priority which application uses to execute it fuction, see more information below.

### Usage:
`subtitlemanager command -subcommands...`

#### Copy Subcommands
* -src `Set your root path.`
* -dst `Set a path to copy the files.`
* -only `Only will copy files inside root path, ignoring all subfolders` 

#### Delete Subcommands
* -path `Set path which have the files that you want to delete.`
* -ignore `Will ignore all files inside this path.`
* -only `Only will delete files inside path, ignoring all subfolders`  

#### Categorize Subcommands
* -src `Set your root path.`
* -dst `Set path to where your files will be moved and categorized.`

#### Download Subcommands
* -path `Defines the path where your files are.`
* -score `Defines the minimum rating.`
* -lang `Defines the language. (Default: 'eng')`
    - You can see others language [ISO 639-2 Code](http://www.loc.gov/standards/iso639-2/php/code_list.php).
* -multi `Defines multiple languages.`
* -scan `It scans the folder and download the subtitles of all the files that it finds.`
* -max `Sets the maximum download queue.`

#### Query Subcommands
* -path `Defines the path to save all downloaded subtitles.`
* -name `Defines query's name.`
* -season `Defines query's season.`
* -episode `Defines query's episode.`
* -score `Defines a minimum rating.`
* -lang `Defines the language. (Default: 'eng')`
    - You can see others language [ISO 639-2 Code](http://www.loc.gov/standards/iso639-2/php/code_list.php).
* -multi `Defines multiple languages.`

### Additional:
If any problems occur, or you have any suggestions please send me an email: samuelpalmeira@outlook.com
	
