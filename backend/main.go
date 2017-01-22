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
	"time"

	"strconv"

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

var PoolShortIDToLongID map[int]string
var UserIDToPoolID map[string]string

func main() {
	// Initialize the pool (single for demonstration)
	var d *DI
	d = new(DI)

	// setup the router and paths
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
	r.Post("/join_pool", d.joinPool)

	r.FileServer("/files", http.Dir("/home/coolbry95/gosrc/src/github.com/coolbry95/partydj/website"))

	// Authenticate the users spotify account
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	go http.ListenAndServe(":6060", r)

	// wait for auth to complete
	d.client = <-ch
	d.pool.PlaylistID = spotify.ID("0nXlYUH7zBAzubO9Yub4rR") // change to spartahack playlist

	// Initialize the short id to long id map
	PoolShortIDToLongID = make(map[int]string)
	// hard code the one join digit for this sample pool to 121
	PoolShortIDToLongID[121] = "0nXlYUH7zBAzubO9Yub4rR"

	// initialize the user map
	UserIDToPoolID = make(map[string]string)

	userID, err := d.client.CurrentUser()
	if err != nil {
		log.Println(err)
	}

	d.pool.UserID = userID.ID

	// Begin the play loop when the server begins for spartahack
	d.BeginPlayLoop()

	block := make(chan struct{})
	<-block
}

func (d *DI) BeginPlayLoop() {
	//fmt.Println("Beginning play loop...")
	// Start the loop after 10 seconds of pool creation
	time.Sleep(10 * time.Second)

	go func() {
		song := heap.Pop(&d.pool).(*pool.Song)
		heap.Push(&d.pool, song)
		// song duration is in milliseconds
		songTime := time.Duration(((song.Duration)/1000) + 15) // 15 seconds + song length
		fmt.Print(songTime * time.Second)
		time.Sleep(songTime * time.Second)
		d.pool.AddNextSong(&d.client)
		//fmt.Println("Now playing new song...")
		d.BeginPlayLoop()
		//return // do this to kill the parent goroutine?
	}()
}

func (d *DI) joinPool(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost")

	userIDNumber := r.PostFormValue("userId")
	poolIDNumberStr := r.PostFormValue("poolShortId")

	// Both must be present
	if len(userIDNumber) == 0 || len(poolIDNumberStr) == 0 {
		w.WriteHeader(http.StatusPartialContent)
		fmt.Println("One of the parameters were empty: ", "userId:",
			userIDNumber, "poolShortId:", poolIDNumberStr)
		return
	}

	poolIDNumber, convErr := strconv.Atoi(poolIDNumberStr)

	// Make sure the poolIDNumber is a valid integer, otherwise
	// respond with Status partial content
	if convErr != nil {
		w.WriteHeader(http.StatusPartialContent)
		fmt.Println("Conversion error: ", convErr)
		return
	}

	// Check if pool exists, if not respond with status not found
	if _, ok := PoolShortIDToLongID[poolIDNumber]; ok != true {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	UserIDToPoolID[userIDNumber] = poolIDNumberStr
	// this is the long pool id that links to the long id
	w.Header().Add("long_pool_id", PoolShortIDToLongID[poolIDNumber])
	// Let the user know the request was accepted
	w.WriteHeader(http.StatusOK)
}

func (d *DI) addSong(w http.ResponseWriter, r *http.Request) {
	songID := chi.URLParam(r, "songID")
	//poolID := chi.URLParam(r, "poolID")

	s := &pool.Song{ID: spotify.ID(songID)}
	heap.Push(&d.pool, s)
}

func (d *DI) upVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost")

	// TODO check for user ID to see if they already voted
	songID := chi.URLParam(r, "songId")
	userID := r.PostFormValue("userId")

	if len(userID) == 0 || len(songID) == 0 {
		w.WriteHeader(http.StatusPartialContent)
		return
	}

	if !d.pool.HasUserVoted(userID, songID) {
		d.pool.UpVote(spotify.ID(songID), userID)
		d.pool.UpdateSpotifyPlaylist(&d.client, d.pool.PlaylistID)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	w.WriteHeader(http.StatusOK)
}

func (d *DI) downVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost")

	// TODO check for user ID to see if they already voted
	songID := chi.URLParam(r, "songId")
	userID := r.PostFormValue("userId")
	//poolID := chi.URLParam(r, "poolID")

	if len(userID) == 0 || len(songID) == 0 {
		w.WriteHeader(http.StatusPartialContent)
		return
	}

	if !d.pool.HasUserVoted(userID, songID) {
		d.pool.DownVote(spotify.ID(songID), userID)
		d.pool.UpdateSpotifyPlaylist(&d.client, d.pool.PlaylistID)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (d *DI) createPool(w http.ResponseWriter, r *http.Request) {
}

func (d *DI) getPool(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost")

	userid, err := d.client.CurrentUser()
	if err != nil {
		log.Println(err)
	}

	// TODO: instead of using existing playlist we require a new playlist to be created
	// playlist, err := d.client.CreatePlaylistForUser(userid.ID, playlistName, true)
	playlist, err := d.client.GetPlaylistTracks(userid.ID, "0nXlYUH7zBAzubO9Yub4rR")
	//spartahackPlaylistName := "Sparthack Sample Playlist!"
	//newPlaylist, err := d.client.CreatePlaylistForUser(userid.ID, spartahackPlaylistName, true)
	if err != nil {
		log.Println(err)
	}

	d.pool.SongHeap = make([]*pool.Song, 0, 100)
	for i, track := range playlist.Tracks {
		d.pool.SongHeap = append(d.pool.SongHeap, pool.TrackToSong(&track.Track, i))
	}
	d.pool.UserToVoteMap = make(map[string][]string)

	//TODO: only call this function only after the the current song finishes

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
