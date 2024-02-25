package login

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/davidroman0O/gogog/data"
	"github.com/davidroman0O/gogog/types"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Allowing you to log into your account",
		Long:  `It will opens a browser, you will have to login and it will find the cookies then quit`,
		Run: func(cmd *cobra.Command, args []string) {
			// if data.Has[types.GogAuthenticationChrome]() {
			// 	fmt.Println("You already have your authentication information")
			// 	return
			// }

			if err := data.Ping(); err != nil {
				slog.Error("Couldn't reach the api", err)
				panic(err)
			}

			log.Println("Please login then close the window")

			gogClient := data.NewGogClient()

			dir, err := os.MkdirTemp("", "chromedp-example")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir)

			opts := append(
				chromedp.DefaultExecAllocatorOptions[3:],
				chromedp.NoFirstRun,
				chromedp.NoDefaultBrowserCheck,
				visible,
			)

			allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

			ctx, cancel := chromedp.NewContext(allocCtx)

			defer func() {
				cancel()
				if r := recover(); r != nil {
					panic(r)
				}
			}()

			if err = chromedp.Run(ctx, chromedp.Navigate("https://www.gog.com/")); err != nil {
				log.Fatal(err)
			}

			var obj []*types.Cookie
			signInFound := false
			for !signInFound {
				// "Sign In" is no longer found, assume user has signed in
				// Retrieve the cookies
				var cookies []*network.Cookie
				if err := chromedp.Run(ctx,
					chromedp.ActionFunc(func(ctx context.Context) error {
						var err error
						cookies, err = network.GetCookies().Do(ctx)
						return err
					}),
				); err != nil {
					log.Fatal(err)
					time.Sleep(1 * time.Second)
					continue
				}

				var bytesCookies []byte

				bytesCookies, err = json.Marshal(cookies)
				if err != nil {
					log.Fatal(err)
					time.Sleep(1 * time.Second)
					continue
				}

				err = json.Unmarshal(bytesCookies, &obj)
				if err != nil {
					log.Fatal(err)
					time.Sleep(1 * time.Second)
					continue
				}

				if err := data.SetCookies(gogClient.Client, obj, types.Hostname); err != nil {
					log.Println("failed set cookies", err)
					time.Sleep(1 * time.Second)
					continue
				}

				user, ok, err := data.FetchUser(gogClient.Client, types.Hostname)
				if err != nil {
					log.Println("failed check cookies", err)
					time.Sleep(1 * time.Second)
					continue
				}

				if !ok {
					log.Println("1s - please login")
					time.Sleep(1 * time.Second)
					continue
				}

				auth := types.GogAuthenticationChrome{
					Cookies: obj,
					User:    user,
				}

				if err := data.PostAccount(auth); err != nil {
					log.Fatalf("couldn't save your authentication data %v", err)
				}

				signInFound = true
			}

			if signInFound {
				log.Println("You can now use `gogog` directly")
			}
		},
	}
}

func visible(a *chromedp.ExecAllocator) {
	chromedp.Flag("headless", false)(a)
	chromedp.Flag("disable-gpu", false)(a)
	chromedp.Flag("hide-scrollbars", false)(a)
	chromedp.Flag("mute-audio", false)(a)
	chromedp.Flag("disable-background-networking", false)
	chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess")
	chromedp.Flag("disable-background-timer-throttling", false)
	chromedp.Flag("disable-backgrounding-occluded-windows", false)
	chromedp.Flag("disable-breakpad", false)
	chromedp.Flag("disable-client-side-phishing-detection", false)
	chromedp.Flag("disable-default-apps", false)
	chromedp.Flag("disable-dev-shm-usage", false)
	chromedp.Flag("disable-extensions", false)
	chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees")
	chromedp.Flag("disable-hang-monitor", false)
	chromedp.Flag("disable-ipc-flooding-protection", false)
	chromedp.Flag("disable-popup-blocking", false)
	chromedp.Flag("disable-prompt-on-repost", false)
	chromedp.Flag("disable-renderer-backgrounding", false)
	chromedp.Flag("disable-sync", false)
	chromedp.Flag("force-color-profile", "srgb")
	chromedp.Flag("metrics-recording-only", false)
	chromedp.Flag("safebrowsing-disable-auto-update", false)
	chromedp.Flag("enable-automation", false)
	chromedp.Flag("password-store", "basic")
	chromedp.Flag("use-mock-keychain", false)
}
