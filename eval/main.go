// Command eval is a chromedp example demonstrating how to evaluate javascript
// and retrieve the result.
package main

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	// var res []string
	var location string
	// var title string
	err := chromedp.Run(ctx,
		// chromedp.Navigate(`https://www.google.com/`),
		chromedp.Navigate(`http://127.0.0.1:8000/docs`),
		// chromedp.Evaluate(`Object.keys(window);`, &res),
		// chromedp.Evaluate(`document.title;`, &title),
		chromedp.Evaluate(`document.location.href;`, &location),
	)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("window object keys: %v", res)
	println(location)
}
