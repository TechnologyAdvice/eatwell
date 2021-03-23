package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile := "./go-quickstart.json"
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

//structs for data
type Response struct {
	Order []Order `json:"orders"`
}

type Order struct {
	Address ShipAddress `json:"shipping_address"`
	ID      int         `json:"id"`
}

type ShipAddress struct {
	Address1  string  `json:"address1"`
	Address2  string  `json:"address2"`
	City      string  `json:"city"`
	State     string  `json:"province"`
	Country   string  `json:"country"`
	Zip       string  `json:"zip"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Order_Num struct {
	Count int `json:"count"`
}

func main() {

	//call google stuff to authenticate with their api
	ctx := context.Background()
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	sheetsService, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	//get id for the spreadsheet to print to
	spreadsheetId := "1YkeduplR7Dnqwkv5Jp5YcodxWHen0X5Sjm7G6sfBU0M"
	//rangeData := "sheet1!A1:M400"
	//valueInputOption := "" // TODO: Update placeholder value.

	//gets order count to use as limit
	count, err := http.Get("https://d3e1e0e9f9d7003f026155ed0cf3f35d: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@eatwellnash.myshopify.com/admin/api/2021-01/orders/count.json")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	order_count, err := ioutil.ReadAll(count.Body)
	if err != nil {
		log.Fatal(err)
	}
	var orderCount Order_Num
	json.Unmarshal(order_count, &orderCount)
	fmt.Println(orderCount.Count) // prints the count

	//handling stuff for shopify and calling the API
	var req *http.Response
	var responseData []byte
	var er error
	req, er = http.Get("https://d3e1e0e9f9d7003f026155ed0cf3f35d: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@eatwellnash.myshopify.com/admin/api/2021-01/orders.json?limit=250&fields=id,shipping_address")
	if er != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, er = ioutil.ReadAll(req.Body)
	if er != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))

	var responseObject Response

	json.Unmarshal(responseData, &responseObject)
	fmt.Println(len(responseObject.Order))

	fmt.Print("\n")
	fmt.Println("Address 1, Address 2, City, State, Country, Zip, Latitude, Longitude")

	//Loops through the struct and assigns it to the variable
	var n int
	n = 0
	var lastid int
	var url string
	for i := 0; i < orderCount.Count; i++ {

		if i == 249 || i == 499 || i == 749 {

			y := len(responseObject.Order)
			fmt.Println(y)
			lastid = responseObject.Order[y-2].ID
			url = fmt.Sprintf("https://d3e1e0e9f9d7003f026155ed0cf3f35d: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@eatwellnash.myshopify.com/admin/api/2021-01/orders.json?limit=250&since_id=%d&fields=id,shipping_address", lastid)
			req, er = http.Get(url)
			if er != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}
			responseData, err = ioutil.ReadAll(req.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(responseData))
			json.Unmarshal(responseData, &responseObject)
			fmt.Println(len(responseObject.Order))
			n = 0
		}

		sheetnum := fmt.Sprintf("sheet1!A%d:M40000", i+2) //turn into string??
		rangeData := sheetnum
		fmt.Println(responseObject.Order[n].Address.Address1, ",", responseObject.Order[n].Address.Address2, ",", responseObject.Order[n].Address.City, ",", responseObject.Order[n].Address.State, ",", responseObject.Order[n].Address.Country, ",", responseObject.Order[n].Address.Zip, ",", responseObject.Order[n].Address.Latitude, ",", responseObject.Order[n].Address.Longitude)

		values := [][]interface{}{{responseObject.Order[n].Address.Address1, responseObject.Order[n].Address.Address2, responseObject.Order[n].Address.City, responseObject.Order[n].Address.State, responseObject.Order[n].Address.Country, responseObject.Order[n].Address.Zip, responseObject.Order[n].Address.Latitude, responseObject.Order[n].Address.Longitude}}
		fmt.Println(rangeData)
		fmt.Println(values)
		rb := &sheets.BatchUpdateValuesRequest{
			ValueInputOption: "USER_ENTERED",
		}
		rb.Data = append(rb.Data, &sheets.ValueRange{
			Range:  rangeData,
			Values: values,
		})
		_, err = sheetsService.Spreadsheets.Values.BatchUpdate(spreadsheetId, rb).Context(ctx).Do()
		if err != nil {
			log.Fatal(err)

		}

		//lastid = responseObject.Order[n].ID
		fmt.Println(len(responseObject.Order))
		//fmt.Println(responseObject.Order[n].ID)
		fmt.Println(n)
		time.Sleep(1 * time.Second)
		n = n + 1

	}
}
