package atcoder

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/x-color/atchk/internal"
	"github.com/x-color/atchk/internal/atcoder/contest"
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
	cache  *internal.Cache
	client *http.Client
}

func NewAtcoder() *atcoder {
	at := &atcoder{
		cache:  &internal.Cache{},
		client: &http.Client{},
	}
	at.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return at
}

func (at *atcoder) SetCache(cache *internal.Cache) {
	at.cache = cache
}

func (at *atcoder) Login(user, password string) error {
	res, err := http.Get(getLoginURL())
	if err != nil {
		return err
	}

	cookies := make(map[string]string, 2)
	for _, c := range res.Cookies() {
		cookies[c.Name] = c.Value
	}
	at.cache.Cookies = cookies

	token, err := getCsrfToken(res)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Add("username", user)
	values.Add("password", password)
	values.Add("csrf_token", token)

	req, err := at.newRequest("POST", getLoginURL(), strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	res, err = at.client.Do(req)
	if err != nil {
		return err
	}

	if url, err := res.Location(); err != nil || url.Path == "/login" {
		return errors.New("authorization error")
	}

	for _, c := range res.Cookies() {
		at.cache.Cookies[c.Name] = c.Value
	}

	return nil
}

func (at *atcoder) IsLoggedIn() bool {
	req, err := at.newRequest("GET", getSubmitURL("abc001"), nil)
	if err != nil {
		return false
	}

	res, err := at.client.Do(req)
	if err != nil || res.StatusCode == 302 {
		return false
	}

	return true
}

func (at *atcoder) Logout() {
	at.cache.Cookies = nil
}

func (at *atcoder) getSamplesFromCache(contest, task string) contest.Samples {
	if contest == at.cache.Contest {
		if task == at.cache.Task {
			return at.cache.Samples
		}
	}
	return nil
}

func (at *atcoder) getSamplesFromWeb(contest, task string) (contest.Samples, error) {
	req, err := at.newRequest("GET", getTasksURL(contest), nil)
	if err != nil {
		return nil, err
	}

	res, err := at.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 302 {
		return nil, fmt.Errorf("Can not access to task page\nPlease execute `atchk login`")
	}

	taskURL, err := getTaskURL(res, task)
	if err != nil {
		return nil, err
	}

	req, err = at.newRequest("GET", taskURL, nil)
	if err != nil {
		return nil, err
	}

	res, err = at.client.Do(req)
	if err != nil {
		return nil, err
	}

	return getSamples(res)
}

func (at *atcoder) getSamples(contest, task string) (contest.Samples, error) {
	samples := at.getSamplesFromCache(contest, task)
	if samples != nil {
		return samples, nil
	}
	samples, err := at.getSamplesFromWeb(contest, task)
	if err != nil {
		return nil, err
	}
	at.cache.Contest = contest
	at.cache.Task = task
	at.cache.Samples = samples
	// This error is no effection for main process
	_ = at.SaveCache()
	return samples, nil
}

func (at *atcoder) Samples(contest, task string) (int, error) {
	samples, err := at.getSamples(contest, task)
	if err != nil {
		return 0, err
	}
	return len(samples), nil
}

func (at *atcoder) Test(i int, cmdList []string) (string, error) {
	return at.cache.Samples[i].Test(cmdList)
}

func (at *atcoder) Submit() bool {
	return true
}

func (at *atcoder) SaveCache() error {
	return at.cache.Update()
}

func (at *atcoder) GetLangList() ([]*language, error) {
	req, err := at.newRequest("GET", getSubmitURL("abc001"), nil)
	if err != nil {
		return nil, err
	}

	res, err := at.client.Do(req)
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
		return nil, fmt.Errorf("Can not get lang list")
	}

	sort.Slice(langs, func(i, j int) bool {
		return langs[i].name < langs[j].name
	})
	return langs, nil
}

func (at *atcoder) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range at.cache.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  strings.ToUpper(k),
			Value: v,
		})
	}
	return req, nil
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

func getSamples(res *http.Response) (contest.Samples, error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	inputs, outputs := make([]string, 0), make([]string, 0)
	list := doc.Find(".lang-ja").Find("pre")
	if list.Length() < 1 {
		list = doc.Find("#task-statement").Last().Find("section > pre")
	}
	list.Slice(1, goquery.ToEnd).Each(func(i int, s *goquery.Selection) {
		if i%2 == 0 {
			inputs = append(inputs, s.Text())
		} else {
			outputs = append(outputs, s.Text())
		}
	})

	samples := make(contest.Samples, 0)
	for i := range inputs {
		samples = append(samples, &contest.Sample{
			ID:     i + 1,
			Input:  strings.TrimSpace(inputs[i]),
			Output: strings.TrimSpace(outputs[i]),
		})
	}

	return samples, nil
}

func getLoginURL() string {
	return "https://atcoder.jp/login"
}

func getTasksURL(contest string) string {
	return "https://atcoder.jp/contests/" + contest + "/tasks"
}

func getSubmitURL(contest string) string {
	return "https://atcoder.jp/contests/" + contest + "/submit"
}

func getTaskURL(res *http.Response, task string) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	var url string
	tokens := doc.Find("div#main-container > div > div > div > table > tbody > tr > td.no-break > a")
	tokens.Each(func(i int, s *goquery.Selection) {
		if strings.ToUpper(task) == s.Text() {
			if href, exist := s.Attr("href"); exist {
				url = href
				return
			}
		}
	})

	if url != "" {
		return "https://atcoder.jp" + url, nil
	}
	return "", fmt.Errorf("Can not get task url")
}
