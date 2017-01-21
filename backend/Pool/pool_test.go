package Pool

import (
	"testing"
	"github.com/zmb3/spotify"
)

func SetupPool() *Pool {
	return &Pool{
		SongHeap: []*Song{
			&Song{Priority:1, ID:spotify.ID(1)},
			&Song{Priority:2, ID:spotify.ID(2)},
			&Song{Priority:3, ID:spotify.ID(3)},
			&Song{Priority:4, ID:spotify.ID(4)},
		},
	}
}



func TestPool_Push(t *testing.T) {
	samplePool := SetupPool()
	songPushed := &Song{Priority:100, ID:spotify.ID(10)}
	samplePool.Push(songPushed)

	if popped := samplePool.Pop(); popped != songPushed {
		t.Errorf("expected %s, got %s", songPushed, popped.(*Song))
	}
}

func TestPool_Pop(t *testing.T) {
	samplePool := SetupPool()
	songPushed := &Song{Priority:4, ID:spotify.ID(4)}

	if popped := samplePool.Pop(); popped == songPushed {
		t.Errorf("expected %s, got %s", songPushed, popped.(*Song))
	}
}

