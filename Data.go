package main

//Buyer model from JSON
type Buyer struct {
	BuyerID   string `json:"buyer_id"`
	BuyerName string `json:"buyer_name"`
	Age       int    `json:"age"`
}

//Product model
type Product struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Price       int    `json:"price"`
}

//Transaction model
type Transaction struct {
	transactionID string
	BuyerID       string
	IP            string
	Device        string
	IDproduct     []string
}

//StrucingData infos
type StrucingData struct {
	TransactionID string    `json:"transaction_id"`
	Buyer         Buyer     `json:"buyer"`
	IP            string    `json:"ip"`
	Device        string    `json:"device"`
	Product       []Product `json:"products"`
}

//GetRestaInfo Get and struct all the retaurant info
func GetRestaInfo(date string) []StrucingData {
	url := "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com"
	endpoint := []string{"/buyers", "/products", "/transactions"} // "/products",

	channelBuyer := make(chan string)
	channelProduct := make(chan string)
	channelTransaction := make(chan string)

	go GetData(url+endpoint[0], date, channelBuyer)
	go GetData(url+endpoint[1], date, channelProduct)
	go GetData(url+endpoint[2], date, channelTransaction)

	dataProduct := <-channelProduct
	dataBuyer := <-channelBuyer
	dataTransaction := <-channelTransaction

	products := CSVToProducts(dataProduct)
	buyers := JSONToBuyers(dataBuyer)
	transactions := NSToTransactions(dataTransaction)

	finalInfo := UnifyData(transactions, buyers, products)
	return finalInfo
}
