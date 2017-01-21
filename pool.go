package main

import (
	"fmt"
	"time"

	"github.com/zmb3/spotify"
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

func (s *Song) getScore() int {
	return s.Upvotes - s.Downvotes
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