package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"main/data"
	"main/tools"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	var choice int
	for {
		fmt.Println("Main Menu")
		fmt.Println("1. Get message")
		fmt.Println("2. Send file")
		fmt.Println("3. Quit")
		fmt.Print(">> ")
		fmt.Scanf("%d\n", &choice)
		if choice == 1 {
			getMessage()
		} else if choice == 2 {
			sendFile()
		} else if choice == 3 {
			break
		} else {
			fmt.Println("Invalid choice")
		}
	}
}

func getMessage() {
	client := createTLSClient()
	resp, err := client.Get("https://localhost:9876")
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	tools.ErrorHandler(err)

	fmt.Println("Server : ", string(data))
}

func sendFile() {
	var name string
	var age int

	scanner := bufio.NewReader(os.Stdin)

	fmt.Print("Input name : ")
	name, _ = scanner.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Input Age: ")
	fmt.Scanf("%d\n", &age)

	person := data.Person{Name: name, Age: age}
	jsonData, err := json.Marshal(person)
	tools.ErrorHandler(err)

	temp := new(bytes.Buffer)
	w := multipart.NewWriter(temp)

	personField, err := w.CreateFormField("Person")
	tools.ErrorHandler(err)

	_, err = personField.Write(jsonData)
	tools.ErrorHandler(err)

	file, err := os.Open("./file.txt")
	tools.ErrorHandler(err)
	defer file.Close()

	fileField, err := w.CreateFormFile("File", file.Name())
	tools.ErrorHandler(err)

	_, err = io.Copy(fileField, file)
	tools.ErrorHandler(err)

	err = w.Close()
	tools.ErrorHandler(err)

	req, err := http.NewRequest("POST", "https://localhost:9876/sendFile", temp)
	tools.ErrorHandler(err)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := createTLSClient()
	resp, err := client.Do(req)
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	tools.ErrorHandler(err)

	fmt.Println("Server : ", string(data))
}

func createTLSClient() *http.Client {
	certPool := x509.NewCertPool()
	cert, err := os.ReadFile("../cert.pem")
	if err != nil {
		fmt.Println("Error reading cert.pem:", err)
		os.Exit(1)
	}
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		fmt.Println("Failed to append cert to certPool")
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := tls.Dial(network, addr, tlsConfig)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	return &http.Client{
		Transport: transport,
	}
}
