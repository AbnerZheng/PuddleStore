package puddlestore

import (
	"fmt"
	"net"
	"net/rpc"
)

type PuddleRPCServer struct {
	node     *PuddleNode
	listener net.Listener
	rpc      *rpc.Server
}

func newPuddlestoreRPCServer(puddle *PuddleNode) (server *PuddleRPCServer) {
	server = new(PuddleRPCServer)
	server.node = puddle
	server.rpc = rpc.NewServer()
	listener, _, err := OpenListener()
	server.rpc.RegisterName(listener.Addr().String(), server)
	server.listener = listener

	if err != nil {
		panic("AA")
	}

	go func() {
		for {
			conn, err := server.listener.Accept()
			if err != nil {
				fmt.Printf("(%v) Raft RPC server accept error: %v\n", err)
				continue
			}
			go server.rpc.ServeConn(conn)
		}
	}()

	return
}

func (server *PuddleRPCServer) ConnectImpl(req *ConnectRequest, rep *ConnectReply) error {
	rvreply, err := server.node.connect(req)
	*rep = rvreply
	return err
}

func (server *PuddleRPCServer) LsImpl(req *LsRequest, rep *LsReply) error {
	rvreply, err := server.node.ls(req)
	*rep = rvreply
	return err
}

func (server *PuddleRPCServer) CdImpl(req *CdRequest, rep *CdReply) error {
	rep, err := server.node.cd(req)
	return err
}

func (server *PuddleRPCServer) MkdirImpl(req *MkdirRequest, rep *MkdirReply) error {
	rep, err := server.node.mkdir(req)
	return err
}
