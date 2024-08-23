package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAcccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error) // returns an Slice(array) of existing accounts
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank  sslmode=disable" //remember to disable
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil

}

// we are gonna actually implement the interface we defined earlier.
// ie create acnt , delete , update and get acnt by ID

func (s *PostgresStore) CreateAccount(acc *Account) error { // here we are getting the actual user information to create account(FIRSTNAME AND LASTNAME)
	query := `INSERT INTO ACCOUNT 
	(first_name,last_name,number,balance,created_at)
	VALUES($1,$2,$3,$4,$5)
	
	`
	response, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", response)
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_,err:=s.db.Query("delete from account where id = $1",id)
	return err
}

func (s *PostgresStore) UpdateAcccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {

		rows,err:=s.db.Query("select * from account where id = $1 ",id)
		
		if err!=nil{
			return nil,err
		}
		for rows.Next(){
			return scanIntoAccount(rows)
		}
		return nil,fmt.Errorf("this account id does not exist : %v ",id)
		}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	//fetching all rows from the DB to Display When getting a GET REQUEST
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}  // an array of accounts.

	//iterating Through the ROWS
	for rows.Next() {
		account,err :=scanIntoAccount(rows)  // returns the account object containing all values
		if err!=nil{
			return nil,err
		}
		accounts = append(accounts, account)  // append this account to the array of accounts.
	}
	return accounts, nil

}

//now to init our DB with tables

func (s *PostgresStore) init() error {
	return s.CreateAccountTable()

}

// the query to make an accounts table if it doesnt exist
func (s *PostgresStore) CreateAccountTable() error {
	query := ` CREATE TABLE IF NOT EXISTS ACCOUNT(
	id serial PRIMARY KEY,
	first_name varchar(50),
	last_name varchar(50),
	number serial,
	balance serial,
	created_at timestamp
	)
	`
	_, err := s.db.Exec(query)
	return err

}


func scanIntoAccount( rows *sql.Rows)(*Account,error){
// function which scans the rows of DB and puts into account format
	account := new(Account)
		if err := rows.Scan(
			&account.ID, // basically we are taking all the columnnames and
			//creating an accnt object for each individual account and then appending it to the accounts array
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		); err != nil {
			return nil, err
		}
	return account,nil
}