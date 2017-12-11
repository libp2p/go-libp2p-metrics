package metrics

import (
	"fmt"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"sync"
	"testing"
	"time"
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
