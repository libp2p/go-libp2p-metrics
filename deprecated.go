package metrics

import moved "github.com/libp2p/go-libp2p-core/metrics"

// Deprecated: use github.com/libp2p/go-libp2p-core/metrics.Reporter instead.
type Reporter = moved.Reporter

// Deprecated: use github.com/libp2p/go-libp2p-core/metrics.Stats instead.
type Stats = moved.Stats

// Deprecated: use github.com/libp2p/go-libp2p-core/metrics.BandwidthCounter instead.
type BandwidthCounter = moved.BandwidthCounter

// Deprecated: use github.com/libp2p/go-libp2p-core/metrics.NewBandwidthCounter instead.
func NewBandwidthCounter() *moved.BandwidthCounter {
	return moved.NewBandwidthCounter()
}
