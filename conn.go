package rsocket

import (
	"log"
	"net"
	"syscall"
)

// TCPListener is a TCP network listener baseded on rsocket.
type TCPListener struct {
	ip      string
	port    int
	tcpAddr *net.TCPAddr
	fd      int
}

type Conn struct {
	fd         int
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr
}

// NewTCPListener creates a new TCPListener.
// It binds the listener to the given ip and port.
func NewTCPListener(ip string, port int) (*TCPListener, error) {
	fd, err := Socket(AF_INET, SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	srcAddr := net.ParseIP(ip)
	sa := &syscall.SockaddrInet4{
		Port: port,
	}
	copy(sa.Addr[:], srcAddr.To4())

	if err := Bind(fd, sa); err != nil {
		return nil, err
	}

	if err := Listen(fd, 128); err != nil {
		return nil, err
	}

	localAddr := &net.TCPAddr{
		IP:   srcAddr,
		Port: port,
	}

	return &TCPListener{
		ip:      ip,
		port:    port,
		fd:      fd,
		tcpAddr: localAddr,
	}, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *TCPListener) Accept() (*Conn, error) {
	fd, addr, err := Accept(l.fd)
	if err != nil {
		return nil, err
	}
	sa := addr.(*syscall.SockaddrInet4)
	ip := net.IPv4(sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3])
	port := sa.Port

	remoteAddr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}

	conn := &Conn{
		fd:         fd,
		localAddr:  l.tcpAddr,
		remoteAddr: remoteAddr,
	}

	return conn, nil
}

// Close closes the listener.
func (l *TCPListener) Close() error {
	return Close(l.fd)
}

// Addr returns the listener's network address.
func (l *TCPListener) Addr() net.Addr {
	return l.tcpAddr
}

// DialTCP connects to the address on the named network based on rsocket.
func DialTCP(address string) (*Conn, error) {
	fd, err := Socket(AF_INET, SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	sa := &syscall.SockaddrInet4{
		Port: tcpAddr.Port,
	}
	copy(sa.Addr[:], tcpAddr.IP.To4())

	if err := Connect(fd, sa); err != nil {
		return nil, err
	}

	conn := &Conn{
		fd:         fd,
		localAddr:  nil,
		remoteAddr: tcpAddr,
	}

	return conn, nil
}

// Read reads data from the connection.
func (c *Conn) Read(p []byte) (int, error) {
	return Read(c.fd, p)
}

// Write writes data to the connection.
func (c *Conn) Write(p []byte) (int, error) {
	return Write(c.fd, p)
}

// Close closes the connection.
func (c *Conn) Close() error {
	return Close(c.fd)
}

// LocalAddr returns the local network address.
func (c *Conn) LocalAddr() net.Addr {
	return c.localAddr
}

// RemoteAddr returns the remote network address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

// SetDeadline sets the read and write deadlines associated with the connection.
// not implementation.
func (c *Conn) SetDeadline() error {
	return nil
}

// SetReadDeadline sets the read deadline on the connection.
// not implementation.
func (c *Conn) SetReadDeadline() error {
	return nil
}

// SetWriteDeadline sets the write deadline on the connection.
// not implementation.
func (c *Conn) SetWriteDeadline() error {
	return nil
}
