package main

import (
	"fmt"
	"os"

	"github.com/ipoluianov/xc-share-file/xc_share_file_client"
	"github.com/ipoluianov/xc-share-file/xc_share_file_server"
)

func main() {
	fmt.Println("XC-FILE-SHARE v" + xc_share_file_server.Version)
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
		err := srv.Start()
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		fmt.Println("Press Enter to exit")
		fmt.Scanln()
		return
	}

	if command == "get" {
		if len(os.Args) < 4 {
			fmt.Println("get addr password [filename]")
			return
		}
		fileName := ""
		if len(os.Args) > 4 {
			fileName = os.Args[4]
		}
		xc_share_file_client.GetFile(os.Args[2], os.Args[3], fileName)
		return
	}

	fmt.Println("unknown command", command)
}
