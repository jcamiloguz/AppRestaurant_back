package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type buyers struct {
	IDBuyer   string `json:"id"`
	NameBuyer string `json:"name"`
	AgeBuyer  int    `json:"age"`
}

func main() {
	url := "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com"
	endpoints := []string{"/buyers"} // "/products",
	//  "/transactions"
	date := []string{"1577836800"}

	info := getAllData(url, endpoints, date)
	// info2 := getAllData(url, endpoints, []string{""})
	var buyers []buyers
	// fmt.Println(info)
	json.Unmarshal([]byte(info), &buyers)

	fmt.Println("Buyers : %+v", buyers)

}

func getAllData(url string, endpoints []string, date []string) string {
	data := "2"
	if date[0] == "" {
		date[0] = strconv.FormatInt(time.Now().Unix(), 10)
	}
	data = getDataBuyer(url + endpoints[0] + "?date=" + date[0])
	fmt.Println(url + endpoints[0] + "?date=" + date[0])
	return data
}

func getDataBuyer(url string) string {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}
