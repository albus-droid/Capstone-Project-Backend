// integration_test.go
// ----------------------------------------------------------------------------
// Run against a live backend (docker-compose up -d must be running)
// ----------------------------------------------------------------------------
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

var baseURL = func() string {
	if u := os.Getenv("API_BASE_URL"); u != "" {
		return u
	}
	return "http://localhost:8080"
}()

func TestUserRegisterLoginProfile(t *testing.T) {
	// Register
	reg := map[string]string{
		"name":     "Test User",
		"email":    "testuser@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(reg)
	res, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Register request failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Register: expected 201, got %d", res.StatusCode)
	}

	// Login
	login := map[string]string{"email": reg["email"], "password": reg["password"]}
	body, _ = json.Marshal(login)
	res, err = http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Login: expected 200, got %d", res.StatusCode)
	}
	var lr map[string]string
	if err := json.NewDecoder(res.Body).Decode(&lr); err != nil {
		t.Fatalf("Login decode failed: %v", err)
	}
	token, ok := lr["token"]
	if !ok || token == "" {
		t.Fatal("Login: missing token")
	}

	// Profile
	req, _ := http.NewRequest("GET", baseURL+"/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Profile request failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Profile: expected 200, got %d", res.StatusCode)
	}
}

func TestSellerFlow(t *testing.T) {
	// Register
	sreg := map[string]string{
		"name":     "Test Seller",
		"email":    "seller@example.com",
		"phone":    "+15550001111",
		"password": "sellerpass",
	}
	body, _ := json.Marshal(sreg)
	res, err := http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Seller register failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Seller register: expected 201, got %d", res.StatusCode)
	}

	// Duplicate → 409
	res, _ = http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(body))
	if res.StatusCode != http.StatusConflict {
		t.Fatalf("Seller duplicate: expected 409, got %d", res.StatusCode)
	}

	// Login
	login := map[string]string{"email": sreg["email"], "password": sreg["password"]}
	body, _ = json.Marshal(login)
	res, err = http.Post(baseURL+"/sellers/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Seller login failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Seller login: expected 200, got %d", res.StatusCode)
	}

	// List all
	res, err = http.Get(baseURL + "/sellers")
	if err != nil {
		t.Fatalf("List sellers failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("List sellers: expected 200, got %d", res.StatusCode)
	}
}

func TestListingCRUD(t *testing.T) {
	// Create listing
	lreq := map[string]any{
		"sellerId":    "00000000-0000-0000-0000-000000000000",
		"title":       "Test Item",
		"description": "Desc",
		"price":       12.34,
		"available":   true,
	}
	body, _ := json.Marshal(lreq)
	res, err := http.Post(baseURL+"/listings", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Create listing failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create listing: expected 201, got %d", res.StatusCode)
	}
	var lr map[string]string
	if err := json.NewDecoder(res.Body).Decode(&lr); err != nil {
		t.Fatalf("Create decode failed: %v", err)
	}
	id, ok := lr["id"]
	if !ok || id == "" {
		t.Fatal("Create listing: missing id")
	}

	// Get by ID
	res, err = http.Get(baseURL + "/listings/" + id)
	if err != nil {
		t.Fatalf("Get listing failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Get listing: expected 200, got %d", res.StatusCode)
	}

	// List all
	res, err = http.Get(baseURL + "/listings")
	if err != nil {
		t.Fatalf("List listings failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("List listings: expected 200, got %d", res.StatusCode)
	}

	// Update
	upp := map[string]any{"price": 9.99}
	body, _ = json.Marshal(upp)
	req, _ := http.NewRequest("PUT", baseURL+"/listings/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Update listing failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Update listing: expected 200, got %d", res.StatusCode)
	}

	// Delete
	req, _ = http.NewRequest("DELETE", baseURL+"/listings/"+id, nil)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Delete listing failed: %v", err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("Delete listing: expected 204, got %d", res.StatusCode)
	}
}

func TestOrderLifecycle(t *testing.T) {
	// Register & login a user
	ureg := map[string]string{"name": "OUser", "email": "ouser@ex.com", "password": "pw"}
	body, _ := json.Marshal(ureg)
	http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
	body, _ = json.Marshal(map[string]string{"email": ureg["email"], "password": ureg["password"]})
	res, _ := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(body))
	var lr map[string]string
	json.NewDecoder(res.Body).Decode(&lr)
	userToken := lr["token"]

	// Create an order
	oreq := map[string]any{
		"listingIds": []string{"l1"},
		"sellerId":   "s1",
		"total":      45.67,
	}
	body, _ = json.Marshal(oreq)
	req, _ := http.NewRequest("POST", baseURL+"/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+userToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Create order failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create order: expected 201, got %d", res.StatusCode)
	}
	var or map[string]any
	json.NewDecoder(res.Body).Decode(&or)
	oid, _ := or["id"].(string)

	// Accept the order
	req, _ = http.NewRequest("PATCH", baseURL+"/orders/"+oid+"/accept", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Accept order failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Accept order: expected 200, got %d", res.StatusCode)
	}

	// Verify status accepted
	req, _ = http.NewRequest("GET", baseURL+"/orders/"+oid, nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Get order failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Get order: expected 200, got %d", res.StatusCode)
	}
	var or2 map[string]any
	json.NewDecoder(res.Body).Decode(&or2)
	if or2["status"] != "accepted" {
		t.Fatalf("Expected status accepted, got %v", or2["status"])
	}

	// Complete the order
	req, _ = http.NewRequest("PATCH", baseURL+"/orders/"+oid+"/complete", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Complete order failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Complete order: expected 200, got %d", res.StatusCode)
	}

	// Verify status completed
	req, _ = http.NewRequest("GET", baseURL+"/orders/"+oid, nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Get order (completed) failed: %v", err)
	}
	json.NewDecoder(res.Body).Decode(&or2)
	if or2["status"] != "completed" {
		t.Fatalf("Expected status completed, got %v", or2["status"])
	}
}
