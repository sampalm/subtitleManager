# Subtitle Manager

Subtitle Manager is an application to manage movie captions, series, animes, etc more efficiently. You can also download the captions directly to the application and organize them by categories.

The application is still in the testing version and also does not have all of your functions. This is a project for learning purposes only, I don't have the goal of turning it into something beyond that.

I appreciate any suggestions or tips to improve application performance or improve the code itself.

## Getting Started
If you want to change something in the code you can use: 
```
go get github.com/sampalm/subtitleManager
```

Or you also can download the executable version **[Outdated]**.

## How to use

To use the application you must use flags to define directories and different functions of the application. For more information use `-help` tag to see a set of instructions explaining the main functions of the application and of each flag or option. You can also combine different flags and options to make multiples task at once, but there are an rank of priority which application uses to execute it fuction, see more information below.

### Usage:
`subtitlemanager [flags] [options...]`

#### Base Flag
* -p `Set your root path. (Required)`
* -e `Set subtitle extension.`
* -v `Set subtitle version.` 
* -help `Show all flags and options that are available to use.` 

#### Download Subtitle
* -dd `Download subtitles to all found files inside your root path.`
    * -rate `Defines the minimum rating.`
    * -lang `Defines the language. (Default: English/eng)`
    * -mlang `Defines multiple languages.`
    * -force `Ignores all confirm messages.`
    * -sl `Allow you to choose only one subtitle to download.`
#### Search Subtitles
* -search `Set search mode on which enable search options.`
    * -sn `Defines the name to search.` 
    * -ss `Defines the season to search.`
    * -se `Defines the episode to search. (Require -ss flag)`
    * -force `Ignores all confirm messages.`
#### Manage your Subtitles
* -org `Organize all subtitles by title and season if exists.`
* -move `Move all subtitles to this path.`
* -d `Delete all files inside your root path.`
* -only `Execute all tasks only into your root path, will ignore all subfolders.`


#### Combining flags:
You can combine flags to perform restricted or concise tasks in this application, the main thing you need to know before use this method is the raking of priorities that the application use to execute the flags.
	
### Additional:
If any problems occur, or you have any suggestions please send me an email: samuelpalmeira@outlook.com
	
