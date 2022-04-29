package models

type File struct {
	Name    string
	Size    int64
	Content []byte
}

type ClientMessage struct {
	Command   string
	Argumment string
	File      *File
}

type ServerMessage struct {
	Message string
	File    *File
}
