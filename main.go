package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Specify when build compile
var (
	token    string
	endpoint string
	username string
	hookurl  string
)

type Repo struct {
	Name string `json:"name"`
}

type OpenedPullReq struct {
	HtmlUrl string `json:"html_url"`
	Title   string `json:"title"`
	User    User   `json:"user"`
	Repo    Repo   `json:"repo"`
}

type User struct {
	LoginName string `json:"login"`
}

type WebHookPostPayload struct {
	Text      string `json:"text,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
}

func PostMessage(payload *WebHookPostPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(hookurl, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(t))
	}

	return nil
}

func Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint+path, nil)
	if err != nil {
		return nil, err
	}
	setToken(req)
	c := http.Client{}
	return c.Do(req)
}

func GetRepos(path string) ([]Repo, error) {
	res, err := Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	repos := []Repo{}
	err = json.NewDecoder(res.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}
	return repos, nil

}

func GetOpenedPullReq(repo string) ([]OpenedPullReq, error) {
	res, err := Get("/repos/" + username + "/" + repo + "/pulls")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	prs := []OpenedPullReq{}
	err = json.NewDecoder(res.Body).Decode(&prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func setToken(req *http.Request) {
	req.Header.Add("Authorization", "token "+token)
}

func main() {
	repos, err := GetRepos("/users/" + username + "/repos?per_page=300")
	if err != nil {
		log.Fatal(err)
	}

	allPR := []OpenedPullReq{}
	for _, repo := range repos {
		prs, err := GetOpenedPullReq(repo.Name)
		if err != nil {
			log.Fatal(err)
		}
		allPR = append(allPR, prs...)
	}

	msg := ""
	for _, pr := range allPR {

		msg = msg + fmt.Sprintf("[\n"+
			pr.Repo.Name+"\n"+
			"url:%s\n"+
			"title:%s\n"+
			"user:%s\n"+
			"]\n",
			pr.HtmlUrl,
			pr.Title,
			pr.User,
		)
	}

	var payload = &WebHookPostPayload{
		Channel:   "#gazou-dev",
		Username:  "PullReqBot",
		IconEmoji: ":ghost:",
		Text:      msg,
	}

	err = PostMessage(payload)
	if err != nil {
		log.Fatal(err)
	}
}
