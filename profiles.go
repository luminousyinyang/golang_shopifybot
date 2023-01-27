package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)


func ProfileCommands() *ProfileCommand {
	gc := &ProfileCommand{
		fs: flag.NewFlagSet("profile", flag.ContinueOnError),
	}
	// if I wish to add flags later in dev
	// gc.fs.BoolVar(&gc.name, "create", false, "create a new profile -create")
	// gc.fs.BoolVar(&gc.name, "edit", false, "edit a current profile -edit")
	// gc.fs.BoolVar(&gc.name, "delete", false, "delete a current profile -delete")


	return gc
}

type Profile struct {
	id string
	firstName string
	lastName string
	email string
	addyLn1 string
	addyLn2 string
	city string
	state string
	zip string
	tele string
	cardNum string
	cardName string
	cardMonth string
	cardYear string
	cardCVV string
}


func ProfileCreate(){

	file, err := os.Open(profilesCSV)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	record := csvRead(file)

	var id, firstName, lastName, email, addyLn1, addyLn2, city, state, zip, tele, cardNum, cardName, cardMonth, cardYear, cardCVV string
	scanner := bufio.NewScanner(os.Stdin)

	prevId := 0
	if len(record) > 0 {
		prevId, err = strconv.Atoi(record[len(record)-1][0])
		if err != nil {
			panic(err)
		}
	} 
	id = strconv.Itoa(prevId + 1)
	fmt.Print("First Name: ")
	scanner.Scan()
	firstName = scanner.Text()
	fmt.Print("Last Name: ")
	scanner.Scan()
	lastName = scanner.Text()
	fmt.Print("Email: ")
	scanner.Scan()
	email = scanner.Text()
	fmt.Print("Address Line 1: ")
	scanner.Scan()
	addyLn1 = scanner.Text()
	fmt.Print("Address Line 2: ")
	scanner.Scan()
	addyLn2 = scanner.Text()
	fmt.Print("City: ")
	scanner.Scan()
	city = scanner.Text()
	fmt.Print("State: ")
	scanner.Scan()
	state = scanner.Text()
	fmt.Print("Zip Code: ")
	scanner.Scan()
	zip = scanner.Text()
	fmt.Print("Telephone Number: ")
	scanner.Scan()
	tele = scanner.Text()
	fmt.Print("Card Number: ")
	scanner.Scan()
	cardNum = scanner.Text()
	fmt.Print("Full Card Name: ")
	scanner.Scan()
	cardName = scanner.Text()
	fmt.Print("Card Expiration Month: ")
	scanner.Scan()
	cardMonth = scanner.Text()
	fmt.Print("Card Expiration Year: ")
	scanner.Scan()
	cardYear = scanner.Text()
	fmt.Print("Card CVV: ")
	scanner.Scan()
	cardCVV = scanner.Text()
	
	profileCsv, err := os.OpenFile(profilesCSV, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer profileCsv.Close()


	profile := Profile{id: id, firstName: firstName, lastName: lastName, email: email, addyLn1: addyLn1, addyLn2: addyLn2, city: city, state: state, zip: zip, tele: tele, cardNum: cardNum, cardName: cardName, cardMonth: cardMonth, cardYear: cardYear, cardCVV: cardCVV}


	writer := csv.NewWriter(profileCsv)
	line := []string{profile.id, profile.firstName, profile.lastName, profile.email, profile.addyLn1, profile.addyLn2, profile.city, profile.state, profile.zip, profile.tele, profile.cardNum, profile.cardName, profile.cardMonth, profile.cardYear, profile.cardCVV}
	err = writer.Write(line)
	if err != nil {
		panic(err)
	}
	writer.Flush()

	fmt.Println()
	fmt.Println("Profile created!...")
	fmt.Printf("id: %s, firstName %s, lastName %s, email %s, addyLn1: %s, addyLn2: %s, city: %s, state: %s, zip: %s, tele: %s, cardNum: %s, cardName: %s, cardMonth: %s, cardYear: %s, cardCVV: %s", id, firstName, lastName, email, addyLn1, addyLn2, city, state, zip, tele, cardNum, cardName, cardMonth, cardYear, cardCVV)
	fmt.Println()
}


func ProfileEdit() {
	file, err := os.Open(profilesCSV)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	record := csvRead(file)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Print out all profiles into the console? Y/N")
	scanner.Scan()
	if (scanner.Text() == "Y") {
	// SEPERATING LINE BETWEEN INPUTS AND PROFILE OUTPUTS
	fmt.Println("--------------------------------")
	var profiles []Profile
	for idx, item := range record {
		profile := Profile{id: item[0], firstName: item[1], lastName: item[2], email: item[3], addyLn1: item[4], addyLn2: item[5], city: item[6], state: item[7], zip: item[8], tele: item[9], cardNum: item[10], cardName: item[11], cardMonth: item[12], cardYear: item[13], cardCVV: item[14]}
		profiles = append(profiles, profile)

		fmt.Println()
		fmt.Printf("id: %s, firstName %s, lastName %s, email %s, addyLn1: %s, addyLn2: %s, city: %s, state: %s, zip: %s, tele: %s, cardNum: %s, cardName: %s, cardMonth: %s, cardYear: %s, cardCVV: %s",
		profiles[idx].id, profiles[idx].firstName, profiles[idx].lastName, profiles[idx].email, profiles[idx].addyLn1, profiles[idx].addyLn2, profiles[idx].city,
		profiles[idx].state, profiles[idx].zip, profiles[idx].tele, profiles[idx].cardNum, profiles[idx].cardName, profiles[idx].cardMonth, profiles[idx].cardYear, profiles[idx].cardCVV)
		fmt.Println()
	}
}
	// indices of the csv array (column# basically)
	firstName, lastName, email, addyLn1, addyLn2, city, state, zip, tele, cardNum, cardName, cardMonth, cardYear, cardCVV := 1,2,3,4,5,6,7,8,9,10,11,12,13,14


	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")

	fmt.Println("Enter profile id you would like to edit: ")
	scanner.Scan()
	editId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}
	fmt.Println("Enter profile property you would like to edit (can't edit id): ")
	scanner.Scan()
	editProp := scanner.Text()
	fmt.Println("Enter value to update profile property to: ")
	scanner.Scan()
	editPropVal := scanner.Text()

	switch editProp {
	case "id":
		log.Fatalln("can't edit id")
	case "firstName":
		record[editId-1][firstName] = editPropVal
	case "lastName":
		record[editId-1][lastName] = editPropVal
	case "email":
		record[editId-1][email] = editPropVal
	case "addyLn1":
		record[editId-1][addyLn1] = editPropVal
	case "addyLn2":
		record[editId-1][addyLn2] = editPropVal
	case "city":
		record[editId-1][city] = editPropVal
	case "state":
		record[editId-1][state] = editPropVal
	case "zip":
		record[editId-1][zip] = editPropVal
	case "tele":
		record[editId-1][tele] = editPropVal
	case "cardNum":
		record[editId-1][cardNum] = editPropVal
	case "cardName":
		record[editId-1][cardName] = editPropVal
	case "cardMonth":
		record[editId-1][cardMonth] = editPropVal
	case "cardYear":
		record[editId-1][cardYear] = editPropVal
	case "cardCVV":
		record[editId-1][cardCVV] = editPropVal
	default:
		log.Fatalln("invalid property")
	}

	file1, err := os.Create(profilesCSV)
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	writer := csv.NewWriter(file1)
	for _, rec := range record {
		err := writer.Write(rec)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
	fmt.Printf("Successfully updated profile %d: ", editId)
	fmt.Println()
}


func ProfileDelete() {
	file, err := os.Open(profilesCSV)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	
	record := csvRead(file)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Print out all profiles into the console? Y/N")
	scanner.Scan()
	if (scanner.Text() == "Y") {
	// SEPERATING LINE BETWEEN INPUTS AND PROFILE OUTPUTS
	fmt.Println("--------------------------------")
	var profiles []Profile
	for idx, item := range record {
		profile := Profile{id: item[0], firstName: item[1], lastName: item[2], email: item[3], addyLn1: item[4], addyLn2: item[5], city: item[6], state: item[7], zip: item[8], tele: item[9], cardNum: item[10], cardName: item[11], cardMonth: item[12], cardYear: item[13], cardCVV: item[14]}
		profiles = append(profiles, profile)

		fmt.Println()
		fmt.Printf("id: %s, firstName %s, lastName %s, email %s, addyLn1: %s, addyLn2: %s, city: %s, state: %s, zip: %s, tele: %s, cardNum: %s, cardName: %s, cardMonth: %s, cardYear: %s, cardCVV: %s",
		profiles[idx].id, profiles[idx].firstName, profiles[idx].lastName, profiles[idx].email, profiles[idx].addyLn1, profiles[idx].addyLn2, profiles[idx].city,
		profiles[idx].state, profiles[idx].zip, profiles[idx].tele, profiles[idx].cardNum, profiles[idx].cardName, profiles[idx].cardMonth, profiles[idx].cardYear, profiles[idx].cardCVV)
		fmt.Println()
	}
}

	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")

	fmt.Println("Enter [id] of profile you want to delete: ")
	scanner.Scan()
	deleteId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}
	if len(record) < deleteId {
		log.Fatalln("id does not exist, either no entries in profiles or out of bounds id")
	}

	record = append(record[:deleteId-1], record[deleteId:]...)

	for _, re := range record[deleteId-1:] {
		temp, err := strconv.Atoi(re[0])
		if err != nil {
			panic(err)
		}
		re[0] = strconv.Itoa(temp - 1)
	}


	file1, err := os.Create(profilesCSV)
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	writer := csv.NewWriter(file1)
	for _, re := range record {
		err := writer.Write(re)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()

	fmt.Println()
	fmt.Println("Successfully deleted profile id #", deleteId)
}