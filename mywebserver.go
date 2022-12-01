package mywebserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type ControllersSet interface {
	Controllers() []Controller
}

type Controller interface {
	Path() string
	Name() string
	DoAction(string)
}

func StartServer(cs Controller) {
	ctx := context.Background()
	router := mux.NewRouter()
	srv := &http.Server{
		Addr:              `0.0.0.0:8080`,
		ReadTimeout:       time.Millisecond * 200,
		WriteTimeout:      time.Millisecond * 200,
		IdleTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Millisecond * 200,
		Handler:           router,
	}

	router.HandleFunc(cs.Path(), func(http.ResponseWriter, *http.Request) {
		cs.DoAction(`Artist invited`)
	})
	router.HandleFunc(cs.Path(), func(http.ResponseWriter, *http.Request) {
		cs.DoAction(`Painting taken off the exhibition`)
	})
	router.HandleFunc(cs.Path(), func(http.ResponseWriter, *http.Request) {
		cs.DoAction(`Exhibition opened`)
	})
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, _ = rw.Write([]byte("Hello!"))
	})

	http.Handle("/", router)

	go func() {
		log.Println(`Web Server started`)
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	<-done

	log.Println(`Web Server is shutting down`)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(ctx, err)
	}
}
