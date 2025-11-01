package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
)

const (
	listenAddr             = ":8080"
	signInWithEmailLinkURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithEmailLink"
	sendOobCodeURL         = "https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode"
	redirectURL            = "http://localhost:8080/finishSignIn"
	defaultEmail           = "miguel.chavez@theksquaregroup.com"
)

func main() {
	app := fiber.New()
	app.Post("/sendOobCode", fiberSendOobCodeHandler)
	app.Post("/signupPasswordless", fiberSignupPasswordlessHandler)
	app.Get("/finishSignIn", fiberFinishSignInHandler)

	log.Println("Fiber server started on ", listenAddr)
	log.Fatal(app.Listen(listenAddr))
}

func getApiKey() string {
	apiKey := os.Getenv("FIREBASE_IDENTITY_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing FIREBASE_IDENTITY_API_KEY env var")
	}

	return apiKey
}

// Fiber handler for /sendOobCode
func fiberSendOobCodeHandler(c *fiber.Ctx) error {
	email := c.Query("email", "")
	if email == "" {
		email = defaultEmail
	}

	encodedEmail := url.QueryEscape(email)
	continueURL := fmt.Sprintf("%s?email=%s", redirectURL, encodedEmail)
	if err := sendSignInEmail(email, continueURL); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to send OOB code by email: " + err.Error())
	}

	return c.SendString("Email with OOB code sent to: " + email)
}

func sendSignInEmail(email, continueURL string) error {
	firebaseAPIKey := getApiKey()
	url := fmt.Sprintf("%s?key=%s", sendOobCodeURL, firebaseAPIKey)

	reqBody := SendSignInLinkRequest{
		RequestType:     "EMAIL_SIGNIN",
		Email:           email,
		ContinueURL:     continueURL,
		HandleCodeInApp: true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to call Firebase API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("firebase API error: %s", string(body))
	}

	var result SendSignInLinkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse Firebase response: %v", err)
	}

	fmt.Println("Email sign-in link sent to:", result.Email)
	return nil
}

func signInWithEmailLink(email, oobCode string) (*SignInWithEmailLinkResponse, error) {
	firebaseAPIKey := getApiKey()
	url := fmt.Sprintf("%s?key=%s", signInWithEmailLinkURL, firebaseAPIKey)
	reqBody := SignInWithEmailLinkRequest{
		Email:   email,
		OobCode: oobCode,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call Firebase API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("firebase API error: %s", string(body))
	}

	var result SignInWithEmailLinkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse Firebase response: %v", err)
	}

	return &result, nil
}

// Fiber handler for /finishSignIn
func fiberFinishSignInHandler(c *fiber.Ctx) error {
	email := c.Query("email", "")
	oobCode := c.Query("oobCode", "")
	if email == "" {
		email = defaultEmail
		// return c.Status(fiber.StatusBadRequest).SendString("Missing email")
	}

	if oobCode == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing oobCode")
	}

	res, err := signInWithEmailLink(email, oobCode)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("sign-in failed: " + err.Error())
	}

	return c.JSON(res)
}

// Fiber handler for /signupPasswordless
func fiberSignupPasswordlessHandler(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing email")
	}

	encodedEmail := url.QueryEscape(req.Email)
	continueURL := fmt.Sprintf("%s?email=%s", redirectURL, encodedEmail)
	if err := sendSignInEmail(req.Email, continueURL); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to send signup OOB code: " + err.Error())
	}

	return c.SendString("Signup OOB code sent")
}
