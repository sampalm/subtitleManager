# Subtitle Manager

Subtitle Manager is an application to manage movie captions, series, animes, etc more efficiently. You can also download the captions directly to the application and organize them by categories.

The application is still in the testing version and also does not have all of your functions. This is a project for learning purposes only, I don't have the goal of turning it into something beyond that.

I appreciate any suggestions or tips to improve application performance or improve the code itself.

## Setup
If you want to change something in the code you can use: 
```
go get github.com/sampalm/subtitleManager
```

Or you also can download the executable version **[Outdated]**.

## How to use

To use the application you must use flags to define directories and different functions of the application. For more information use `-h` tag to see a set of instructions explaining the main functions of the application and of each flag or option. You can also combine different flags and options to make multiples task at once, but there are an rank of priority which application uses to execute it fuction, see more information below at **#Combining options and flags** topic.

### Usage:
`subtitlemanager [flags] [options...]`

#### Flags [rank] [compatibility]:
* **-p [5] [all]**	
    * Execute the program with "-p" flag to set the main path that contains the subtitles. eg: "./subs/" *(Required)*
* **-e [6] [all]**	
    * Use it to set the file extension. eg: ".sub" *(Optional)*
* **-v [7] [all]**	
    * Use it to set the subtitle version. eg: "720p-WEB". *(Required)*
* **-m [3] [2-4-5]**	
    * Use it to set the folder to where the files will be moved. eg: "./my-subs/". *(Optional)*
	    * If you set this flag no files will be deleted, only the files that matched will be moved. 

#### Options [rank] [compatibility]:
* **-d [4] [2-5-6]**	
    * Delete all files inside the path set by flag[-p]. (Optional)
* **-h [1] [none]**	
    * Show all flags and options that are available to use. (Optional)
* **-only [8] [all]**	
    * Execute search only into the main path, it will ignore all subfolders and files within. (Optional)
* **-org [2] [3-4-5]**	
    * Organize all files by title and creates subfolders to each season. (Optional)

#### Combining options and flags:
You can combine both to perform restricted or concise tasks in this application, the main thing you need to know before use this method is the raking of priorities that the application use to execute the flags. As you can see the **[rank]** determine order of execution in the flow of the application and **[compatibility]** informs what are the possible combinations between flags and options which means that probably if you try to use other flags besides these compatible will not have any effect.
	
### Additional:
If any problems occur, or you have any suggestions please send me an email: samuelpalmeira@outlook.com
	
