package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/SantiagoZuluaga/fileserver/config"
	"github.com/SantiagoZuluaga/fileserver/models"
)

const MAX_FILE_SIZE = 100000000

func GetFile(channel string, name string, content []byte) error {

	err := os.MkdirAll(fmt.Sprintf("files/%s", channel), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	file, err := os.Create(fmt.Sprintf("./files/%s/%s", channel, name))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func LoadFile(path string) (string, int64, []byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", 0, nil, err
	}

	if fileInfo.Size() > MAX_FILE_SIZE {
		return "", 0, nil, errors.New("max file size is 100mb")
	}

	buffer := make([]byte, fileInfo.Size())
	if _, err := file.Read(buffer); err != nil {
		return "", 0, nil, err
	}
	file.Close()

	return fileInfo.Name(), fileInfo.Size(), buffer, nil
}

func RunTCPClient() {

	connection, err := net.Dial("tcp", fmt.Sprintf("%s:%s", config.GetConfig().Host, config.GetConfig().Port))
	if err != nil {
		fmt.Println(err)
		connection.Close()
		return
	}
	defer connection.Close()

	clientReader := bufio.NewReader(os.Stdin)

	for {

		input, err := clientReader.ReadString('\n')
		if err != nil && err == io.EOF {
			fmt.Println(err)
			continue
		}

		inputParsed := strings.TrimSpace(string(input))
		args := strings.Split(inputParsed, " ")

		command := strings.TrimSpace(args[0])

		switch command {
		case "/username":
			err := json.NewEncoder(connection).Encode(models.ClientMessage{
				Command:   command,
				Argumment: args[1],
			})
			if err != nil {
				fmt.Println(err)
				continue
			}

			var msg models.ServerMessage
			if err := json.NewDecoder(connection).Decode(&msg); err != nil {
				if err == io.EOF {
					fmt.Println("Disconnected from server")
					connection.Close()
					return
				}
				fmt.Println(err)
				continue
			}
			fmt.Println("> ", msg.Message)

		case "/subscribe":
			err := json.NewEncoder(connection).Encode(models.ClientMessage{
				Command:   command,
				Argumment: args[1],
			})
			if err != nil {
				fmt.Println(err)
				continue
			}

			var msg models.ServerMessage
			if err := json.NewDecoder(connection).Decode(&msg); err != nil {
				if err == io.EOF {
					fmt.Println("Disconnected from server")
					connection.Close()
					return
				}
				fmt.Println(err)
				continue
			}
			fmt.Println("> ", msg.Message)

			for {
				var msg models.ServerMessage
				if err := json.NewDecoder(connection).Decode(&msg); err != nil {
					if err == io.EOF {
						fmt.Println("Disconnected from server")
						connection.Close()
						return
					}
					fmt.Println(err)
					continue
				}

				err := GetFile(args[1], msg.File.Name, msg.File.Content)
				if err != nil {
					fmt.Println("Error getting file")
					continue
				}

				fmt.Println("> ", msg.Message)
			}

		case "/channels":
			err := json.NewEncoder(connection).Encode(models.ClientMessage{
				Command: command,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}

			var msg models.ServerMessage
			if err := json.NewDecoder(connection).Decode(&msg); err != nil {
				if err == io.EOF {
					fmt.Println("Disconnected from server")
					connection.Close()
					return
				}
				fmt.Println(err)
				continue
			}
			fmt.Println("> ", msg.Message)

		case "/send":

			name, size, bytes, err := LoadFile(args[2])
			if err != nil {
				fmt.Println(err)
				continue
			}

			if err := json.NewEncoder(connection).Encode(models.ClientMessage{
				Command:   command,
				Argumment: args[1],
				File: &models.File{
					Name:    name,
					Size:    size,
					Content: bytes,
				},
			}); err != nil {
				fmt.Println(err)
				continue
			}

			var msg models.ServerMessage
			if err := json.NewDecoder(connection).Decode(&msg); err != nil {
				if err == io.EOF {
					fmt.Println("Disconnected from server")
					connection.Close()
					return
				}
				fmt.Println(err)
				continue
			}
			fmt.Println("> ", msg.Message)
		case "/exit":
			fmt.Println("Bye, come back soon")
			connection.Close()
			return
		default:
			fmt.Println("Invalid command.\nCommands available:\n1. /username USERNAME\n2. /send CHANNEL FILE_NAME\n3. /suscribe CHANNEL_NAME\n4. /quit")
		}
	}
}
