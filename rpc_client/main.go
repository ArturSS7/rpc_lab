package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"strings"
)

type File struct {
	Path string
	Data []byte
	Credentials Credentials
}

type Path struct {
	Path        string
	Credentials Credentials
}

func printMenu() {
	fmt.Print("RPC Server Menu\n[1] Press \"1\" to list directories\n[2] Press \"2\" to upload file\n[3] Press \"3\" to view file\n[4] Press \"4\" to exit\n")
}

type Credentials struct {
	Username string
	Password string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify server ip as 1 argument")
		return
	}
	serverAddress := os.Args[1]

	client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	reader := bufio.NewReader(os.Stdin)
	creds := Credentials{}
	fmt.Println("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.Replace(username, "\n", "", -1)
	creds.Username = username

	fmt.Println("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.Replace(password, "\n", "", -1)
	creds.Password = password

	var reply string
	err = client.Call("Action.Auth", creds, &reply)
	if err != nil {
		log.Println(err)
		return
	}

	printMenu()
	for {
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "1":
			fmt.Print("Enter directory path\n>> ")
			path, _ := reader.ReadString('\n')
			path = strings.Replace(path, "\n", "", -1)
			dir := Path{
				Path:        path,
				Credentials: creds,
			}
			var reply string
			err := client.Call("Action.ListDir", dir, &reply)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(reply)
		case "2":
			fmt.Print("Enter filename to upload\n>> ")
			filename, _ := reader.ReadString('\n')
			filename = strings.Replace(filename, "\n", "", -1)
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Print("Enter path where file should be uploaded\n>> ")
			path, _ := reader.ReadString('\n')
			path = strings.Replace(path, "\n", "", -1)
			file := File{
				Path: path,
				Data: data,
				Credentials: creds,
			}
			var reply string
			err = client.Call("Action.UploadFile", file, &reply)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(reply)
		case "3":
			fmt.Print("Enter filename\n>> ")
			filename, _ := reader.ReadString('\n')
			filename = strings.Replace(filename, "\n", "", -1)
			var reply []byte
			path := Path{
				Path:        filename,
				Credentials: creds,
			}
			err := client.Call("Action.ViewFile", path, &reply)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Printf("File data: \n%s", reply)
		case "4":
			fmt.Println("Bye!")
			return
		}
	}
}