package thread

import (
	"testing"
	"time"
)

func TestTicker1(t *testing.T) {
	sl := NewThread(10, 10*time.Second)
	sl.Run()
	ticker := NewTicker(sl, t)

	t.Log("start")
	i := 0
	ticker.Start(time.Second, func() {
		t.Log("i: ", i)
		i++

	})

	time.Sleep(10 * time.Second)

}

func TestTicker2(t *testing.T) {
	sl := NewThread(10, 10*time.Second)
	sl.Run()
	ticker := NewTicker(sl, t)

	ticker.Start(time.Second, func() {
		t.Log("ticker1")
	})

	ticker.Start(time.Second*2, func() {
		t.Log("ticker2")
	})

	time.Sleep(10 * time.Second)

	p := ticker.Release()
	p.SyncThen()
}
