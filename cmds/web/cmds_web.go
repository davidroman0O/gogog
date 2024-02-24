package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/davidroman0O/gogog/views"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "web",
		Short: "server web",
		Long:  `.`,
		Run: func(cmd *cobra.Command, args []string) {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)

			server := &http.Server{
				Addr:         ":8080",
				ReadTimeout:  10 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			mux := http.NewServeMux()
			mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Add("Content-Type", "text/html")
				if err := views.Hello("myself").Render(context.Background(), w); err != nil {
					panic(err)
				}
			})

			server.Handler = mux

			go func() {
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
