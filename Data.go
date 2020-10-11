package main

//Buyer model from JSON
type Buyer struct {
	Uid       string   `json:"uid"`
	BuyerID   string   `json:"id"`
	BuyerName string   `json:"name"`
	Age       int      `json:"age"`
	DType     []string `json:"dgraph.type"`
}

//Product model
type Product struct {
	Uid         string   `json:"uid"`
	ProductID   string   `json:"product_id"`
	ProductName string   `json:"product_name"`
	Price       int      `json:"price"`
	DType       []string `json:"dgraph.type"`
}

//Transaction model
type Transaction struct {
	Uid           string   `json:"uid"`
	TransactionID string   `json:"transaction_id"`
	BuyerID       string   `json:"buyer_id"`
	IP            string   `json:"ip"`
	Device        string   `json:"device"`
	IDproduct     []string `json:"products_id"`
	DType         []string `json:"dgraph.type"`
}

//StrucingData infos
type StrucingData struct {
	Uid           string    `json:"uid"`
	TransactionID string    `json:"transaction_id"`
	Buyer         Buyer     `json:"buyer"`
	IP            string    `json:"ip"`
	Device        string    `json:"device"`
	Product       []Product `json:"products"`
	DType         []string  `json:"dgraph.type"`
}
type TransactionsRsp struct {
	Transaction_id string
	Ip             string
	Device         string
	Products_id    []string
}
type RespTransaction struct {
	Transactions []byte
}

//GetRestaInfo Get and struct all the retaurant info
func GetRestaInfo(date string) ([]Buyer, []Product, []Transaction) {
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

	return buyers, products, transactions
}
