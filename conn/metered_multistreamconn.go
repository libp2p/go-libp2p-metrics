package meterconn

import (
	metrics "github.com/libp2p/go-libp2p-metrics"
	tpt "github.com/libp2p/go-libp2p-transport"
	smux "github.com/libp2p/go-stream-muxer"
)

type meteredStream struct {
	smux.Stream

	mc *MeteredMultiStreamConn
}

var _ smux.Stream = &meteredStream{}

func newMeteredStream(mc *MeteredMultiStreamConn, stream smux.Stream) *meteredStream {
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

type MeteredMultiStreamConn struct {
	tpt.MultiStreamConn

	mesRecv metrics.MeterCallback
	mesSent metrics.MeterCallback
}

var _ tpt.MultiStreamConn = &MeteredMultiStreamConn{}

func newMeteredMultiStreamConn(base tpt.MultiStreamConn, rcb metrics.MeterCallback, scb metrics.MeterCallback) *MeteredMultiStreamConn {
	return &MeteredMultiStreamConn{
		MultiStreamConn: base,
		mesRecv:         rcb,
		mesSent:         scb,
	}
}

func (ms *MeteredMultiStreamConn) OpenStream() (smux.Stream, error) {
	s, err := ms.MultiStreamConn.OpenStream()
	if err != nil {
		return nil, err
	}
	return newMeteredStream(ms, s), nil
}

func (ms *MeteredMultiStreamConn) AcceptStream() (smux.Stream, error) {
	s, err := ms.MultiStreamConn.AcceptStream()
	if err != nil {
		return nil, err
	}
	return newMeteredStream(ms, s), nil
}
