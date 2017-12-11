package metrics

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

func BenchmarkBandwidthCounter(b *testing.B) {
	b.StopTimer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bwc := NewBandwidthCounter()
		round(bwc, b)
	}
}

func round(bwc *BandwidthCounter, b *testing.B) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 0; i < 1000; i++ {
		p := peer.ID(fmt.Sprintf("peer-%d", i))
		for j := 0; j < 10; j++ {
			proto := protocol.ID(fmt.Sprintf("bitswap-%d", j))
			go func() {
				defer wg.Done()
				<-start

				for i := 0; i < 1000; i++ {
					bwc.LogSentMessage(100)
					bwc.LogSentMessageStream(100, proto, p)
					time.Sleep(1 * time.Millisecond)
				}
			}()
		}
	}

	b.StartTimer()
	close(start)
	wg.Wait()
	b.StopTimer()
}

func TestBandwidthCounter(t *testing.T) {
	bwc := NewBandwidthCounter()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(200)
	for i := 0; i < 100; i++ {
		p := peer.ID(fmt.Sprintf("peer-%d", i))
		for j := 0; j < 2; j++ {
			proto := protocol.ID(fmt.Sprintf("bitswap-%d", j))
			go func() {
				defer wg.Done()
				<-start

				t := time.NewTicker(100 * time.Millisecond)
				defer t.Stop()

				for i := 0; i < 40; i++ {
					bwc.LogSentMessage(100)
					bwc.LogRecvMessage(50)
					bwc.LogSentMessageStream(100, proto, p)
					bwc.LogRecvMessageStream(50, proto, p)
					<-t.C
				}
			}()
		}
	}

	close(start)
	time.Sleep(2*time.Second + 500*time.Millisecond)
	for i := 0; i < 100; i++ {
		stats := bwc.GetBandwidthForPeer(peer.ID(fmt.Sprintf("peer-%d", i)))
		if !approxEq(stats.RateOut, 2000, 200) {
			t.Errorf("expected rate 1000 (±200), got %f", stats.RateOut)
		}

		if !approxEq(stats.RateIn, 1000, 100) {
			t.Errorf("expected rate 500 (±100), got %f", stats.RateIn)
		}
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
	for i := 0; i < 100; i++ {
		stats := bwc.GetBandwidthForPeer(peer.ID(fmt.Sprintf("peer-%d", i)))
		if stats.TotalOut != 8000 {
			t.Errorf("expected total 8000, got %d", stats.TotalOut)
		}

		if stats.TotalIn != 4000 {
			t.Errorf("expected total 4000, got %d", stats.TotalIn)
		}
	}
}

func approxEq(a, b, err float64) bool {
	return math.Abs(a-b) < err
}
