# Schema
****
> Graphql
``` 
type buyer {
	buyerID:String! @id
	buyerName:String!
	age:Int!
}
****
type product {
	productID:String! @id
	productName:String!
	price:Int!
}
type transaction {
	transactionID:String! @id 
	buyer:buyer! @search (by:buyerID)
	ip:Int!
	device:String!
	product:[product!]!
}
```
****
> Graphql+(Dgraph)
``` 
type Buyer {
    id
    name
    age
}
type Product {
    id
    name
    price
}
type Transaction {
    id
    buyer
    ip
    device
  	products
   
}

# Define Directives and index

name: string @index(term) @lang .
age: int @index(int) .
price: int @index(int) .
products: [uid] @count .
buyer: [uid] .
device: string .
id: string  @index(term) .

```