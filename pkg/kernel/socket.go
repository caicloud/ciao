// Package kernel is mostly copied from by https://github.com/gopherdata/gophernotes.
package kernel

import (
	"fmt"
	"log"
	"sync"
	"time"

	zmq "github.com/pebbe/zmq4"
)

// RunWithSocket invokes the `run` function after acquiring the `Socket.Lock` and releases the lock when done.
func (s *Socket) RunWithSocket(run func(socket *zmq.Socket) error) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	return run(s.Socket)
}

// prepareSockets sets up the ZMQ sockets through which the kernel
// will communicate.
func prepareSockets(connInfo ConnectionInfo) (SocketGroup, error) {

	// Initialize the context.
	context, err := zmq.NewContext()
	if err != nil {
		return SocketGroup{}, err
	}

	// Initialize the socket group.
	var sg SocketGroup

	// Create the shell socket, a request-reply socket that may receive messages from multiple frontend for
	// code execution, introspection, auto-completion, etc.
	sg.ShellSocket.Socket, err = context.NewSocket(zmq.ROUTER)
	sg.ShellSocket.Lock = &sync.Mutex{}
	if err != nil {
		return sg, err
	}

	// Create the control socket. This socket is a duplicate of the shell socket where messages on this channel
	// should jump ahead of queued messages on the shell socket.
	sg.ControlSocket.Socket, err = context.NewSocket(zmq.ROUTER)
	sg.ControlSocket.Lock = &sync.Mutex{}
	if err != nil {
		return sg, err
	}

	// Create the stdin socket, a request-reply socket used to request user input from a front-end. This is analogous
	// to a standard input stream.
	sg.StdinSocket.Socket, err = context.NewSocket(zmq.ROUTER)
	sg.StdinSocket.Lock = &sync.Mutex{}
	if err != nil {
		return sg, err
	}

	// Create the iopub socket, a publisher for broadcasting data like stdout/stderr output, displaying execution
	// results or errors, kernel status, etc. to connected subscribers.
	sg.IOPubSocket.Socket, err = context.NewSocket(zmq.PUB)
	sg.IOPubSocket.Lock = &sync.Mutex{}
	if err != nil {
		return sg, err
	}

	// Create the heartbeat socket, a request-reply socket that only allows alternating recv-send (request-reply)
	// calls. It should echo the byte strings it receives to let the requester know the kernel is still alive.
	sg.HBSocket.Socket, err = context.NewSocket(zmq.REP)
	sg.HBSocket.Lock = &sync.Mutex{}
	if err != nil {
		return sg, err
	}

	// Bind the sockets.
	address := fmt.Sprintf("%v://%v:%%v", connInfo.Transport, connInfo.IP)
	sg.ShellSocket.Socket.Bind(fmt.Sprintf(address, connInfo.ShellPort))
	sg.ControlSocket.Socket.Bind(fmt.Sprintf(address, connInfo.ControlPort))
	sg.StdinSocket.Socket.Bind(fmt.Sprintf(address, connInfo.StdinPort))
	sg.IOPubSocket.Socket.Bind(fmt.Sprintf(address, connInfo.IOPubPort))
	sg.HBSocket.Socket.Bind(fmt.Sprintf(address, connInfo.HBPort))

	// Set the message signing key.
	sg.Key = []byte(connInfo.Key)

	return sg, nil
}

// startHeartbeat starts a go-routine for handling heartbeat ping messages sent over the given `hbSocket`. The `wg`'s
// `Done` method is invoked after the thread is completely shutdown. To request a shutdown the returned `shutdown` channel
// can be closed.
func startHeartbeat(hbSocket Socket, wg *sync.WaitGroup) (shutdown chan struct{}) {
	quit := make(chan struct{})

	// Start the handler that will echo any received messages back to the sender.
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Create a `Poller` to check for incoming messages.
		poller := zmq.NewPoller()
		poller.Add(hbSocket.Socket, zmq.POLLIN)

		for {
			select {
			case <-quit:
				return
			default:
				// Check for received messages waiting at most 500ms for once to arrive.
				pingEvents, err := poller.Poll(500 * time.Millisecond)
				if err != nil {
					log.Fatalf("Error polling heartbeat channel: %v\n", err)
				}

				// If there is at least 1 message waiting then echo it.
				if len(pingEvents) > 0 {
					hbSocket.RunWithSocket(func(echo *zmq.Socket) error {
						// Read a message from the heartbeat channel as a simple byte string.
						pingMsg, err := echo.RecvBytes(0)
						if err != nil {
							log.Fatalf("Error reading heartbeat ping bytes: %v\n", err)
							return err
						}

						// Send the received byte string back to let the front-end know that the kernel is alive.
						if _, err = echo.SendBytes(pingMsg, 0); err != nil {
							log.Printf("Error sending heartbeat pong bytes: %b\n", err)
							return err
						}

						return nil
					})
				}
			}
		}
	}()

	return quit
}
