package client

import (
	"errors"
	"io"
	"net/rpc"

	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/log"
)

type Client struct {
	*rpc.Client
}

// New creates a new client for a kopach_worker.
// Note that any kind of connection can be used here, other than the StdConn
func New(conn io.ReadWriteCloser) *Client {
	return &Client{rpc.NewClient(conn)}
}

// The following are all blocking calls as they are all triggers rather than
// queries and should return immediately the message is received.
// If deadlines are needed, set them on the connection,
// for StdConn this shouldn't be required as usually if the server is running
// worker will be too, a deadline would be needed for a network connection,
// or alternatively as with the Controller just spew messages over UDP

// NewJob is a delivery of a new job for the worker, this starts a miner
func (c *Client) NewJob(job *job.Container) (err error) {
	//log.DEBUG("sending new job")
	var reply bool
	err = c.Call("Worker.NewJob", job, &reply)
	if err != nil {
		log.ERROR(err)
		return
	}
	if reply != true {
		err = errors.New("new job command not acknowledged")
	}
	return
}

func (c *Client) Pause() (err error) {
	//log.DEBUG("sending pause")
	var reply bool
	err = c.Call("Worker.Pause", 1, &reply)
	if err != nil {
		log.ERROR(err)
		return
	}
	if reply != true {
		err = errors.New("pause command not acknowledged")
	}
	return
}

func (c *Client) Stop() (err error) {
	log.DEBUG("stop working (exit)")
	var reply bool
	err = c.Call("Worker.Stop", 1, &reply)
	if err != nil {
		log.ERROR(err)
		return
	}
	if reply != true {
		err = errors.New("stop command not acknowledged")
	}
	return
}

func (c *Client) SendPass(pass string) (err error) {
	log.DEBUG("sending dispatch password")
	var reply bool
		err = c.Call("Worker.SendPass", pass, &reply)
	if err != nil {
		log.ERROR(err)
		return
	}
	if reply != true {
		err = errors.New("send pass command not acknowledged")
	}
	return
}
