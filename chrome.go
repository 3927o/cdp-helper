package cdpHelper

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/storage"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

var DefaultExecAllocatorOptions = map[string]interface{}{
	"force-color-profile":                    "srgb",
	"safebrowsing-disable-auto-update":       true,
	"no-first-run":                           true,
	"hide-scrollbars":                        true,
	"mute-audio":                             true,
	"enable-features":                        "NetworkService,NetworkServiceInProcess",
	"disable-backgrounding-occluded-windows": true,
	"enable-automation":                      false,
	"disable-prompt-on-repost":               true,
	"disable-sync":                           true,
	"password-store":                         "basic",
	"disable-dev-shm-usage":                  true,
	"disable-features":                       "site-per-process,Translate,BlinkGenPropertyTrees",
	"disable-background-networking":          true,
	"disable-breakpad":                       true,
	"disable-renderer-backgrounding":         true,
	"disable-default-apps":                   false,
	"disable-hang-monitor":                   true,
	"no-default-browser-check":               true,
	"disable-ipc-flooding-protection":        true,
	"disable-popup-blocking":                 true,
	"metrics-recording-only":                 true,
	"use-mock-keychain":                      true,
	"headless":                               true,
	"disable-background-timer-throttling":    true,
	"disable-client-side-phishing-detection": true,
	"disable-extensions":                     true,
}

var DefaultPCDevice chromedp.Device = device.Info{
	Name:      "PC",
	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36",
	Width:     1920,
	Height:    1080,
	Scale:     1.000,
	Landscape: false,
	Mobile:    false,
	Touch:     false,
}

// CreateRemoteBrowserWithContext
func CreateRemoteBrowserWithContext(parent context.Context, remote string, stealth bool, options map[string]interface{}) (context.Context, context.CancelFunc) {
	remoteUrl := constructRemoteUrl(remote, options)

	if stealth {
		remoteUrl = remoteUrl + "&stealth"
	}

	c, _ := chromedp.NewRemoteAllocator(parent, remoteUrl, func(ra *chromedp.RemoteAllocator) {
		// prevent modifying the url
		chromedp.NoModifyURL(ra)
	})

	ctx, cancel := chromedp.NewContext(c)
	return ctx, cancel
}

// CreateRemoteBrowser create a local browser
// not return cancel func because you can close it in a more elegant way: chromedp.Cancel()
func CreateRemoteBrowser(remote string, options map[string]interface{}) context.Context {
	ctx, _ := CreateRemoteBrowserWithContext(context.Background(), remote, true, options)
	return ctx
}

func CreateBrowserWithContext(parent context.Context, stealth bool, options map[string]interface{}) (context.Context, context.CancelFunc) {
	execAllocatorOptions := convertExecAllocatorOption(options)

	if stealth {
		execAllocatorOptions = append(execAllocatorOptions,
			chromedp.Flag("disable-blink-features", "AutomationControlled"),
		)
	}

	c, _ := chromedp.NewExecAllocator(parent, execAllocatorOptions...)

	ctx, cancel := chromedp.NewContext(c)

	return ctx, cancel
}

// CreateBrowser create a local browser
// not return cancel func because you can close it in a more elegant way: chromedp.Cancel()
func CreateBrowser(options map[string]interface{}) context.Context {
	ctx, _ := CreateBrowserWithContext(context.Background(), true, options)
	return ctx
}

func SetCookies(ctx context.Context, cookies []*network.CookieParam) error {
	err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return storage.SetCookies(cookies).Do(ctx)
	}))
	return err
}

func GetCookies(ctx context.Context) (cookies []*network.Cookie, err error) {
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		cookies, err = storage.GetCookies().Do(ctx)
		return err
	}))
	return
}

func CaptureScreenshot(ctx context.Context, filename string) error {
	if ctx == nil {
		return errors.New("context is nil")
	}
	c := chromedp.FromContext(ctx)
	if c.Target == nil {
		return errors.New("given context is a invalid chromedp context")
	}
	screenBytes, err := page.CaptureScreenshot().Do(cdp.WithExecutor(ctx, c.Target))
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, screenBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func convertExecAllocatorOption(options map[string]interface{}) []chromedp.ExecAllocatorOption {
	var cdpOptions []chromedp.ExecAllocatorOption
	for k, v := range options {
		cdpOptions = append(cdpOptions, chromedp.Flag(k, v))
	}
	return cdpOptions
}

func constructRemoteUrl(remote string, options map[string]interface{}) string {
	remoteUrl, err := url.Parse(remote)
	if err != nil {
		panic(err)
	}
	urlQuery := remoteUrl.Query()
	for k, v := range options {
		urlQuery.Add(k, fmt.Sprintf("%v", v))
	}
	remoteUrl.RawQuery = urlQuery.Encode()
	return remoteUrl.String()
}
