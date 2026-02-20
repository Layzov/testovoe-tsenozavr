package cookies

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func GetCookies() (string, error){
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var cookies []*network.Cookie

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate("https://samokat.ru"),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			for {
				cookies, err = network.GetCookies().Do(ctx)
				if err != nil {
					break
				}
				if len(cookies) > 0 {
					break
				}
			}
			return err
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	cookieHeader := buildCookieHeader(cookies)
	return cookieHeader, nil
}

func buildCookieHeader(cookies []*network.Cookie) string {
	var parts []string
	for _, c := range cookies {
		if c.Name == "sberid_auto_error_pause" {
			continue
		}
		parts = append(parts, c.Name+"="+c.Value)
	}
	return strings.Join(parts, "; ")
}