package pool

import (
	"testing"
	"github.com/zmb3/spotify"
	"container/heap"
	"fmt"
)

func Setuppool() *Pool {
	pool := &Pool{
		SongHeap: []*Song{
			&Song{Priority:1, ID:spotify.ID("1")},
			&Song{Priority:2, ID:spotify.ID("2")},
			&Song{Priority:3, ID:spotify.ID("3")},
			&Song{Priority:4, ID:spotify.ID("4")},
		},
	}

	heap.Init(pool)
	return pool
}

func Testpool_Push(t *testing.T) {
	samplepool := Setuppool()
	songPushed := &Song{Priority:100, ID:spotify.ID("10")}
	//fmt.Println("Pushing ", songPushed.ID)
	samplepool.Push(songPushed)


	if popped := samplepool.Pop(); popped.(*Song).ID != songPushed.ID {
		t.Errorf("expected %s, got %s", songPushed, popped.(*Song))
	}
}

func Testpool_Pop(t *testing.T) {
	samplepool := Setuppool()
	songPushed := &Song{Priority:4, ID:spotify.ID("4")}

	if popped := samplepool.Pop(); popped.(*Song).ID == songPushed.ID {
		t.Errorf("expected %s, got %s", songPushed, popped.(*Song))
	}
}

func Testpool_UpVote(t *testing.T) {
	samplepool := Setuppool()
	targetSongId := spotify.ID("3")

	fmt.Printf("song %s is %d\n", targetSongId, samplepool.SongHeap[2].Priority)

	samplepool.UpVote(targetSongId)
	samplepool.UpVote(targetSongId)

	fmt.Printf("song %s is %d\n", targetSongId, samplepool.SongHeap[2].Priority)

	if popped := samplepool.Pop(); popped.(*Song).ID != targetSongId {
		t.Errorf("expected %s, got %s", targetSongId.String(), popped.(*Song))

	}

}
