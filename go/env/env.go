package env

import (
	"encoding/json"
	"io/ioutil"
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
	SALESFORCE struct {
		ClientId     string `json:"clientId" form:"clientId"`
		ClientSecret string `json:"clientSecret" form:"clientSecret"`
		GrantType    string `json:"grantType" form:"grantType"`
		Username     string `json:"username" form:"username"`
		Password     string `json:"password" form:"password"`
		LoginUrl     string `json:"loginUrl" form:"loginUrl"`
	}
}

func NewEnv() (*ENV, error) {
	// intellij だと GOROOT の path がプロジェクトの top になる
	// なので、GOROOT からの相対パスで import すると、
	env, err := ioutil.ReadFile("./env.json")
	//fmt.Printf("env: %v \n", env)

	if err != nil {
		return nil, err
	}

	var environment ENV
	err = json.Unmarshal(env, &environment)

	return &environment, nil
}
