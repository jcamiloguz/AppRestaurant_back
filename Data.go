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

type transactions struct {
	IDBuyer   string `json:"id"`
	NameBuyer string `json:"name"`
	AgeBuyer  int    `json:"age"`
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

	// dataTransactions := <-channelTransaction
	dataProduct := <-channelProduct
	// dataBuyer := <-channelBuyer

	// buyers := JsonToBuyers(dataBuyer)
	produsts := CsvToProducts(dataProduct)
	// fmt.Println(buyers)
	// fmt.Println(products)
	// fmt.Println(dataProduct)
	// fmt.Println(dataTransactions)
	timeUsed := time.Since(start)
	fmt.Printf("Tiempo de ejecucion %s\n", timeUsed)

}
func JsonToBuyers(data string) []buyer {
	var buyers []buyer
	json.Unmarshal([]byte(data), &buyers)
	return buyers
}
func CsvToProducts(data string) []product {
	all := strings.Split(strings.Replace(data, "\r\n", "\n", -1), "\n")
	a := regexp.MustCompile(`(?m)^(?P<id>\w+)\'(?P<name>[ \'\w\+\&\-\"\%\&\.À-ÿ]+)\'(?P<price>\d+)$`)
	var products []product
	// fmt.Println(data)
	for _, line := range all {
		info := a.FindStringSubmatch(line)
		if len(info) > 1 {
			product := product{info[1], info[2], info[3]}
			products = append(products, product)
		}
	}
	return products
}

//  "/transactions"
// info := getAllData(url, endpoints, date)
// info2 := getAllData(url, endpoints, []string{""})
// fmt.Println(info)
// var buyers []buyer
// json.Unmarshal([]byte(info), &buyers)

// fmt.Println("Buyers : %+v", buyers[1])
// func getAllData(url string, endpoints []string, date []string) string {
// 	data := "2"
// 	if date[0] == "" {
// 		date[0] = strconv.FormatInt(time.Now().Unix(), 10)
// 	}
// 	data = getDataBuyer(url + endpoints[0] + "?date=" + date[0])
// 	fmt.Println(url + endpoints[0] + "?date=" + date[0])
// 	return data
// }

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
