// integration_test.go
// ----------------------------------------------------------------------------
// Run: go test -v -timeout 30s integration_test.go
// ----------------------------------------------------------------------------
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

var baseURL = func() string {
	if u := os.Getenv("API_BASE_URL"); u != "" {
		return u
	}
	return "http://127.0.0.1:8000"
}()

// decode JSON response into v, fail on error
func mustDecode(t *testing.T, res *http.Response, v interface{}) {
	t.Helper()
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		t.Fatalf("decode %T: %v", v, err)
	}
}

// =================================================================================
// 1. Users
// =================================================================================

func TestUsers(t *testing.T) {
	email := fmt.Sprintf("user+%d@ex.com", time.Now().UnixNano())
	pass := "p@ssw0rd"

	// 1.1 Register
	{
		req := map[string]string{"name": "Alice Example", "email": email, "password": pass}
		body, _ := json.Marshal(req)
		res, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST /users/register error: %v", err)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("POST /users/register: expected 201, got %d", res.StatusCode)
		}
	}

	// 1.2 Login
	var loginResp struct{ Token string `json:"token"` }
	{
		req := map[string]string{"email": email, "password": pass}
		body, _ := json.Marshal(req)
		res, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST /users/login error: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("POST /users/login: expected 200, got %d", res.StatusCode)
		}
		mustDecode(t, res, &loginResp)
		if loginResp.Token == "" {
			t.Fatal("login: missing token")
		}
	}

	// 1.3 Profile
	{
		req, _ := http.NewRequest("GET", baseURL+"/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("GET /users/profile error: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET /users/profile: expected 200, got %d", res.StatusCode)
		}
		var prof struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		mustDecode(t, res, &prof)
		if prof.Email != email {
			t.Fatalf("profile email=%q; want %q", prof.Email, email)
		}
	}
}

// =================================================================================
// 2. Sellers
// =================================================================================

func TestSellers(t *testing.T) {
	email := fmt.Sprintf("seller+%d@ex.com", time.Now().UnixNano())
	pass := "hunter2!"

	// 2.1 Register Seller
	{
		req := map[string]string{
			"name":     "Bob’s Burgers",
			"email":    email,
			"phone":    "+1-555-1234",
			"password": pass,
		}
		body, _ := json.Marshal(req)
		res, err := http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST /sellers/register error: %v", err)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("POST /sellers/register: expected 201, got %d", res.StatusCode)
		}
	}

	// 2.1 Duplicate → 409
	{
		req := map[string]string{"name": "Bob’s Burgers", "email": email, "phone": "+1-555-1234", "password": pass}
		body, _ := json.Marshal(req)
		res, _ := http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(body))
		if res.StatusCode != http.StatusConflict {
			t.Fatalf("POST /sellers/register duplicate: expected 409, got %d", res.StatusCode)
		}
	}

	// 2.2 Seller Login
	var loginResp struct{ Token string `json:"token"` }
	{
		req := map[string]string{"email": email, "password": pass}
		body, _ := json.Marshal(req)
		res, err := http.Post(baseURL+"/sellers/login", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST /sellers/login error: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("POST /sellers/login: expected 200, got %d", res.StatusCode)
		}
		mustDecode(t, res, &loginResp)
		if loginResp.Token == "" {
			t.Fatal("seller login: missing token")
		}
	}

	// 2.4 List All Sellers
	var sellers []struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Verified bool   `json:"verified"`
	}
	{
		res, err := http.Get(baseURL + "/sellers")
		if err != nil {
			t.Fatalf("GET /sellers error: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET /sellers: expected 200, got %d", res.StatusCode)
		}
		mustDecode(t, res, &sellers)
		found := false
		for _, s := range sellers {
			if s.Email == email {
				found = true
			}
		}
		if !found {
			t.Fatalf("GET /sellers did not include %s", email)
		}
	}

	// 2.3 Get Seller by ID
	{
		var sellerID string
		for _, s := range sellers {
			if s.Email == email {
				sellerID = s.ID
			}
		}
		res, err := http.Get(baseURL + "/sellers/" + sellerID)
		if err != nil {
			t.Fatalf("GET /sellers/%s error: %v", sellerID, err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET /sellers/%s: expected 200, got %d", sellerID, res.StatusCode)
		}
		var one struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Verified bool   `json:"verified"`
		}
		mustDecode(t, res, &one)
		if one.ID == "" || one.Email != email {
			t.Fatalf("GET /sellers/%s returned %+v", sellerID, one)
		}
	}
}

// =================================================================================
// 3. Listings
// =================================================================================

type listingInfo struct {
	ID          string  `json:"id"`
	SellerID    string  `json:"sellerId"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
}

func TestListings(t *testing.T) {
	// First register a seller to get an ID
	email := fmt.Sprintf("lstsell+%d@ex.com", time.Now().UnixNano())
	b, _ := json.Marshal(map[string]string{
		"name":     "List Seller",
		"email":    email,
		"phone":    "+1-222-3333",
		"password": "pw",
	})
	http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))

	// fetch seller ID
	res, _ := http.Get(baseURL + "/sellers")
	var sellers []struct{ ID, Email string }
	mustDecode(t, res, &sellers)
	var sellerID string
	for _, s := range sellers {
		if s.Email == email {
			sellerID = s.ID
		}
	}

	// 3.1 Create Listing
	var createResp struct{ Message, ID string }
	{
		req := map[string]interface{}{
			"sellerId":    sellerID,
			"title":       "Fresh Apples",
			"description": "Crisp and sweet",
			"price":       2.99,
			"available":   true,
		}
		body, _ := json.Marshal(req)
		res, err := http.Post(baseURL+"/listings", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("POST /listings error: %v", err)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("POST /listings: expected 201, got %d", res.StatusCode)
		}
		mustDecode(t, res, &createResp)
		if createResp.ID == "" {
			t.Fatal("listing create: missing id")
		}
	}

	// 3.2 Get Listing by ID
	{
		res, err := http.Get(baseURL + "/listings/" + createResp.ID)
		if err != nil {
			t.Fatalf("GET /listings/%s error: %v", createResp.ID, err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET /listings/%s: expected 200, got %d", createResp.ID, res.StatusCode)
		}
		var info listingInfo
		mustDecode(t, res, &info)
		if info.ID != createResp.ID || info.SellerID != sellerID {
			t.Fatalf("GET /listings/%s returned %+v", createResp.ID, info)
		}
	}

	// 3.3 List Listings (no filter)
	var all []listingInfo
	{
		res, err := http.Get(baseURL + "/listings")
		if err != nil {
			t.Fatalf("GET /listings error: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET /listings: expected 200, got %d", res.StatusCode)
		}
		mustDecode(t, res, &all)
		found := false
		for _, L := range all {
			if L.ID == createResp.ID {
				found = true
			}
		}
		if !found {
			t.Fatalf("GET /listings did not include %s", createResp.ID)
		}
	}

	// 3.3 List Listings (filter by sellerId)
	{
		res, err := http.Get(baseURL + "/listings?sellerId=" + sellerID)
		if err != nil {
			t.Fatalf("GET /listings?sellerId error: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET /listings?sellerId: expected 200, got %d", res.StatusCode)
		}
		var filtered []listingInfo
		mustDecode(t, res, &filtered)
		if len(filtered) == 0 {
			t.Fatal("filtered listings returned zero results")
		}
	}

	// 3.4 Update Listing
	{
		payload := map[string]interface{}{"price": 3.49, "available": false}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", baseURL+"/listings/"+createResp.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("PUT /listings/%s error: %v", createResp.ID, err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("PUT /listings/%s: expected 200, got %d", createResp.ID, res.StatusCode)
		}
	}

	// 3.5 Delete Listing
	{
		req, _ := http.NewRequest("DELETE", baseURL+"/listings/"+createResp.ID, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("DELETE /listings/%s error: %v", createResp.ID, err)
		}
		if res.StatusCode != http.StatusNoContent {
			t.Fatalf("DELETE /listings/%s: expected 204, got %d", createResp.ID, res.StatusCode)
		}
	}
}

// =================================================================================
// 4. Orders
// =================================================================================

func TestOrders(t *testing.T) {
	// Register & login a user
	email := fmt.Sprintf("orderusr+%d@ex.com", time.Now().UnixNano())
	pass := "pw"
	{
		b, _ := json.Marshal(map[string]string{"name": "Order User", "email": email, "password": pass})
		http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(b))
	}
	var loginResp struct{ Token string `json:"token"` }
	{
		b, _ := json.Marshal(map[string]string{"email": email, "password": pass})
		res, _ := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(b))
		mustDecode(t, res, &loginResp)
	}

	// Seed seller & listing
	sellerEmail := fmt.Sprintf("orderslr+%d@ex.com", time.Now().UnixNano())
	{
		b, _ := json.Marshal(map[string]string{"name": "Order Seller", "email": sellerEmail, "phone": "+3000", "password": "pw"})
		http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))
	}
	// fetch seller ID
	res, _ := http.Get(baseURL + "/sellers")
	var sellers []struct{ ID, Email string }
	mustDecode(t, res, &sellers)
	var sid string
	for _, s := range sellers {
		if s.Email == sellerEmail {
			sid = s.ID
		}
	}

	// create listing for order
	var lr struct{ ID string `json:"id"` }
	{
		b, _ := json.Marshal(map[string]interface{}{
			"sellerId":    sid,
			"title":       "Order Item",
			"description": "desc",
			"price":       15.0,
			"available":   true,
		})
		res, _ := http.Post(baseURL+"/listings", "application/json", bytes.NewReader(b))
		mustDecode(t, res, &lr)
	}

	// 4.1 Create Order
	var or struct {
		ID        string   `json:"id"`
		UserEmail string   `json:"user_email"`
		SellerID  string   `json:"sellerId"`
		ListingIDs []string `json:"listingIds"`
		Total     float64  `json:"total"`
		CreatedAt int64    `json:"createdAt"`
		Status    string   `json:"status"`
	}
	{
		b, _ := json.Marshal(map[string]interface{}{
			"listingIds": []string{lr.ID},
			"sellerId":   sid,
			"total":      15.0,
		})
		req, _ := http.NewRequest("POST", baseURL+"/orders", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST /orders error: %v", err)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("POST /orders: expected 201, got %d", res.StatusCode)
		}
		mustDecode(t, res, &or)
		if or.Status != "pending" {
			t.Fatalf("new order status=%q; want pending", or.Status)
		}
	}

	// 4.3 List My Orders
	{
		req, _ := http.NewRequest("GET", baseURL+"/orders", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET /orders: expected 200, got %d", res.StatusCode)
		}
		var list []struct{ ID, Status string }
		mustDecode(t, res, &list)
		found := false
		for _, x := range list {
			if x.ID == or.ID {
				found = true
			}
		}
		if !found {
			t.Fatalf("GET /orders missing %s", or.ID)
		}
	}

	// 4.4 Accept Order
	{
		req, _ := http.NewRequest("PATCH", baseURL+"/orders/"+or.ID+"/accept", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("PATCH /orders/%s/accept: expected 200, got %d", or.ID, res.StatusCode)
		}
	}

	// 4.2 Get Order by ID (accepted)
	{
		req, _ := http.NewRequest("GET", baseURL+"/orders/"+or.ID, nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, _ := http.DefaultClient.Do(req)
		mustDecode(t, res, &or)
		if or.Status != "accepted" {
			t.Fatalf("order status after accept=%q; want accepted", or.Status)
		}
	}

	// 4.5 Complete Order
	{
		req, _ := http.NewRequest("PATCH", baseURL+"/orders/"+or.ID+"/complete", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("PATCH /orders/%s/complete: expected 200, got %d", or.ID, res.StatusCode)
		}
	}

	// Get Order by ID (completed)
	{
		req, _ := http.NewRequest("GET", baseURL+"/orders/"+or.ID, nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, _ := http.DefaultClient.Do(req)
		mustDecode(t, res, &or)
		if or.Status != "completed" {
			t.Fatalf("order status after complete=%q; want completed", or.Status)
		}
	}
}
