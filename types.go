package main

import "math/rand"

//we are creating Types of the USER.

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Number    int64  `json:"number"`
	balance   int64  `json:"balance"` // `json:"id"` says that whenever ID is mentioned, the json  encoding should be id
}

//initializing new Account

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000)),
	}

}
