package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/SantiagoZuluaga/fileserver/config"
	"github.com/SantiagoZuluaga/fileserver/models"
)

type Channel struct {
	Name    string
	Members map[net.Addr]*User
}

type User struct {
	Connection     net.Conn
	Username       string
	CurrentChannel *Channel
}

var channels = make(map[string]*Channel)
var users = make(map[string]*User)

func (channel *Channel) broadcastFile(sender *User, file models.File) {
	for addr, member := range channel.Members {
		if sender.Connection.RemoteAddr() != addr {
			err := json.NewEncoder(member.Connection).Encode(models.ServerMessage{
				Message: fmt.Sprintf("%s sent a file: %s", sender.Username, file.Name),
				File:    &file,
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func (user *User) connectionHandler() {

	users[user.Connection.RemoteAddr().String()] = user
	fmt.Println("User connected: " + user.Connection.RemoteAddr().String())

	for {

		var msg models.ClientMessage
		err := json.NewDecoder(user.Connection).Decode(&msg)
		if err != nil {
			if err == io.EOF {
				fmt.Println("User disconnected: " + user.Connection.RemoteAddr().String())
				delete(users, user.Connection.RemoteAddr().String())
				delete(channels[user.CurrentChannel.Name].Members, user.Connection.RemoteAddr())
				return
			}
			fmt.Println(err)
			continue
		}

		switch msg.Command {
		case "/username":
			fmt.Println("Username")
			if msg.Argumment == "" {
				err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
					Message: "Username is required",
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			}

			err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
				Message: fmt.Sprintf("Welcome: %s", msg.Argumment),
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "/send":
			fmt.Println("Send")
			if msg.Argumment == "" {
				err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
					Message: "Channel is required",
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			}

			channelName := msg.Argumment
			channel := channels[channelName]
			if channel == nil {
				err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
					Message: "Channel not found",
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			}

			err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
				Message: "File received",
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
			channel.broadcastFile(user, *msg.File)

		case "/subscribe":
			fmt.Println("Subscribe")
			if msg.Argumment == "" {
				err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
					Message: "Channel is required",
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			}

			channelName := msg.Argumment
			channel := channels[channelName]
			if channel == nil {
				channel = &Channel{
					Name:    channelName,
					Members: make(map[net.Addr]*User),
				}

				channels[channelName] = channel
			}

			if channel.Members[user.Connection.RemoteAddr()] == nil {

				if user.CurrentChannel != nil {
					delete(channels[user.CurrentChannel.Name].Members, user.Connection.RemoteAddr())
				}

				user.CurrentChannel = channel
				channel.Members[user.Connection.RemoteAddr()] = user

				err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
					Message: fmt.Sprintf("Welcome to channel: %s", channelName),
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
				continue
			}

			err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
				Message: "You are already on this channel",
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "/channels":
			fmt.Println("Channels")
			var channelsTemp []string

			for name := range channels {
				channelsTemp = append(channelsTemp, fmt.Sprintf("- %s (%d users)", name, len(channels[name].Members)))
			}

			err := json.NewEncoder(user.Connection).Encode(models.ServerMessage{
				Message: fmt.Sprintf("Channels available: \n%s", strings.Join(channelsTemp, "\n")),
			})
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func RunTCPServer() {

	fmt.Println("SERVER FILE SERVER")

	server, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.TCP_HOST, config.TCP_PORT))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.Close()

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			return
		}
		defer connection.Close()

		user := &User{
			Connection: connection,
			Username:   "Anonymous",
		}

		go user.connectionHandler()
	}
}
