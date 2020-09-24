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

type buyer struct {
	IDBuyer   string `json:"id"`
	NameBuyer string `json:"name"`
	AgeBuyer  int    `json:"age"`
}
type product struct {
	IDProduct    string
	NameProduct  string
	PriceProduct string
}

type transaction struct {
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

	go getData(url+endpoint[0], date, channelBuyer)
	go getData(url+endpoint[1], date, channelProduct)
	go getData(url+endpoint[2], date, channelTransaction)

	// dataBuyer := <-channelBuyer
	// dataProduct := <-channelProduct
	dataTransactions := <-channelTransaction

	// buyers := JsonToBuyers(dataBuyer)
	// products := CsvToProducts(dataProduct)
	transactions := NfToTransactions(dataTransactions)
	fmt.Println(transactions)
	timeUsed := time.Since(start)
	fmt.Printf("Tiempo de ejecucion %s\n", timeUsed)

}
func getData(url string, date string, channel chan string) {
	url = url + "?date=" + date
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	channel <- string(body)
}
func JsonToBuyers(data string) []buyer {
	var buyers []buyer
	json.Unmarshal([]byte(data), &buyers)
	return buyers
}
func CsvToProducts(data string) []product {
	all := strings.Split(strings.Replace(data, "\r\n", "\n", -1), "\n")

	regex := regexp.MustCompile(`(?m)^(?P<id>\w+)\'(?P<name>[ \'\w\+\&\-\"\%\&\.À-ÿ]+)\'(?P<price>\d+)$`)

	var products []product
	for _, line := range all {
		info := regex.FindStringSubmatch(line)
		if len(info) > 1 {
			product := product{info[1], info[2], info[3]}
			products = append(products, product)
		}
	}
	return products
}
func NfToTransactions(data string) []transaction {
	cleaner := regexp.MustCompile(`[^a-zA-Z0-9\.\(\)\,\#]+`)
	dataClean := cleaner.ReplaceAllString(data, "")

	filter := regexp.MustCompile(`(?m)[\s]?(?P<id>#\w{12,12})\s?(?P<buyerid>\w{8,8})\s?(?P<ip>\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3})\s?(?P<device>\w{3,10})\s?\((?P<transactions>[\w{8,8},?]{1,})\)\s?`)
	info := filter.FindAllStringSubmatch(dataClean, -1)

	var transactions []transaction
	for _, params := range info {
		products := strings.Split(params[5], ",")
		transaction := transaction{params[1], params[2], params[3], params[4], products}
		transactions = append(transactions, transaction)
	}
	return transactions
}
