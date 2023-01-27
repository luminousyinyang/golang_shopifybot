# golang_shopifybot

### Preface
This is a basic shopify bot that was built with GoLang. This bot is fully powered by HTTP requests so it doesn't need to use a browser and is extremely lightweight. All the profile and task data is stored in CSV files to ensure quick and easy access to the program. Can create and edit profiles/tasks directly in the program instead of accessing the CSV if you wish to.

### How to use
There is 1 major thing you need to do to get this started up, which should be super easy since its just 2 lines in the main.go file. Keep in mind these edits are just showcasing where the CSV file is in your system so the program can find it. 

Line #s to edit in main.go: 16,17. Variables: tasksCSV AND profilesCSV

once edited you can build or install the program with their respective GoLang commands. "go build" or "go install"

To use the program you must enter a subcommand while starting up the program or right after starting it up. To edit profiles there is a "profile" subcommand, to edit or start tasks there is a "task" subcommand. 

Hope this program is helpful. Enjoy :)
