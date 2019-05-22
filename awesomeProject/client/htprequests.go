package client

func GetListOfHTTPRequests() []string{
	var ListOfHTTPRequests []string
	ListOfHTTPRequests = append(ListOfHTTPRequests, "http://localhost:8080/accounts/")
	ListOfHTTPRequests = append(ListOfHTTPRequests, "http://localhost:8080/account/crud/")
	ListOfHTTPRequests = append(ListOfHTTPRequests, "http://localhost:8080/account/score/")
	ListOfHTTPRequests = append(ListOfHTTPRequests, "http://localhost:8080/account/transactions/")

	return ListOfHTTPRequests
}

func GetListsOfMethods() []string{
	var ListOfMethods []string
	ListOfMethods = append(ListOfMethods, "GET")
	ListOfMethods = append(ListOfMethods, "POST")
	ListOfMethods = append(ListOfMethods, "PUT")
	ListOfMethods = append(ListOfMethods, "DELETE")

	return  ListOfMethods
}