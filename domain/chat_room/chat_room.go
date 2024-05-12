package chatroom

import (
	"context"
	"fmt"

	"github.com/ChikaKakazu/CLI-Chat-Client/pb"
	"google.golang.org/grpc"
)

func CreateChatRoom() {
	// TODO: ipアドレスを環境変数から取得する
	conn, err := grpc.Dial("host.docker.internal:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect to server", err)
		return
	}
	defer conn.Close()

	c := pb.NewChatRoomServiceClient(conn)

	response, err := c.CreateRoom(context.Background(), &pb.CreateRoomRequest{RoomName: "test_room"})
	if err != nil {
		fmt.Println("Failed to create room", err)
		return
	}
	fmt.Println("Room created", response)

	JoinChatRoom(c)
}

func JoinChatRoom(c pb.ChatRoomServiceClient) {
	joinRes, err := c.JoinRoom(context.Background(), &pb.JoinRoomRequest{RoomName: "test_room"})
	if err != nil {
		fmt.Println("Failed to join room", err)
	}
	fmt.Println("Joined room", joinRes)
}
