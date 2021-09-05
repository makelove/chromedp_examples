
chromedp Request response

chromedp 怎样拦截 Request和response ？

chromedp intercept Request response


- [How to intercept and modify the Response? #409](https://github.com/chromedp/chromedp/issues/409)
```js
chromedp.ListenTarget(ctx, func(event interface{}) {
	switch ev := event.(type) {
        case *fetch.EventRequestPaused:
            go func() {
                c := chromedp.FromContext(ctx)
                ctx := cdp.WithExecutor(ctx, c.Target)
                fmt.Printf("[fetch.EventRequestPaused]: %s\n", ev.Request.URL)
                status := 200
                body := "123456"
                headers := make([]*fetch.HeaderEntry, 0)
                headers = append(headers, &fetch.HeaderEntry{Name: "Connection", Value: "closed"})
                headers = append(headers, &fetch.HeaderEntry{Name: "Content-Length", Value: "6"})
                headers = append(headers, &fetch.HeaderEntry{Name: "Content-Type", Value: "text/javascript"})
                fetch.FulfillRequest(ev.RequestID, int64(status)).
                        WithBody(base64.StdEncoding.EncodeToString([]byte(body))).
                        WithResponseHeaders(headers).
                        WithResponsePhrase("OK").
                        Do(ctx)
            }()
            
	}
})

```

