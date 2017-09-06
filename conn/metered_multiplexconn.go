package meterconn

import (
	metrics "github.com/libp2p/go-libp2p-metrics"
	tpt "github.com/libp2p/go-libp2p-transport"
	smux "github.com/libp2p/go-stream-muxer"
)

type meteredStream struct {
	smux.Stream

	mc *MeteredMultiplexConn
}

var _ smux.Stream = &meteredStream{}

func newMeteredStream(mc *MeteredMultiplexConn, stream smux.Stream) *meteredStream {
	return &meteredStream{
		Stream: stream,
		mc:     mc,
	}
}

func (ms *meteredStream) Read(b []byte) (int, error) {
	n, err := ms.Stream.Read(b)
	ms.mc.mesRecv(int64(n))
	return n, err
}

func (ms *meteredStream) Write(b []byte) (int, error) {
	n, err := ms.Stream.Write(b)
	ms.mc.mesSent(int64(n))
	return n, err
}

type MeteredMultiplexConn struct {
	tpt.MultiplexConn

	mesRecv metrics.MeterCallback
	mesSent metrics.MeterCallback
}

var _ tpt.MultiplexConn = &MeteredMultiplexConn{}

func newMeteredMultiplexConn(base tpt.MultiplexConn, rcb metrics.MeterCallback, scb metrics.MeterCallback) *MeteredMultiplexConn {
	return &MeteredMultiplexConn{
		MultiplexConn: base,
		mesRecv:       rcb,
		mesSent:       scb,
	}
}

func (ms *MeteredMultiplexConn) OpenStream() (smux.Stream, error) {
	s, err := ms.MultiplexConn.OpenStream()
	if err != nil {
		return nil, err
	}
	return newMeteredStream(ms, s), nil
}

func (ms *MeteredMultiplexConn) AcceptStream() (smux.Stream, error) {
	s, err := ms.MultiplexConn.AcceptStream()
	if err != nil {
		return nil, err
	}
	return newMeteredStream(ms, s), nil
}
