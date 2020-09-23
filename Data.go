package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type buyer struct {
	IDBuyer   string `json:"id"`
	NameBuyer string `json:"name"`
	AgeBuyer  int    `json:"age"`
}
type product struct {
	IDBuyer   string `json:"id"`
	NameBuyer string `json:"name"`
	AgeBuyer  int    `json:"age"`
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

	dataTransactions := <-channelTransaction
	dataProduct := <-channelProduct
	dataBuyer := <-channelBuyer

	fmt.Println(dataBuyer)
	fmt.Println(dataProduct)
	fmt.Println(dataTransactions)
	timeUsed := time.Since(start)
	fmt.Printf("Tiempo de ejecucion %s\n", timeUsed)

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
