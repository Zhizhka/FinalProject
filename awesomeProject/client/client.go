package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Date struct {
	Day string `json:"day"`
	Month string `json:"month"`
	Year string `json:"year"`
}

type Transaction struct {
	ID string `json:"id"`
	Sum string `json:"sum"`
	Date *Date `json:"date"`
}

type Account struct {
	ID string `json:"id"`
	Score string `json:"score"`
	Transactions []Transaction `json:"transactions"`
}

type AccountsList struct {
	Accounts []Account
}

var ListOfHTTPRequests = GetListOfHTTPRequests()
var ListOfMethods = GetListsOfMethods()

func Menu(){
	fmt.Println("Hello user, glad to see you!")
	fmt.Println("\nPlease, choose action:\n1 - get list of all your accounts\n2 - use CRUD functions\n3 - exit")
	var chose int
	fmt.Scan(&chose)
	client  := &http.Client{}
	if chose == 1{
		GoRequest(client, MakeRequest(ListOfMethods[0], ListOfHTTPRequests[0], nil))
		Menu()
	}else if chose == 2{
		CRUDMenu(client)
		Menu()
	}else if chose == 3{
		return
	}else{
		fmt.Println("\nWrong value, please, try again")
		Menu()
	}
}

func CRUDMenu(client *http.Client){
	fmt.Println("\nPlease, choose action:\n1 - get account by ID\n2 - create account\n3 - update account")
	fmt.Println("\n4 - delete account\n5 - get back")
	var chose int
	fmt.Scan(&chose)
	if chose == 1{
		fmt.Println("\nEnter id of account")
		var accountid string
		fmt.Scan(&accountid)
		GoRequest(client, MakeRequest(ListOfMethods[0], ListOfHTTPRequests[1] + accountid, nil))
		MenuOfWorkWithAccount(client, accountid)
		CRUDMenu(client)
	}else if chose == 2{
		account := MenuOfCreatingAccount()
		accountJson, _ := json.Marshal(account)
		body := bytes.NewReader(accountJson)
		GoRequest(client, MakeRequest(ListOfMethods[1], ListOfHTTPRequests[1], body))
		CRUDMenu(client)
	}else if chose == 3{
		fmt.Println("\nEnter id of account")
		var accountid string
		fmt.Scan(&accountid)
		account := MenuOfUpdatingAccount(client, accountid)
		accountJson, _ := json.Marshal(account)
		body := bytes.NewReader(accountJson)
		GoRequest(client, MakeRequest(ListOfMethods[2], ListOfHTTPRequests[1] + accountid, body))
		CRUDMenu(client)
	}else if chose == 4{
		fmt.Println("\nEnter id of account")
		var accountid string
		fmt.Scan(&accountid)
		GoRequest(client, MakeRequest(ListOfMethods[3], ListOfHTTPRequests[1] + accountid, nil))
		CRUDMenu(client)
	}else if chose == 5{
		return
	}else {
		fmt.Println("\nWrong value, please, try again")
		CRUDMenu(client)
	}
}

func MenuOfCreatingAccount() Account{
	var id string
	var score string
	fmt.Println("\nNow you should create new account")
	fmt.Println("\nPlease, enter your ID:")
	fmt.Scan(&id)
	fmt.Println("\nPlease, enter your Score:")
	fmt.Scan(&score)
	account := Account{
		ID: id,
		Score: score,
		Transactions: make([]Transaction, 0),
	}

	return account
}

func MenuOfUpdatingAccount(client *http.Client, accountid string) Account{
	response, _ := client.Do(MakeRequest(ListOfMethods[0], ListOfHTTPRequests[1] + accountid, nil))
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s", data)
	response, _ = client.Do(MakeRequest(ListOfMethods[0], ListOfHTTPRequests[1] + accountid, nil))
	decoder := json.NewDecoder(response.Body)
	var account Account
	err := decoder.Decode(&account)
	if err != nil {
		panic(err)
	}
	var id string
	var score string
	var answer int
	fmt.Println("\nNow you should update this account")
	fmt.Println("\nPlease, enter new ID:")
	fmt.Scan(&id)
	fmt.Println("\nPlease, enter new Score:")
	fmt.Scan(&score)
	fmt.Println("\nDecide 1 - update transactions, else - not")
	fmt.Scan(&answer)
	var transactions []Transaction
	if answer == 1{
		transactions = UpdateTransactions(len(account.Transactions))
	}else {
		transactions = account.Transactions
	}
	account.ID = id
	account.Score = score
	account.Transactions = transactions

	return account
}

func UpdateTransactions(len int) []Transaction {
	fmt.Println("Please, update list of transactions")
	var transactions []Transaction
	transactions = make([]Transaction, 0)
	for i := 0; i < len; i++{
		var transaction Transaction
		var id string
		var sum string
		fmt.Println("\nPlease, enter new id")
		fmt.Scan(&id)
		fmt.Println("\nPlease, enter new Sum:")
		fmt.Scan(&sum)
		fmt.Println("Please, enter new date: year, month, day")
		var date Date
		var year string
		var month string
		var day string
		fmt.Scan(&year)
		fmt.Scan(&month)
		fmt.Scan(&day)
		date.Day = day
		date.Month = month
		date.Year = year
		transaction.ID = id
		transaction.Sum = sum
		transaction.Date = &date
		transactions = append(transactions, transaction)
	}

	return  transactions
}

func MenuOfWorkWithAccount(client *http.Client, accountid string){
	fmt.Println("\nNow, please, choose:\n1 - action with account score\n2 - action with account transactions\n3 - exit")
	var answer int
	fmt.Scan(&answer)
	if answer == 1{
		MenuOfWorkWithScore(client, accountid)
		MenuOfWorkWithAccount(client, accountid)
	}else if answer == 2{
		MenuOfWorkWithTransactions(client, accountid)
		MenuOfWorkWithAccount(client, accountid)
	}else if answer == 3{
		return
	}else{
		fmt.Println("\nThis is uncorrect value")
		MenuOfWorkWithAccount(client, accountid)
	}
}

func MenuOfWorkWithScore(client *http.Client, accountid string)  {
	fmt.Println("\nNow, please, choose:\n1 - get score\n2 - add money to score\n3 - delete score\n4 - exit")
	fmt.Println("5 - convert money from one account to enouther")
	var answer int
	fmt.Scan(&answer)
	if answer == 1{
		GoRequest(client, MakeRequest(ListOfMethods[0], ListOfHTTPRequests[2] + accountid, nil))
		MenuOfWorkWithScore(client, accountid)
	}else if answer == 2{
		fmt.Println("\nPlease, enter the sum which you want to add to score")
		var sumText string
		fmt.Scan(&sumText)
		GoRequest(client, MakeRequest(ListOfMethods[2], ListOfHTTPRequests[2] + accountid + "." + sumText, nil))
		MenuOfWorkWithScore(client, accountid)
	}else if answer == 3{
		GoRequest(client, MakeRequest(ListOfMethods[3], ListOfHTTPRequests[2] + accountid, nil))
		MenuOfWorkWithScore(client, accountid)
	}else if answer == 4{
		return
	}else if answer == 5{
		var sumText string
		var id string
		fmt.Println("\nPlease, enter the id of account of convertation")
		fmt.Scan(&id)
		fmt.Println("\nPlease, enter the sum which you want to convert")
		fmt.Scan(&sumText)
		GoRequest(client, MakeRequest(ListOfMethods[1], ListOfHTTPRequests[2] + accountid + "." + id + "." + sumText,
			nil))
		MenuOfWorkWithScore(client, accountid)
	} else{
		fmt.Println("\nUncorrect value")
		MenuOfWorkWithScore(client, accountid)
	}
}

func MenuOfWorkWithTransactions(client *http.Client, accountid string){
	fmt.Println("\nNow, please, choose:\n1 - get list of transactions by date\n2 - delete list of transactions\n3 - exit")
	var answer int
	fmt.Scan(&answer)
	if answer == 1{
		fmt.Println("\nPlease, enter the date to see all transactions under it: year, month, day")
		var year string
		var month string
		var day string
		fmt.Scan(&year)
		fmt.Scan(&month)
		fmt.Scan(&day)
		GoRequest(client, MakeRequest(ListOfMethods[0],
			ListOfHTTPRequests[3] + accountid + "." + day + "." + month + "." + year, nil))
		MenuOfWorkWithTransactions(client, accountid)
	}else if answer == 2{
		GoRequest(client, MakeRequest(ListOfMethods[3], ListOfHTTPRequests[3] + accountid, nil))
		MenuOfWorkWithTransactions(client, accountid)
	}else if answer == 3{
		return
	}else{
		fmt.Println("\nUncorect value")
		MenuOfWorkWithTransactions(client, accountid)
	}
}

func MakeRequest(method string, url string, body io.Reader) *http.Request{
	request, _ := http.NewRequest(method, url, body)
	request.Header.Set("X-Custom-Header", "myvalue")
	request.Header.Set("Content-Type", "application/json")

	return request
}

func GoRequest(client *http.Client, request *http.Request){
	response, _ := client.Do(request)
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	fmt.Printf("%s", string(data))
}