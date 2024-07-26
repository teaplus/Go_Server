package main

import (
	"fmt"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Price  float64
}

type Customer struct {
	ID    string
	Name  string
	Email string
}

type Order struct {
	OrderID    string
	Customer   Customer
	Books      []Book
	TotalPrice float64
}

func (o *Order) CalculateTotal() {
	var total float64
	for _, book := range o.Books {
		fmt.Println(book.Price, total)
		total += book.Price
	}
	o.TotalPrice = total

}

func (o *Order) AddBook(book Book) {
	o.Books = append(o.Books, book)
	o.CalculateTotal()
}

func NewBook(id int, title, author string, price float64) Book {
	return Book{ID: id, Title: title, Author: author, Price: price}
}

func NewCustomer(id, name, email string) Customer {
	return Customer{ID: id, Name: name, Email: email}
}

func NewOrder(orderID string, customer Customer) Order {
	return Order{OrderID: orderID, Customer: customer, Books: []Book{}}
}

func main() {
	//create Book
	book1 := NewBook(1, "The Go Programming Language", "Alan A. A. Donovan", 35.99)
	book2 := NewBook(2, "Learning Go", "Jon Bodner", 29.99)

	//create customer
	customer := NewCustomer("1", "John Doe", "john@example.com")

	//create Order
	order := NewOrder("1", customer)

	//add book to order
	order.AddBook(book1)
	order.AddBook(book2)

	fmt.Printf("Order ID: %s\n", order.OrderID)
	fmt.Printf("Customer: %s\n", order.Customer.Name)
	fmt.Println("Books in Order:")
	for _, book := range order.Books {
		fmt.Printf("- %s by %s: $%.2f\n", book.Title, book.Author, book.Price)
	}
	fmt.Printf("Total Price: $%.2f\n", order.TotalPrice)
}
