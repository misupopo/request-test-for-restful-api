package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"request-test-for-restful-api/env"
	"request-test-for-restful-api/salesforce"
	"strconv"
	"strings"
)

func main() {
	environment, err := env.NewEnv()

	if err != nil {
		fmt.Printf("New environment error: %v \n", err)
	}

	fmt.Printf("environment: %v \n", environment.SAP.Host)

	// sap へのリクエストテスト
	//requestSap(environment)

	// salesforce へのリクエストテスト
	rsf := requestSalesforce(environment)
	requestSalesforceApi(rsf)
}

func requestSap(environment *env.ENV) {
	method := "GET"
	service := "sap/opu/odata4/sap/api_bank/srvd_a2x/sap/api_bank_2/0001/Bank"

	requestUrl := fmt.Sprintf("http://%v:%v/%v", environment.SAP.Host, environment.SAP.PORT, service)

	body, err := json.Marshal(map[string]string{})
	if err != nil {
		fmt.Printf("API request error: %v \n", err)
	}

	req, err := http.NewRequest(method, requestUrl, strings.NewReader(string(body)))
	if err != nil {
		fmt.Printf("HTTP NewRequest error: %v \n", err)
	}

	csrfToken := environment.SAP.XCSRFToken

	req.SetBasicAuth(environment.SAP.Username, environment.SAP.Password)
	req.Header.Add("x-csrf-token", csrfToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	fmt.Sprintf("BankCountry eq '%v' and Bank eq '%v'", environment.SAP.Bank.Country, environment.SAP.Bank.BankId)

	params := map[string]string{
		"sap-client": strconv.Itoa(environment.SAP.Client),
		"$filter":    fmt.Sprintf("BankCountry eq '%v' and Bank eq '%v'", environment.SAP.Bank.Country, environment.SAP.Bank.BankId),
	}

	parameter := req.URL.Query()
	for k, v := range params {
		parameter.Add(k, v)
	}
	req.URL.RawQuery = parameter.Encode()

	// なぜか下記は必要がなくなった
	//unescapedQuery, err := url.QueryUnescape(req.URL.RawQuery)
	//fmt.Printf("req.URL.RawQuery: %v \n", req.URL.RawQuery)
	//
	//// クエリパラメータの文字列を SAP が求める形式に変更する
	//// スペースを + に変更する
	//unescapedQuery = strings.ReplaceAll(unescapedQuery, " ", "+")
	//fmt.Printf("unescapedQuery: %v \n", unescapedQuery)
	//req.URL.RawQuery = unescapedQuery

	j, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: j,
	}

	fmt.Printf("start request \n")

	response, err := client.Do(req)

	fmt.Printf("complete request \n")

	if err != nil {
		fmt.Printf("request returns error: %v", err)
	}

	fmt.Printf("response: %v \n", response)
	fmt.Printf("statusCode: %v \n", response.StatusCode)
	fmt.Printf("Body: %v \n", response.Body)
}

func requestSalesforce(environment *env.ENV) *salesforce.OAuthInfo {
	form := url.Values{}
	form.Add("grant_type", environment.SALESFORCE.GrantType)
	form.Add("client_id", environment.SALESFORCE.ClientId)
	form.Add("client_secret", environment.SALESFORCE.ClientSecret)
	form.Add("username", environment.SALESFORCE.Username)
	form.Add("password", environment.SALESFORCE.Password)

	resp, err := http.PostForm(environment.SALESFORCE.LoginUrl, form)
	if err != nil {
		fmt.Printf("failed to fetch response: %v", err)
	}
	defer safeClose(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP %s: %s", resp.Status, body)
	}

	fmt.Printf("body: %v \n", string(body))

	var oauthResp salesforce.OAuthInfo
	if err := json.Unmarshal(body, &oauthResp); err != nil {
		fmt.Printf("failed to unmarshal json to response struct: %v", err)
	}

	fmt.Printf("oauthResp AccessToken: %v \n", oauthResp.AccessToken)
	fmt.Printf("oauthResp InstanceUrl: %v \n", oauthResp.InstanceUrl)

	return &oauthResp
}

func safeClose(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			fmt.Printf("failed to close: %v", err)
		}
	}
}

func requestSalesforceApi(salesforceInfo *salesforce.OAuthInfo) {
	method := "GET"

	// https://latona--test.sandbox.my.salesforce.com/services/apexrest/ContractRelatedList/doGetContractRelatedList
	requestUrl := fmt.Sprintf("%v/services/apexrest/%v/%v", salesforceInfo.InstanceUrl, "ContractRelatedList", "doGetContractRelatedList")

	body, err := json.Marshal(map[string]string{})
	if err != nil {
		fmt.Printf("API request error: %v \n", err)
	}

	fmt.Printf("requestUrl: %v \n", requestUrl)
	fmt.Printf("strings.NewReader(string(body)): %v \n", strings.NewReader(string(body)))

	req, err := http.NewRequest(method, requestUrl, strings.NewReader(string(body)))
	if err != nil {
		fmt.Printf("HTTP NewRequest error: %v \n", err)
	}

	requiredHeaders := http.Header{}
	requiredHeaders.Add("Content-Type", "application/json")

	req.Header = requiredHeaders

	req.Header.Add("Authorization", "Bearer "+salesforceInfo.AccessToken)

	fmt.Printf("salesforceInfo.AccessToken: %v \n", salesforceInfo.AccessToken)

	client := &http.Client{}

	response, err := client.Do(req)

	fmt.Printf("response: %v \n", response)
	fmt.Printf("statusCode: %v \n", response.StatusCode)

	responseBody, _ := io.ReadAll(response.Body)

	// responseBody は []byte なので string に変換する
	fmt.Printf("body: %v \n", string(responseBody))
}
