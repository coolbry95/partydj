package pool

import (
	"time"
	"container/heap"
	"fmt"

	"github.com/zmb3/spotify"
)

type Pool struct {
	PlaylistID spotify.ID `json:playlistid`
	UserID     string     `json:userid`
	// TimeStarted
	SongHeap []*Song `json:songheap`
}

type Song struct {
	ID        spotify.ID
	Upvotes   int       `json:upvotes`
	Downvotes int       `json:downvotes`
	Priority  int       `json:priority`
	index     int       `json:index` // index into the priorty queue
	TimeAdded time.Time `json:timeadded`
}

func (s *Song) String() string {
	return s.ID.String() + ", Priority: " + fmt.Sprintf("%d", s.Priority)
}

func (p *Pool) UpVote(id spotify.ID) {
	if song := p.FindSong(id); song != nil {
		song.Upvotes++
		song.Priority++
		p.update(song, song.Priority)
	} else {
		fmt.Println("(UpVote) DID NOT FIND SONG")
	}
}

func (p *Pool) UpVote(id spotify.ID) {
	if song := p.FindSong(id); song != nil {
		song.Downvotes++
		song.Priority--
		p.update(song, song.Priority)
	} else {
		fmt.Println("(UpVote) DID NOT FIND SONG")
	}
}

func (p *Pool) FindSong(id spotify.ID) *Song {
	for i := range p.SongHeap {
		if p.SongHeap[i].ID == id {
			//fmt.Printf("(UpVote) Updated priority of %s is %d\n", v.ID, v.Priority)
			return p.SongHeap[i]
		}
	}
	return nil
}

func (p *Pool) Len() int { return len(p.SongHeap) }

func (p *Pool) Less(i, j int) bool {
	return p.SongHeap[i].Priority > p.SongHeap[j].Priority
}

func (p *Pool) Swap(i, j int) {
	p.SongHeap[i], p.SongHeap[j] = p.SongHeap[j], p.SongHeap[i]
	p.SongHeap[i].Priority = i
	p.SongHeap[j].Priority = j
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
