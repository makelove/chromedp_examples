
## chromedp 入门

- 参考
    - [chromedp入门](https://www.codercto.com/a/114358.html)
        -  [原文](https://mp.weixin.qq.com/s?__biz=MjM5MDM5NTY0OA==&mid=2648697048&idx=1&sn=3986d965a6e819625f9cfdc614f9dc7e)

    - [Go使用chromedp库操作headless-chrome爬取"JS画出来的"网站](https://zhuanlan.zhihu.com/p/139261122)
        - 可以看到chromedp是严重依赖于go的context包的，如果不熟悉context的使用最好先去研究一下。
        - 现在网上很多例子都是只爬取一个页面的例子，然而实际中我们经常需要按一定规则爬取整个网站，这时候就涉及到Chrome实例的复用问题，总不能跟其他博客说的那样每爬取一个网页都要销毁、打开Chrome吧，效率太低了。


## 显示Chrome窗口
```js
chromedp.Flag("headless", false),
context.Background()//后台运行
```

## 指定浏览器路径？默认是Chrome？
- 自动找到
```js
//TODO
```



## 初始化chromedp的上下文，后续这个页面都使用这个上下文进行操作
```js
chromedp.NewContext()
```
## 运行一个chrome的一系列操作
```js
chromedp.Run() 
```


## 窗口大小
```js
//TODO "50,400"不正确？
chromedp.Flag("window-size", "50,400"),
```

## 设置UserAgent
```js
chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
```

## 禁用图片加载
```js
 chromedp.Flag("blink-settings", "imagesEnabled=false"),
```



## 数据文件夹，保存用户数据
```js
chromedp.UserDataDir(dir),
```

## 使用代理，认证username+password
```js
//TODO

```


## 设置超时时间
//TODO
在实际应用中连续爬取多个网页，只需要每次在chromeCtx的基础上，派生一个带有超时功能的子Context即可：
```js
timeoutCtx, cancel := context.WithTimeout(chromeContext, 20 * time.Second)
```
这样就不会销毁Chrome进程再重新创建一遍了。

## 导航
```js
 chromedp.Navigate(`https://space.bilibili.com/180948619/dynamic`),
```

##  等候某个元素可见，再继续执行。
```js
chromedp.WaitVisible(`body`, chromedp.BySearch),
```

## 监听 Request和Response
```js
 chromedp.ListenTarget(ctx, func(ev interface{}) {
```

## 使用代理，认证username+password
```js
//TODO
chromedp.UserDataDir(dir),
```



//TODO
获取page
第一个页面

## 截图,元素截图

- 全屏 ?
```js
chromedp.Screenshot(sel, res, chromedp.NodeVisible),
```



chromedp.Tasks{} 一系列Action组成的任务

## 模拟手机设备 emulate
- Run task

```js
chromedp.Emulate(device.IPhone7landscape),
```

## 执行JS代码
- Run task

```js
chromedp.Evaluate(`document.location.href;`, &location)
```

## 提取HTML 获取元素的outer html
```js
chromedp.OuterHTML() 
chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
```

