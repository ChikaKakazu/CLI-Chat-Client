package chat

import (
	"context"
	"fmt"

	"github.com/ChikaKakazu/CLI-Chat-Client/pb"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
)

type View interface {
	SetRoot(p tview.Primitive, fullscreen bool) *tview.Application
	ChatPage(c *ChatClient, roomName, userName string) tview.Primitive
}

type ChatClient struct {
	ChatClient pb.ChatRoomServiceClient
}

func NewChatClient() *ChatClient {
	// TODO: ipアドレスを環境変数から取得する
	conn, err := grpc.Dial("host.docker.internal:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect to server", err)
		return nil
	}

	return &ChatClient{
		ChatClient: pb.NewChatRoomServiceClient(conn),
	}
}

func (c *ChatClient) CreateChatRoomAndJoin(v View, roomName, userName string) {
	c.CreateChatRoom(roomName, userName)
	// c.JoinChatRoom(v, roomName, userName)
	v.SetRoot(v.ChatPage(c, roomName, userName), true)
}

func (c *ChatClient) CreateChatRoom(roomName, userName string) {
	_, err := c.ChatClient.CreateRoom(context.Background(), &pb.CreateRoomRequest{RoomName: roomName})
	if err != nil {
		fmt.Println("Failed to create room", err)
		return
	}
}

func (c *ChatClient) JoinChatRoom(v View, roomName, userName string) {
	_, err := c.ChatClient.JoinRoom(context.Background(), &pb.JoinRoomRequest{RoomName: roomName})
	if err != nil {
		fmt.Println("Failed to join room", err)
	}

	v.SetRoot(v.ChatPage(c, roomName, userName), true)
}

func (c *ChatClient) SendMessage(chatClient pb.ChatRoomService_ChatClient, roomName, userName, message string) error {
	err := chatClient.Send(&pb.ChatMessage{
		RoomName: roomName,
		UserName: userName,
		Message:  message,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *ChatClient) ReceiveMessages(chatClient pb.ChatRoomService_ChatClient, roomName, userName string,
	chatBox *tview.TextView, app *tview.Application) {
	go func() {
		for {
			in, err := chatClient.Recv()
			if err != nil {
				fmt.Println("Failed to receive message", err)
				return
			}
			app.QueueUpdateDraw(func() {
				chatBox.Write([]byte(in.UserName + " > " + in.Message + "\n"))
				chatBox.ScrollToEnd()
			})
		}
	}()
}

func (c *ChatClient) ListRooms() (*pb.ListRoomsResponse, error) {
	res, err := c.ChatClient.ListRooms(context.Background(), &pb.ListRoomsRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
