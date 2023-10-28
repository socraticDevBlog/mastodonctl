package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

const (
	APP_NAME                 string = "mastodonctl"
	APP_DESCRIPTION          string = "commandline client for a Mastodon social media user"
	APP_VERSION              string = "0.2.0"
	MSG_EXPECT_BAD_BEHAVIORS string = "Some Commands might not work properly"
)

type Conf struct {
	ResultsDisplayCount int
	ApiUrl              string
	AuthToken           string
}

// This function returns a Configuration struct
// it retrieves configuration values either by:
//
// 1- environment variables
//   - RESULT_DISPLAY_COUNT is set to an integer indicating how many results should be fetched
//   - API_URL is set to the base URL of a Mastodon server
//   - AUTH_TOKEN is set to the Bearer token required to interact with Mastodon APIs
//
// 2- a configuration file located next the the binary file
//
// ONLY the Authorization Token value is required to properly operate mastodonctl
func FetchConf() Conf {

	conf := Conf{ResultsDisplayCount: 10, ApiUrl: "https://mastodon.social", AuthToken: ""}

	// 1- start by fetching conf from environment variables
	//
	//   return early if Auth Token value is found
	//
	displayCount := os.Getenv("RESULT_DISPLAY_COUNT")
	if len(displayCount) > 0 {
		int_val, err := strconv.ParseInt(displayCount, 6, 12)
		if err == nil {
			conf.ResultsDisplayCount = int(int_val)
		}
	}
	apiUrl := os.Getenv("API_URL")
	if len(apiUrl) > 0 {
		conf.ApiUrl = apiUrl
	}
	authToken := os.Getenv("AUTH_TOKEN")
	if len(authToken) > 0 {
		conf.AuthToken = authToken

		// return early if gotten Token from environment variable
		return conf
	}

	// 2 - get conf values from configuration file if no environment variables are set
	//
	var configFilepath string
	// not ideal but less evil way
	// (binary executed away from its directory
	// doesn't know how to locate conf.json file!)
	if os.Getenv("MASTODONCTL_CONFIG_FILEPATH") != "" {
		configFilepath = os.Getenv("MASTODONCTL_CONFIG_FILEPATH")
	} else {
		// only works when app is executed from within its directory!
		configFilepath = "conf.json"
	}

	configs_file, err := os.Open(configFilepath)
	if os.IsNotExist(err) {
		fmt.Println("Program is unable to fetch Auth Token!!")
		fmt.Println(MSG_EXPECT_BAD_BEHAVIORS)
	} else {
		defer configs_file.Close()
		decoder := json.NewDecoder(configs_file)
		err := decoder.Decode(&conf)
		if err != nil {
			log.Fatal(err)
		}
	}
	if conf.ResultsDisplayCount == 0 {
		conf.ResultsDisplayCount = 10
	}

	return conf
}

func main() {
	conf := FetchConf()

	statusCmd := cli.Command{
		Name:  "status",
		Usage: "retrieve status info by ID",
		Action: func(c *cli.Context) error {
			statusId := c.Args().Get(0)
			if len(statusId) <= 0 {
				fmt.Println("Error: must provide a status ID!")
				return nil
			}

			var authToken string
			if len(conf.AuthToken) > 0 {
				authToken = fmt.Sprintf("Bearer %s", conf.AuthToken)
			} else {
				fmt.Println(MSG_EXPECT_BAD_BEHAVIORS)
			}

			status, err := GetStatus(InStatus{
				Id:        statusId,
				AuthToken: authToken,
				ApiUrl:    conf.ApiUrl,
			})
			if err != nil {
				log.Fatal(err)
			}
			headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgHiBlue).SprintfFunc()

			tbl := table.New("user", "content", "favorited count")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			tbl.AddRow(status.Account.Username, status.Content[0:64], status.FavouritesCount)

			tbl.Print()

			return nil

		},
	}
	accountsCmd := cli.Command{
		Name:  "accounts",
		Usage: "Retrieve Mastodon Accounts infos by username",
		Action: func(c *cli.Context) error {
			var token_val string
			if len(conf.AuthToken) > 0 {
				token_val = fmt.Sprintf("Bearer %s", conf.AuthToken)
			} else {
				fmt.Println(MSG_EXPECT_BAD_BEHAVIORS)
			}

			userName := c.Args().Get(0)
			accounts, err := GetAccounts(InAccounts{
				Username:     userName,
				AuthToken:    token_val,
				ApiUrl:       conf.ApiUrl,
				ResultsCount: conf.ResultsDisplayCount,
			})
			if err != nil {
				log.Fatal(err)
			}

			headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgHiBlue).SprintfFunc()

			tbl := table.New("id", "username", "displayname", "URL", "follower count", "following count")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			for _, r := range accounts {
				tbl.AddRow(r.ID, r.UserName, r.DisplayName, r.URL, r.FollowersCount, r.FollowingCount)
			}

			if len(accounts) == 0 {
				fmt.Println("No results?  Are you sure you have provided a valid APi auth token in conf.json file? or have you NOT provided a username to search?")
			}

			tbl.Print()

			return nil
		},
	}
	hashtagCmd := cli.Command{
		Name:  "hashtag",
		Usage: "Will get latest post informations about a specific hashtag - append searched word after the command",
		Action: func(c *cli.Context) error {
			hashtag := c.Args().Get(0)

			if len(hashtag) <= 0 {
				fmt.Println("Error: must provide a hashtag value to look for!")
				return nil
			}

			results, err := GetHashtag(InTopics{Hashtag: hashtag, ApiUrl: conf.ApiUrl, ResultsCount: conf.ResultsDisplayCount})

			if err != nil {
				log.Fatal(err)
			}

			headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgHiBlue).SprintfFunc()

			tbl := table.New("hashtag", "username", "media url")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			for _, r := range results {
				tbl.AddRow(r.Hashtag, r.Username, r.MediaURL)
			}

			tbl.Print()

			return nil
		},
	}

	app := cli.NewApp()
	app.Name = APP_NAME
	app.Usage = APP_DESCRIPTION
	app.Authors = append(app.Authors, &cli.Author{Name: "socraticDev", Email: "socraticdev@gmail.com"})
	app.Version = APP_VERSION
	app.Commands = append(app.Commands, &statusCmd, &accountsCmd, &hashtagCmd)

	app.Run(os.Args)
}
