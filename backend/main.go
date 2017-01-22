// This example demonstrates how to authenticate with Spotify.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"container/heap"
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
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserLibraryModify, spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistModifyPublic)
	ch    = make(chan spotify.Client)
	state = "stateless"
)

type DI struct {
	client spotify.Client
	pool   pool.Pool
}

func main() {

	var d *DI
	d = new(DI)

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
	})
	r.Get("/callback", completeAuth)
	r.Get("/getpool", d.getPool)
	r.Post("/createpool", d.createPool)
	r.Post("/add_song/:poolID/:songID", d.addSong)
	r.Post("/upvote/:poolID/:songID", d.upVote)
	r.Post("/downvote/:poolID/:songID", d.downVote)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	go http.ListenAndServe(":6060", r)

	// wait for auth to complete
	d.client = <-ch

	d.pool.PlaylistID = spotify.ID("0nXlYUH7zBAzubO9Yub4rR")

	userID, err := d.client.CurrentUser()
	if err != nil {
		log.Println(err)
	}

	d.pool.UserID = userID.ID

	block := make(chan struct{})
	<-block

}

func (d *DI) addSong(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "songID")
	//poolID := chi.URLParam(r, "poolID")

	s := &pool.Song{ID: spotify.ID(songID)}
	heap.Push(&d.pool, s)

}

func (d *DI) upVote(w http.ResponseWriter, r *http.Request) {
	// TODO check for user ID to see if they already voted
	songID := chi.URLParam(r, "songID")

	d.pool.UpVote(spotify.ID(songID))
	d.pool.UpdateSpotifyPlaylist(&d.client, d.pool.PlaylistID)
}

func (d *DI) downVote(w http.ResponseWriter, r *http.Request) {
	// TODO check for user ID to see if they already voted
	songID := chi.URLParam(r, "songID")
	//poolID := chi.URLParam(r, "poolID")

	d.pool.DownVote(spotify.ID(songID))
	d.pool.UpdateSpotifyPlaylist(&d.client, d.pool.PlaylistID)
}

func (d *DI) createPool(w http.ResponseWriter, r *http.Request) {
}

func (d *DI) getPool(w http.ResponseWriter, r *http.Request) {
	userid, err := d.client.CurrentUser()
	if err != nil {
		log.Println(err)
	}

	// TODO: instead of using existing playlist we require a new playlist to be created
	// playlist, err := d.client.CreatePlaylistForUser(userid.ID, playlistName, true)
	playlist, err := d.client.GetPlaylistTracks(userid.ID, "0nXlYUH7zBAzubO9Yub4rR")
	if err != nil {
		log.Println(err)
	}

	d.pool.SongHeap = make([]*pool.Song, 0, 10)
	for i, track := range playlist.Tracks {
		d.pool.SongHeap = append(d.pool.SongHeap, pool.TrackToSong(&track.Track.SimpleTrack, i))
	}

	tracks, err := d.client.GetTracks("7ccI9cStQbQdystvc6TvxD", "2U9v51tNOoLRhwrU6j83uU")
	if err != nil {
		fmt.Println(err)
	}

	for i, track := range tracks {
		s := pool.TrackToSong(&track.SimpleTrack, (i+1)*-100)
		//fmt.Println(d.pool)
		heap.Push(&d.pool, s)
		//fmt.Println(d.pool)
	}

	//TODO: only call this function only after the the current song finishes
	//d.pool.AddNextSong(&d.client)

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
	ch <- client
	http.Redirect(w, r, "/getpool", 301)
}
