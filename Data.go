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

//TransactionsRsp response Transaction struct
type TransactionsRsp struct {
	Ip          string   `json:"ip"`
	Device      string   `json:"device"`
	Products_id []string `json:"products_id"`
}

//BuyerRsp response Buyer struct
type BuyerRsp struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//ProductRsp response Product struct
type ProductRsp struct {
	Product_id string `json:"product_id"`
	Name       string `json:"product_name"`
	Price      int    `json:"price"`
}

//Details Model
type Details struct {
	Buyer       []BuyerRsp        `json:"buyer`
	Transaction []TransactionsRsp `json:"transaction"`
	Buyers      []BuyerRsp        `json:"buyers"`
}

//History Model
type History struct {
	Products []ProductRsp `json:"products"`
}

//DetailResponse Main model
type DetailResponse struct {
	Buyers   Details `json:"Details"`
	Products History `json:"History"`
}

//StrucingData likely the best way to sruct the mutation, Redundancy Degraph problems
type StrucingData struct {
	Uid           string    `json:"uid"`
	TransactionID string    `json:"transaction_id"`
	Buyer         Buyer     `json:"buyer"`
	IP            string    `json:"ip"`
	Device        string    `json:"device"`
	Product       []Product `json:"products"`
	DType         []string  `json:"dgraph.type"`
}
