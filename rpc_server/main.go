package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc"
	"os"
)

type Action int

type Path struct {
	Path        string
	Credentials Credentials
}

func (a *Action) ListDir(dirPath *Path, reply *string) error {
	err := checkCredentials(dirPath.Credentials)
	if err != nil {
		return err
	}
	f, err := os.Open(dirPath.Path)
	if err != nil {
		return err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}
	dirList := ""
	for _, file := range files {
		dirList += file.Name() + "\n"
	}
	*reply = dirList
	return nil
}

type File struct {
	Credentials Credentials
	Path string
	Data []byte
}

func(a *Action) UploadFile(file *File, reply *string) error {
	f, err := os.Create(file.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(file.Data)
	if err != nil {
		return err
	}
	*reply = "File uploaded"
	return nil
}

func (a *Action) ViewFile(filepath *Path, reply *[]byte) error {
	err := checkCredentials(filepath.Credentials)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(filepath.Path)
	if err != nil {
		return err
	}
	*reply = data
	return nil
}

type Credentials struct {
	Username string
	Password string
}

type InvalidCredentials struct {}

func (i *InvalidCredentials) Error() string {
	return "Invalid credentials"
}

func checkCredentials(credentials Credentials) error {
	if credentials.Username == "test" && credentials.Password == "test" {
		return nil
	}
	return &InvalidCredentials{}
}

func (a *Action) Auth(credentials Credentials, reply *string) error{
	err := checkCredentials(credentials)
	if err != nil {
		return err
	}
	*reply = "Success"
	return nil
}

func main() {
	action := new(Action)
	server := rpc.NewServer()
	server.Register(action)
	server.HandleHTTP("/", "/debug")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Println(err)
	}
}