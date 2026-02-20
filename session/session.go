package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var BrowserCtx context.Context

func GetSession(parent context.Context) (string, string, context.Context, context.CancelFunc, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		// chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-infobars", true),
		// chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(parent, opts...)

	browserCtx, browserCancel := chromedp.NewContext(allocCtx)

	// экспортируем живой browserCtx
	BrowserCtx = browserCtx

	cancelAll := func() {
		browserCancel()
		allocCancel()
	}

	workCtx, workCancel := context.WithTimeout(browserCtx, 90*time.Second)
	defer workCancel()

	var token string
	var cookies string

	tokenCh := make(chan struct{})


	err := chromedp.Run(workCtx,
		network.Enable(),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// сохраняем Actx для использования в других частях кода
			browserCtx = ctx

			// listener остаётся внутри ActionFunc, как у тебя было,
			// но внутри горутины используем browserCtx (не ctx),
			// чтобы чтение тела не зависело от workCtx.
			chromedp.ListenTarget(browserCtx, func(ev interface{}) {
				if resp, ok := ev.(*network.EventResponseReceived); ok {
					if strings.Contains(resp.Response.URL, "/api/auth/session") {

						go func(id network.RequestID) {
							// ВАЖНО: используем browserCtx вместо ctx
							body, err := network.GetResponseBody(id).Do(ctx)
							if err != nil {
								return
							}

							var parsed map[string]interface{}
							if json.Unmarshal(body, &parsed) == nil {

								if t, ok := parsed["accessToken"].(string); ok {
									token = t
								}

								if data, ok := parsed["data"].(map[string]interface{}); ok {
									if t, ok := data["accessToken"].(string); ok {
										token = t
									}
								}

								if token != "" {
									select {
									case <-tokenCh:
									default:
										close(tokenCh)
									}
								}
							}
						}(resp.RequestID)
					}
				}
			})

			return nil
		}),

		chromedp.Navigate("https://samokat.ru"),

		//антибот
		chromedp.Sleep(12*time.Second),
		chromedp.MouseClickXY(500, 500),

		//форс-запрос токена
		chromedp.Evaluate(`(async () => {
			const r = await fetch("/api/auth/session", {credentials:"include"});
			return await r.text();
		})()`, nil),

		// ждём ответ
		chromedp.ActionFunc(func(ctx context.Context) error {
			select {
			case <-tokenCh:
				return nil
			case <-time.After(20 * time.Second):
				return fmt.Errorf("token not received")
			}
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			c, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}

			var parts []string
			for _, ck := range c {
				if ck.Name == "sberid_auto_error_pause"{
					continue
				}
				parts = append(parts, ck.Name+"="+ck.Value)
			}
			cookies = strings.Join(parts, "; ")
			return nil
		}),
	)

	if err != nil {
		// если ошибка — закрываем созданные контексты
		cancelAll()
		return "", "", nil, nil, err
	}

	// УСПЕШНО: возвращаем token, cookies, сам browserCtx и функцию cancelAll, caller должен вызвать cancelAll()
	return token, cookies, browserCtx, cancelAll, nil
}
