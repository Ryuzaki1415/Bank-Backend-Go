package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//we are creating Types of the USER.

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstname"`
	LastName          string    `json:"lastname"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"encryptedpassword"`
	Balance           int64     `json:"balance"` // `json:"id"` says that whenever ID is mentioned, the json  encoding should be id
	CreatedAt         time.Time `json:"time"`
}

//initializing new Account

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		//ID:        rand.Intn(1000),
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(10000)),
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
	}, nil

}

// creating a type for create_account Request ie when the user requests to create this account, they have to send data in this foramt.

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LasttName string `json:"lastName"`
	Password  string `json:"password"`
}

// creating a type for the money transfer request that the user sends.
type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

type LoginRequest struct {
	Number   int    `json:"number"`
	Password string `json:"password"`
}


type LoginResponse struct{
	Number int64 `json:"number"`
	Token string  `json:"token"`
}

func (a *Account) ValidatePassword(pw string)bool{
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword),[]byte(pw)) == nil 
}





//5456