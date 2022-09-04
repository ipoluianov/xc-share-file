package xc_share_file_client

import (
	"crypto/rsa"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ipoluianov/gomisc/crypt_tools"
	"github.com/ipoluianov/xchg/xchg_connections"
	"github.com/ipoluianov/xchg/xchg_network"
)

func GetFile(publicAddress string, password string, destFile string) {
	var err error
	var privateKey *rsa.PrivateKey
	privateKey, err = crypt_tools.GenerateRSAKey()
	if err != nil {
		return
	}
	privateKeyBS := crypt_tools.RSAPrivateKeyToDer(privateKey)
	clientPrivateKey := base32.StdEncoding.EncodeToString(privateKeyBS)

	client := xchg_connections.NewClientConnection(xchg_network.NewNetworkFromInternet(), publicAddress, clientPrivateKey, password, nil)

	var version string
	version, err = getVersion(client)
	if err != nil {
		fmt.Println("getVersion ERROR:", err)
		return
	}
	fmt.Println("Version:", version)

	var fileName string
	fileName, err = getFileName(client)
	if err != nil {
		fmt.Println("getFileName ERROR:", err)
		return
	}
	fmt.Println("FileName:", fileName)

	if len(destFile) == 0 {
		_, err = os.Stat(fileName)
		if err == nil {
			fmt.Println("ERROR: local file already exists")
			return
		}
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Println("ERROR:", err)
			return
		}
		destFile = fileName
	}

	var fileSize int
	fileSize, err = getFileSize(client)
	if err != nil {
		fmt.Println("getFileSize ERROR:", err)
		return
	}

	//fileContent := make([]byte, 0)

	var file *os.File
	file, err = os.OpenFile(destFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("create local file ERROR:", err)
		return
	}
	defer file.Close()

	receviedBytes := 0

	blockSize := 1 * 1024
	for receviedBytes < fileSize {
		var block []byte
		block, err = getFileContent(client, receviedBytes, blockSize)
		if err != nil {
			fmt.Println("getFileContent ERROR:", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		file.Write(block)
		receviedBytes += len(block)
		fmt.Printf("received %d bytes; %.2f %%\r\n", receviedBytes, float64(receviedBytes)/float64(fileSize)*100)
	}

	fmt.Println("Complete", receviedBytes, "bytes")
	client.Call("thank-you", nil)
}

func getVersion(client *xchg_connections.ClientConnection) (version string, err error) {
	var bs []byte
	bs, err = client.Call("get-version", nil)
	if err != nil {
		return
	}
	version = string(bs)
	return
}

func getFileName(client *xchg_connections.ClientConnection) (fileName string, err error) {
	var bs []byte
	bs, err = client.Call("get-file-name", nil)
	if err != nil {
		return
	}
	fileName = string(bs)
	return
}

func getFileSize(client *xchg_connections.ClientConnection) (fileSize int, err error) {
	var bs []byte
	bs, err = client.Call("get-file-size", nil)
	if err != nil {
		return
	}
	if len(bs) != 8 {
		err = errors.New("wrong get-file-size response")
		return
	}
	fileSize = int(binary.LittleEndian.Uint64(bs))
	return
}

func getFileContent(client *xchg_connections.ClientConnection, offset int, size int) (fileContent []byte, err error) {
	var bs []byte
	parameter := make([]byte, 16)
	binary.LittleEndian.PutUint64(parameter[0:], uint64(offset))
	binary.LittleEndian.PutUint64(parameter[8:], uint64(size))
	bs, err = client.Call("get-file-content", parameter)
	if err != nil {
		return
	}
	if len(bs) < 16 {
		err = errors.New("wrong get-file-content response")
		return
	}

	fileContent = bs[16:]
	return
}
