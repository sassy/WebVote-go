package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

func getIdList() []string {
	var idList = make([]string, 0)
	urls := []string{
		"http://www.mache.tv/m/twitter_pickup.php?per_screen_name=mikichi0223&e=2016flash",
		"http://www.mache.tv/m/twitter_pickup.php?per_screen_name=mikichi0223&e=2016flash&page=2",
		"http://www.mache.tv/m/twitter_pickup.php?per_screen_name=mikichi0223&e=2016flash&page=3",
		"http://www.mache.tv/m/twitter_pickup.php?per_screen_name=mikichi0223&e=2016flash&page=4",
	}

	for _, url := range urls {
		re, _ := regexp.Compile("^https://twitter.com/mikichi0223/status/([0-9]+)")
		doc, _ := goquery.NewDocument(url)
		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			match := re.FindStringSubmatch(href)
			if len(match) == 2 {
				fmt.Println(match[1])
				idList = append(idList, match[1])
			}
		})
	}
	return idList
}

func vote(id string, proxy string) {
	vote_url := "http://www.mache.tv/m/twitter_pickup/vote.php"
	values := url.Values{}
	values.Add("e_id", "2016flash")
	values.Add("mode", "vote_confirm")
	values.Add("p_id", "14")
	values.Add("tweet_id", id)

	proxyUrl, _ := url.Parse(proxy)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}, Jar: jar}
	res, err := client.PostForm(vote_url, values)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromResponse(res)
	onclick, _ := doc.Find("a").First().Attr("onclick")

	fmt.Println(onclick)
	re, _ := regexp.Compile("vote_hash=([a-z0-9]+)")
	match := re.FindStringSubmatch(onclick)
	fmt.Println(match[1])

	time.Sleep(time.Second)

	values2 := url.Values{}
	values2.Add("e_id", "2016flash")
	values2.Add("mode", "vote_complete")
	values2.Add("p_id", "14")
	values2.Add("tweet_id", id)
	values2.Add("vote_hash", match[1])
	res2, err := client.PostForm(vote_url, values2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res2.Body.Close()
	body, _ := ioutil.ReadAll(res2.Body)
	fmt.Println(string(body))
}

func main() {
	//proxyのURL
	proxys := []string{
		"http://0.0.0.0:0000", //proxy url
	}

	idList := getIdList()
	for _, proxy := range proxys {
		fmt.Println(proxy)
		for _, id := range idList {
			vote(id, proxy)
		}
		time.Sleep(time.Minute * 10)
	}
}
