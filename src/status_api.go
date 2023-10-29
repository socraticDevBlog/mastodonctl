package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Status struct {
	Account struct {
		Acct           string    `json:"acct"`
		Avatar         string    `json:"avatar"`
		AvatarStatic   string    `json:"avatar_static"`
		Bot            bool      `json:"bot"`
		CreatedAt      time.Time `json:"created_at"`
		Discoverable   bool      `json:"discoverable"`
		DisplayName    string    `json:"display_name"`
		Emojis         []any     `json:"emojis"`
		Fields         []any     `json:"fields"`
		FollowersCount int       `json:"followers_count"`
		FollowingCount int       `json:"following_count"`
		Group          bool      `json:"group"`
		Header         string    `json:"header"`
		HeaderStatic   string    `json:"header_static"`
		ID             string    `json:"id"`
		LastStatusAt   string    `json:"last_status_at"`
		Locked         bool      `json:"locked"`
		Noindex        bool      `json:"noindex"`
		Note           string    `json:"note"`
		Roles          []any     `json:"roles"`
		StatusesCount  int       `json:"statuses_count"`
		URI            string    `json:"uri"`
		URL            string    `json:"url"`
		Username       string    `json:"username"`
	} `json:"account"`
	Application struct {
		Name    string `json:"name"`
		Website any    `json:"website"`
	} `json:"application"`
	Bookmarked bool `json:"bookmarked"`
	Card       struct {
		AuthorName       string `json:"author_name"`
		AuthorURL        string `json:"author_url"`
		Blurhash         string `json:"blurhash"`
		Description      string `json:"description"`
		EmbedURL         string `json:"embed_url"`
		Height           int    `json:"height"`
		Html             string `json:"html"`
		Image            string `json:"image"`
		ImageDescription string `json:"image_description"`
		Language         string `json:"language"`
		ProviderName     string `json:"provider_name"`
		ProviderURL      string `json:"provider_url"`
		PublishedAt      any    `json:"published_at"`
		Title            string `json:"title"`
		Type             string `json:"type"`
		URL              string `json:"url"`
		Width            int    `json:"width"`
	} `json:"card"`
	Content            string    `json:"content"`
	CreatedAt          time.Time `json:"created_at"`
	EditedAt           any       `json:"edited_at"`
	Emojis             []any     `json:"emojis"`
	Favourited         bool      `json:"favourited"`
	FavouritesCount    int       `json:"favourites_count"`
	Filtered           []any     `json:"filtered"`
	ID                 string    `json:"id"`
	InReplyToAccountID any       `json:"in_reply_to_account_id"`
	InReplyToID        any       `json:"in_reply_to_id"`
	Language           string    `json:"language"`
	MediaAttachments   []any     `json:"media_attachments"`
	Mentions           []any     `json:"mentions"`
	Muted              bool      `json:"muted"`
	Pinned             bool      `json:"pinned"`
	Poll               any       `json:"poll"`
	Reblog             any       `json:"reblog"`
	Reblogged          bool      `json:"reblogged"`
	ReblogsCount       int       `json:"reblogs_count"`
	RepliesCount       int       `json:"replies_count"`
	Sensitive          bool      `json:"sensitive"`
	SpoilerText        string    `json:"spoiler_text"`
	Tags               []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"tags"`
	URI        string `json:"uri"`
	URL        string `json:"url"`
	Visibility string `json:"visibility"`
}

type InStatus struct {
	ApiUrl    string
	AuthToken string
	Id        string
}

func GetStatus(in InStatus) (Status, error) {
	uri := fmt.Sprintf("%s/api/v1/statuses/%s", in.ApiUrl, in.Id)
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	req.Header.Set("Content-Type", "application/json") // => your content-type

	req.Header.Add("Authorization", in.AuthToken)

	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(response.Body)

	var res Status
	if err := json.Unmarshal(body, &res); err != nil { // Parse []byte to go struct pointer
		log.Fatal("Can not unmarshal JSON")
	}

	return res, nil
}
