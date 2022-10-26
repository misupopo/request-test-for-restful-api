package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
)

type ENV struct {
	SAP struct {
		Host     string `json:"host" form:"host"`
		PORT     string `json:"port" form:"port"`
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
		Bank     struct {
			Country string `json:"country" form:"country"`
			BankId  string `json:"bankId" form:"bankId"`
		}
		Client     int    `json:"client" form:"client"`
		XCSRFToken string `json:"xCSRFToken" form:"xCSRFToken"`
	}
}

func main() {
	// intellij だと GOROOT の path がプロジェクトの top になる
	// なので、GOROOT からの相対パスで import すると、
	env, err := ioutil.ReadFile("./env.json")
	//fmt.Printf("env: %v \n", env)

	var environment ENV
	err = json.Unmarshal(env, &environment)

	fmt.Printf("environment: %v \n", environment.SAP.Host)

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

	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("request returns error: %v", err)
	}

	fmt.Printf("response: %v \n", response)
	fmt.Printf("statusCode: %v \n", response.StatusCode)
	fmt.Printf("Body: %v \n", response.Body)
}
