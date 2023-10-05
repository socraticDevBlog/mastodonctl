package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"gopkg.in/urfave/cli.v1"
)

type Conf struct {
	ResultsDisplayCount int
	ApiUrl              string
	AuthToken           string
}

func (conf *Conf) DefaultConf() {
	if conf.ResultsDisplayCount == 0 {
		conf.ResultsDisplayCount = 10
	}

	if conf.ApiUrl == "" {
		conf.ApiUrl = "https://mastodon.social"
	}
}

// todo: move to Utils - PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func main() {
	const (
		APP_NAME                 string = "mastodonctl"
		APP_DESCRIPTION          string = "commandline client for a Mastodon social media user"
		APP_VERSION              string = "0.1.0"
		MSG_EXPECT_BAD_BEHAVIORS string = "Some Commands might not work properly"
	)

	conf := Conf{}

	var configFilepath string
	// not ideal but less evil way
	// (binary executed away from its directory
	// doesn't know how to locate conf.json file!)
	if os.Getenv("CONFIG_FILEPATH") != "" {
		configFilepath = os.Getenv("CONFIG_FILEPATH")
	} else {
		// only works when app is executed from within its directory!
		configFilepath = "conf.json"
	}

	configs_file, err := os.Open(configFilepath)
	if os.IsNotExist(err) {
		fmt.Println("Program is unable to open configuration file: conf.json ...")
		fmt.Println(MSG_EXPECT_BAD_BEHAVIORS)
		conf.DefaultConf()
	} else {
		defer configs_file.Close()
		decoder := json.NewDecoder(configs_file)
		err := decoder.Decode(&conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	app := cli.NewApp()
	app.Name = APP_NAME
	app.Usage = APP_DESCRIPTION

	app.Authors = append(app.Authors, cli.Author{Name: "socraticDev", Email: "socraticdev@gmail.com"})
	app.Version = APP_VERSION
	app.Commands = []cli.Command{
		{
			Name:  "accounts",
			Usage: "Retrieve Mastodon Accounts infos by username",
			Action: func(c *cli.Context) error {
				userName := c.Args().Get(0)

				var token_val string
				if len(conf.AuthToken) > 0 {
					token_val = fmt.Sprintf("Bearer %s", conf.AuthToken)
				} else {
					fmt.Println(MSG_EXPECT_BAD_BEHAVIORS)
				}

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
		}, {
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
		},
	}

	app.Run(os.Args)
}
