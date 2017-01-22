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
		UserToVoteMap: make(map[string][]string),
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
        targetUserId := "1234"
	// Avoid the priority being equivalent
	samplePool.UpVote(targetSongId, targetUserId)
	samplePool.UpVote(targetSongId, targetUserId)

	if popped := heap.Pop(samplePool); popped.(*Song).ID != targetSongId {
		t.Errorf("expected %s, got %s", targetSongId.String(), popped.(*Song).ID)
	}
}

func TestPool_DownVote(t *testing.T) {
	samplePool := SetupPool()
	targetSongId := spotify.ID("4")
	targetUserId := "1234"
	newLargestSongId := spotify.ID("3")

	// Avoid the priority being equivalent
	samplePool.DownVote(targetSongId, targetUserId)
	samplePool.DownVote(targetSongId, targetUserId)

	if popped := heap.Pop(samplePool); popped.(*Song).ID != newLargestSongId {
		t.Errorf("expected %s, got %s", targetSongId.String(), popped.(*Song).ID)
	}
}

func TestPool_HasUserVoted(t *testing.T) {
	samplePool := SetupPool()
	targetUserID := "1234"
	targetSongID := spotify.ID("3")

	if voted := samplePool.HasUserVoted(targetUserID, targetSongID.String()); voted{
		t.Errorf("expected false (user not found), got true")
	}

	samplePool.UpVote(targetSongID, targetUserID)
	if voted := samplePool.HasUserVoted(targetUserID, targetSongID.String()); !voted{
		t.Errorf("expected true (user voted on %s)", targetSongID.String())
	}

	if voted := samplePool.HasUserVoted(targetUserID, "5678"); voted{
		t.Errorf("expected false (user did not vote on 5678)")
	}

	targetSongIDDown := spotify.ID("33")
	samplePool.DownVote(targetSongID, targetUserID)
	if voted := samplePool.HasUserVoted(targetUserID, targetSongIDDown.String()); voted{
		t.Errorf("expected false (user did vote on 33)")
	}

	if voted := samplePool.HasUserVoted(targetUserID, "34"); voted{
		t.Errorf("expected false (user did not vote on 34)")
	}
}