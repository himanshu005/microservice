package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/himanshu005/microservice/handler"
	"golang.org/x/net/context"
)

func main() {
	l := log.New(os.Stdout, "product-api :", log.LstdFlags)
	hh := handler.NewHello(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	sig := <-sigChan

	l.Println("Recieved terminate,gracefull shoutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
