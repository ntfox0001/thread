package thread

import (
	"testing"
)

func BenchmarkThread1(b *testing.B) {
	thread := NewThread(10, 10)
	thread.Run()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := i
		thread.AsyncExecRt(func(resolve Resolve) {
			resolve.R(v * v)
		})
	}

	thread.Close()
}

func BenchmarkFixedSelectLoop(b *testing.B) {
	sl := NewThread(10, 10)
	sl.Run()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := i
		sl.AsyncExecRt(func(resolve Resolve) {
			resolve.R(v * v)
		})
	}

	sl.Close()
}

func BenchmarkSelectLoop(b *testing.B) {
	sl := NewThread(10, 10)
	sl.Run()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := i
		sl.AsyncExecRt(func(resolve Resolve) {
			resolve.R(v * v)
		})
	}

	sl.Close()
}
