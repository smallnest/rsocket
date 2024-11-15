# rsocket
`rsockets` is a protocol over RDMA that supports a socket-level API for applications. rsocket APIs are intended to match the behavior of corresponding socket calls, except where noted. rsocket functions match the name and function signature of socket calls, with the exception that all function calls are prefixed with an 'r'.

This project encapsulates rsocket, provides an interface compatible with `net.Conn`, can be directly replaced with `net.Conn`, provides a more convenient interface, can directly use the encapsulated interface of rsocket, or use the interface of `net.Conn`. This project is based on `RDMA`, so the performance will be better than net.Conn.

![GitHub](https://img.shields.io/github/license/smallnest/rsocket) [![GoDoc](https://godoc.org/github.com/smallnest/rsocket?status.png)](http://godoc.org/github.com/smallnest/rsocket)  


**This project is developing so don't use it in production environments.**

## Usage

Some examples of using rsocket are provided in the `examples` directory. The examples contains a simple TCP/UDP echo server and client.

## Reference

- [rsocket(7) - Linux man page](https://linux.die.net/man/7/rsocket)
- [rsocket](https://github.com/linux-rdma/rdma-core/blob/master/librdmacm/docs/rsocket)