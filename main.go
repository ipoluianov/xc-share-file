package main

import (
	"fmt"
	"os"
	"xc_share_file/xc_share_file_client"
	"xc_share_file/xc_share_file_server"
)

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
