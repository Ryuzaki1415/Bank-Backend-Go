package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// here we create api server and handlers

//first thing to do is to create a type of API server

type APIServer struct {
	listenAddr string
	//db to be added later
	store Storage
}

// now to actually create an instance of this server

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{ // here we are instantiating an API server with the listening port and returning the address of the location of the created APISERVER
		listenAddr: listenAddr,
		store:      store,
	}
}

// creation of HANDLERS!!!!!
// now to create the central handler which is gonna facilitate the creation, deletion,transfer money  and reading APIs
//prefixing with "handle" is preffered

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	//the purpose of handle account is that we cannot directly specify which methods are acceptable . We define them here with if statements.
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("METHOD NOT ALLOWED!!!! %s", r.Method)

}

// func to get all available accounts  GET/account

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

// function to get accnt by ID
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error { // main account handling
	//account := NewAccount("DHEERAJ", "UNNI")
	id := mux.Vars(r)["id"] // getting id from the request argument
	fmt.Println("the selected ID is ", id)
	return WriteJSON(w, http.StatusOK, &Account{}) // return an empty account

}

// NOW WE ARE GOING TO DEFINE THE LOGIC TO HANDLE ACCOUNT CREATION
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error { // FOR POST request

	//thr Request shld be of the type we defined.

	createAccountReq := new(CreateAccountRequest)
	//decoding the user sent JSON and matching it with our Type.
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	//creation of a new account with the data that we got from the user.
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LasttName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error { // FOR DELETE REQUEST
	return nil

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil

}

// Now create a function to start our server!!!   // install gorilla/mux for making routers!!!
//remember to convert func to HTTP handler

//handleFunc  accepts a handler

func (s *APIServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount))             //see that we have wrapped the function handleAccount and converted it to a HTTP handler.
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleGetAccountByID)) // this is so that we can retrieve accounts by ID.
	log.Println("JSON API SERVER RUNNING ON PORT", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router) //Running the server.
}

// we are co=nverting the functions to handlers using this code

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil { //if there is any problem in converting to handlers
			WriteJSON(w, http.StatusBadRequest, APIError{
				Error: err.Error(),
			})
		}
	}
}

// now we need to write JSON.. creating a function for that

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json") // we need to add headers before we Writeheader
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// we  use this to handle the errors
type APIError struct {
	Error string
}
