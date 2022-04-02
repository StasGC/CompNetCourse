package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
	"os"
)

const host = "127.0.0.1"
const port = ":21"
const user = "TestUser"
const password = "123456"

func main() {
	client, err := ftp.Dial(host + port)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Login(user, password)
	if err != nil {
		log.Fatal(err)
	}

	fileName := "/Temp/FTP/test_pushing.txt"
	localPathPush := "test_push.txt"
	localPathGet := "test_get.txt"

	if err := GetDirectoryListing(client, "/Temp/FTP"); err != nil {
		log.Fatal(err)
	}

	if err := PushFile(client, localPathPush, fileName); err != nil {
		log.Fatal(err)
	}

	if err := RequestFile(client, fileName, localPathGet); err != nil {
		log.Fatal(err)
	}

	if err := client.Quit(); err != nil {
		log.Fatal(err)
	}
}

func PushFile(client *ftp.ServerConn, localPath string, fileName string) error {
	f, err := os.Open(localPath)
	if err != nil {
		return err
	}

	if err := client.Stor(fileName, f); err != nil {
		return err
	}
	return nil
}

func RequestFile(client *ftp.ServerConn, filePath string, localPath string) error {
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}

	reader, err := client.Retr(filePath)
	if err != nil {
		return err
	}

	buf := make([]byte, 2048)
	n, err := reader.Read(buf)
	if err != nil {
		return err
	}
	buf = buf[:n]
	fmt.Print(string(buf))

	if _, err := f.Write(buf); err != nil {
		return err
	}

	return nil
}

func GetDirectoryListing(client *ftp.ServerConn, path string) error {
	entries, err := client.List(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fmt.Printf("%20s	%10s	%v\n", entry.Name, entry.Type, entry.Time)
	}

	return nil
}
