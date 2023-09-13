package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}

func (app *App) Initialize(DbUser string, DbPassword string, DbName string) error {
	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	
	app.Router = mux.NewRouter().StrictSlash(true)
	app.HandleRoutes()
	return nil
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))

}

func (app *App) HandleRoutes(){
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}){
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string){
	error_message := map[string]string{"error": err}
	sendResponse(w, statusCode, error_message)
}


func (app *App) getProducts(w http.ResponseWriter, r *http.Request){
	products, err := getProducts(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id,err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	product := Product{Id: id}
	err = product.getProduct(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "Product not found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	sendResponse(w, http.StatusOK, product)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request){
	var product Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	if err := product.createProduct(app.DB); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusCreated, product)
} 

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	var product Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	product.Id = id
	if err := product.updateProduct(app.DB); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, product)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	product := Product{Id: id}
	if err := product.deleteProduct(app.DB); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"result": "success"})
}