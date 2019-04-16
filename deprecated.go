package metrics

import moved "github.com/libp2p/go-libp2p-core/metrics"

// Deprecated: use github.com/libp2p/go-libp2p/metrics.Reporter instead.
type Reporter = moved.Reporter

// Deprecated: use github.com/libp2p/go-libp2p/metrics.Stats instead.
type Stats = moved.Stats

// Deprecated: use github.com/libp2p/go-libp2p/metrics.BandwidthCounter instead.
type BandwidthCounter = moved.BandwidthCounter

// Deprecated: use github.com/libp2p/go-libp2p/metrics.NewBandwidthCounter instead.
var NewBandwidthCounter = moved.NewBandwidthCounter
