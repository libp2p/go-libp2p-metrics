package meterconn

import (
	metrics "github.com/libp2p/go-libp2p-metrics"
	tpt "github.com/libp2p/go-libp2p-transport"
)

type MeteredDuplexConn struct {
	tpt.DuplexConn

	mesRecv metrics.MeterCallback
	mesSent metrics.MeterCallback
}

var _ tpt.DuplexConn = &MeteredDuplexConn{}

func newMeteredDuplexConn(base tpt.DuplexConn, rcb metrics.MeterCallback, scb metrics.MeterCallback) *MeteredDuplexConn {
	return &MeteredDuplexConn{
		DuplexConn: base,
		mesRecv:    rcb,
		mesSent:    scb,
	}
}

func (mc *MeteredDuplexConn) Read(b []byte) (int, error) {
	n, err := mc.DuplexConn.Read(b)
	mc.mesRecv(int64(n))
	return n, err
}

func (mc *MeteredDuplexConn) Write(b []byte) (int, error) {
	n, err := mc.DuplexConn.Write(b)
	mc.mesSent(int64(n))
	return n, err
}
