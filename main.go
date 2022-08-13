package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"os"
	"strings"
	"xc_share_file/xc_share_file_client"
	"xc_share_file/xc_share_file_server"

	"github.com/btcsuite/btcutil/base58"
)

func genAddr(substr string, prefix string) {
	bs := make([]byte, 8)
	for i := 0; i < 4000000000; i++ {
		rand.Read(bs)
		base32.StdEncoding.DecodeString("123")
		bsSHA := sha256.Sum256(bs)
		address := base58.Encode(bsSHA[:32])
		if (i % 100000) == 0 {
			fmt.Println(prefix, i)
		}
		a := strings.ToLower(address)
		if strings.Contains(a, substr) {
			index := strings.Index(a, substr)
			address = address[0:index] + "." + address[index:index+len(substr)] + "." + address[index+len(substr):]
			fmt.Println("Addr:", address, i)
			panic(address)
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("LEGEND")
		return
	}

	command := os.Args[1]

	if command == "share" {
		if len(os.Args) != 4 {
			fmt.Println("... share FILE PASSWORD")
			return
		}
		srv := xc_share_file_server.NewXcFileShareServer(os.Args[2], os.Args[3])
		publicKey, err := srv.Start()
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		fmt.Println("File is shared. PublicKey:")
		fmt.Println(publicKey)
		fmt.Println("Press Enter to exit")
		fmt.Scanln()
		return
	}

	if command == "get" {
		if len(os.Args) != 4 {
			fmt.Println("need address")
			return
		}
		xc_share_file_client.GetFile(os.Args[2], os.Args[3])
		return
	}

	fmt.Println("unknown command", command)
}
