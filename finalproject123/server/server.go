package server

import (
	"context"
	"encoding/json"
	"finalproject123/serverfunctional"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)
var accountsList = serverfunctional.AccountsList{
	Accounts: make([]serverfunctional.Account, 0),
}

func Run() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)

	handler := http.NewServeMux()

	handler.HandleFunc("/account/crud/", Logger(HandleAccountCRUD))
	handler.HandleFunc("/accounts/", Logger(HandleAccounts))
	handler.HandleFunc("/account/score/", Logger(HandleAccountScore))
	handler.HandleFunc("/account/transactions/", Logger(HandleAccountListOfTransactions))

	server := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}


	go func() {
		log.Printf("Listening on http://%s\n", server.Addr)
		log.Fatal(server.ListenAndServe())
	}()

	graceful(server, 5*time.Second)
}

func graceful(server *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Printf("\nShutdown with timeout: %s\n", timeout)

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
}
// HANDLE FUNCTIONS
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request){

		log.Printf("server [net/http] method [%s]  connection from [%v]", request.Method,
			request.RemoteAddr)

		next.ServeHTTP(responseWriter, request)
	}
}

func HandleAccounts(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	if request.Method == http.MethodGet {
		HandleGetAccounts(responseWriter, request)
	} else {
		HandleMethodIsNotAllowed(responseWriter, request)

	}
}

func HandleGetAccounts(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
	accounts, _ := json.Marshal(accountsList.GetAccounts())

	responseWriter.Write(accounts)
}

func HandleAccountCRUD(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	if request.Method == http.MethodGet {
		HandleGetAccount(responseWriter, request)

	} else if request.Method == http.MethodPost {
		HandleAddAccount(responseWriter, request)

	} else if request.Method == http.MethodPut {
		HandleUpdateAccount(responseWriter, request)

	} else if request.Method == http.MethodDelete {
		HandleDeleteAccount(responseWriter, request)

	} else {
		HandleMethodIsNotAllowed(responseWriter, request)

	}
}

func HandleMethodIsNotAllowed(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusBadRequest)
	msg, _ := json.Marshal(fmt.Sprintf("Method %s not allowed", request.Method))
	responseWriter.Write(msg)
}

func HandleGetAccount(responseWriter http.ResponseWriter, request *http.Request) {
	accountid := strings.Replace(request.URL.Path, "/account/crud/", "", 1)

	account := accountsList.GetAccountByID(accountid)

	if account == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("Account with id %s not found", accountid))

		responseWriter.Write(error)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	accountidJson, _ := json.Marshal(account)

	responseWriter.Write(accountidJson)
}

func HandleAddAccount(responseWriter http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var account serverfunctional.Account

	err := decoder.Decode(&account)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		error, _ := json.Marshal(fmt.Sprintf("Bad request. %v", err))

		responseWriter.Write(error)
		return
	}


	err = accountsList.AddAccount(account)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		error, _ := json.Marshal(fmt.Sprintf("Bad request. %v", err))

		responseWriter.Write(error)
		return
	}

	HandleGetAccounts(responseWriter, request)
}

func HandleUpdateAccount(responseWriter http.ResponseWriter, request *http.Request) {
	accountid := strings.Replace(request.URL.Path, "/account/crud/", "", 1)

	decoder := json.NewDecoder(request.Body)

	var account serverfunctional.Account

	err := decoder.Decode(&account)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		error, _ := json.Marshal(fmt.Sprintf("Bad request. %v", err))

		responseWriter.Write(error)
		return
	}

	account.ID = accountid

	err = accountsList.SetAccount(account)

	if err != nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("%v", err))

		responseWriter.Write(error)

		return
	}

	HandleGetAccount(responseWriter, request)
}

func HandleDeleteAccount(responseWriter http.ResponseWriter, request *http.Request) {
	accountid := strings.Replace(request.URL.Path, "/account/crud/", "", 1)
	err := accountsList.DeleteAccount(accountid)

	if err != nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("%v", err))

		responseWriter.Write(error)

		return
	}

	HandleGetAccounts(responseWriter, request)
}

func HandleAccountScore(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	if request.Method == http.MethodGet {
		HandleGetAccountScore(responseWriter, request)

	}else if request.Method == http.MethodPost{
		HandleFromScoreToScore(responseWriter, request)

	} else if request.Method == http.MethodPut{
		HandleAddToScore(responseWriter, request)

	} else if request.Method == http.MethodDelete{
		HandleDeleteAccountScore(responseWriter, request)

	} else {
		HandleMethodIsNotAllowed(responseWriter, request)

	}
}

func HandleGetAccountScore(responseWriter http.ResponseWriter, request *http.Request){
	accountid := strings.Replace(request.URL.Path, "/account/score/", "", 1)

	account := accountsList.GetAccountScore(accountid)

	if account == "" {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("Account with id %s not found", accountid))

		responseWriter.Write(error)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	accountidJson, _ := json.Marshal(account)

	responseWriter.Write(accountidJson)
}

func HandleFromScoreToScore(responseWriter http.ResponseWriter, request *http.Request)  {
	accountidandsum := strings.Replace(request.URL.Path, "/account/score/", "", 1)

	pair := strings.SplitN(accountidandsum, ".", 3)

	list := accountsList.FromScoreToScore(pair)

	if list == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("Account with id %s not found", pair[0]))

		responseWriter.Write(error)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	accountidJson, _ := json.Marshal(list)

	responseWriter.Write(accountidJson)
}

func HandleAddToScore(responseWriter http.ResponseWriter, request *http.Request){
	accountidandsum := strings.Replace(request.URL.Path, "/account/score/", "", 1)

	pair := strings.SplitN(accountidandsum, ".", 2)

	list := accountsList.AddToScore(pair)

	if list == "" {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("Account with id %s not found", pair[0]))

		responseWriter.Write(error)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	accountidJson, _ := json.Marshal(list)

	responseWriter.Write(accountidJson)
}

func HandleDeleteAccountScore(responseWriter http.ResponseWriter, request *http.Request){
	accountid := strings.Replace(request.URL.Path, "/account/score/", "", 1)

	account := accountsList.DeleteFromScore(accountid)

	if account == "" {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("Account with id %s not found", accountid))

		responseWriter.Write(error)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	accountidJson, _ := json.Marshal(account)

	responseWriter.Write(accountidJson)
}

func HandleAccountListOfTransactions(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	if request.Method == http.MethodGet {
		HandleGetAccountListOfTransactionsByDate(responseWriter, request)

	} else if request.Method == http.MethodDelete{
		HandleDeleteAccountListOfTransactions(responseWriter, request)
	} else {
		HandleMethodIsNotAllowed(responseWriter, request)

	}
}

func HandleGetAccountListOfTransactionsByDate(responseWriter http.ResponseWriter, request *http.Request){
	accountidanddate := strings.Replace(request.URL.Path, "/account/transactions/", "", 1)

	pair := strings.SplitN(accountidanddate, ".", 4)

	list := accountsList.GetAccountListOfTransactionsByDate(pair)

	if list == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("Account with id %s not found", pair[0]))

		responseWriter.Write(error)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	accountidJson, _ := json.Marshal(list)

	responseWriter.Write(accountidJson)
}

func HandleDeleteAccountListOfTransactions(responseWriter http.ResponseWriter, request *http.Request){
	accountid := strings.Replace(request.URL.Path, "/account/transactions/", "", 1)

	list := accountsList.DeleteAllTransactions(accountid)

	if list == nil {
		responseWriter.WriteHeader(http.StatusNotFound)
		error, _ := json.Marshal(fmt.Sprintf("Account with id %s not found", accountid))

		responseWriter.Write(error)

		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	accountidJson, _ := json.Marshal(list)

	responseWriter.Write(accountidJson)
}