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
	//create the handler t
	hh := handler.NewProduct(l)

	//create serv max and add handler
	sm := http.NewServeMux()
	sm.Handle("/", hh)

	s := &http.Server{
		Addr:         ":9090",           // bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the error log for handler
		IdleTimeout:  120 * time.Second, // set idle timeout for request
		ReadTimeout:  1 * time.Second,   // set read time out for request and
		WriteTimeout: 1 * time.Second,   // set write time out for request
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
