package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gophkeeper/client/internal/api"
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/models"
	"net/http"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GophKeeper",
	Long:  `Login to the GophKeeper server with your username and password to obtain an authentication token.`,
	Run: func(cmd *cobra.Command, args []string) {
		login, _ := cmd.Flags().GetString("login")
		password, _ := cmd.Flags().GetString("password")

		if login == "" || password == "" {
			fmt.Println("Error: Login and password cannot be empty.")
			cmd.Help()
			return
		}

		user := models.User{
			Login:    login,
			Password: password,
		}

		client := api.NewClient()
		resp, err := client.Request(http.MethodPost, "/api/user/login", user)
		if err != nil {
			fmt.Printf("Error sending login request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			fmt.Printf("Login failed: %s\n", buf.String())
			return
		}

		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Printf("Error decoding login response: %v\n", err)
			return
		}

		token, ok := result["token"]
		if !ok || token == "" {
			fmt.Println("Login failed: no token received.")
			return
		}

		// Store the token locally
		// TODO: Implement secure storage of token, for now just in config.
		err = config.SaveToken(token)
		if err != nil {
			fmt.Printf("Error saving token: %v\n", err)
			return
		}

		fmt.Println("Login successful! Token saved.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("login", "l", "", "User login/username")
	loginCmd.Flags().StringP("password", "p", "", "User password")
	loginCmd.MarkFlagRequired("login")
	loginCmd.MarkFlagRequired("password")
}
