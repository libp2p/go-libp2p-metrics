package meterconn

import (
	"bytes"
	"net"
	"time"

	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockDuplexConn struct {
	dataToRead  bytes.Buffer
	dataWritten bytes.Buffer
}

var _ tpt.DuplexConn = &mockDuplexConn{}

func (c *mockDuplexConn) Read(b []byte) (int, error)       { return c.dataToRead.Read(b) }
func (c *mockDuplexConn) Write(p []byte) (int, error)      { return c.dataWritten.Write(p) }
func (c *mockDuplexConn) Close() error                     { panic("not implemented") }
func (c *mockDuplexConn) LocalAddr() net.Addr              { panic("not implemented") }
func (c *mockDuplexConn) LocalMultiaddr() ma.Multiaddr     { panic("not implemented") }
func (c *mockDuplexConn) RemoteAddr() net.Addr             { panic("not implemented") }
func (c *mockDuplexConn) RemoteMultiaddr() ma.Multiaddr    { panic("not implemented") }
func (c *mockDuplexConn) SetDeadline(time.Time) error      { panic("not implemented") }
func (c *mockDuplexConn) SetReadDeadline(time.Time) error  { panic("not implemented") }
func (c *mockDuplexConn) SetWriteDeadline(time.Time) error { panic("not implemented") }
func (c *mockDuplexConn) Transport() tpt.Transport         { panic("not implemented") }

type counter struct {
	count int64
}

func (c *counter) Count(n int64) { c.count += n }

var _ = Describe("DuplexConn", func() {
	var (
		conn  *mockDuplexConn
		mconn *MeteredDuplexConn

		readCounter  *counter
		writeCounter *counter
	)

	BeforeEach(func() {
		readCounter = &counter{}
		writeCounter = &counter{}

		conn = &mockDuplexConn{}
		conn.dataToRead.Write([]byte("foobar"))
		mconn = newMeteredDuplexConn(conn, readCounter.Count, writeCounter.Count)
	})

	It("counts data read", func() {
		n, _ := mconn.Read(make([]byte, 3))
		Expect(n).To(Equal(3))
		n, _ = mconn.Read(make([]byte, 2))
		Expect(n).To(Equal(2))
		Expect(readCounter.count).To(BeEquivalentTo(5))
		Expect(writeCounter.count).To(BeZero())
	})

	It("counts data written", func() {
		_, err := mconn.Write([]byte("foo"))
		Expect(err).ToNot(HaveOccurred())
		_, err = mconn.Write([]byte("bar"))
		Expect(err).ToNot(HaveOccurred())
		Expect(readCounter.count).To(BeZero())
		Expect(writeCounter.count).To(BeEquivalentTo(6))
	})
})
