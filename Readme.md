## CDP-Helper
A simple tool to help you with your chromedp project

### Features
- connect to a remote chrome instance
- enable stealth mode to avoid automation detection

### Usage
#### proxy、headless、user-agent options
```go
options := map[string]interface{}{
    "proxy-server": "https://proxy.com",
    "headless": true,
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)",
}

ctx := CreateBrowser(options)
```

### API
```go
// CreateBrowser create a browser instance
func CreateRemoteBrowser(remote string, options map[string]interface{}) context.Context

// CreateBrowserWithContext create a browser instance with context
func CreateBrowserWithContext(parent context.Context, stealth bool, options map[string]interface{}) (context.Context, context.CancelFunc)

// CreateRemoteBrowser create a local browser
func CreateBrowser(remote string, options map[string]interface{}) context.Context

// CreateRemoteBrowserWithContext create a local browser with context
func CreateBrowserWithContext(parent context.Context, remote string, stealth bool, options map[string]interface{}) (context.Context, context.CancelFunc)

func SetCookies(ctx context.Context, cookies []*network.CookieParam) error

func GetCookies(ctx context.Context) (cookies []*network.Cookie, err error)

func CaptureScreenshot(ctx context.Context, filename string) error
```
