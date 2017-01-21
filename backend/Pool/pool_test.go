package Pool

import (
	"testing"
	"github.com/zmb3/spotify"
	"container/heap"
	"fmt"
)

func SetupPool() *Pool {
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

func TestPool_Push(t *testing.T) {
	samplePool := SetupPool()
	songPushed := &Song{Priority:100, ID:spotify.ID("10")}
	//fmt.Println("Pushing ", songPushed.ID)
	samplePool.Push(songPushed)


	if popped := samplePool.Pop(); popped.(*Song).ID != songPushed.ID {
		t.Errorf("expected %s, got %s", songPushed, popped.(*Song))
	}
}

func TestPool_Pop(t *testing.T) {
	samplePool := SetupPool()
	songPushed := &Song{Priority:4, ID:spotify.ID("4")}

	if popped := samplePool.Pop(); popped.(*Song).ID == songPushed.ID {
		t.Errorf("expected %s, got %s", songPushed, popped.(*Song))
	}
}

func TestPool_UpVote(t *testing.T) {
	samplePool := SetupPool()
	targetSongId := spotify.ID("3")

	fmt.Printf("song %s is %d\n", targetSongId, samplePool.SongHeap[2].Priority)

	samplePool.UpVote(targetSongId)
	samplePool.UpVote(targetSongId)

	fmt.Printf("song %s is %d\n", targetSongId, samplePool.SongHeap[2].Priority)

	if popped := samplePool.Pop(); popped.(*Song).ID != targetSongId {
		t.Errorf("expected %s, got %s", targetSongId.String(), popped.(*Song))

	}

}