package util

import (
	"log"
	"net"
	"net/rpc"
)

// MakeConnection Creates and returns a TCP connection between localAddr and remoteAddr
func MakeConnection(localAddr string, remoteAddr string) *net.TCPConn {
	localTcpAddr, err := net.ResolveTCPAddr("tcp", localAddr)
	CheckErr(err, "Could not resolve address: "+localAddr)
	remoteTcpAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	CheckErr(err, "Could not resolve address: "+remoteAddr)
	conn, err := net.DialTCP("tcp", localTcpAddr, remoteTcpAddr)
	CheckErr(err, "Could not connect "+localAddr+" to "+remoteAddr)
	return conn
}

// MakeClient Creates an RPC client given between a local and remote address
func MakeClient(localAddr string, remoteAddr string) (*net.TCPConn, *rpc.Client) {
	conn := MakeConnection(localAddr, remoteAddr)
	client := rpc.NewClient(conn)
	return conn, client
}

type RPCEndPoint struct {
	Addr   string
	Client *rpc.Client
}

func (e *RPCEndPoint) Call(methodName string, args interface{}, reply interface{}) error {
	if e.Client == nil {
		client, err := Connect(e.Addr)
		if err != nil {
			return err
		}
		e.Client = client
	}

	err := e.Client.Call(methodName, args, reply)
	if err != nil {
		log.Println(err)
		e.Client = nil
		return err
	}
	return nil
}

func Connect(address string) (*rpc.Client, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return client, nil
}