package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/Jeffail/tunny"
	"github.com/davidroman0O/gogog/data"
	"github.com/davidroman0O/gogog/db"
	"github.com/davidroman0O/gogog/handlers"
	"github.com/davidroman0O/gogog/types"
	"github.com/davidroman0O/gogog/web"
	"github.com/k0kubun/pp/v3"
	"github.com/spf13/cobra"
)

func MiddlewareAccount(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total, err := data.CountAccounts()
		if err != nil {
			fmt.Println("err", err)
			http.Redirect(w, r, fmt.Sprintf("/error?error=%v", err), http.StatusTemporaryRedirect)
			return
		}
		if total == 0 {
			fmt.Println("redirecting")
			http.Redirect(w, r, "/onboarding", http.StatusTemporaryRedirect)
			return
		}
		handler(w, r)
	}
}

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "server",
		Long:  `.`,
		Run: func(cmd *cobra.Command, args []string) {

			if err := db.Initialize(
				db.WithDBConfig(),
			); err != nil {
				panic(err)
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)

			if err := data.Initialize("./tmp"); err != nil {
				panic(err)
			}

			numCPUs := runtime.NumCPU()

			tmpStateDownload := false

			poolGameCache := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
				var result []byte

				switch msg := payload.(type) {
				case []types.Account:
					gogClient := data.NewGogClient()

					for _, account := range msg {

						if err := data.SetCookies(gogClient.Client, account.Cookies, types.Hostname); err != nil {
							return err
						}

						params, err := types.NewSearchParams()

						if err != nil {
							return err
						}

						products, err := data.Search(gogClient.Client, types.Hostname, *params)
						if err != nil {
							return err
						}

						for _, v := range products {
							data.GamesDB().Save(&v)
							slog.Info("save %v", v.Title)
						}
					}

					tmpStateDownload = false

				default:
					slog.Info("it's not an account array")
				}

				return result
			})

			defer poolGameCache.Close()

			server := &http.Server{
				Addr:         ":8000",
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			mux := http.NewServeMux()

			mux.HandleFunc("GET /", MiddlewareAccount(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/" {
					http.Redirect(w, r, "/404", http.StatusFound)
					return
				}
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				if err := web.Index().Render(context.Background(), w); err != nil {
					panic(err)
				}
			}))

			mux.HandleFunc("GET /onboarding", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.Header().Add("Content-Type", "text/html")
				if err := web.PageOnboarding().Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /error", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.Header().Add("Content-Type", "text/html")
				err := r.URL.Query().Get("error")
				if err := web.PageError(fmt.Errorf(err)).Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /accounts", MiddlewareAccount(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.Header().Add("Content-Type", "text/html")
				accounts, err := data.GetAccounts()
				if err != nil {
					http.Redirect(w, r, "/404", http.StatusFound)
					return
				}
				if err := web.PageAccounts(accounts).Render(context.Background(), w); err != nil {
					panic(err)
				}
			}))

			mux.HandleFunc("GET /games/table", MiddlewareAccount(handlers.GetGamesTablefunc))

			mux.HandleFunc("GET /games", MiddlewareAccount(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.Header().Add("Content-Type", "text/html")
				// data, err := data.GetGames()
				// if err != nil {
				// 	http.Redirect(w, r, "/404", http.StatusFound)
				// 	return
				// }
				// pp.Println(data)
				if err := web.PageGames(tmpStateDownload).Render(context.Background(), w); err != nil {
					panic(err)
				}
			}))

			mux.HandleFunc("GET /backups", MiddlewareAccount(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.Header().Add("Content-Type", "text/html")
				if err := web.PageBackups().Render(context.Background(), w); err != nil {
					panic(err)
				}
			}))

			mux.HandleFunc("GET /downloads", MiddlewareAccount(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.Header().Add("Content-Type", "text/html")
				if err := web.PageDownloads().Render(context.Background(), w); err != nil {
					panic(err)
				}
			}))

			mux.HandleFunc("GET /404", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				if err := web.PageNotFound().Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				accounts, err := data.GetAccounts()
				if err != nil {
					w.WriteHeader(400)
					w.Header().Add("Content-Type", "text/plain; charset=utf-8")
					w.Write([]byte("bad request"))
				}
				w.Header().Add("Content-Type", "application/json")
				bytes, err := json.Marshal(accounts)
				if err != nil {
					w.WriteHeader(400)
					w.Header().Add("Content-Type", "text/plain; charset=utf-8")
					w.Write([]byte("bad request"))
				}
				w.Write(bytes)
			})

			mux.HandleFunc("POST /api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
				var state types.GogAuthenticationChrome

				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(400)
					w.Header().Add("Content-Type", "text/plain; charset=utf-8")
					w.Write([]byte("bad request"))
					return
				}

				if err := json.Unmarshal(body, &state); err != nil {
					w.WriteHeader(400)
					w.Header().Add("Content-Type", "text/plain; charset=utf-8")
					w.Write([]byte("bad request"))
					return
				}

				if err := data.CreateAccountFromSignIn(state); err != nil {
					w.WriteHeader(500)
					w.Header().Add("Content-Type", "text/plain; charset=utf-8")
					w.Write([]byte("internal server error"))
					return
				}

				w.WriteHeader(201)
				w.Write([]byte("created"))
			})

			// TODO @droman: do more tests with REAL cookies get from the website using the procedure with the extension
			mux.HandleFunc("POST /api/v1/cookies", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.Header().Add("Content-Type", "text/plain; charset=utf-8")

				// Get the uploaded file
				file, handler, err := r.FormFile("file")
				if err != nil {
					// Handle the error
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("couldn't read the file"))
					return
				}
				defer file.Close()

				if !strings.Contains(handler.Filename, ".json") {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("file should be a json file"))
					return
				}

				// TODO @droman: deserve a whole handler for that

				dataFile := make([]byte, handler.Size)
				if _, err = file.Read(dataFile); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("couldn't read the bytes of file"))
					return
				}

				if len(dataFile) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("file is empty"))
					return
				}

				var obj []types.Cookie
				var datamap map[string]interface{}
				var dataCookies []byte

				if err = json.Unmarshal(dataFile, &datamap); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("couldn't unmarshal the file"))
					pp.Println(dataFile)
					return
				}

				if dataCookies, err = json.Marshal(datamap["Cookies"]); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("couldn't marshal the cookies"))
					return
				}

				if err := json.Unmarshal(dataCookies, &obj); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("couldn't unmarshal the cookies"))
					return
				}

				gogClient := data.NewGogClient()

				if err := data.SetCookies(gogClient.Client, obj, types.Hostname); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("couldn't set cookies"))
					return
				}

				user, ok, err := data.FetchUser(gogClient.Client, types.Hostname)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("couldn't fetch user"))
					return
				}

				if !ok {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("couldn't fetch user informations"))
					return
				}

				auth := types.GogAuthenticationChrome{
					Cookies: obj,
					User:    user,
				}

				if err := data.PostAccount(auth); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("couldn't save your authentication data"))
					return
				}

				// TODO @droman: check if the account already exists but for now...
				accounts, err := data.GetAccounts()
				if err != nil {
					w.WriteHeader(400)
					w.Header().Add("Content-Type", "text/plain; charset=utf-8")
					w.Write([]byte("bad request"))
				}

				go poolGameCache.Process(accounts)
				tmpStateDownload = true

				w.WriteHeader(200)
				w.Write([]byte("file uploaded: " + handler.Filename))
			})

			mux.HandleFunc("GET /api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				w.Header().Add("Content-Type", "text/plain; charset=utf-8")
				slog.Info("ping-pong")
				w.Write([]byte("pong"))
			})

			// TODO @droman: make a flag or env var
			mux.Handle("GET /static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./web/static/js"))))
			mux.Handle("GET /static/svg/", http.StripPrefix("/static/svg/", http.FileServer(http.Dir("./web/static/svg"))))
			mux.Handle("GET /static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("./web/static/css"))))

			server.Handler = mux

			go func() {
				slog.Info("server started")
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()

			<-c

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				// Handle error
				panic(err)
			}

			os.Exit(0)
		},
	}
}
