package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type InAccounts struct {
	Username     string
	AuthToken    string
	ApiUrl       string
	ResultsCount int
}

func GetAccounts(in InAccounts) ([]Account, error) {
	uri := fmt.Sprintf("%s/api/v2/search?q=%s&resolve=true?limit=%d", in.ApiUrl, in.Username, in.ResultsCount)

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	req.Header.Set("Content-Type", "application/json") // => your content-type

	req.Header.Add("Authorization", in.AuthToken)

	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(response.Body)

	var res Users
	if err := json.Unmarshal(body, &res); err != nil { // Parse []byte to go struct pointer
		log.Fatal("Can not unmarshal JSON")
	}

	var ret []Account
	for _, r := range res.Accounts {
		var n Account
		n.ID = r.ID
		n.UserName = r.Username
		n.DisplayName = r.DisplayName
		n.URL = r.URL
		n.FollowersCount = r.FollowersCount
		n.FollowingCount = r.FollowingCount
		ret = append(ret, n)
	}

	return ret, err
}

type Account struct {
	ID             string
	UserName       string
	DisplayName    string
	URL            string
	FollowersCount int
	FollowingCount int
}

type Users struct {
	Accounts []struct {
		ID             string        `json:"id"`
		Username       string        `json:"username"`
		Acct           string        `json:"acct"`
		DisplayName    string        `json:"display_name"`
		Locked         bool          `json:"locked"`
		Bot            bool          `json:"bot"`
		Discoverable   bool          `json:"discoverable"`
		Group          bool          `json:"group"`
		CreatedAt      time.Time     `json:"created_at"`
		Note           string        `json:"note"`
		URL            string        `json:"url"`
		Avatar         string        `json:"avatar"`
		AvatarStatic   string        `json:"avatar_static"`
		Header         string        `json:"header"`
		HeaderStatic   string        `json:"header_static"`
		FollowersCount int           `json:"followers_count"`
		FollowingCount int           `json:"following_count"`
		StatusesCount  int           `json:"statuses_count"`
		LastStatusAt   string        `json:"last_status_at"`
		Noindex        bool          `json:"noindex"`
		Emojis         []interface{} `json:"emojis"`
		Fields         []interface{} `json:"fields"`
	} `json:"accounts"`
	Statuses []interface{} `json:"statuses"`
	Hashtags []struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		History []struct {
			Day      string `json:"day"`
			Accounts string `json:"accounts"`
			Uses     string `json:"uses"`
		} `json:"history"`
		Following bool `json:"following"`
	} `json:"hashtags"`
}
