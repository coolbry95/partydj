package pool

import (
	"container/heap"
	"fmt"
	"time"

	"github.com/zmb3/spotify"
)

type Pool struct {
	PlaylistID spotify.ID `json:"playlistid"`
	PlaylistName string `'json:"playlist_name`
	UserID     string     `json:"userid"`
	// TimeStarted
	SongHeap []*Song `json:"songheap"`
	UserToVoteMap map[string][]string
}

type Song struct {
	ID        spotify.ID
	Upvotes   int       `json:"upvotes"`
	Downvotes int       `json:"downvotes"`
	Priority  int       `json:"priority"`
	index     int       `json:"index"` // index into the priorty queue
	TimeAdded time.Time `json:"timeadded"`

	// meta
	Album    string                 `json:"albumname"`
	Images   []spotify.Image        `json:"images"`
	Artists  []spotify.SimpleArtist `json:"artists"`
	Duration int                    `json"duration"`
	Name     string                 `json:"name"`
}

func (s *Song) String() string {
	return s.ID.String() + ", Priority: " + fmt.Sprintf("%d", s.Priority)
}

func (p *Pool) UpVote(id spotify.ID, userID string) {
	for i := range p.SongHeap {
		if p.SongHeap[i].ID == id {
			//fmt.Printf("(UpVote) Updated priority of %s is %d\n", v.ID, v.Priority)
			p.SongHeap[i].Upvotes++
			p.SongHeap[i].Priority++
			p.update(p.SongHeap[i], p.SongHeap[i].Priority)
			p.UserToVoteMap[userID] = append(p.UserToVoteMap[userID], p.SongHeap[i].ID.String())
			return
		}
	}
	fmt.Println("(UpVote) DID NOT FIND SONG")
}

func (p *Pool) DownVote(id spotify.ID, userID string) {
	//song := p.FindSong(id)
	for i := range p.SongHeap {
		if p.SongHeap[i].ID == id {
			//fmt.Printf("(UpVote) Updated priority of %s is %d\n", v.ID, v.Priority)
			p.SongHeap[i].Downvotes++
			p.SongHeap[i].Priority--
			p.update(p.SongHeap[i], p.SongHeap[i].Priority)
			p.UserToVoteMap[userID] = append(p.UserToVoteMap[userID], p.SongHeap[i].ID.String())
			return
		}
	}
	fmt.Println("(DownVote) DID NOT FIND SONG")
}

func (p *Pool) Len() int { return len(p.SongHeap) }

func (p *Pool) Less(i, j int) bool {
	return p.SongHeap[i].Priority > p.SongHeap[j].Priority
}

func (p *Pool) Swap(i, j int) {
	p.SongHeap[i], p.SongHeap[j] = p.SongHeap[j], p.SongHeap[i]
	p.SongHeap[i].index = i
	p.SongHeap[j].index = j
}

func (p *Pool) Push(x interface{}) {
	n := len(p.SongHeap)
	item := x.(*Song)
	item.index = n
	p.SongHeap = append(p.SongHeap, item)
}

func (p *Pool) Pop() interface{} {
	old := p.SongHeap
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	p.SongHeap = old[:n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *Pool) update(item *Song, priority int) {
	//fmt.Println("(update) Requested priority: ", priority)
	//fmt.Println("(update) Old priority of " + item.ID.String(), item.Priority)
	item.Priority = priority
	heap.Fix(pq, item.index)
	//fmt.Println("(update) Updated priority of " + item.ID.String(), item.Priority)
}

func (p *Pool) copyPool() *Pool {
	ps := &Pool{
		PlaylistID: p.PlaylistID,
		UserID:     p.UserID,
		SongHeap:   make([]*Song, len(p.SongHeap)),
	}
	copy(ps.SongHeap, p.SongHeap)
	return ps
}

func (p *Pool) getSecondSong() *Song {
	poolCopy := p.copyPool()

	firstSong := heap.Pop(poolCopy)
	song := heap.Pop(poolCopy)
	heap.Push(p, firstSong)
	heap.Push(p, song)
	return song.(*Song)
}

func (p *Pool) UpdateSpotifyPlaylist(c *spotify.Client, playlistId spotify.ID) {
	playlist, err := c.GetPlaylistTracks(p.UserID, playlistId)

	if err != nil {
		fmt.Println("ERROR: (UpdateSpotifyPlaylist) ", err.Error())
	}

	trackToRemoveId := playlist.Tracks[1].Track.ID
	trackToAddId := p.getSecondSong().ID

	//fmt.Println("To Remove: ", trackToRemoveId)
	//fmt.Println("To Add: ", trackToAddId)
	if trackToAddId != trackToRemoveId {
		//fmt.Println("(UpdateSpotifyPlaylist) Next song is not the correct one!")
		newPlayListId, _ := c.RemoveTracksFromPlaylist(p.UserID, playlistId, trackToRemoveId)
		newPlayListId, _ = c.AddTracksToPlaylist(p.UserID, spotify.ID(newPlayListId), trackToAddId)
	} else {
		//fmt.Println("(UpdateSpotifyPlaylist) Next song is correct!!")
	}
}

func (p *Pool) AddNextSong(c *spotify.Client) {
	firstSong := heap.Pop(p)
	nextSong := heap.Pop(p)
	toBeNextSong := heap.Pop(p)

	//fmt.Println("first song: ", firstSong,"next song: ", nextSong, "to be next song: ", toBeNextSong)
	_, err := c.RemoveTracksFromPlaylist(p.UserID, p.PlaylistID, firstSong.(*Song).ID)

	if err != nil {
		fmt.Println("remove error: ", err)
	}

	heap.Push(p, toBeNextSong)
	heap.Push(p, nextSong)

	c.AddTracksToPlaylist(p.UserID, p.PlaylistID, toBeNextSong.(*Song).ID)
}

func (p *Pool) HasUserVoted(userId string, songID string) bool {
	_, ok := p.UserToVoteMap[userId]
	if ok{
		votedSongs := p.UserToVoteMap[userId]
		for i := range votedSongs{
			if votedSongs[i] == songID{
				return true
			}
		}
	}
	return false
}

func TrackToSong(track *spotify.FullTrack, priority int) *Song {
	return &Song{
		ID: track.ID,
		Name: track.Name,
		Duration: track.Duration,
		Album: track.Album.Name,
		Images: track.Album.Images,
		Artists: track.Artists,
		Priority: priority,
	}
}
