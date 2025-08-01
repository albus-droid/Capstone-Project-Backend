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

// mustDecode unmarshals JSON or fails the test
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
		payload := map[string]string{"name": "Alice Example", "email": email, "password": pass}
		b, _ := json.Marshal(payload)
		res, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(b))
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
		payload := map[string]string{"email": email, "password": pass}
		b, _ := json.Marshal(payload)
		res, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(b))
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
		payload := map[string]string{
			"name":     "Bob’s Burgers",
			"email":    email,
			"phone":    "+1-555-1234",
			"password": pass,
		}
		b, _ := json.Marshal(payload)
		res, err := http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))
		if err != nil {
			t.Fatalf("POST /sellers/register error: %v", err)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("POST /sellers/register: expected 201, got %d", res.StatusCode)
		}
	}

	// 2.1 Duplicate → 409
	{
		payload := map[string]string{"name": "Bob’s Burgers", "email": email, "phone": "+1-555-1234", "password": pass}
		b, _ := json.Marshal(payload)
		res, _ := http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))
		if res.StatusCode != http.StatusConflict {
			t.Fatalf("POST /sellers/register duplicate: expected 409, got %d", res.StatusCode)
		}
	}

	// 2.2 Seller Login
	var loginResp struct{ Token string `json:"token"` }
	{
		payload := map[string]string{"email": email, "password": pass}
		b, _ := json.Marshal(payload)
		res, err := http.Post(baseURL+"/sellers/login", "application/json", bytes.NewReader(b))
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

func TestListings(t *testing.T) {
    // 1) Register a seller
    sellerEmail := fmt.Sprintf("lstsell+%d@ex.com", time.Now().UnixNano())
    reg := map[string]string{
        "name":     "ListSeller",
        "email":    sellerEmail,
        "phone":    "+1000",
        "password": "pw",
    }
    b, _ := json.Marshal(reg)
    res, err := http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))
    if err != nil || res.StatusCode != http.StatusCreated {
        t.Fatalf("POST /sellers/register failed: %v / %d", err, res.StatusCode)
    }

    // 2) Log in to get a token
    var loginResp struct{ Token string `json:"token"` }
    {
        creds := map[string]string{"email": sellerEmail, "password": "pw"}
        b, _ := json.Marshal(creds)
        res, err := http.Post(baseURL+"/sellers/login", "application/json", bytes.NewReader(b))
        if err != nil {
            t.Fatalf("POST /sellers/login failed: %v", err)
        }
        if res.StatusCode != http.StatusOK {
            t.Fatalf("POST /sellers/login: expected 200, got %d", res.StatusCode)
        }
        mustDecode(t, res, &loginResp)
        if loginResp.Token == "" {
            t.Fatal("seller login: missing token")
        }
    }
    authHeader := "Bearer " + loginResp.Token

    // 3) Fetch the seller’s ID
    res, err = http.Get(baseURL + "/sellers")
    if err != nil {
        t.Fatalf("GET /sellers: %v", err)
    }
    var sellers []struct{ ID, Email string }
    mustDecode(t, res, &sellers)
    var sellerID string
    for _, s := range sellers {
        if s.Email == sellerEmail {
            sellerID = s.ID
            break
        }
    }
    if sellerID == "" {
        t.Fatal("could not find new seller in /sellers")
    }

    // 4) POST /listings
    var createResp struct{ Message, ID string }
    {
        payload := map[string]interface{}{
			"sellerId":    sid,         // seller’s UUID string
			"title":       "OrderItem",
			"description": "Freshly made order item",
			"price":       15.0,
			"available":   true,
			"portionSize": 1,           // int, the size of each portion
			"leftSize":    10,          // int, how many portions are available
		}

        b, _ := json.Marshal(payload)
        req, _ := http.NewRequest("POST", baseURL+"/listings", bytes.NewReader(b))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", authHeader)
        res, err := http.DefaultClient.Do(req)
        if err != nil {
            t.Fatalf("POST /listings: %v", err)
        }
        if res.StatusCode != http.StatusCreated {
            t.Fatalf("POST /listings: expected 201, got %d", res.StatusCode)
        }
        mustDecode(t, res, &createResp)
        if createResp.ID == "" {
            t.Fatal("listing create: missing id")
        }
    }

    // 5) GET /listings/:id
    {
        url := fmt.Sprintf("%s/listings/%s", baseURL, createResp.ID)
        req, _ := http.NewRequest("GET", url, nil)
        req.Header.Set("Authorization", authHeader)
        res, err := http.DefaultClient.Do(req)
        if err != nil {
            t.Fatalf("GET /listings/%s: %v", createResp.ID, err)
        }
        if res.StatusCode != http.StatusOK {
            t.Fatalf("GET /listings/%s: expected 200, got %d", createResp.ID, res.StatusCode)
        }
    }

    // 6) GET /listings (all)
    {
        req, _ := http.NewRequest("GET", baseURL+"/listings", nil)
        req.Header.Set("Authorization", authHeader)
        res, err := http.DefaultClient.Do(req)
        if err != nil {
            t.Fatalf("GET /listings: %v", err)
        }
        if res.StatusCode != http.StatusOK {
            t.Fatalf("GET /listings: expected 200, got %d", res.StatusCode)
        }
    }

    // 7) PUT /listings/:id
    {
        update := map[string]interface{}{"price": 5.5, "available": false}
        b, _ := json.Marshal(update)
        url := fmt.Sprintf("%s/listings/%s", baseURL, createResp.ID)
        req, _ := http.NewRequest("PUT", url, bytes.NewReader(b))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", authHeader)
        res, err := http.DefaultClient.Do(req)
        if err != nil {
            t.Fatalf("PUT /listings/%s: %v", createResp.ID, err)
        }
        if res.StatusCode != http.StatusOK {
            t.Fatalf("PUT /listings/%s: expected 200, got %d", createResp.ID, res.StatusCode)
        }
    }

    // 8) DELETE /listings/:id
    {
        url := fmt.Sprintf("%s/listings/%s", baseURL, createResp.ID)
        req, _ := http.NewRequest("DELETE", url, nil)
        req.Header.Set("Authorization", authHeader)
        res, err := http.DefaultClient.Do(req)
        if err != nil {
            t.Fatalf("DELETE /listings/%s: %v", createResp.ID, err)
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
		payload := map[string]string{"name": "OrderUser", "email": email, "password": pass}
		b, _ := json.Marshal(payload)
		http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(b))
	}
	var loginResp struct{ Token string `json:"token"` }
	{
		payload := map[string]string{"email": email, "password": pass}
		b, _ := json.Marshal(payload)
		res, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(b))
		if err != nil {
			t.Fatalf("POST /users/login: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("POST /users/login: expected 200, got %d", res.StatusCode)
		}
		mustDecode(t, res, &loginResp)
	}

	// Register a seller
	sellerEmail := fmt.Sprintf("orderslr+%d@ex.com", time.Now().UnixNano())
	{
		payload := map[string]string{"name": "OrderSeller", "email": sellerEmail, "phone": "+3000", "password": "pw"}
		b, _ := json.Marshal(payload)
		http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))
	}
	res, _ := http.Get(baseURL + "/sellers")
	var sellers []struct{ ID, Email string }
	mustDecode(t, res, &sellers)
	var sid string
	for _, s := range sellers {
		if s.Email == sellerEmail {
			sid = s.ID
		}
	}

	// Create a listing for order with leftSize=10, portionSize=1
	var lr struct{ ID string `json:"id"` }
	{
		payload := map[string]interface{}{
			"sellerId":    sid,
			"title":       "OrderItem",
			"description": "d",
			"price":       15.0,
			"available":   true,
			"portionSize": 1,
			"leftSize":    10,
		}
		b, _ := json.Marshal(payload)
		res, _ := http.Post(baseURL+"/listings", "application/json", bytes.NewReader(b))
		mustDecode(t, res, &lr)
	}

	// Create Order
	var or struct {
		ID         string   `json:"id"`
		UserEmail  string   `json:"user_email"`
		SellerID   string   `json:"sellerId"`
		ListingIDs []string `json:"listingIds"`
		Total      float64  `json:"total"`
		CreatedAt  int64    `json:"createdAt"`
		Status     string   `json:"status"`
	}
	{
		payload := map[string]interface{}{
			"listingIds": []string{lr.ID},
			"sellerId":   sid,
			"total":      15.0,
		}
		b, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+"/orders", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST /orders: %v", err)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("POST /orders: expected 201, got %d", res.StatusCode)
		}
		mustDecode(t, res, &or)
		if or.Status != "pending" {
			t.Fatalf("new order status=%q; want pending", or.Status)
		}
	}

	// Fetch listing before accepting order
	var before, after listingResp
	{
		req, _ := http.NewRequest("GET", baseURL+"/listings/"+lr.ID, nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != http.StatusOK {
			t.Fatalf("GET /listings/%s before accept: %v / %d", lr.ID, err, res.StatusCode)
		}
		mustDecode(t, res, &before)
	}

	// Accept Order
	{
		req, _ := http.NewRequest("PATCH", baseURL+"/orders/"+or.ID+"/accept", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("PATCH /orders/%s/accept: expected 200, got %d", or.ID, res.StatusCode)
		}
	}

	// Fetch listing after accepting order
	{
		req, _ := http.NewRequest("GET", baseURL+"/listings/"+lr.ID, nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, err := http.DefaultClient.Do(req)
		if err != nil || res.StatusCode != http.StatusOK {
			t.Fatalf("GET /listings/%s after accept: %v / %d", lr.ID, err, res.StatusCode)
		}
		mustDecode(t, res, &after)
		if after.LeftSize != before.LeftSize-1 {
			t.Fatalf("leftSize not decremented after accept: before=%d, after=%d", before.LeftSize, after.LeftSize)
		}
	}

	// Check order status after accept
	{
		req, _ := http.NewRequest("GET", baseURL+"/orders/"+or.ID, nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, _ := http.DefaultClient.Do(req)
		mustDecode(t, res, &or)
		if or.Status != "accepted" {
			t.Fatalf("order status after accept=%q; want accepted", or.Status)
		}
	}

	// Complete Order
	{
		req, _ := http.NewRequest("PATCH", baseURL+"/orders/"+or.ID+"/complete", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		res, _ := http.DefaultClient.Do(req)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("PATCH /orders/%s/complete: expected 200, got %d", or.ID, res.StatusCode)
		}
	}

	// Final GET to confirm completion
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