package rsocket

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -lrdmacm -libverbs
#include <rdma/rsocket.h>
*/
import "C"
import (
	"syscall"
	"unsafe"
)

// Socket domain constants
const (
	AF_INET  = C.AF_INET
	AF_INET6 = C.AF_INET6
)

// Socket type constants
const (
	SOCK_STREAM = C.SOCK_STREAM
	SOCK_DGRAM  = C.SOCK_DGRAM
)

// Protocol constants
const (
	IPPROTO_TCP = C.IPPROTO_TCP
	IPPROTO_UDP = C.IPPROTO_UDP
)

// RDMA specific socket options
const (
	SOL_RDMA    = C.SOL_RDMA
	RDMA_SQSIZE = C.RDMA_SQSIZE
	RDMA_RQSIZE = C.RDMA_RQSIZE
	RDMA_INLINE = C.RDMA_INLINE
	RDMA_ROUTE  = C.RDMA_ROUTE
)

// Socket creates a new RDMA socket
func Socket(domain, typ, protocol int) (int, error) {
	fd := C.rsocket(C.int(domain), C.int(typ), C.int(protocol))
	if fd < 0 {
		return -1, syscall.Errno(-fd)
	}
	return int(fd), nil
}

// Bind binds the socket to the given address
func Bind(fd int, sa syscall.Sockaddr) error {
	ptr, len, err := sockaddrToAny(sa)
	if err != nil {
		return err
	}
	if rc := C.rbind(C.int(fd), (*C.struct_sockaddr)(unsafe.Pointer(ptr)), C.socklen_t(len)); rc < 0 {
		return syscall.Errno(-rc)
	}
	return nil
}

// Listen marks the socket as a passive socket
func Listen(fd int, backlog int) error {
	if rc := C.rlisten(C.int(fd), C.int(backlog)); rc < 0 {
		return syscall.Errno(-rc)
	}
	return nil
}

// Accept accepts a connection on the given socket
func Accept(fd int) (int, syscall.Sockaddr, error) {
	var (
		addr syscall.RawSockaddrAny
		len  = C.socklen_t(syscall.SizeofSockaddrAny)
	)
	nfd := C.raccept(C.int(fd), (*C.struct_sockaddr)(unsafe.Pointer(&addr)), &len)
	if nfd < 0 {
		return -1, nil, syscall.Errno(-nfd)
	}
	sa, err := anyToSockaddr(&addr)
	if err != nil {
		return -1, nil, err
	}
	return int(nfd), sa, nil
}

// Connect connects the socket to a remote address
func Connect(fd int, sa syscall.Sockaddr) error {
	ptr, len, err := sockaddrToAny(sa)
	if err != nil {
		return err
	}
	if rc := C.rconnect(C.int(fd), (*C.struct_sockaddr)(unsafe.Pointer(ptr)), C.socklen_t(len)); rc < 0 {
		return syscall.Errno(-rc)
	}
	return nil
}

// Read reads data from the socket
func Read(fd int, p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n := C.rread(C.int(fd), unsafe.Pointer(&p[0]), C.size_t(len(p)))
	if n < 0 {
		return 0, syscall.Errno(-n)
	}
	return int(n), nil
}

// Write writes data to the socket
func Write(fd int, p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n := C.rwrite(C.int(fd), unsafe.Pointer(&p[0]), C.size_t(len(p)))
	if n < 0 {
		return 0, syscall.Errno(-n)
	}
	return int(n), nil
}

// Close closes the socket
func Close(fd int) error {
	if rc := C.rclose(C.int(fd)); rc < 0 {
		return syscall.Errno(-rc)
	}
	return nil
}

// SetSockOpt sets a socket option
func SetSockOpt(fd, level, opt int, value unsafe.Pointer, len uint32) error {
	if rc := C.rsetsockopt(C.int(fd), C.int(level), C.int(opt), value, C.socklen_t(len)); rc < 0 {
		return syscall.Errno(-rc)
	}
	return nil
}

// GetSockOpt gets a socket option
func GetSockOpt(fd, level, opt int, value unsafe.Pointer, len *uint32) error {
	l := C.socklen_t(*len)
	if rc := C.rgetsockopt(C.int(fd), C.int(level), C.int(opt), value, &l); rc < 0 {
		return syscall.Errno(-rc)
	}
	*len = uint32(l)
	return nil
}

// sockaddrToAny converts a syscall.Sockaddr to a syscall.RawSockaddrAny
func sockaddrToAny(sa syscall.Sockaddr) (*syscall.RawSockaddrAny, uint32, error) {
	if sa == nil {
		return nil, 0, syscall.EINVAL
	}

	switch sa := sa.(type) {
	case *syscall.SockaddrInet4:
		raw := syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
			Port:   uint16((sa.Port >> 8) | ((sa.Port & 0xff) << 8)), // network byte order
		}
		copy(raw.Addr[:], sa.Addr[:])
		return (*syscall.RawSockaddrAny)(unsafe.Pointer(&raw)), syscall.SizeofSockaddrInet4, nil

	case *syscall.SockaddrInet6:
		raw := syscall.RawSockaddrInet6{
			Family:   syscall.AF_INET6,
			Port:     uint16((sa.Port >> 8) | ((sa.Port & 0xff) << 8)), // network byte order
			Flowinfo: sa.ZoneId,
		}
		copy(raw.Addr[:], sa.Addr[:])
		return (*syscall.RawSockaddrAny)(unsafe.Pointer(&raw)), syscall.SizeofSockaddrInet6, nil

	default:
		return nil, 0, syscall.EAFNOSUPPORT
	}
}

// anyToSockaddr converts a syscall.RawSockaddrAny to a syscall.Sockaddr
func anyToSockaddr(rsa *syscall.RawSockaddrAny) (syscall.Sockaddr, error) {
	if rsa == nil {
		return nil, syscall.EINVAL
	}

	switch rsa.Addr.Family {
	case syscall.AF_INET:
		pp := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := &syscall.SockaddrInet4{
			Port: int(pp.Port<<8 | pp.Port>>8), // network byte order
		}
		copy(sa.Addr[:], pp.Addr[:])
		return sa, nil

	case syscall.AF_INET6:
		pp := (*syscall.RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := &syscall.SockaddrInet6{
			Port:   int(pp.Port<<8 | pp.Port>>8), // network byte order
			ZoneId: pp.Scope_id,
		}
		copy(sa.Addr[:], pp.Addr[:])
		return sa, nil

	default:
		return nil, syscall.EAFNOSUPPORT
	}
}
