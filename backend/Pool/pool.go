package Pool

import (
	"time"

	"container/heap"
	"github.com/zmb3/spotify"
	"fmt"
)

type Pool struct {
	PlaylistID spotify.ID
	UserID     string
	// TimeStarted
	SongHeap []*Song
}

type Song struct {
	ID        spotify.ID
	Upvotes   int
	Downvotes int
	Priority  int
	index     int // index into the priorty queue
	TimeAdded time.Time
}

func (s *Song) String() string {
	return "(ID: " + s.ID.String() + " Priority: " + fmt.Sprintf("%d)", s.Priority)
}

func (p *Pool) UpVote(id spotify.ID) {
	if song := p.FindSong(id); song != nil {
		song.Upvotes++
		p.update(song, song.Priority + 1)
	} else {
		fmt.Println("(UpVote) DID NOT FIND SONG")
	}
}

func (p *Pool) FindSong(id spotify.ID) *Song {
	for i := range p.SongHeap{
		if p.SongHeap[i].ID == id{
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
	p.SongHeap[i].index = i
	p.SongHeap[j].index = j
}

func (p *Pool) Push(x interface{}) {
	n := len(p.SongHeap)
	item := x.(*Song)
	item.index = n
	p.SongHeap = append(p.SongHeap, item)
	fmt.Println(p.SongHeap)
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
	item.Priority = priority
	heap.Fix(pq, item.index)
}


