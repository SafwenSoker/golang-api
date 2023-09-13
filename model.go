package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type Product struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
}

var Products []Product

func getProducts(db *sql.DB) ([]Product, error) {
	query := "SELECT id, name, quantity, price from products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.Id, &product.Name, &product.Quantity, &product.Price)
		if err != nil {
			return nil, err
		}
		Products = append(Products, product)
	}
	return Products, nil
}

func (p *Product) getProduct(db *sql.DB) error 	  {
	query := "SELECT name, quantity, price from products where id=?"
	err := db.QueryRow(query, &p.Id).Scan(&p.Name, &p.Quantity, &p.Price)
	fmt.Println(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) createProduct(db *sql.DB) error {
	query := "INSERT INTO products(name, quantity, price) VALUES(?, ?, ?)"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}	
	p.Id = int(id)
	return nil
}

func (p *Product) updateProduct(db *sql.DB) error {
	query := "UPDATE products SET name=?, quantity=?, price=? WHERE id=?"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price, p.Id)
	rowsAffected, err := result.RowsAffected()	
	if rowsAffected == 0 {
		return errors.New("Product not found")
	}
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) deleteProduct(db *sql.DB) error {
	query := "DELETE FROM products WHERE id=?"
	result, err := db.Exec(query, p.Id)
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("Product not found")
	}
	if err != nil {
		return err
	}
	return nil
}