package view

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ChikaKakazu/CLI-Chat-Client/domain/chat"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type View struct {
	App        *tview.Application
	ChatClient *chat.ChatClient
}

func NewView() *View {
	return &View{
		App:        tview.NewApplication(),
		ChatClient: chat.NewChatClient(),
	}
}

func (v *View) Run() {
	v.createUserName()
}

func (v *View) SetRoot(p tview.Primitive, focus bool) *tview.Application {
	return v.App.SetRoot(p, focus)
}

func (v *View) createUserName() {
	// ユーザー名入力画面
	userNameForm := tview.NewForm()
	userNameForm.AddInputField("ユーザー名: ", "", 20, nil, nil)
	userNameForm.AddButton("次へ", func() {
		name := userNameForm.GetFormItem(0).(*tview.InputField).GetText()
		if strings.TrimSpace(name) != "" {
			v.App.SetRoot(v.roomSelectPage(name), true)
		}
	})

	userNameForm.SetBorder(true).SetTitle("ユーザー名を入力してください").SetTitleAlign(tview.AlignCenter)

	// アプリケーションのルートをユーザー名入力画面に設定
	if err := v.App.SetRoot(userNameForm, true).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

func (v *View) roomSelectPage(userName string) tview.Primitive {
	/**************************
	* チャットルーム作成画面
	***************************/
	createRoomForm := tview.NewForm()
	createRoomForm.AddInputField("チャットルーム名: ", "", 30, nil, nil)
	createRoomForm.AddButton("作成", func() {
		roomName := createRoomForm.GetFormItem(0).(*tview.InputField).GetText()
		if strings.TrimSpace(roomName) != "" {
			v.ChatClient.CreateChatRoomAndJoin(v, roomName, userName)
		}
	})
	createRoomForm.SetBorder(true).SetTitle("新しいチャットルームを作成")

	/**************************
	* チャットルーム一覧画面
	***************************/
	// 部屋リストと作成フォーム
	roomsResp, err := v.ChatClient.ListRooms()
	if err != nil {
		log.Fatalf("Failed to list rooms: %v", err)
	}
	rooms := roomsResp.GetRoomNames()

	roomList := tview.NewList()
	for _, room := range rooms {
		roomList.AddItem(room, "", 0, func(roomName string) func() {
			return func() {
				v.ChatClient.JoinChatRoom(v, roomName, userName)
			}
		}(room))
	}
	roomList.SetBorder(true).SetTitle("部屋を選択してください")

	// レイアウト
	flex := tview.NewFlex().
		AddItem(createRoomForm, 0, 1, true).
		AddItem(roomList, 0, 1, false)

	// タブキーでフォーカスの移動
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if v.App.GetFocus() == roomList {
				v.App.SetFocus(createRoomForm)
			} else {
				v.App.SetFocus(roomList)
			}
			return nil
		}
		return event
	})

	return flex
}

func (v *View) ChatPage(c *chat.ChatClient, roomName, userName string) tview.Primitive {
	chatBox := tview.NewTextView().SetDynamicColors(true)
	inputField := tview.NewInputField().SetLabel("メッセージ: ")

	chatClient, err := c.ChatClient.Chat(context.Background())
	if err != nil {
		log.Fatalf("Failed to chat: %v", err)
	}

	c.ReceiveMessages(chatClient, roomName, userName, chatBox, v.App)

	// 最初にメッセージを送信してチャットに参加
	go func() {
		if err := c.SendMessage(chatClient, roomName, userName, fmt.Sprintf("%sさんがチャットに参加しました", userName)); err != nil {
			log.Printf("Failed to send join message: %v", err)
		}
	}()

	// メッセージ送受信
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			message := inputField.GetText()
			if strings.TrimSpace(message) != "" {
				go func() {
					if err := c.SendMessage(chatClient, roomName, userName, message); err != nil {
						log.Printf("Failed to send message: %v", err)
					}
				}()
				inputField.SetText("")
			}
		}
	})

	chatBox.SetBorder(true).SetTitle(roomName)
	inputField.SetBorder(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatBox, 0, 7, false).
		AddItem(inputField, 3, 1, true)

	c.ReceiveMessages(chatClient, roomName, userName, chatBox, v.App)

	go func() {
		if err := c.SendMessage(chatClient, roomName, userName, fmt.Sprintf("%sさんがチャットに参加しました", userName)); err != nil {
			log.Printf("Failed to send join message: %v", err)
		}
	}()

	return flex
}
