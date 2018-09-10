// Package kernel is mostly copied from by https://github.com/gopherdata/gophernotes.
package kernel

import (
	"sync"

	zmq "github.com/pebbe/zmq4"
)

// ConnectionInfo stores the contents of the kernel connection
// file created by Jupyter.
type ConnectionInfo struct {
	SignatureScheme string `json:"signature_scheme"`
	Transport       string `json:"transport"`
	StdinPort       int    `json:"stdin_port"`
	ControlPort     int    `json:"control_port"`
	IOPubPort       int    `json:"iopub_port"`
	HBPort          int    `json:"hb_port"`
	ShellPort       int    `json:"shell_port"`
	Key             string `json:"key"`
	IP              string `json:"ip"`
}

// LanguageInfo holds information about the language that this kernel executes code in.
type LanguageInfo struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	MIMEType          string `json:"mimetype"`
	FileExtension     string `json:"file_extension"`
	PygmentsLexer     string `json:"pygments_lexer"`
	CodeMirrorMode    string `json:"codemirror_mode"`
	NBConvertExporter string `json:"nbconvert_exporter"`
}

// HelpLink stores data to be displayed in the help menu of the notebook.
type HelpLink struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

// Info holds information about the kubeflow kernel, for kernel_info_reply messages.
type Info struct {
	ProtocolVersion       string       `json:"protocol_version"`
	Implementation        string       `json:"implementation"`
	ImplementationVersion string       `json:"implementation_version"`
	LanguageInfo          LanguageInfo `json:"language_info"`
	Banner                string       `json:"banner"`
	HelpLinks             []HelpLink   `json:"help_links"`
}

// Socket wraps a zmq socket with a lock which should be used to control write access.
type Socket struct {
	Socket *zmq.Socket
	Lock   *sync.Mutex
}

// SocketGroup holds the sockets needed to communicate with the kernel,
// and the key for message signing.
type SocketGroup struct {
	ShellSocket   Socket
	ControlSocket Socket
	StdinSocket   Socket
	IOPubSocket   Socket
	HBSocket      Socket
	Key           []byte
}

// shutdownReply encodes a boolean indication of shutdown/restart.
type shutdownReply struct {
	Restart bool `json:"restart"`
}

// MsgHeader encodes header info for ZMQ messages.
type MsgHeader struct {
	MsgID           string `json:"msg_id"`
	Username        string `json:"username"`
	Session         string `json:"session"`
	MsgType         string `json:"msg_type"`
	ProtocolVersion string `json:"version"`
	Timestamp       string `json:"date"`
}

// ComposedMsg represents an entire message in a high-level structure.
type ComposedMsg struct {
	Header       MsgHeader
	ParentHeader MsgHeader
	Metadata     map[string]interface{}
	Content      interface{}
}

// msgReceipt represents a received message, its return identities, and
// the sockets for communication.
type msgReceipt struct {
	Msg        ComposedMsg
	Identities [][]byte
	Sockets    SocketGroup
}

// bundledMIMEData holds data that can be presented in multiple formats. The keys are MIME types
// and the values are the data formatted with respect to its MIME type. All bundles should contain
// at least a "text/plain" representation with a string value.
type BundledMIMEData map[string]interface{}

type Data struct {
	Data      BundledMIMEData `json:"data"`
	Metadata  BundledMIMEData `json:"metadata"`
	Transient BundledMIMEData `json:"transient"`
}

// InvalidSignatureError is returned when the signature on a received message does not
// validate.
type InvalidSignatureError struct{}
