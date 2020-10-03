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

//FindBuyer Find a buyer from  id
func FindBuyer(a []Buyer, x string) int {
	for i, n := range a {
		if x == n.BuyerID {
			return i
		}
	}
	return len(a) - 1
}

//FindProduct find a product from products
func FindProduct(a []Product, x string) int {
	for i, n := range a {
		if x == n.ProductID {
			return i
		}
	}
	return len(a) - 1
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
	for i, b := range buyers {
		buyers[i].Uid = "_:" + b.BuyerID
		buyers[i].DType = []string{"Buyers"}
	}
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
			product := Product{"_:" + info[1], info[1], info[2], price, []string{"Product"}}
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

//UnifyData unify the transactions, buyers and products in one struct
func UnifyData(transactions []Transaction, buyers []Buyer, products []Product) []StrucingData {
	finalInfo := []StrucingData{}
	for _, t := range transactions {
		transacProcs := StrucingData{}
		id := strings.Replace(t.transactionID, "#", "", -1)
		transacProcs.Uid = "_:" + id
		transacProcs.TransactionID = t.transactionID
		transacProcs.IP = t.IP
		b := FindBuyer(buyers, t.BuyerID)
		transacProcs.Buyer = buyers[b]
		for _, pro := range t.IDproduct {
			d := FindProduct(products, pro)
			transacProcs.Product = append(transacProcs.Product, products[d])

		}
		transacProcs.Device = t.Device
		transacProcs.DType = []string{"Transaction"}
		finalInfo = append(finalInfo, transacProcs)
	}
	return finalInfo
}
