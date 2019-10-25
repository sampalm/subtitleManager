# Subtitle Manager

Manage your movie captions, series, animes in an easier way. You can also download the captions directly from the application and organize them.

The application is still in development but I'm always looking to update it when its possible or I'm disposal to. This is a project for learning purposes only.

I appreciate any suggestions or tips to improve application performance or improve the code itself.

## Getting Started
If you want to change something in the code you can use: 
```
go get github.com/sampalm/subtitleManager
```

## How to use

To make use you must choose a flag to define directories and different functions. For more information use `-help` tag to see a set of instructions explaining the main commands available and its own sub-commands. Each command has a set of subcommands as you can see below:

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
* -path `Set the path where to put the downloaded files.`
* -score `Set the minimum rating.`
* -lang `Set the language. (Default: 'eng')`
    - You can see others language [ISO 639-2 Code](http://www.loc.gov/standards/iso639-2/php/code_list.php).
* -multi `Set multiple languages separate by comma.`
* -scan `It scans the folder and download the subtitles of all the files that scan finds.`

#### Query Subcommands
* -path `Set the path to save all downloaded subtitles.`
* -name `Set query's name.`
* -season `Set query's season.`
* -episode `Set query's episode.`
* -score `Set a minimum rating.`
* -lang `Defines the language. (Default: 'eng')`
    - You can see others language [ISO 639-2 Code](http://www.loc.gov/standards/iso639-2/php/code_list.php).
* -multi `Defines multiple languages.`

