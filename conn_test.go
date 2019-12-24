package medusa

import (
	"bufio"
	"errors"
	"fmt"
	"github/wziww/medusa/config"
	"github/wziww/medusa/encrpt"
	"net"
	"os"
	"sync"
	"testing"
)

var encryptor encrpt.Encryptor
var ErrNetClosing = errors.New("use of closed network connection")

func TestMain(m *testing.M) {
	config.Init()
	password := []byte(config.C.Base.Password)
	encryptor = encrpt.InitEncrypto(&password, config.C.Base.Crypto)
	m.Run()
}
func Init() (*net.TCPListener, *net.TCPConn) {
	var server *net.TCPListener
	var serverAddress *net.TCPAddr
	var client *net.TCPConn
	serverAddress, resoveErr := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if resoveErr != nil {
		fmt.Println(resoveErr)
		os.Exit(0)
	}
	// 服务启动
	var serverError error
	server, serverError = net.ListenTCP("tcp", serverAddress)
	if serverError != nil {
		os.Exit(0)
	}
	serverAddress, resoveErr = net.ResolveTCPAddr("tcp", server.Addr().String())
	if resoveErr != nil {
		fmt.Println(resoveErr)
		os.Exit(0)
	}
	var clientError error
	client, clientError = net.DialTCP("tcp", nil, serverAddress)
	if clientError != nil {
		fmt.Println(clientError)
		os.Exit(0)
	}
	return server, client
}

type rwc struct {
	buf   []byte
	mutex sync.Mutex
}

func (r *rwc) Read(p []byte) (n int, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	n = len(r.buf)
	if n > len(p) {
		n = len(p)
		r.buf = r.buf[:len(p)]
		for i := range p {
			p[i] = r.buf[i]
		}
		return
	}
	for i := range r.buf {
		p[i] = r.buf[i]
	}
	return
}

func (r *rwc) Write(p []byte) (n int, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.buf = append(r.buf, p...)
	return len(p), nil
}
func (r *rwc) Close() (err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.buf = r.buf[:0]
	return nil
}
func TestWriteRead(t *testing.T) {
	fakelocalConn := &rwc{}
	flT := &TCPConn{
		L:         "",
		R:         "",
		Reader:    bufio.NewReader(fakelocalConn),
		Writer:    fakelocalConn,
		Closer:    fakelocalConn,
		Encryptor: encryptor,
	}
	s := "hello World"
	flT.EncodeWrite([]byte(s))
	_, buf, _ := flT.DecodeRead()
	if string(buf) != s {
		t.Fatal(s, "!=", buf, "fail to Write and Read")
	}
}

func TestCopy(t *testing.T) {
	server, client := Init()
	flT := &TCPConn{
		L:         "",
		R:         "",
		Reader:    bufio.NewReader(client),
		Writer:    client,
		Closer:    client,
		Encryptor: encryptor,
	}
	conn, connError := server.Accept()
	if connError != nil {
		t.Fatal(connError)
	}
	connT := &TCPConn{
		Reader:    bufio.NewReader(conn),
		Writer:    conn,
		Closer:    conn,
		Encryptor: encryptor,
	}
	s := "hello World"
	flT.Write([]byte(s))
	buf := make([]byte, bufSize)
	encodeCopy(connT, connT, buf)
	decodeCopy(flT, flT)
	n, err := connT.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != s {
		t.Fatal(s, "!=", buf, "fail to Copy")
	}
}

func TestCloseEncodeCopy(t *testing.T) {
	server, client := Init()
	flT := &TCPConn{
		L:         "",
		R:         "",
		Reader:    bufio.NewReader(client),
		Writer:    client,
		Closer:    client,
		Encryptor: encryptor,
	}
	conn, connError := server.Accept()
	if connError != nil {
		t.Fatal(connError)
	}
	connT := &TCPConn{
		Reader:    bufio.NewReader(conn),
		Writer:    conn,
		Closer:    conn,
		Encryptor: encryptor,
	}
	s := "hello World"
	flT.Write([]byte(s))
	connT.Close()
	buf := make([]byte, bufSize)
	c, err := encodeCopy(connT, connT, buf)
	if c != false && err != ErrNetClosing {
		t.Fatal(err)
	}
}

func TestCloseEncodeCopy2(t *testing.T) {
	server, client := Init()
	flT := &TCPConn{
		L:         "",
		R:         "",
		Reader:    bufio.NewReader(client),
		Writer:    client,
		Closer:    client,
		Encryptor: encryptor,
	}
	conn, connError := server.Accept()
	if connError != nil {
		t.Fatal(connError)
	}
	connT := &TCPConn{
		Reader:    bufio.NewReader(conn),
		Writer:    conn,
		Closer:    conn,
		Encryptor: encryptor,
	}
	s := "hello World"
	flT.Write([]byte(s))
	flT.Close()
	buf := make([]byte, bufSize)
	c, err := encodeCopy(connT, flT, buf)
	if c != false && err != ErrNetClosing {
		t.Fatal(err)
	}
}
func TestCloseDecodeCopy(t *testing.T) {
	server, client := Init()
	flT := &TCPConn{
		L:         "",
		R:         "",
		Reader:    bufio.NewReader(client),
		Writer:    client,
		Closer:    client,
		Encryptor: encryptor,
	}
	conn, connError := server.Accept()
	if connError != nil {
		t.Fatal(connError)
	}
	connT := &TCPConn{
		Reader:    bufio.NewReader(conn),
		Writer:    conn,
		Closer:    conn,
		Encryptor: encryptor,
	}
	s := "hello World"
	flT.Write([]byte(s))
	flT.Close()
	buf := make([]byte, bufSize)
	encodeCopy(connT, connT, buf)
	c, err := decodeCopy(flT, flT)
	if c != false && err != ErrNetClosing {
		t.Fatal(err)
	}
}
func TestCloseDecodeCopy2(t *testing.T) {
	server, client := Init()
	flT := &TCPConn{
		L:         "",
		R:         "",
		Reader:    bufio.NewReader(client),
		Writer:    client,
		Closer:    client,
		Encryptor: encryptor,
	}
	conn, connError := server.Accept()
	if connError != nil {
		t.Fatal(connError)
	}
	connT := &TCPConn{
		Reader:    bufio.NewReader(conn),
		Writer:    conn,
		Closer:    conn,
		Encryptor: encryptor,
	}
	s := "hello World"
	flT.Write([]byte(s))
	buf := make([]byte, bufSize)
	encodeCopy(connT, connT, buf)
	connT.Close()
	c, err := decodeCopy(flT, connT)
	if c != false && err != ErrNetClosing {
		t.Fatal(err)
	}
}

func TestEncodeCopy(t *testing.T) {
	server, client := Init()
	flT := &TCPConn{
		L:         "",
		R:         "",
		Reader:    bufio.NewReader(client),
		Writer:    client,
		Closer:    client,
		Encryptor: encryptor,
	}
	conn, connError := server.Accept()
	if connError != nil {
		t.Fatal(connError)
	}
	connT := &TCPConn{
		Reader:    bufio.NewReader(conn),
		Writer:    conn,
		Closer:    conn,
		Encryptor: encryptor,
	}
	s := "hello World"
	flT.Write([]byte(s))
	connT.Close()
	connT.EncodeCopy(flT)
}
