# golang_shopifybot

### Preface
This is a basic shopify bot that was built with GoLang. This bot is fully powered by HTTP requests so it doesn't need to use a browser and is extremely lightweight. All the profile and task data is stored in CSV files to ensure quick and easy access to the program. Can create and edit profiles/tasks directly in the program instead of accessing the CSV if you wish to.

### How to use
There is 1 major thing you need to do to get this started up, for ease of explanation I will provide the file and line #s you need to edit. Keep in mind all these edits is just showcasing where the CSV file is in your system so the program can find it. 

File + Line #s to edit:
File: taskRun.go, Lines: 37,69,90,132,161
File: profiles.go Lines: 48,110,136,217,237,289
File: taskRunHelpFuncs.go Line: 30

OR instead of finding those lines individually, you can just CTRL+F search with "/Users/userid/path/to" to find those lines. All you need to do on those lines is edit the file path for the CSV files.

once edited you can build or install the program with their respective GoLang commands. "go build" or "go install"

To use the program you must enter a subcommand while starting up the program or right after starting it up. To edit profiles there is a "profile" subcommand, to edit or start tasks there is a "task" subcommand. 

Hope this program is helpful. Enjoy :)
