package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var accounts []Account

//Account is used to hold account details
type Account struct {
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	MobileNumber int64  `json:"mobilenumber"`
	Password     string `json:"password"`
}

//CreateAccount is used to create an account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	account := &Account{}
	//read request body from front end
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	//unmarshall body into the object called account
	err = json.Unmarshal(body, account)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//validate account details
	for _, acc := range accounts {
		if acc.FirstName == account.FirstName {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(" Sorry, Account already exists"))
			return
		}
	}
	accounts = append(accounts, *account)
	msg := fmt.Sprintf("Hello %s, we", account.FirstName)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(msg))
	fmt.Println(account)

}

//ListAccounts is used to list all the accounts
func ListAccounts(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(accounts)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

//GetDetailsByName is used to get specific account by their ID
func GetDetailsByName(w http.ResponseWriter, r *http.Request) {
	queryParams := mux.Vars(r)

	for _, acc := range accounts {
		if acc.FirstName == queryParams["name"] {

			resp, err := json.Marshal(acc)

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
			return
		}

	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("no Account found"))
}

//DeleteAccount is used to delete an account
func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	a := []Account{}
	queryParams := mux.Vars(r)
	for _, acc := range accounts {
		if acc.FirstName != queryParams["name"] {
			a = append(a, acc)
		}
	}
	accounts = a
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Account deleted"))
}

//UpdateAccount is used to update an account
func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var isUpdated bool
	account := &Account{}
	a := []Account{}
	queryParams := mux.Vars(r)
	for _, acc := range accounts {
		if acc.FirstName == queryParams["name"] {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}

			//unmarshall body into the object called account
			err = json.Unmarshal(body, account)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			acc.FirstName = account.FirstName
			acc.LastName = account.LastName
			acc.MobileNumber = account.MobileNumber
			acc.Password = account.Password
			isUpdated = true
		}
		a = append(a, acc)
	}

	if isUpdated {

		accounts = a
		msg := fmt.Sprintf("hello %s, your account is updated", account.FirstName)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))

	} else {
		msg := fmt.Sprintln("no such account")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg))
	}

}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/account", CreateAccount).Methods(http.MethodPost) // to create an account

	router.HandleFunc("/account", ListAccounts).Methods(http.MethodGet) //to get the details ofall account

	router.HandleFunc("/account/{name}", GetDetailsByName).Methods(http.MethodGet) //details of specific account
	router.HandleFunc("/account/{name}", DeleteAccount).Methods(http.MethodDelete) //Delete an account
	router.HandleFunc("/account/{name}", UpdateAccount).Methods(http.MethodPut)    //to update an account
	http.ListenAndServe(":8080", router)
	fmt.Println("hey we are listening in: 8080")

}
