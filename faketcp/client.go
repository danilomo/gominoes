package faketcp

import (
	"bytes"
	"net"
	"time"
)

type fakeClient struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

// FakeClient aaa
func FakeClient(str string) net.Conn {
	return &fakeClient{
		readBuffer:  bytes.NewBufferString(str),
		writeBuffer: bytes.NewBuffer(make([]byte, 1_000_000)),
	}
}

func (fc *fakeClient) Read(b []byte) (n int, err error) {
	return fc.readBuffer.Read(b)
}

func (fc *fakeClient) Write(b []byte) (n int, err error) {
	return fc.writeBuffer.Write(b)
}

func (fc *fakeClient) Close() error {
	return nil
}

func (fc *fakeClient) LocalAddr() net.Addr {
	return nil
}

func (fc *fakeClient) RemoteAddr() net.Addr {
	return nil
}

func (fc *fakeClient) SetDeadline(t time.Time) error {
	return nil
}

func (fc *fakeClient) SetReadDeadline(t time.Time) error {
	return nil
}

func (fc *fakeClient) SetWriteDeadline(t time.Time) error {
	return nil
}
