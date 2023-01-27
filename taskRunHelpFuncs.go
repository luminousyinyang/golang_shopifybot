package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func csvRead(file *os.File) ([][]string) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	taskRecord, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return taskRecord
}
// Gets the profile specified in the task.
func getTaskProf(taskRecords [][]string, taskId int) ([]string){
	profFile, err := os.Open("/Users/userid/path/to/csv/profiles.csv")
	if err != nil {
		panic(err)
	}
	
	profRecords := csvRead(profFile)
	
	taskProfId, err := strconv.Atoi(taskRecords[taskId-1][3])
	if err != nil {
		panic(err)
	}
	taskProfId--

	taskProf := profRecords[taskProfId]
	return taskProf
}

func findTaskVar(jsonURL string, taskSize string) (int) {
	resp, err := http.Get(jsonURL)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	
	type ShopVarResp struct {
		Inner struct {
			InnerArray []struct {
				Size1 string `json:"option1"`
				Size2 string `json:"option2"`
				Size3 string `json:"option3"`
				Var1 int `json:"id"`
			} `json:"variants"`
		} `json:"product"`
	}

	var response ShopVarResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	var taskVariant int
	for _, val := range response.Inner.InnerArray {
		if val.Size1 == taskSize || val.Size2 == taskSize || val.Size3 == taskSize {
			taskVariant = val.Var1
			break
		}
	}
	return taskVariant
}

func taskATC(client *http.Client, siteURL string, taskVariant int) {
	atcReqBuff := bytes.NewBuffer([]byte(`{
		"id": ` + strconv.Itoa(taskVariant) + `,
		"quantity": 1
	}`))
	req, err := http.NewRequest("POST", siteURL + "cart/add.json", atcReqBuff)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	

	if resp.StatusCode != http.StatusOK {
		panic("Failed to add item to cart")
	} else {
		println("Added item to Cart. HTTP Status code: ", resp.StatusCode)
	}
}

func getAuthTok(body []byte) (authTok string){
	regex := regexp.MustCompile(`"authenticity_token" value="([^"]*)"`)
	authToken := regex.FindAllString(string(body), -1)[0]
	regex = regexp.MustCompile(`[^"][A-Za-z0-9-_]+`)
	authToken1 := regex.FindAllString(authToken, -1)[2]

	return authToken1
}
//************************ 
func submitContactInfo(client *http.Client, getCheckURL, authToken string, profile Profile) (authTok string){
	data := url.Values{}
	data.Add("utf8", "✓")
	data.Add("_method", "patch")
	data.Add("authenticity_token", authToken)
	data.Add("previous_step", "contact_information")
	data.Add("step", "shipping_method")
	data.Add("checkout[email]", profile.email)
	data.Add("checkout[shipping_address][first_name]", profile.firstName)
	data.Add("checkout[shipping_address][last_name]", profile.lastName)
	data.Add("checkout[shipping_address][company]", "")
	data.Add("checkout[shipping_address][address1]", profile.addyLn1)
	data.Add("checkout[shipping_address][address2]", profile.addyLn2)
	data.Add("checkout[shipping_address][city]", profile.city)
	data.Add("checkout[shipping_address][country]", "United States")
	data.Add("checkout[shipping_address][province]", profile.state)
	data.Add("checkout[shipping_address][zip]", profile.zip)
	data.Add("checkout[shipping_address][phone]", profile.tele)
	data.Add("checkout[remember_me]", "0")
	data.Add("checkout[client_details][browser_width]", "979")
	data.Add("checkout[client_details][browser_height]", "631")
	data.Add("checkout[client_details][javascript_enabled]", "1")
	data.Add("button", "")
	
	dataStr := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", getCheckURL, dataStr)
		if err != nil {
			panic(err)
		}
		
		req.Header.Set("Referer", getCheckURL + "?step=contact_information")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.8")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		limitedReader := &io.LimitedReader{resp.Body, 100000}
		body, err := ioutil.ReadAll(limitedReader)
		if err != nil {
			panic(err)
		}
		
		authToken1 := getAuthTok(body)
		
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Failed to submit contact info. HTTP Status code: ", resp.StatusCode)
			panic(err)
		} else {
			fmt.Println("Submitted Contact Info. HTTP Status code: ", resp.StatusCode)
			return authToken1
		}
}
//************************** 
func getShipRate(client *http.Client, siteURL string, profile Profile) (shipRateReturn string) {
	shipMethodFind := siteURL + "cart/shipping_rates.json?shipping_address[zip]=" + profile.zip + "&shipping_address[country]=" + "United States" + "&shipping_address[province]=" + profile.state
	req, err := http.NewRequest("GET", shipMethodFind, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	type shipMethod struct {
		Rates []struct {
			Name string `json:"name"`
			Price string `json:"price"`
		} `json:"shipping_rates"`
	}

	var ship shipMethod
	err = json.Unmarshal(body, &ship)
	if err != nil {
		panic(err)
	}
	shipName := ship.Rates[0].Name
	shipName = strings.Replace(shipName, " ", "%20", -1)
	shipPrice := ship.Rates[0].Price
	shipRate := "shopify-" + shipName + "-" + shipPrice
	return shipRate
}

func submitShipInfo(client *http.Client, getCheckURL string, shipRate string, authToken string) (authTok, payGate string){
	dataShip := url.Values{}
	dataShip.Add("_method", "patch")
	dataShip.Add("authenticity_token", authToken)
	dataShip.Add("previous_step", "shipping_method")
	dataShip.Add("step", "payment_method")
	dataShip.Add("checkout[shipping_rate][id]", shipRate)
	dataShip.Add("checkout[client_details][browser_width]", "979")
	dataShip.Add("checkout[client_details][browser_height]", "631")
	dataShip.Add("checkout[client_details][javascript_enabled]", "1")
	dataShip.Add("checkout[client_details][color_depth]", "30")
	dataShip.Add("checkout[client_details][java_enabled]", "false")
	dataShip.Add("checkout[client_details][browser_tz]", "300")
	
	dataStrShip := strings.NewReader(dataShip.Encode())

	req, err := http.NewRequest("POST", getCheckURL, dataStrShip)
		if err != nil {
			panic(err)
		}
		req.PostForm = dataShip
		req.Header.Set("Referer", getCheckURL + "?step=shipping_method")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.8")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("Submitted Shipping Info. HTTP Status code: ", resp.StatusCode)

		return getPayGateway(client, resp.Request.URL.String())
}

func getPayGateway(client *http.Client, paymentURL string) (authTok, payGate string) {
		req, err := http.NewRequest("GET", paymentURL, nil)
		if err != nil {
			panic(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		regex := regexp.MustCompile(`class="input-radio".+`)
		pay := regex.FindAllString(string(body), -1)
		regex = regexp.MustCompile(`\d+`)
		pay_gateway := regex.FindAllString(pay[0], -1)[0]
		
		authToken1 := getAuthTok(body)

		return authToken1, pay_gateway
}
//*********************** 
func fetchPayTok(client *http.Client, profile Profile) (payTok string) {
	paytokenlink := "https://elb.deposit.shopifycs.com/sessions"

		paytokenBuff := bytes.NewBuffer([]byte(`{
			"credit_card": {
				"number": "` + profile.cardNum + `",
				"name": "`+ profile.cardName +`",
				"month": "` +profile.cardMonth +`",
				"year": "` + profile.cardYear + `",
				"verification_value": "` + profile.cardCVV +`"
			}
		}`))



		tokenReq, err := http.NewRequest("POST", paytokenlink, paytokenBuff)
		if err != nil {
			panic(err)
		}
		tokenReq.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(tokenReq)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		
		regex := regexp.MustCompile(`\w+-[A-Za-z0-9]+`)
		payToken := regex.FindAllString(string(body), -1)[0]

		return payToken
}
// ********************** 
func submitPayment(client *http.Client, authToken, payToken, pay_gateway, shipRate, getCheckURL string, profile Profile) {
	dataPay := url.Values{}
	dataPay.Add("utf8", "✓")
	dataPay.Add("_method", "patch")
	dataPay.Add("authenticity_token", authToken)
	dataPay.Add("previous_step", "payment_method")
	dataPay.Add("step", "")
	dataPay.Add("s", payToken)
	dataPay.Add("checkout[payment_gateway]", pay_gateway)
	dataPay.Add("checkout[credit_card][vault]", "false")
	dataPay.Add("checkout[billing_address][first_name]", profile.firstName)
	dataPay.Add("checkout[billing_address][last_name]", profile.lastName)
	dataPay.Add("checkout[billing_address][address1]", profile.addyLn1)
	dataPay.Add("checkout[billing_address][address2]", profile.addyLn2)
	dataPay.Add("checkout[billing_address][city]", profile.city)
	dataPay.Add("checkout[billing_address][country]", "United States")
	dataPay.Add("checkout[billing_address][province]", profile.state)
	dataPay.Add("checkout[billing_address][zip]", profile.zip)
	dataPay.Add("checkout[billing_address][phone]", profile.tele)
	dataPay.Add("checkout[shipping_rate][id]", shipRate)
	dataPay.Add("complete", "1")
	dataPay.Add("checkout[client_details][browser_width]", "979")
	dataPay.Add("checkout[client_details][browser_height]", "631")
	dataPay.Add("checkout[client_details][javascript_enabled]", "1")
	dataPay.Add("g-recaptcha-response", "")
	dataPay.Add("button", "")
	
	dataStrPay := strings.NewReader(dataPay.Encode())

	req, err := http.NewRequest("POST", getCheckURL, dataStrPay)
		if err != nil {
			panic(err)
		}
		req.Header.Set("Referer", getCheckURL + "?step=payment_method")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.8")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Failed to complete checkout. Status code: ", resp.StatusCode)
		} else {
			fmt.Println("Submitted payment info, check card!")
		}
}

func readAllTasks(taskRecord [][]string) {
	for _, task := range taskRecord {
		t := Task{id: task[0], url: task[1], size: task[2], profile: task[3]}
		
		fmt.Println()
		fmt.Printf("{id: %s, url: %s, size: %s, profile: %s}", t.id, t.url, t.size, t.profile)
		fmt.Println()
	}

	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")
	fmt.Println("--------------------------------")
}