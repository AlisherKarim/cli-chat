// utils/api.go
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Login sends a POST request to the login endpoint
func Login(username, password string) (string, error) {
    url := "http://localhost:8080/api/v1/auth/login"
    
    // Create the request payload
    payload := map[string]string{
        "username": username,
        "password": password,
    }
    
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("error marshaling JSON: %v", err)
    }

    // Make the POST request
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Check for a successful status code
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body) // Read the response body
        return "", fmt.Errorf("failed to login: %s", body)
    }

    // Read and return the response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(responseBody), nil
}

// Register sends a POST request to the registration endpoint
func Register(username, email, password string) (string, error) {
    url := "http://localhost:8080/api/v1/auth/register" // Update the URL for registration
    
    // Create the request payload
    payload := map[string]string{
        "username": username,
        "email": email,
        "password": password,
    }
    
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("error marshaling JSON: %v", err)
    }

    // Make the POST request
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Check for a successful status code
    if resp.StatusCode != http.StatusCreated { // Assuming 201 is the success code for registration
        body, _ := io.ReadAll(resp.Body) // Read the response body
        return "", fmt.Errorf("failed to register: %s", body)
    }

    // Read and return the response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(responseBody), nil
}


func RequestChatList(username string) (string, error) {
    url := "http://localhost:8080/api/v1/chats"
    
    // Create the GET request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", fmt.Errorf("error creating request: %v", err)
    }
    
    // Set headers
    req.Header.Set("Content-Type", "application/json")
    // req.Header.Set("Authorization", "Bearer your_access_token_here") // Example authorization header
    
    // Create an HTTP client and set a timeout
    client := &http.Client{
        Timeout: 5 * time.Second,
    }
    
    // Send the request
    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error while making GET request: %v", err)
    }
    defer resp.Body.Close()
    
    // Check for a successful status code
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("failed to request: %s", body)
    }
    
    // Read and return the response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response body: %v", err)
    }

    return string(responseBody), nil
}
