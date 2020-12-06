package account

import (
	"encoding/json"
	"net/http"

	"github.com/brushed-charts/backend/institutions/oanda/src/util"
	"github.com/pkg/errors"
)

var (
	oandaAccountURLPath = "/v3/accounts"
)

type account struct {
	ID   string   `json:"id"`
	Tags []string `json:"tags"`
}
type accountList struct {
	Accounts []account `json:"accounts"`
}

// GetAccountID retrieve one ID in the IDs returned by
// oanda servers
func GetAccountID(client *http.Client, token, url string) (string, error) {
	request, err := makeAccountIDRequest(token, url)
	if err != nil {
		return "", err
	}

	response, err := sendRequest(client, request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	accounts, err := getAccountListFromBody(response)
	if err != nil {
		return "", err
	}

	firstAccountID := accounts[0].ID

	return firstAccountID, nil
}

func makeAccountIDRequest(token, url string) (*http.Request, error) {
	urlWithPath := url + oandaAccountURLPath
	req, err := util.MakeBearerGetRequest(urlWithPath, token)
	return req, err
}

func sendRequest(client *http.Client, request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)
	if err != nil || util.IsHTTPResponseError(response) {
		body := util.TryReadingResponseBody(response)
		err := errors.New("Error when fetching accountID from oanda\n" + body)
		return nil, err
	}

	return response, nil
}

func getAccountListFromBody(response *http.Response) ([]account, error) {
	var accList accountList
	if response == nil {
		return []account{}, nil
	}
	err := json.NewDecoder(response.Body).Decode(&accList)
	if err != nil {
		return []account{}, errors.New("Account -- Error during JSON parsing\n" + err.Error())
	}
	return accList.Accounts, nil
}