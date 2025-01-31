package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Nerzal/gocloak/v13"
)

/* TODO:
- make the IP address dynamic
- change to token authorization after first login; refresh token automatically with previous token
- change to config parameters or secrets to mount in K8s or alternatively, switch to certificates
*/

func main() {
	log.Printf("Reading environment variables:\n KEYCLOAK, CLIENTID, REALM, USER, PASSWORD, SECRET, LIGHTHOUSE")
	clientID := os.Getenv("CLIENTID")
	realm := os.Getenv("REALM")
	username := os.Getenv("USER")
	secret := os.Getenv("SECRET")
	password := os.Getenv("PASSWORD")
	keycloak := os.Getenv("KEYCLOAK")
	lighthouse := os.Getenv("LIGHTHOUSE")
	advertiseAddress := os.Getenv("ADV_ADDRESS")
	advertiseName := os.Getenv("ADV_NAME")
	ctx := context.Background()

	log.Print("clientID: " + clientID + "| realm: " + realm + "| username: " + username + "| keycloak: " + keycloak + "| lighthouse: " + lighthouse + "| advertising address: " + advertiseAddress + "| advertising name: " + advertiseName)

	for {

		ref_token, err := LoginUser(username, password, keycloak, secret, clientID, realm, ctx)
		if err != nil {
			log.Print(err)
		}

		localVarPath := lighthouse + "/controller/"

		var localVarPostBody = []byte(`{"name": "` + advertiseName + `", "address": "` + advertiseAddress + `"}`)

		req, err := http.NewRequest("POST", localVarPath, bytes.NewBuffer(localVarPostBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("api_key", ref_token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case 201:
			log.Println("Registered successfully")
		case 202:
			log.Println("Updated successfully")
		default:
			log.Println("Could not register controller:")
			body, _ := io.ReadAll(resp.Body)
			fmt.Println("response code: "+resp.Status+" with body: ", string(body))
		}

		time.Sleep(13 * time.Second)
	}
}

// LoginUser - Logs user into the system
func LoginUser(username string, password string, server string, secret string, clientID string, realm string, ctx context.Context) (string, error) {
	client := gocloak.NewClient(server)
	token, err := client.Login(ctx, clientID, secret, realm, username, password)
	if err != nil {
		fmt.Println("Could not log in to IAM: ", err)
		return "none", err
	} else {
		// filter the refresh token from the token
		refresh_token := token.RefreshToken
		return refresh_token, err
	}
}
