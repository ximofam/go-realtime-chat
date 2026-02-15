package chat

type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
	SendAt  int64  `json:"send_at"`
}

const (
	TypeChat           = "chat"
	TypeConnectUser    = "connect_user"
	TypeDisconnectUser = "disconnect_user"
)

type WSResponse struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
