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



{  "set":[
  {
	"uid": "_:wqe123",
  "id":"wqe123",
  "dgraph.type":"buyer",
  "name":"pedro",
  "age":"32"
},
  {
	"uid": "_:213gj",
  "id":"213gj",
  "dgraph.type":"buyer",
  "name":"juan",
  "age":"18"
}
]
}



{  "set":[
  {
	"uid": "_:231dsd",
  "id":"231dsd",
  "dgraph.type":"product",
  "name":"pizza",
  "price":"32133"
},
  {
	"uid": "_:dsvger3",
  "id":"dsvger3",
  "dgraph.type":"product",
  "name":"hotdog",
  "price":"12318"
}
]
}

{  "set":[
  {
	"uid": "_:dsaad",
  "id":"dsaad",
  "dgraph.type":"trasaction",
  "ip":"2313213",
  "device":"mac",
  "buyerid":"wqe123",
  "products":["231dsd","dsvger3"]
}
]
}


{  
  var(func: eq(id, "wqe123")) {
    name
  age
    BUY as id
  
  }
trans(func: eq(buyerid, val(BUY)) {
  ip
  device
}
}