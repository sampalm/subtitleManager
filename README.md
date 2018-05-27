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

To use the application you must use flags to define directories and different functions of the application. For more information use `-h` tag to see a set of instructions explaining the main functions of the application and of each flag or option. You can also combine different flags and options to make multiples task at once, but there are an rank of priority which application uses to execute it fuction, see more information below at **#Combining options and flags** topic.

### Usage:
`subtitlemanager [flags] [options...]`

#### Flags [rank] [compatibility]:
* **-p [-] [all]**	
    * Execute the program with "-p" flag to set the main path that contains the subtitles. eg: "./subs/" *(Required)*
    * See some examples of what you can combine with this flag:
        * `-p=./your-path/ -m=./new-path/`
            -  Moving all files inside `your-path` to `new-path`.
        * `-p=./your-path/ -e=.sub`
            -  Get only files with `.sub` extension.
        * `-p=./your-path/ -v=720p`
            - Get only files that are from `720p` version.    
* **-e [-] [all]**	
    * Use it to set the file extension. eg: ".sub" *(Optional)*
* **-v [-] [all]**	
    * Use it to set the subtitle version. eg: "720p-WEB". *(Required)*
* **-move [5] [all]**	
    * Use it to set the folder to where the files will be moved. eg: "./my-subs/". *(Optional)*
	    * If you set this flag no files will be deleted, only the files that matched will be moved.
* **-rate [-] [3]**
    * Set a minimum rating to download subtitles. (Optional)
* **-lang [-] [3]**
    * Set the language that the subtitles will be downloaded. (Optional)
* **-mlang [-] [3]**
    * Set all the languages that the subtitles will be downloaded. (Optional)
* **-sn [-] [2]**
    * Search for subtitles with this name to download. (Optional)
* **-ss [7] [2]**
    * Search for subtitles in this season to download. (Required to -se)
* **-se [8] [2-7]**
    * Search for subtitles of this episode to download. (Optional)

#### Options [rank] [compatibility]:
* **-d [6] [4-6]**	
    * Delete all files inside the path set by flag[-p]. (Optional)
    * See some examples of what you can combine with this flag:
        * `-p=./your-path/ -d -only`
            - Will delete files in `your-path` but will ignore any sub-folder and files inside it.
        * `-p=./your-path/ -d -e=.srt`
            - Will delete all files inside `your-path` that have `.srt` as extension.
* **-help [1] [none]**	
    * Show all flags and options that are available to use. (Optional)
* **-only [-] [all]**	
    * Execute search only into the main path, it will ignore all subfolders and files within. (Optional)
* **-org [4] [-]**	
    * Organize all files by title and creates subfolders to each season. (Optional)
* **-search [2] [-]**
    * Set search mode on and enable seach options. (Require to -s* flags)
* **-force [9] [2-3]**
    * Will force all downloads and ignores confirm messages. (Optional)
* **-sl [10] [2]**
    * Will allow you to choose only one subtitle to download. (Optional)
* **-dd [3] [-]**
    * Download subtitles to all found files. (Optional)
    * See some examples of what you can combine with this flag:
        * `-p=./your-path/ -dd -m=./new-path/`
            - All your downloaded subtitles will be placed in `new-path` folder.
        * `-p=./your-path/ -dd -lang=spa`
            - Only will download `spanish` subtitles. You can see others language [ISO 639-2 Code here](http://www.loc.gov/standards/iso639-2/php/code_list.php).


#### Combining options and flags:
You can combine both to perform restricted or concise tasks in this application, the main thing you need to know before use this method is the raking of priorities that the application use to execute the flags. As you can see the **[rank]** determine order of execution in the flow of the application and **[compatibility]** informs what are the possible combinations between flags and options which means that probably if you try to use other flags besides these compatible will not have any effect.
	
### Additional:
If any problems occur, or you have any suggestions please send me an email: samuelpalmeira@outlook.com
	
