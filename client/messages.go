package client

type SendMessage struct {
	Error   bool
	Message interface{}
}

type ErrorMessage struct {
	Host  string
	Error string
	Name  string
}

type Message struct {
	Host     string
	Name     string
	Platform string
	Data     interface{}
}

// ReceiveMessage : specify the host to filter by
type ReceiveMessage struct {
	FilterBy string
}
