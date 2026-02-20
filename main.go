package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testovoe/internal/requests/categories"
	"testovoe/session"
	"time"

	"github.com/chromedp/chromedp"
)

type Showcase struct {
	ShowcaseID string `json:"showcaseId"`
	StoreID    string `json:"storeId"`
	Type       string `json:"type"`
	Title      string `json:"title"`
}

func main() {
	// получаем сессию (token, cookies) и живой browserCtx + cancel
	t, c, browserCtx, cancelAll, err := session.GetSession(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer cancelAll()

	// теперь можно использовать browserCtx без context canceled:

	fmt.Println("TOKEN:", t)
	fmt.Println("COOKIES:", c)

	// ------ теперь получаем showcase через браузер (fetch внутри browserCtx) ------
	lat := "56.4601736"
	lon := "84.9616012"

	var jsonResult string

	fetchJS := fmt.Sprintf(`(async () => {
		let r = await fetch("https://api-web.samokat.ru/showcases/list?lat=%s&lon=%s", {
			method: "GET",
			credentials: "include",
			headers: {
				"Authorization": "Bearer %s",
				"X-Application-Platform": "web",
				"Accept": "application/json"
			}
		});

		if (r.status === 307 || r.redirected) {
			r = await fetch(r.url, {
				credentials: "include",
				headers: {
					"Authorization": "Bearer %s",
					"X-Application-Platform": "web",
					"Accept": "application/json"
				}
			});
		}

		if (!r.ok) {
			return JSON.stringify({
				error:true,
				status:r.status,
				text:await r.text()
			});
		}

		return JSON.stringify(await r.json());
	})()`, lat, lon, t, t)

	// используем отдельный короткий timeout для этого запроса
	ctx, cancel := context.
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Evaluate(fetchJS, &jsonResult)); err != nil {
		log.Fatal("fetch showcases error:", err)
	}

	// парсим результат
	type showcaseResp []struct {
		ShowcaseID string `json:"showcaseId"`
	}

	var parsed showcaseResp
	if err := json.Unmarshal([]byte(jsonResult), &parsed); err != nil {
		log.Fatal("unmarshal showcases error:", err, "\nbody:", jsonResult)
	}

	if len(parsed) == 0 {
		log.Fatal("no showcases: " + jsonResult)
	}

	showcaseID := parsed[0].ShowcaseID
	fmt.Println("SHOWCASE:", showcaseID)

	// дальше — используем твой categories модуль, передаём token/cookies как раньше
	cats, err := categories.GetCategories(showcaseID, t, c)
	if err != nil {
		log.Fatal("GetCategories error:", err)
	}

	for _, cat := range cats {
		fmt.Println(cat.Name, cat.Slug, cat.ProductsCount)
	}

	fmt.Println("DONE")
}
