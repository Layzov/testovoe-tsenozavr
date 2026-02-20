package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func DoRequest(url string, token string, cookies string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization":       {"Bearer " + token},
		"Cookie":              {cookies},
		"User-Agent":          {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36 OPR/127.0.0.0"},
		"Accept":              {"*/*"},
		"X-Application-Platform": {"web"},
		// "Connection":          {"keep-alive"},
		"Accept-Encoding":     {"gzip, deflate, br"},
		"Cache-Control":       {"no-cache"},
		// "Host":                {"api-web.samokat.ru"},
		// "Accept-Language":     {"ru,en-US;q=0.9,en;q=0.8,es;q=0.7,de;q=0.6"},
		"Content-Type":        {"application/json"},
		"Origin": 				{"https://samokat.ru"},
		"Referer": 				{"https://samokat.ru/"},
	}

	client := &http.Client{
	Timeout: 30 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// не следуем автоматически
		return http.ErrUseLastResponse
	},
}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 307 {
		redirectURL := resp.Header.Get("Location")
		req2, _ := http.NewRequest("GET", redirectURL, nil)
		req2.Header = req.Header

		resp2, err := client.Do(req2)
		if err != nil {
			panic(err)
		}
		defer resp2.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("STATUS2:", resp.Status)
		fmt.Println("BODY2", string(bodyBytes))
		fmt.Println("LOCATION2:", resp.Header.Get("Location"))

		return io.ReadAll(resp2.Body)
	}
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("STATUS:", resp.Status)
	fmt.Println("BODY:", string(bodyBytes))
	fmt.Println("LOCATION:", resp.Header.Get("Location"))

	return io.ReadAll(resp.Body)
}