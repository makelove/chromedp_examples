/**
* @File    :   main.go
* @Time    :   2021/08/08 20:55:11
* @Author  :   GH
* @Desc    :   
*/

package main

import (
	"context"
	"github.com/chromedp/chromedp/runner"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 本期启动chrome的一些参数相当于执行了 shell 命令
	// C:\Users\mojotv.cn\AppData\Local\Google\Chrome\Application\chrome.exe --no-default-browser-check=true --no-sandbox=true --window-size=1280,900
	// 如果需要更多参数详解chrome浏览器参数的文档
	runnerOps := chromedp.WithRunnerOptions(
		//我的windows10电脑使用chromedp默认配置导致找不到chrome.exe
		//这行代码可以注释掉,如果找不到自己的chrome.exe 请像我一样制定chrome.exe路径
		//一下配置都不是必选的
		//更多参数详解文档 https://blog.csdn.net/wanwuguicang/article/details/79751571
		runner.Path(`C:\Users\mojotv.cn\AppData\Local\Google\Chrome\Application\chrome.exe`),
		//启动chrome的时候不检查默认浏览器
		runner.Flag("no-default-browser-check", true),
		//启动chrome 不适用沙盒, 性能优先
		runner.Flag("no-sandbox", true),
		//设置浏览器窗口尺寸,
		runner.WindowSize(1280, 1024),
		//设置浏览器的userage
		runner.UserAgent(`Mozilla/5.0 (iPhone; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.25 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1`),
	)
	//在普通模式的情况下启动chrome程序,并且建立共代码和chrome程序的之间的连接(https://127.0.0.1:9222)
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf), runnerOps)
	if err != nil {
		log.Fatal(err)
	}

	var siteHref, title, iFrameCode string
	err = c.Run(ctxt, visitMojoTvDotCn("https://mojotv.cn/2018/12/10/how-to-create-a-https-proxy-serice-in-100-lines-of-code.html", &siteHref, &title, &iFrameCode))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("`%s` (%s),html:::%s", title, siteHref, iFrameCode)

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func visitMojoTvDotCn(url string, elementHref, pageTitle, iFrameHtml *string) chromedp.Tasks {
	//临时放图片buf
	var buf []byte
	return chromedp.Tasks{
		//跳转到页面
		chromedp.Navigate(url),
		//chromedp.Sleep(2 * time.Second),
		//等待博客正文显示
		chromedp.WaitVisible(`#post`, chromedp.ByQuery),
		//滑动页面到google adsense 广告
		chromedp.ScrollIntoView(`ins`, chromedp.ByQuery),
		chromedp.Screenshot(`#post`, &buf, chromedp.ByQuery, chromedp.NodeVisible),
		//等待2s
		chromedp.Sleep(2 * time.Second),
        //截图到文件
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			//保存图片到mojotv_local.png
			return ioutil.WriteFile("mojotv_local.png", buf, 0644)
		}),
		//滑动页面到#copyright
		chromedp.ScrollIntoView(`#copyright`, chromedp.ByID),
		//等待mojotv google广告展示出来
		chromedp.WaitVisible(`#post__title`, chromedp.ByID),
		chromedp.Sleep(2 * time.Second),

		//获取我的google adsense 广告代码
		chromedp.InnerHTML(`#post__title`, iFrameHtml, chromedp.ByID),
		//跳转到我的bilibili网站
		chromedp.Sleep(5 * time.Second),

		chromedp.Click("#copyright > a:nth-child(3)", chromedp.NodeVisible),
		//等待则个页面显现出来
		chromedp.WaitVisible(`#page`, chromedp.ByQuery),
		//在chrome浏览器页面里执行javascript
		chromedp.Evaluate(`document.title`, pageTitle),
		chromedp.Screenshot(`#page`, &buf, chromedp.ByQuery, chromedp.NodeVisible),
		chromedp.Sleep(5 * time.Second),

		//截取bili网页图片
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("bili_local.png", buf, 0644)
		}),
		//获取bilibili网页的标题
		chromedp.JavascriptAttribute(`a`, "href", elementHref, chromedp.ByQuery),
	}
}
