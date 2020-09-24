package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//Buyer model from JSON
type Buyer struct {
	IDBuyer   string `json:"id"`
	NameBuyer string `json:"name"`
	AgeBuyer  int    `json:"age"`
}

//Product model
type Product struct {
	IDProduct    string
	NameProduct  string
	PriceProduct string
}

//Transaction model
type Transaction struct {
	IDTransaction string
	IDBuyer       string
	IP            string
	Device        string
	IDproduct     []string
}

func main() {
	start := time.Now()
	url := "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com"
	endpoint := []string{"/buyers", "/products", "/transactions"} // "/products",
	date := "1577836800"

	channelBuyer := make(chan string)
	channelProduct := make(chan string)
	channelTransaction := make(chan string)

	go GetData(url+endpoint[0], date, channelBuyer)
	go GetData(url+endpoint[1], date, channelProduct)
	go GetData(url+endpoint[2], date, channelTransaction)

	// dataBuyer := <-channelBuyer
	// dataProduct := <-channelProduct
	dataTransactions := <-channelTransaction

	// buyers := CSVnTProductss(dataBucsv
	// products := CSVToProductcts(dataProduct)
	transactions := NSToTransactions(dataTransactions)
	fmt.Println(transactions)
	timeUsed := time.Since(start)
	fmt.Printf("Tiempo de ejecucion %s\n", timeUsed)

}

//GetData from a url passinga a Date param
func GetData(url string, date string, channel chan string) {
	url = url + "?date=" + date
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	// buyers := CSVnTProductss(dataBucsv
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	channel <- string(body)
}

//JSONToBuyers take a JSON string a pass a in []Buyer model
func JSONToBuyers(data string) []Buyer {
	var buyers []Buyer
	json.Unmarshal([]byte(data), &buyers)
	return buyers
}

//CSVToProducts take a csv string a pass a in []Product model
func CSVToProducts(data string) []Product {
	all := strings.Split(strings.Replace(data, "\r\n", "\n", -1), "\n")

	regex := regexp.MustCompile(`(?m)^(?P<id>\w+)\'(?P<name>[ \'\w\+\&\-\"\%\&\.À-ÿ]+)\'(?P<price>\d+)$`)

	var products []Product
	for _, line := range all {
		info := regex.FindStringSubmatch(line)
		if len(info) > 1 {
			product := Product{info[1], info[2], info[3]}
			products = append(products, product)
		}
	}
	return products
}

//NSToTransactions take a non-standard string a pass a in []Transaction model
func NSToTransactions(data string) []Transaction {
	cleaner := regexp.MustCompile(`[^a-zA-Z0-9\.\(\)\,\#]+`)
	dataClean := cleaner.ReplaceAllString(data, "")

	filter := regexp.MustCompile(`(?m)[\s]?(?P<id>#\w{12,12})\s?(?P<buyerid>\w{8,8})\s?(?P<ip>\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3})\s?(?P<device>\w{3,10})\s?\((?P<transactions>[\w{8,8},?]{1,})\)\s?`)
	info := filter.FindAllStringSubmatch(dataClean, -1)

	var transactions []Transaction
	for _, params := range info {
		products := strings.Split(params[5], ",")
		transaction := Transaction{params[1], params[2], params[3], params[4], products}
		transactions = append(transactions, transaction)
	}
	return transactions
}
