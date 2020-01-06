package medusa

import (
	"bufio"
	"errors"
	"fmt"
	"github/wziww/medusa/config"
	"github/wziww/medusa/encrypt"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

var encryptor encrypt.Encryptor
var ErrNetClosing = errors.New("use of closed network connection")

func TestMain(m *testing.M) {
	config.Init()
	password := []byte(config.C.Base.Password)
	encryptor = encrypt.InitEncrypto(&password, config.C.Base.Crypto, config.C.Base.Padding, nil)
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
		Encryptor: &encryptor,
	}
	s := []byte("hello World")
	b := make([]byte, len(s)+(*flT.Encryptor).Ivlen())
	flT.ssEncodeWrite(s)
	flT.Read(b)
	s2 := (*flT.Encryptor).Decode(b[(*flT.Encryptor).Ivlen():], b[:(*flT.Encryptor).Ivlen()])
	if string(s2) != string(s) {
		t.Fatal(s, "!=", s2, "fail to Write and Read")
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
		Encryptor: &encryptor,
	}
	conn, connError := server.Accept()
	if connError != nil {
		t.Fatal(connError)
	}
	connT := &TCPConn{
		Reader:    bufio.NewReader(conn),
		Writer:    conn,
		Closer:    conn,
		Encryptor: &encryptor,
	}
	s := []byte("hello World")
	flT.Write(s)
	buf := make([]byte, len(s)+(*connT.Encryptor).Ivlen())
	go func() {
		connT.SSEncodeCopy(flT)
	}()
	select {
	case <-time.After(100 * time.Millisecond):
	}
	connT.Read(buf)
	s2 := (*flT.Encryptor).Decode(buf[(*flT.Encryptor).Ivlen():], buf[:(*flT.Encryptor).Ivlen()])
	if string(s2) != string(s) {
		t.Fatal(s, "!=", s2, "fail to Write and Read")
	}
}

func TestDecodeCopy(t *testing.T) {
	server, client := Init()
	flT := &TCPConn{
		L:         "",
		R:         "",
		Reader:    bufio.NewReader(client),
		Writer:    client,
		Closer:    client,
		Encryptor: &encryptor,
	}
	conn, connError := server.Accept()
	if connError != nil {
		t.Fatal(connError)
	}
	connT := &TCPConn{
		Reader:    bufio.NewReader(conn),
		Writer:    conn,
		Closer:    conn,
		Encryptor: &encryptor,
	}
	s := []byte("hello World")
	flT.ssEncodeWrite(s)
	buf := make([]byte, len(s)+(*connT.Encryptor).Ivlen())
	go func() {
		connT.SSDecodeCopy(connT, nil)
	}()
	select {
	case <-time.After(100 * time.Millisecond):
	}
	b, _ := flT.Read(buf)
	fmt.Println(string(b))
	s2 := buf[(*flT.Encryptor).Ivlen():]
	if string(s2) != string(s) {
		t.Fatal(s, "!=", s2, "fail to Write and Read")
	}
}
