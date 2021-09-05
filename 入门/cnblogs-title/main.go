/**
* @File    :   main.go
* @Time    :   2021/09/04 17:14:39
* @Author  :   GH
* @Desc    :   

实践
打开 https://www.cnblogs.com/ 的首页，然后获取所有文章的标题和链接
*/

package main

import (
 "context"
 "fmt"
 "log"

 "github.com/chromedp/cdproto/cdp"
 "github.com/chromedp/chromedp"
)

func main() {

     opts := append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.DisableGPU,
        chromedp.NoDefaultBrowserCheck,
        chromedp.Flag("headless", false),
        chromedp.Flag("ignore-certificate-errors", true),
        chromedp.WindowSize(1920, 1080),
        chromedp.Flag("blink-settings", "imagesEnabled=false"),
        // chromedp.UserDataDir(dir),
    )
    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()

    ctx, cancel := chromedp.NewContext(
        allocCtx,
        chromedp.WithLogf(log.Printf),
    )
    defer cancel()

    var nodes []*cdp.Node
    err := chromedp.Run(ctx,
          chromedp.Navigate("https://www.cnblogs.com/"),
        // chromedp.Navigate("file:///Users/play/Desktop/cnblogs.com/index.html"),
        chromedp.WaitVisible(`#footer`, chromedp.ByID),
        //   chromedp.Nodes(`a[@class="title"]`, &nodes),
        chromedp.Nodes(`a.post-item-title`, &nodes),
        //   chromedp.Nodes(`a`, &nodes),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("get nodes:", len(nodes))
    // print titles
	for _, node := range nodes {
		fmt.Println(node.Children[0].NodeValue, ":", node.AttributeValue("href"))
        println("------------------------")
	}
}
