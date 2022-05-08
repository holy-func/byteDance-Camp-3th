package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"github.com/holy-func/async"
)

type Query struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
}
type CY struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

func (d *CY) output() {
	fmt.Println("UK:", d.Dictionary.Prons.En, "US:", d.Dictionary.Prons.EnUs)
	for _, item := range d.Dictionary.Explanations {
		fmt.Println(strings.TrimSpace(item))
	}
}
func (y *YouDao) output() {
	for _, item := range strings.Split(y.Data.Entries[0].Explain, ";") {
		fmt.Println(item)
	}
}

type YouDao struct {
	Result struct {
		Msg  string `json:"msg"`
		Code int    `json:"code"`
	} `json:"result"`
	Data struct {
		Entries []struct {
			Explain string `json:"explain"`
			Entry   string `json:"entry"`
		} `json:"entries"`
		Query    string `json:"query"`
		Language string `json:"language"`
		Type     string `json:"type"`
	} `json:"data"`
}

func queryWordYD(word string) *async.GoPromise {
	return async.Promise(func(resolve, reject async.Handler) {
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://dict.youdao.com/suggest?num=5&ver=3.0&doctype=json&cache=false&le=en&q=%s", word), nil)
		if err != nil {
			reject(err)
		}
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
		req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="100", "Microsoft Edge";v="100"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		resp, err := client.Do(req)
		if err != nil {
			reject(err)
		}
		defer resp.Body.Close()
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			reject(err)
		}
		var yd YouDao
		json.Unmarshal(bodyText, &yd)
		resolve(&yd)
	})
}

type queryResponse interface {
	output()
}

func output(ret queryResponse, word string) {
	fmt.Println(word)
	ret.output()
}
func queryWordCY(word string) *async.GoPromise {
	return async.Promise(func(resolve, reject async.Handler) {
		client := &http.Client{}
		query, _ := json.Marshal(Query{TransType: "en2zh", Source: word})
		var data = bytes.NewReader(query)
		req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
		if err != nil {
			reject(err)
		}
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
		req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "cross-site")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
		req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
		req.Header.Set("app-name", "xy")
		req.Header.Set("os-type", "web")
		req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="100", "Microsoft Edge";v="100"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		resp, err := client.Do(req)
		if err != nil {
			reject(err)
		}
		defer resp.Body.Close()
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			reject(err)
		}
		var cy CY
		json.Unmarshal(bodyText, &cy)
		resolve(&cy)
	})
}
func main() {
	if len(os.Args) != 2 {
		fmt.Println("hi~,try with any word")
	} else {
		word := os.Args[1]
		ret, err := async.Any(queryWordCY(word), queryWordYD(word)).UnsafeAwait()
		if err == nil {
			switch v := ret.(type) {
			case *CY:
				output(v, word)
			case *YouDao:
				output(v, word)
			}
		} else {
			fmt.Println("please try again!")
		}
	}
}
