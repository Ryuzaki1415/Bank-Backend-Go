package main

import (
	"math/rand"
	"time"
)

//we are creating Types of the USER.

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Number    int64  `json:"number"`

	Balance   int64     `json:"balance"` // `json:"id"` says that whenever ID is mentioned, the json  encoding should be id
	CreatedAt time.Time `json:"time`
}

//initializing new Account

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		//ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000)),
		CreatedAt: time.Now().UTC(),
	}

}

// creating a type for create_account Request ie when the user requests to create this account, they have to send data in this foramt.

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LasttName string `json:"lastName"`
}
