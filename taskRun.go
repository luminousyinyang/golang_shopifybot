package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strconv"
)

func TaskCommands() *ProfileCommand {
	taskc := &ProfileCommand{
		fs: flag.NewFlagSet("task", flag.ContinueOnError),
	}
// if I wish to add flags later in dev
	// taskc.fs.BoolVar(&taskc.name, "create", false, "create a new task -create")
	// taskc.fs.BoolVar(&taskc.name, "start", false, "start an existing task -edit")
	// taskc.fs.BoolVar(&taskc.name, "edit", false, "edit a current task property -delete")

	return taskc
}

type Task struct {
	id string
	url string
	size string
	profile string
}
	
func TaskCreate() {
	fileRD, err := os.Open(tasksCSV)
	if err != nil {
		panic(err)
	}
	defer fileRD.Close()

	taskRecord := csvRead(fileRD)

	var id, url, size, profile string
	scanner := bufio.NewScanner(os.Stdin)

	prevId := 0
	// gets prev id to make sure task ids go in sequential order
	if len(taskRecord) > 0 {
		prevId, err = strconv.Atoi(taskRecord[len(taskRecord) - 1][0])
		if err != nil {
			panic(err)
		}
	}

	id = strconv.Itoa(prevId + 1)
	fmt.Print("Task Product URL: ")
	scanner.Scan()
	url = scanner.Text()
	fmt.Print("Task Product Size: ")
	scanner.Scan()
	size = scanner.Text()
	fmt.Print("Task Profile id: ")
	scanner.Scan()
	profile = scanner.Text()


	fileAPWR, err := os.OpenFile(tasksCSV, os.O_APPEND|os.O_CREATE|os.O_WRONLY , 0644)
	if err != nil {
		panic(err)
	}
	defer fileAPWR.Close()

	taskWriter := csv.NewWriter(fileAPWR)
	err = taskWriter.Write([]string{id, url, size, profile})
	if err != nil {
		panic(err)
	}
	taskWriter.Flush()

	fmt.Println()
	fmt.Println("Task successfully created...")
	
	fmt.Printf("{id: %s, url: %s, size: %s, profile: %s}", id, url, size, profile)
	fmt.Println()
}

func TaskEdit() {
	fileRD, err := os.Open(tasksCSV)
	if err != nil {
		panic(err)
	}
	defer fileRD.Close()

	taskRecord := csvRead(fileRD)


	urlidx, sizeidx, profileidx := 1,2,3
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Print out all tasks into the console? Y/N")
	scanner.Scan()
	if (scanner.Text() == "Y") {
		readAllTasks(taskRecord)
	} 

	fmt.Print("Enter task id you would like to edit: ")
	scanner.Scan()
	editId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}
	fmt.Print("Enter task property you would like to edit (can't edit id): ")
	scanner.Scan()
	editProp := scanner.Text()
	fmt.Print("Enter value to update task property to: ")
	scanner.Scan()
	editPropVal := scanner.Text()

	switch editProp {
	case "id":
		log.Fatalln("cant edit id")
	case "url":
		taskRecord[editId-1][urlidx] = editPropVal
	case "size":
		taskRecord[editId-1][sizeidx] = editPropVal
	case "profile":
		taskRecord[editId-1][profileidx] = editPropVal
	}

	fileWR, err := os.Create(tasksCSV)
	if err != nil {
		panic(err)
	}
	defer fileWR.Close()

	writer := csv.NewWriter(fileWR)
	for _, task := range taskRecord {
		err = writer.Write(task)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()

	fmt.Printf("Successfully updated profile %d", editId)
	fmt.Println()
}

func TaskStart() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Input task id you would like to start: ")
	scanner.Scan()
	taskId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	file, err := os.Open(tasksCSV)
	if err != nil {
		panic(err)
	}

	taskRecords := csvRead(file)

	taskProf := getTaskProf(taskRecords, taskId)
	profile := Profile{taskProf[0], taskProf[1], taskProf[2], taskProf[3], taskProf[4], taskProf[5], taskProf[6], taskProf[7], taskProf[8], taskProf[9], taskProf[10], taskProf[11], taskProf[12], taskProf[13], taskProf[14]}
	fmt.Println()
	fmt.Println("Profile info: ", profile)

	taskSize := taskRecords[taskId-1][2]
	taskProdJSON := taskRecords[taskId-1][1] + ".json"

	taskVariant := findTaskVar(taskProdJSON, taskSize)

	regex := regexp.MustCompile(`(https:\/\/)(www.)?[a-zA-z0-9]+\.[a-z]+\/`)
	siteURL := regex.FindAllString(taskProdJSON, -1)[0]
	fmt.Println("Site URL :",siteURL)
	

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	taskATC(client, siteURL, taskVariant)

	req, _ := http.NewRequest("GET", siteURL + "checkout.json", nil)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	limitedReader := &io.LimitedReader{resp.Body, 100000}
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		panic(err)
	}
	
	authToken := getAuthTok(body)
	getCheckURL := resp.Request.URL.String()
	
	authToken = submitContactInfo(client, getCheckURL, authToken, profile)

		
	shipRate := getShipRate(client, siteURL, profile)

	authToken, pay_gateway := submitShipInfo(client, getCheckURL, shipRate, authToken)

	payToken := fetchPayTok(client, profile)


	submitPayment(client, authToken, payToken, pay_gateway, shipRate, getCheckURL, profile)

}