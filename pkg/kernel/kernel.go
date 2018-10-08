// Package kernel is mostly copied from by https://github.com/gopherdata/gophernotes.
package kernel

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/caicloud/ciao/pkg/manager"
	zmq "github.com/pebbe/zmq4"
)

const (
	kernelStarting = "starting"
	kernelBusy     = "busy"
	kernelIdle     = "idle"
)

var (
	// ExecCounter is incremented each time we run user code in the notebook.
	// It cannot be defined in Kernel struct, or it will be always 1.
	ExecCounter = 0
)

// Kernel is a kubeflow kernel.
type Kernel struct {
	// ProtocolVersion defines the Jupyter protocol version.
	ProtocolVersion string
	// Version defines the gophernotes version.
	Version        string
	ConnectionFile string
	Manager        *manager.Manager
}

// New creates a new kernel instance.
func New(protocolVersion, version, connectionFile string, manager *manager.Manager) *Kernel {
	return &Kernel{
		ProtocolVersion: protocolVersion,
		Version:         version,
		ConnectionFile:  connectionFile,
		Manager:         manager,
	}
}

// RunKernel is the main entry point to start the kernel.
func (k Kernel) RunKernel() {
	// Parse the connection info.
	var connInfo ConnectionInfo

	connData, err := ioutil.ReadFile(k.ConnectionFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal(connData, &connInfo); err != nil {
		log.Fatal(err)
	}

	// Set up the ZMQ sockets through which the kernel will communicate.
	sockets, err := prepareSockets(connInfo)
	if err != nil {
		log.Fatal(err)
	}

	// TODO connect all channel handlers to a WaitGroup to ensure shutdown before returning from runKernel.

	// Start up the heartbeat handler.
	startHeartbeat(sockets.HBSocket, &sync.WaitGroup{})

	// TODO gracefully shutdown the heartbeat handler on kernel shutdown by closing the chan returned by startHeartbeat.

	poller := zmq.NewPoller()
	poller.Add(sockets.ShellSocket.Socket, zmq.POLLIN)
	poller.Add(sockets.StdinSocket.Socket, zmq.POLLIN)
	poller.Add(sockets.ControlSocket.Socket, zmq.POLLIN)

	// msgParts will store a received multipart message.
	var msgParts [][]byte

	// Start a message receiving loop.
	for {
		polled, err := poller.Poll(-1)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range polled {

			// Handle various types of messages.
			switch socket := item.Socket; socket {

			// Handle shell messages.
			case sockets.ShellSocket.Socket:
				msgParts, err = sockets.ShellSocket.Socket.RecvMessageBytes(0)
				if err != nil {
					log.Println(err)
				}

				msg, ids, err := WireMsgToComposedMsg(msgParts, sockets.Key)
				if err != nil {
					log.Println(err)
					return
				}

				k.handleShellMsg(msgReceipt{msg, ids, sockets})

				// TODO Handle stdin socket.
			case sockets.StdinSocket.Socket:
				sockets.StdinSocket.Socket.RecvMessageBytes(0)

				// Handle control messages.
			case sockets.ControlSocket.Socket:
				msgParts, err = sockets.ControlSocket.Socket.RecvMessageBytes(0)
				if err != nil {
					log.Println(err)
					return
				}

				msg, ids, err := WireMsgToComposedMsg(msgParts, sockets.Key)
				if err != nil {
					log.Println(err)
					return
				}

				k.handleShellMsg(msgReceipt{msg, ids, sockets})
			}
		}
	}
}

// handleShellMsg responds to a message on the shell ROUTER socket.
func (k Kernel) handleShellMsg(receipt msgReceipt) {
	// Tell the front-end that the kernel is working and when finished notify the
	// front-end that the kernel is idle again.
	if err := receipt.PublishKernelStatus(kernelBusy); err != nil {
		log.Printf("Error publishing kernel status 'busy': %v\n", err)
	}
	defer func() {
		if err := receipt.PublishKernelStatus(kernelIdle); err != nil {
			log.Printf("Error publishing kernel status 'idle': %v\n", err)
		}
	}()

	switch receipt.Msg.Header.MsgType {
	case "kernel_info_request":
		if err := k.sendKernelInfo(receipt); err != nil {
			log.Fatal(err)
		}
	case "complete_request":
		if err := k.handleCompleteRequest(receipt); err != nil {
			log.Fatal(err)
		}
	case "execute_request":
		if err := k.handleExecuteRequest(receipt); err != nil {
			log.Fatal(err)
		}
	case "shutdown_request":
		k.handleShutdownRequest(receipt)
	default:
		log.Println("Unhandled shell message: ", receipt.Msg.Header.MsgType)
	}
}

// sendKernelInfo sends a kernel_info_reply message.
func (k Kernel) sendKernelInfo(receipt msgReceipt) error {
	return receipt.Reply("kernel_info_reply",
		Info{
			ProtocolVersion:       k.ProtocolVersion,
			Implementation:        "Kubeflow",
			ImplementationVersion: k.Version,
			Banner:                fmt.Sprintf("Kubeflow kernel: Kubeflow - v%s", k.Version),
			LanguageInfo: LanguageInfo{
				Name:          "Python (Kubeflow)",
				Version:       runtime.Version(),
				FileExtension: ".py",
			},
			HelpLinks: []HelpLink{
				{Text: "Kubeflow", URL: "https://www.kubeflow.org/"},
			},
		},
	)
}

// handleExecuteRequest runs code from an execute_request method,
// and sends the various reply messages.
func (k Kernel) handleExecuteRequest(receipt msgReceipt) error {

	// Extract the data from the request.
	reqcontent := receipt.Msg.Content.(map[string]interface{})
	code := reqcontent["code"].(string)
	silent := reqcontent["silent"].(bool)

	if !silent {
		ExecCounter++
	}

	// Prepare the map that will hold the reply content.
	content := make(map[string]interface{})
	content["execution_count"] = ExecCounter

	// Tell the front-end what the kernel is about to execute.
	if err := receipt.PublishExecutionInput(ExecCounter, code); err != nil {
		log.Printf("Error publishing execution input: %v\n", err)
	}

	// Redirect the standard out from the REPL.
	oldStdout := os.Stdout
	rOut, wOut, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stdout = wOut

	// Redirect the standard error from the REPL.
	oldStderr := os.Stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stderr = wErr

	var writersWG sync.WaitGroup
	writersWG.Add(2)

	// Forward all data written to stdout/stderr to the front-end.
	go func() {
		defer writersWG.Done()
		jupyterStdOut := JupyterStreamWriter{StreamStdout, &receipt}
		io.Copy(&jupyterStdOut, rOut)
	}()

	go func() {
		defer writersWG.Done()
		jupyterStdErr := JupyterStreamWriter{StreamStderr, &receipt}
		io.Copy(&jupyterStdErr, rErr)
	}()

	// eval
	vals, executionErr := k.doEval(code)

	// Close and restore the streams.
	wOut.Close()
	os.Stdout = oldStdout

	wErr.Close()
	os.Stderr = oldStderr

	// Wait for the writers to finish forwarding the data.
	writersWG.Wait()

	if executionErr == nil {
		// if the only non-nil value is image.Image or Data, render it
		data := renderResults(vals)

		content["status"] = "ok"
		content["user_expressions"] = make(map[string]string)

		if !silent && len(data.Data) != 0 {
			// Publish the result of the execution.
			if err := receipt.PublishExecutionResult(ExecCounter, data); err != nil {
				log.Printf("Error publishing execution result: %v\n", err)
			}
		}
	} else {
		content["status"] = "error"
		content["ename"] = "ERROR"
		content["evalue"] = executionErr.Error()
		content["traceback"] = nil

		if err := receipt.PublishExecutionError(executionErr.Error(), []string{executionErr.Error()}); err != nil {
			log.Printf("Error publishing execution error: %v\n", err)
		}
	}

	// Send the output back to the notebook.
	return receipt.Reply("execute_reply", content)
}

// doEval evaluates the code in the interpreter. This function captures an uncaught panic
// as well as the values of the last statement/expression.
func (k Kernel) doEval(code string) (val []interface{}, err error) {
	job, err := k.Manager.Execute(code)
	if err != nil {
		fmt.Printf("Failed to create job: %v", err)
		return nil, nil
	}
	fmt.Printf("Job %s is created.", job.Name)
	return val, nil
}

// handleShutdownRequest sends a "shutdown" message.
func (k Kernel) handleShutdownRequest(receipt msgReceipt) {
	content := receipt.Msg.Content.(map[string]interface{})
	restart := content["restart"].(bool)

	reply := shutdownReply{
		Restart: restart,
	}

	if err := receipt.Reply("shutdown_reply", reply); err != nil {
		log.Fatal(err)
	}

	log.Println("Shutting down in response to shutdown_request")
	os.Exit(0)
}

func (k Kernel) handleCompleteRequest(receipt msgReceipt) error {
	return nil
}
