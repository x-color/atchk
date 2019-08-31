package atcoder

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/x-color/atchk/internal"
)

var (
	loginURL  = "https://atcoder.jp/login"
	submitURL = "https://atcoder.jp/contests/abc138/submit"
	tasksURL  = ""
)

type language struct {
	id   string
	name string
	mime string
}

func (l *language) String() string {
	return l.name
}

type atcoder struct {
	contests Contests
	config   *internal.Config
}

func NewAtcoder() *atcoder {
	at := &atcoder{}
	at.config = &internal.Config{}
	return at
}

func (at *atcoder) Login(user, password string) error {
	res, err := http.Get(loginURL)
	if err != nil {
		return err
	}

	token, err := getCsrfToken(res)
	if err != nil {
		return err
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
		return err
	}

	if url, err := res.Location(); err != nil || url.Path == "/login" {
		return errors.New("authorization error")
	}

	cookies := make(map[string]string, 2)
	for _, c := range res.Cookies() {
		cookies[c.Name] = c.Value
	}

	at.config.System.Cookies = cookies

	return nil
}

func getCsrfToken(res *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	token, exist := doc.Find("div#main-div > form > input").Attr("value")
	if !exist {
		return "", fmt.Errorf("Can not get token")
	}
	return token, nil
}

func (at *atcoder) IsLoggedIn() bool {
	req, err := http.NewRequest("GET", submitURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range at.config.System.Cookies {
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

func (at *atcoder) Logout() {
	at.config.System.Cookies = nil
}

func (at *atcoder) Submit() bool {
	return true
}

func (at *atcoder) LoadConfig() error {
	return at.config.Read()
}

func (at *atcoder) SaveConfig() error {
	return at.config.Update()
}

func (at *atcoder) SetConfig(key, value string) error {
	return at.config.Set(key, value)
}

func (at *atcoder) GetLangList() ([]*language, error) {
	req, err := http.NewRequest("GET", submitURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range at.config.System.Cookies {
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
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 302 {
		return nil, fmt.Errorf("Faild to access Atcoder\nPlease execute `atchk login`")
	}

	return getLangList(res)
}

func getLangList(res *http.Response) ([]*language, error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	all := true
	langs := make([]*language, 0)
	list := doc.Find("div#select-lang > div > select").First().Find("option")
	list.Each(func(i int, s *goquery.Selection) {
		id, exist1 := s.Attr("value")
		mime, exist2 := s.Attr("data-mime")
		all = all && exist1 && exist2
		langs = append(langs, &language{
			id:   id,
			name: s.Text(),
			mime: mime,
		})
	})

	if !all {
		return nil, fmt.Errorf("Can not get token")
	}

	sort.Slice(langs, func(i, j int) bool {
		return langs[i].name < langs[j].name
	})
	return langs, nil
}

func (at *atcoder) SetLang(lang *language) error {
	if err := at.SetConfig("system.language", lang.name); err != nil {
		return err
	}
	if err := at.SetConfig("system.languageid", lang.id); err != nil {
		return err
	}
	return nil
}

func (at *atcoder) String() string {
	return at.config.String()
}
