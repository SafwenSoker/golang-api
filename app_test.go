package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"bytes"
	"encoding/json"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("Error occured while iniatializing the database")

	}
	createTable()
	m.Run()
}

func createTable(){
	query := `CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		quantity INT NOT NULL,
		price FLOAT NOT NULL
	)`
	_, err := a.DB.Exec(query)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error occured while creating the table")
	}	
}

func clearTable(){
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER TABLE products AUTO_INCREMENT = 1")
}

func addProduct(name string, quantity int, price float64){
	a.DB.Exec("INSERT INTO products(name, quantity, price) VALUES(?, ?, ?)", name, quantity, price)
} 

func TestGetProducts(t *testing.T){
	clearTable()
	addProduct("Product 1", 10, 10.00)
	req, _ := http.NewRequest("GET", "/products", nil)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func checkResponseCode(t *testing.T, expected, actual int){
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
} 

func sendRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func TestCreateProduct(t *testing.T){
	clearTable()
	payload := []byte(`{"name":"Product 1", "quantity":10, "price":10.00}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	response := sendRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "Product 1" {
		t.Errorf("Expected product name to be Product 1. Got %v", m["name"])
	}
	if m["quantity"] != 10.00 {
		t.Errorf("Expected product quantity to be 10. Got %v", m["quantity"])
	}
	if m["price"] != 10.00 {
		t.Errorf("Expected product price to be 10. Got %v", m["price"])
	}
	 
}

func TestDeleteProduct(t *testing.T){
	clearTable()
	addProduct("Product 1", 10, 10.00)
	req, _ := http.NewRequest("GET", "/products", nil)
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
} 

func TestUpdateProduct(t *testing.T){
	clearTable()
	addProduct("Product 1", 10, 10.00)
	payload := []byte(`{"name":"Product 1", "quantity":1, "price":10.00}`)
	req, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	response := sendRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	
	if m["name"] != "Product 1" {
		t.Errorf("Expected product name to be Product 1. Got %v", m["name"])
	}
	if m["quantity"] != 1.00 {
		t.Errorf("Expected product quantity to be 1. Got %v", m["quantity"])
	}
	if m["price"] != 10.00 {
		t.Errorf("Expected product price to be 10. Got %v", m["price"])
	}
}