package mywebserver

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"mvc/app/controllers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer() {
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
	artistC := &controllers.ArtistController{}
	artC := &controllers.ArtController{}
	galleryC := &controllers.ExhibitionController{}

	// CREATE AN ART
	// localhost:8080/createart/blackCat
	router.HandleFunc("/createart/{art}", artC.ArtCreation)

	// CREATE AN ARTIST
	// localhost:8080/createartist/Fillip
	router.HandleFunc("/createartist/{artist}", artistC.Registration)

	// CREATE EXHIBITION
	// localhost:8080/creategallery/Tokio
	router.HandleFunc("/createexhibition/{exhibition}", galleryC.ExhibitionCreation)

	//ASSIGN AN ART TO THE ARTIST (BY NAME)
	// localhost:8080/artist/assign/Fillip/blackCat
	router.HandleFunc("/artist/assign/{artist}/{art}", artC.AssignArt)

	//REGISTRATION AN ARTIST ON THE EXHIBITION
	// localhost:8080/artist/register/Fillip/Tokio
	router.HandleFunc("/artist/register/{artist}/{exhibition}", artistC.ArtistRegistration)

	// DELETE AN ARTIST FROM THE EXHIBITION
	// localhost:8080/artist/delete/Fillip/Tokio
	router.HandleFunc("/artist/delete/{artist}/{gallery}", galleryC.KickArtistOffExhibition)
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
