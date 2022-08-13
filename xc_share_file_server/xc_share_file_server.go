package xc_share_file_server

import (
	"crypto/rsa"
	"encoding/base32"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/ipoluianov/gomisc/crypt_tools"
	"github.com/ipoluianov/xchg/xchg"
	"github.com/ipoluianov/xchg/xchg_connections"
	"github.com/ipoluianov/xchg/xchg_network"
)

type XcShareFileServer struct {
	srv         *xchg_connections.ServerConnection
	fileName    string
	password    string
	fileContent []byte
	maxFileSize int64
}

func NewXcFileShareServer(fileName string, password string) *XcShareFileServer {
	var c XcShareFileServer
	c.fileName = fileName
	c.password = password
	c.maxFileSize = 1024 * 1024
	return &c
}

func (c *XcShareFileServer) Start() (serverAddress string, err error) {
	var fi fs.FileInfo
	fi, err = os.Stat(c.fileName)
	if err != nil {
		return
	}

	if fi.Size() > c.maxFileSize {
		err = errors.New("file is too large. max_size=" + fmt.Sprint(c.maxFileSize) + " bytes")
		return
	}

	c.fileContent, err = ioutil.ReadFile(c.fileName)
	if err != nil {
		return
	}

	if len(c.fileContent) > int(c.maxFileSize) {
		err = errors.New("file is too large. max_size=" + fmt.Sprint(c.maxFileSize) + " bytes")
		return
	}

	var privateKey *rsa.PrivateKey
	privateKey, err = crypt_tools.GenerateRSAKey()
	if err != nil {
		return
	}
	serverPrivateKeyBS := crypt_tools.RSAPrivateKeyToDer(privateKey)
	serverPrivateKey32 := base32.StdEncoding.EncodeToString(serverPrivateKeyBS)
	serverAddress = xchg.AddressForPublicKey(&privateKey.PublicKey)
	network := xchg_network.NewNetworkDefault()
	c.srv = xchg_connections.NewServerConnection(serverPrivateKey32, network)
	c.srv.SetProcessor(c)
	c.srv.Start()
	return
}

func (c *XcShareFileServer) ServerProcessorAuth(authData []byte) (err error) {
	if string(authData) == c.password {
		return nil
	}
	return errors.New(xchg.ERR_XCHG_ACCESS_DENIED)
}

func (c *XcShareFileServer) ServerProcessorCall(function string, parameter []byte) (response []byte, err error) {
	switch function {
	case "version":
		response = []byte("xc-share-file 2.02")
	case "get-file-name":
		response = []byte(c.fileName)
	case "get-file-content":
		response = c.fileContent
	default:
		err = errors.New(ERR_SIMPLE_SERVER_FUNC_IS_NOT_IMPL)
	}

	return
}

const (
	ERR_SIMPLE_SERVER_FUNC_IS_NOT_IMPL = "{ERR_SIMPLE_SERVER_FUNC_IS_NOT_IMPL}"
)
