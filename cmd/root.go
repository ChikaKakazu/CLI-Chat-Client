/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	chatroom "github.com/ChikaKakazu/CLI-Chat-Client/domain/chat_room"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "CLI-ChatApp",
	Short: "A CLI Chat Application",
	Long:  `Hello World`,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.CLI-Chat-Client.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(cmd *cobra.Command, args []string) {
	// ユーザー名の入力
	prompt := promptui.Prompt{
		Label: "Enter your name",
	}
	userName, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed", err)
		return
	}

	fmt.Printf("Hello, %s! Welcome to the chat app.\n", userName)

	// チャットルームのリスト表示または作成
	selectPrompt := promptui.Select{
		Label: "Select a chat room or create a new one",
		Items: []string{"List Chat Rooms", "Create New Chat Room"},
	}
	_, result, err := selectPrompt.Run()
	if err != nil {
		fmt.Println("Select failed", err)
		return
	}

	switch result {
	case "List Chat Rooms":
		fmt.Println("List Chat Rooms")
	case "Create New Chat Room":
		fmt.Println("Create New Chat Room")
		chatroom.CreateChatRoom()
	default:
		fmt.Println("Invalid selection")
	}
}
