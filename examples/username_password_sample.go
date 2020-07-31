// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package main

import (
	"fmt"
	"reflect"

	msalgo "github.com/AzureAD/microsoft-authentication-library-for-go/src/msal"
	log "github.com/sirupsen/logrus"
)

func tryUsernamePasswordFlow(publicClientApp *msalgo.PublicClientApplication) {
	userNameParams := msalgo.CreateAcquireTokenUsernamePasswordParameters(config.Scopes, config.Username, config.Password)
	result, err := publicClientApp.AcquireTokenByUsernamePassword(userNameParams)
	if err != nil {
		log.Fatal(err)
	}
	accessToken := result.GetAccessToken()
	log.Info("Access token is: " + accessToken)
}

func acquireByUsernamePasswordPublic() {
	config := createConfig("config.json")
	// Creating the Public Client Application
	pcaParams := createPCAParams(config.ClientID, config.Authority)
	publicClientApp, err := msalgo.CreatePublicClientApplication(pcaParams)
	if err != nil {
		log.Fatal(err)
	}
	publicClientApp.SetCacheAccessor(cacheAccessor)
	var userAccount msalgo.IAccount
	accounts := publicClientApp.GetAccounts()
	for _, account := range accounts {
		if account.GetUsername() == config.Username {
			userAccount = account
		}
	}
	if reflect.ValueOf(userAccount).IsNil() {
		log.Info("No valid account found")
		tryUsernamePasswordFlow(publicClientApp)
	} else {
		silentParams := msalgo.CreateAcquireTokenSilentParameters(config.Scopes, userAccount)
		result, err := publicClientApp.AcquireTokenSilent(silentParams)
		if err != nil {
			log.Info(err)
			tryUsernamePasswordFlow(publicClientApp)
		} else {
			fmt.Println("Access token is " + result.GetAccessToken())
		}
	}
}