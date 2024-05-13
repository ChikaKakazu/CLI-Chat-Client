package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/ChikaKakazu/CLI-Chat-Client/pb"
	"google.golang.org/grpc"
)

type ChatClient struct {
	c pb.ChatRoomServiceClient
}

func NewChatClient() *ChatClient {
	// TODO: ipアドレスを環境変数から取得する
	conn, err := grpc.Dial("host.docker.internal:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect to server", err)
		return nil
	}

	return &ChatClient{
		c: pb.NewChatRoomServiceClient(conn),
	}
}

func (c *ChatClient) CreateChatRoom(userName string) {
	response, err := c.c.CreateRoom(context.Background(), &pb.CreateRoomRequest{RoomName: "test_room"})
	if err != nil {
		fmt.Println("Failed to create room", err)
		return
	}
	fmt.Println("Room created", response)

	c.JoinChatRoom("test_room", userName)
}

func (c *ChatClient) JoinChatRoom(roomName, userName string) {
	joinRes, err := c.c.JoinRoom(context.Background(), &pb.JoinRoomRequest{RoomName: roomName})
	if err != nil {
		fmt.Println("Failed to join room", err)
	}
	fmt.Println("Joined room", joinRes)

	c.Chat(roomName, userName)
}

func (c *ChatClient) Chat(roomName, userName string) {
	chatRes, err := c.c.Chat(context.Background())
	if err != nil {
		fmt.Println("Failed to chat", err)
		return
	}

	go func() {
		for {
			in, err := chatRes.Recv()
			if err != nil {
				fmt.Println("Failed to receive message", err)
				return
			}
			fmt.Printf("%s: %s\n", in.UserName, in.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		msg := scanner.Text()
		if err := chatRes.Send(&pb.ChatMessage{
			RoomName: roomName,
			UserName: userName,
			Message:  msg,
		}); err != nil {
			fmt.Println("Failed to send message", err)
			return
		}
	}
}

func (c *ChatClient) ListRooms() (*pb.ListRoomsResponse, error) {
	res, err := c.c.ListRooms(context.Background(), &pb.ListRoomsRequest{})
	if err != nil {
		return nil, err
	}

	// for _, roomName := range res.GetRoomNames() {
	// 	fmt.Println("・", roomName)
	// }
	return res, nil
}
