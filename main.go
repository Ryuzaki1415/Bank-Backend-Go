package main

import (
	"flag"
	"fmt"
	"log"
)
//seed =  6064
func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}
	fmt.Println("NEW ACCOUNT =>>>",acc.Number)
	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "DJDJDJDJ", "UNNNNIIII", "MAAAANUNITDDDD")
}
func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.init(); err != nil {
		log.Fatal(err)
	}

	//seed stuff
	if *seed {
		fmt.Println("SEEDING THE DATABASE")
		seedAccounts(store)
	}

	server := NewAPIServer(":3451", store)
	server.Run()

}
