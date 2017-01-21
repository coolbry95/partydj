package pool

import (
	"container/heap"
	"github.com/zmb3/spotify"
	"testing"
	//"fmt"
)

func SetupPool() *Pool {
	pool := &Pool{
		SongHeap: []*Song{
			&Song{Priority: 1, ID: spotify.ID("1")},
			&Song{Priority: 2, ID: spotify.ID("2")},
			&Song{Priority: 3, ID: spotify.ID("3")},
			&Song{Priority: 4, ID: spotify.ID("4")},
		},
	}

	heap.Init(pool)
	return pool
}

func TestPool_Push(t *testing.T) {
	samplePool := SetupPool()
	songPushed := &Song{Priority: 100, ID: spotify.ID("10")}
	heap.Push(samplePool, songPushed)

	if popped := heap.Pop(samplePool); popped.(*Song).ID != songPushed.ID {
		t.Errorf("expected %s, got %s", songPushed.ID, popped.(*Song).ID)
	}
}

func TestPool_Pop(t *testing.T) {
	samplePool := SetupPool()

	if popped := heap.Pop(samplePool); popped.(*Song).ID != spotify.ID("4") {
		t.Errorf("expected %s, got %s", spotify.ID("4"), popped.(*Song).ID)
	}
}

func TestPool_UpVote(t *testing.T) {
	samplePool := SetupPool()
	targetSongId := spotify.ID("3")

	// Avoid the priority being equivalent
	samplePool.UpVote(targetSongId)
	samplePool.UpVote(targetSongId)

	if popped := heap.Pop(samplePool); popped.(*Song).ID != targetSongId {
		t.Errorf("expected %s, got %s", targetSongId.String(), popped.(*Song).ID)
	}
}

func TestPool_DownVote(t *testing.T) {
	samplePool := SetupPool()
	targetSongId := spotify.ID("4")
	newLargestSongId := spotify.ID("3")

	// Avoid the priority being equivalent
	samplePool.DownVote(targetSongId)
	samplePool.DownVote(targetSongId)

	if popped := heap.Pop(samplePool); popped.(*Song).ID != newLargestSongId {
		t.Errorf("expected %s, got %s", targetSongId.String(), popped.(*Song).ID)
	}
}

func TestPool_FindSong(t *testing.T) {
	samplePool := SetupPool()
	targetSongID := spotify.ID("3")

	if found := samplePool.FindSong(spotify.ID("3")); found.ID != targetSongID {
		t.Errorf("expected %s, got %s", targetSongID, found.ID)
	}
}
