package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//Buyer model from JSON
type Buyer struct {
	BuyerID   string `json:"id"`
	BuyerName string `json:"name"`
	Age       int    `json:"age"`
}

//Product model
type Product struct {
	ProductID   string
	ProductName string
	Price       int
}

//Transaction model
type Transaction struct {
	transactionID string
	BuyerID       string
	IP            string
	Device        string
	IDproduct     []string
}

func FindBuyer(a []Buyer, x string) int {
	for i, n := range a {
		if x == n.BuyerID {
			return i
		}
	}
	return len(a) - 1
}
func FindProduct(a []Product, x string) int {
	for i, n := range a {
		if x == n.ProductID {
			return i
		}
	}
	return len(a) - 1
}

// func main() {
// 	start := time.Now()
// 	url := "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com"
// 	endpoint := []string{"/buyers", "/products", "/transactions"} // "/products",
// 	date := "1577836800"

// 	channelBuyer := make(chan string)
// 	channelProduct := make(chan string)
// 	channelTransaction := make(chan string)

// 	go GetData(url+endpoint[0], date, channelBuyer)
// 	go GetData(url+endpoint[1], date, channelProduct)
// 	go GetData(url+endpoint[2], date, channelTransaction)

// 	dataProduct := <-channelProduct

// 	products := CSVToProducts(dataProduct)
// 	fmt.Println(products)
// 	timeUsed := time.Since(start)
// 	fmt.Printf("Tiempo de ejecucion %s\n", timeUsed)

// }

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

	regex := regexp.MustCompile(`(?m)^(?P<id>\w+)\'\"?(?P<name>[ \'\w\+\&\-\%\&\.À-ÿ]+)\"?\'(?P<price>\d+)$`)

	var products []Product
	for _, line := range all {
		info := regex.FindStringSubmatch(line)
		if len(info) > 1 {
			price, err := strconv.Atoi(info[3])
			if err != nil {
				// handle error
				fmt.Println(err)
				os.Exit(2)
			}
			product := Product{info[1], info[2], price}
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
