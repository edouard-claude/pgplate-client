# Squirrel API PGPlate Client


Sample usage (./cmd/pgplate-client):

```go
package main

import (
	"encoding/json"
	"log"
	"time"

	pgplateclient "github.com/edouard-claude/pgplate-client"
)

func main() {

	client := pgplateclient.Client{
		BaseUrl:     "https://localhost:1324",
		OAuthID:     "00000",
		OAuthSecret: "99999",
	}

	c, err := pgplateclient.NewClient(&client)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	loginStruct := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "email@domain.tld",
		Password: "********",
	}

	login, err := json.Marshal(&loginStruct)
	if err != nil {
		log.Println(err)
	}

	var responseLogin struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Success      bool   `json:"success"`
	}

	c.Fetch("/v1/web/signin", pgplateclient.POST, login, &responseLogin)

	if !responseLogin.Success {
		log.Fatal("Login failed")
	}

	c.Jwt = &responseLogin.AccessToken
	c.RefreshToken = &responseLogin.RefreshToken

	var responseConstest struct {
		Contests []struct {
			ID              string    `json:"id"`
			Number          string    `json:"number"`
			Name            string    `json:"name"`
			StartDate       time.Time `json:"start_date"`
			EndDate         time.Time `json:"end_date"`
			ClosingDate     time.Time `json:"closing_date"`
			Departement     string    `json:"departement"`
			DepartementName string    `json:"departement_name"`
			RecognitionTime time.Time `json:"recognition_time"`
			ExternalURL     string    `json:"external_url"`
			TrialsTypes     []string  `json:"trials_types"`
			Days            []string  `json:"days"`
			Juries          []struct {
				Lic       string `json:"lic"`
				Num       string `json:"num"`
				Function  string `json:"function"`
				Lastname  string `json:"lastname"`
				Firstname string `json:"firstname"`
				FuncName  string `json:"func_name"`
			} `json:"juries"`
			Paddocks []struct {
				Lic       string `json:"lic"`
				Num       string `json:"num"`
				Function  string `json:"function"`
				Lastname  string `json:"lastname"`
				Firstname string `json:"firstname"`
				FuncName  string `json:"func_name"`
			} `json:"paddocks"`
			TrackLeaders []struct {
				Lic       string `json:"lic"`
				Num       string `json:"num"`
				Function  string `json:"function"`
				Lastname  string `json:"lastname"`
				Firstname string `json:"firstname"`
				FuncName  string `json:"func_name"`
			} `json:"track_leaders"`
			Participants []struct {
				ID        string `json:"id"`
				Firstname string `json:"firstname"`
				Lastname  string `json:"lastname"`
				Avatar    string `json:"avatar"`
				Created   bool   `json:"created"`
			} `json:"participants"`
			ParticipantsCount int       `json:"participants_count"`
			FirstStartTime    time.Time `json:"first_start_time"`
			Status            bool      `json:"status"`
			Address           string    `json:"address"`
		} `json:"contests"`
		Success bool `json:"success"`
	}

	c.Fetch("/v1/web/contests/admin/0/10", pgplateclient.GET, nil, &responseConstest)

	if !responseConstest.Success {
		log.Fatal("List contests failed")
	}

	for _, c := range responseConstest.Contests {
		log.Println(c.ID, c.Name, c.Number)
	}

}

```