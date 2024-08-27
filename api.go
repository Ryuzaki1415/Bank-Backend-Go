package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/golang-jwt/jwt/v5"
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
func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error { // fetch an account using ID
	// the id that we get is gonna be a string so we convert it to int.
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, account) // return the selected  account
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)

	}
	return fmt.Errorf("METHOD IS NOT ALLOWED : %s", r.Method)
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
	//creation of a JWT token when we are making a new account.
	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}
	fmt.Println("JWT TOKEN :", tokenString)
	return WriteJSON(w, http.StatusOK, account)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error { // FOR DELETE REQUEST

	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}
	defer r.Body.Close() //defer activates when fn execution ends.
	return WriteJSON(w, http.StatusOK, transferRequest)

}

// Now create a function to start our server!!!   // install gorilla/mux for making routers!!!
//remember to convert func to HTTP handler

//handleFunc  accepts a handler

func (s *APIServer) Run() {

	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount)) //see that we have wrapped the function handleAccount and converted it to a HTTP handler.
	//we have to protect /id
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandler(s.handleGetAccountByID), s.store)) // this is so that we can retrieve accounts by ID.
	router.HandleFunc("/transfer", makeHTTPHandler(s.handleTransfer))
	log.Println("JSON API SERVER RUNNING ON PORT", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router) //Running the server.
}

// we are converting the functions to handlers using this code

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

// function to fetch ID from request ( since its being used in multiple places)

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"] // getting id from the request argument
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid ID given %s !!! ", idStr)
	}
	return id, nil
}

//We are going to implement JWT Authentication.

// this is gonna get executed all the time when we call the getaccountbyID
func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("----------------CALLING MIDDLEWARE-------------------")
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetAccountByID(userID)
		if err != nil {
			permissionDenied(w)
			return

		}
		//34:49
		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "INVALID TOKEN"})
		}
		handlerFunc(w, r)

	}
}

// returns a pointer to JWT token
// the secret is a long string saved in environ.“
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Secret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

}

// function to create a JWT token when account creation is done.

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))

}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjozNDcwLCJleHBpcmVzQXQiOjE1MDAwfQ.wqJ10qI8JOHKv-uqAyydQpJYekVZfw0P7SpeCGMZkM8

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, APIError{Error: "permission denied"})
	return
}
