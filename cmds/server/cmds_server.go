package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/davidroman0O/gogog/data"
	"github.com/davidroman0O/gogog/types"
	"github.com/davidroman0O/gogog/web"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "server",
		Long:  `.`,
		Run: func(cmd *cobra.Command, args []string) {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)

			if err := data.Initialize("./tmp"); err != nil {
				panic(err)
			}

			server := &http.Server{
				Addr:         ":8080",
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			mux := http.NewServeMux()

			mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/" {
					http.Redirect(w, r, "/404", http.StatusFound)
					return
				}
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				if err := web.Index().Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				accounts, err := data.GetAccounts()
				if err != nil {
					http.Redirect(w, r, "/404", http.StatusFound)
					return
				}
				if err := web.PageAccounts(accounts).Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /games", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				if err := web.PageGames().Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /backups", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				if err := web.PageBackups().Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /downloads", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				if err := web.PageDownloads().Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /404", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				if err := web.PageNotFound().Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			mux.HandleFunc("GET /api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
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

			mux.HandleFunc("GET /api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/plain; charset=utf-8")
				slog.Info("ping-pong")
				w.Write([]byte("pong"))
			})

			// TODO @droman: make a flag or env var
			mux.Handle("GET /static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./web/static/js"))))
			mux.Handle("GET /static/svg/", http.StripPrefix("/static/svg/", http.FileServer(http.Dir("./web/static/svg"))))

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
