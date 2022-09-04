package xc_share_file_server

import (
	"crypto/rsa"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/ipoluianov/gomisc/crypt_tools"
	"github.com/ipoluianov/xchg/xchg"
	"github.com/ipoluianov/xchg/xchg_connections"
	"github.com/ipoluianov/xchg/xchg_network"
)

const Version = "0.0.44"

type XcShareFileServer struct {
	srv      *xchg_connections.ServerConnection
	fileName string
	baseName string
	password string
}

func NewXcFileShareServer(fileName string, password string) *XcShareFileServer {
	var c XcShareFileServer
	c.fileName = fileName
	c.password = password
	return &c
}

func (c *XcShareFileServer) Start() (err error) {
	fmt.Println("XcShareFileServer starting ...")
	fmt.Println("FileName:", c.fileName)
	var fi fs.FileInfo
	fi, err = os.Stat(c.fileName)
	if err != nil {
		return
	}

	c.baseName = fi.Name()

	var privateKey *rsa.PrivateKey
	privateKey, err = crypt_tools.GenerateRSAKey()
	if err != nil {
		return
	}
	serverPrivateKeyBS := crypt_tools.RSAPrivateKeyToDer(privateKey)
	serverPrivateKey32 := base32.StdEncoding.EncodeToString(serverPrivateKeyBS)
	serverAddress := xchg.AddressForPublicKey(&privateKey.PublicKey)
	network := xchg_network.NewNetworkFromInternet()
	c.srv = xchg_connections.NewServerConnection()
	c.srv.SetProcessor(c)
	c.srv.Start(serverPrivateKey32, network)

	// waiting for connections
	time.Sleep(1 * time.Second)

	state := c.srv.State()
	connectedPeers := make([]string, 0)
	fmt.Println()
	fmt.Println("----- CONNECTIONS -----")
	for _, p := range state.PeerConnections {
		if p.Init6Received {
			connectedPeers = append(connectedPeers, p.BaseConnection.Host)
			fmt.Println("Connected to ", p.BaseConnection.Host)
		}
	}
	fmt.Println("-----------------------")

	fmt.Println("ADDRESS:")
	fmt.Println()
	fmt.Println(serverAddress)
	fmt.Println()

	fmt.Println("To get file:")
	fmt.Println("xc-share-file get " + serverAddress + " " + c.password)

	return
}

func (c *XcShareFileServer) ServerProcessorAuth(authData []byte) (err error) {
	if string(authData) == c.password {
		fmt.Println("Success authentication")
		return nil
	}
	fmt.Println("WARNING: Wrong password:", string(authData))
	return errors.New(xchg.ERR_XCHG_ACCESS_DENIED)
}

func (c *XcShareFileServer) ServerProcessorCall(function string, parameter []byte) (response []byte, err error) {
	switch function {
	case "get-version":
		response = []byte("xc-share-file v" + Version)
	case "get-file-name":
		response = c.processGetFileName()
	case "get-file-size":
		response, err = c.processGetFileSize()
	case "get-file-content":
		response, err = c.processGetFileContent(parameter)
	case "thank-you":
		fmt.Println("Thank you")
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	default:
		err = errors.New(ERR_SIMPLE_SERVER_FUNC_IS_NOT_IMPL)
	}

	return
}

func (c *XcShareFileServer) processGetFileName() []byte {
	return []byte(c.baseName)
}

func (c *XcShareFileServer) processGetFileSize() ([]byte, error) {
	fs, err := os.Stat(c.fileName)
	if err != nil {
		return nil, err
	}
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(fs.Size()))
	return bs, nil
}

func (c *XcShareFileServer) processGetFileContent(parameter []byte) ([]byte, error) {
	var err error
	if len(parameter) != 16 {
		return nil, errors.New("wrong get-file-content parameter size")
	}
	offset := int(binary.LittleEndian.Uint64(parameter[0:]))
	size := int(binary.LittleEndian.Uint64(parameter[8:]))

	if size < 0 || size > 1024*1024 {
		return nil, errors.New("wrong size")
	}

	if offset < 0 || offset > 0x7FFFFFFF {
		return nil, errors.New("wrong offset")
	}

	var file *os.File

	file, err = os.OpenFile(c.fileName, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nSeek int64
	nSeek, err = file.Seek(int64(offset), 0)
	if err != nil {
		return nil, err
	}
	if int(nSeek) != offset {
		return nil, errors.New("can not seek offset in file")
	}

	var n int
	buffer := make([]byte, size)
	n, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 8+8+n)
	copy(response, parameter[0:16])
	copy(response[16:], buffer[:n])
	return response, nil
}

const (
	ERR_SIMPLE_SERVER_FUNC_IS_NOT_IMPL = "{ERR_SIMPLE_SERVER_FUNC_IS_NOT_IMPL}"
)
