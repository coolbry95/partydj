// This example demonstrates how to authenticate with Spotify.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/coolbry95/partydj/backend/pool"
	"github.com/pressly/chi"
	"github.com/zmb3/spotify"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "https://linode.shellcode.in/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate)
	ch    = make(chan *spotify.Client)
	state = "stateless"
)

type DI struct {
	client *spotify.Client
	pool   pool.Pool
}

func main() {

	var d *DI


	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
	})
	r.Post("/callback", completeAuth)
	r.Get("/pool/:poolID", d.getPool)
	r.Get("/pool", d.getPool)
	//r.Post("/add_song/:poolID/:songID", handle)
	//r.Post("/upvote/:poolID/:songID", handle)
	//r.Post("/down/:poolID/:songID", handle)

	// wait for auth to complete
	d.client = <-ch
	d.pool.PlaylistID = spotify.ID("0nXlYUH7zBAzubO9Yub4rR")
	userID, err := d.client.CurrentUser()
	if err != nil {
		log.Println(err)
	}
	d.pool.UserID = userID.ID
	fmt.Println(d)

	go http.ListenAndServe(":6060", r)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

}

func (d *DI) getPool(w http.ResponseWriter, r *http.Request) {
	userid, err := d.client.CurrentUser()
	if err != nil {
		log.Println(err)
	}
	playlist, err := d.client.GetPlaylistTracks(userid.ID, "0nXlYUH7zBAzubO9Yub4rR")
	if err != nil {
		log.Println(err)
	}
	d.pool.SongHeap = make([]*pool.Song, 10)
	for _, track := range playlist.Tracks {
		var s *pool.Song
		s.ID = track.Track.ID
		d.pool.SongHeap = append(d.pool.SongHeap, s)
	}
	json.NewEncoder(w).Encode(d.pool)
	return
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
	http.Redirect(w, r, "https://linode.shellcode.in/pool", 300)
}
