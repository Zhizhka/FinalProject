package serverfunctional

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// TYPES  AND  VARS
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


// USUAL FUNCTIONS
func (accountList AccountsList) GetAccounts() []Account{
	return accountList.Accounts
}

func (accountList AccountsList) GetAccountByID(id string) *Account {
	for _, account := range accountList.Accounts {
		if account.ID == id {
			return &account
		}
	}

	return nil
}

func (accountList *AccountsList) AddAccount(account Account) error{
	localaccount := accountList.GetAccountByID(account.ID)
	if localaccount != nil {
		return errors.New(fmt.Sprintf("Account with id %s already exists", account.ID))

	}
	accountList.Accounts = append(accountList.Accounts, account)

	return nil
}

func (accountList *AccountsList) SetAccount(account Account) error{
	for i, localaccount := range accountList.Accounts {
		if localaccount.ID == account.ID {

			accountList.Accounts[i] = account

			return nil
		}
	}

	return errors.New(fmt.Sprintf("There is no account with id %s", account.ID))
}

func (accountList *AccountsList) DeleteAccount(id string) error {
	for i, account := range accountList.Accounts {
		if account.ID == id {

			accountList.Accounts = append(accountList.Accounts[:i], accountList.Accounts[i+1:]...)

			return nil
		}
	}

	return errors.New(fmt.Sprintf("There is no account with id %s", id))
}

func (accountList AccountsList) GetAccountScore(id string) string{
	for _, account := range accountList.Accounts {
		if account.ID == id{

			return account.Score
		}
	}

	return ""
}

func (accountList AccountsList) AddToScore(pair []string) string{
	for _, account := range accountList.Accounts {
		if account.ID == pair[0]{
			var transaction = CreateTransaction()
			iscore, _ := strconv.Atoi(account.Score)
			sumForAdd, _ := strconv.Atoi(pair[1])
			iscore += sumForAdd
			sscore := strconv.Itoa(iscore)
			account.Score = sscore
			transaction.ID = strconv.Itoa(len(account.Transactions)+1)
			transaction.Sum = pair[1]
			account.Transactions = append(account.Transactions, transaction)

			return account.Score
		}
	}

	return ""
}

func (accountList AccountsList) FromScoreToScore(pair []string) []Account{
	var accounts []Account
	for _, account1 := range accountList.Accounts {
		if account1.ID == pair[0]{
			var transaction = CreateTransaction()
			iscore, _ := strconv.Atoi(account1.Score)
			sumForDelete, _ := strconv.Atoi(pair[2])
			iscore -= sumForDelete
			sscore := strconv.Itoa(iscore)
			account1.Score = sscore
			transaction.ID = strconv.Itoa(len(account1.Transactions) + 1)
			transaction.Sum = pair[2]
			account1.Transactions = append(account1.Transactions, transaction)
			accounts = append(accounts, account1)
		}
	}

	for _, account2 := range accountList.Accounts {
		if account2.ID == pair[1]{
			var transaction = CreateTransaction()
			iscore, _ := strconv.Atoi(account2.Score)
			sumForAdd, _ := strconv.Atoi(pair[2])
			iscore += sumForAdd
			sscore := strconv.Itoa(iscore)
			account2.Score = sscore
			transaction.ID = strconv.Itoa(len(account2.Transactions) + 1)
			transaction.Sum = pair[2]
			account2.Transactions = append(account2.Transactions, transaction)
			accounts = append(accounts, account2)

			return accounts
		}
	}

	return nil
}

func (accountList AccountsList) DeleteFromScore(id string) string{
	for _, account := range accountList.Accounts {
		if account.ID == id{
			var transaction = CreateTransaction()
			iscore, _ := strconv.Atoi(account.Score)
			iscore = 0
			sscore := strconv.Itoa(iscore)
			account.Score = sscore
			transaction.ID = strconv.Itoa(len(account.Transactions) + 1)
			transaction.Sum = "0"
			account.Transactions = append(account.Transactions, transaction)

			return account.Score
		}
	}

	return ""
}

func CreateTransaction() Transaction{
	var curenttime = time.Now().Format("01-01-2006")
	var transaction Transaction
	var date Date
	datearray := strings.SplitN(curenttime,"-", 3)
	date.Year = string(datearray[2])
	date.Month = string(datearray[1])
	date.Day = string(datearray[0])
	transaction.Date = &date

	return transaction
}

func (accountList AccountsList) GetAccountListOfTransactionsByDate(pair []string) []Transaction{
	runtime.GC()
	for _, account := range accountList.Accounts {
		if account.ID == pair[0]{
			transactions := account.Transactions

			dateYear, _ := strconv.Atoi(pair[3])
			dateYear *= 365
			dateMonth, _ := strconv.Atoi(pair[2])
			dateMonth *= 30
			dateDay, _ := strconv.Atoi(pair[1])
			dateSumOfDays := 0
			dateSumOfDays += dateYear
			dateSumOfDays += dateMonth
			dateSumOfDays += dateDay

			for j := 0; j < len(account.Transactions); j++ {
				for i := 0; i < len(transactions); i++ {
					trYear, _ := strconv.Atoi(transactions[i].Date.Year)
					trYear *= 365
					trMonth, _ := strconv.Atoi(transactions[i].Date.Month)
					trMonth *= 30
					trDay, _ := strconv.Atoi(transactions[i].Date.Day)
					trSumOfDays := 0
					trSumOfDays += trYear
					trSumOfDays += trMonth
					trSumOfDays += trDay

					if trSumOfDays < dateSumOfDays {
						transactions = append(transactions[:i], transactions[i+1:]...)
						break
					}
				}
			}

			return transactions
		}
	}

	return nil
}

func (accountList AccountsList) DeleteAllTransactions(id string) []Transaction{
	for _, account := range accountList.Accounts {
		if account.ID == id{
			transactions := make([]Transaction, 0)
			account.Transactions = transactions

			return account.Transactions
		}
	}

	return nil
}