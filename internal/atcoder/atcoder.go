package atcoder

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Atcoder struct {
	answerFile  string
	taskName    string
	contestName string
	contests    Contests
}

var (
	loginURL  = "https://atcoder.jp/login"
	submitURL = "https://atcoder.jp/contests/abc138/submit"
	tasksURL  = ""
)

func (at *Atcoder) Login(user, password string) (map[string]string, error) {
	res, err := http.Get(loginURL)
	if err != nil {
		return nil, err
	}

	token, err := getCsrfToken(res)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("username", user)
	values.Add("password", password)
	values.Add("csrf_token", token)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for _, c := range res.Cookies() {
		req.AddCookie(c)
	}

	client := http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if url, err := res.Location(); err != nil || url.Path == "/login" {
		return nil, errors.New("authorization error")
	}

	cookies := make(map[string]string, 2)
	for _, c := range res.Cookies() {
		cookies[c.Name] = c.Value
	}

	return cookies, nil
}

func getCsrfToken(res *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	token, exist := doc.Find("div#main-div > form > input").Attr("value")
	if !exist {
		return "", errors.New("can not get token")
	}
	return token, nil
}

func (at *Atcoder) IsLoggedIn(cookies map[string]string) bool {
	req, err := http.NewRequest("GET", submitURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  strings.ToUpper(k),
			Value: v,
		})
	}

	client := http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := client.Do(req)
	if err != nil || res.StatusCode == 302 {
		return false
	}

	return true
}

func (at *Atcoder) Submit() bool {
	return true
}
