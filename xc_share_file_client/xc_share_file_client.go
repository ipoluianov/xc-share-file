package xc_share_file_client

import (
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ipoluianov/gomisc/crypt_tools"
	"github.com/ipoluianov/xchg/xchg"
)

func GetFile(publicAddress string, password string, destFile string) {
	var err error
	var privateKey *rsa.PrivateKey
	privateKey, err = crypt_tools.GenerateRSAKey()
	if err != nil {
		return
	}

	client := xchg.NewPeer(privateKey)

	var version string
	version, err = getVersion(client, publicAddress, password)
	if err != nil {
		fmt.Println("getVersion ERROR:", err)
		return
	}
	fmt.Println("Version:", version)

	var fileName string
	fileName, err = getFileName(client, publicAddress, password)
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
	fileSize, err = getFileSize(client, publicAddress, password)
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

	blockSize := 64 * 1024
	for receviedBytes < fileSize {
		var block []byte
		block, err = getFileContent(client, publicAddress, password, receviedBytes, blockSize)
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
	client.Call(publicAddress, password, "thank-you", nil, time.Second)
}

func getVersion(client *xchg.Peer, publicAddress string, password string) (version string, err error) {
	for i := 0; i < 3; i++ {
		var bs []byte
		bs, err = client.Call(publicAddress, password, "get-version", nil, time.Second)
		if err == nil {
			version = string(bs)
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func getFileName(client *xchg.Peer, publicAddress string, password string) (fileName string, err error) {
	var bs []byte
	bs, err = client.Call(publicAddress, password, "get-file-name", nil, time.Second)
	if err != nil {
		return
	}
	fileName = string(bs)
	return
}

func getFileSize(client *xchg.Peer, publicAddress string, password string) (fileSize int, err error) {
	var bs []byte
	bs, err = client.Call(publicAddress, password, "get-file-size", nil, time.Second)
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

func getFileContent(client *xchg.Peer, publicAddress string, password string, offset int, size int) (fileContent []byte, err error) {
	var bs []byte
	parameter := make([]byte, 16)
	binary.LittleEndian.PutUint64(parameter[0:], uint64(offset))
	binary.LittleEndian.PutUint64(parameter[8:], uint64(size))
	bs, err = client.Call(publicAddress, password, "get-file-content", parameter, time.Second)
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
