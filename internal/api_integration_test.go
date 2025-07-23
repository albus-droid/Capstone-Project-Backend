// integration_test.go
// ----------------------------------------------------------------------------
// Run: go test -v -timeout 30s integration_test.go
// ----------------------------------------------------------------------------
package main

import (
	"bytes"
	"encoding/json"
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

func mustDecode[T any](t *testing.T, res *http.Response, out *T) {
	t.Helper()
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		t.Fatalf("decode %T: %v", out, err)
	}
}

// 1. Users
func TestUserEndpoints(t *testing.T) {
	email := "user+" + time.Now().Format("150405") + "@ex.com"
	pass := "pw1234"

	// Register
	reg := map[string]string{"name": "Alice", "email": email, "password": pass}
	b, _ := json.Marshal(reg)
	res, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST /users/register failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("POST /users/register expected 201, got %d", res.StatusCode)
	}

	// Login
	login := map[string]string{"email": email, "password": pass}
	b, _ = json.Marshal(login)
	res, err = http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST /users/login failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("POST /users/login expected 200, got %d", res.StatusCode)
	}
	var lr struct{ Token string `json:"token"` }
	mustDecode(t, res, &lr)
	if lr.Token == "" {
		t.Fatal("login: empty token")
	}

	// Profile
	req, _ := http.NewRequest("GET", baseURL+"/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+lr.Token)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET /users/profile failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET /users/profile expected 200, got %d", res.StatusCode)
	}
	var prof struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	mustDecode(t, res, &prof)
	if prof.Email != email {
		t.Fatalf("profile email = %q; want %q", prof.Email, email)
	}
}

// 2. Sellers
func TestSellerEndpoints(t *testing.T) {
	email := "seller+" + time.Now().Format("150506") + "@ex.com"
	pass := "passw0rd"

	// Register
	reg := map[string]string{"name": "Bob", "email": email, "phone": "+15550001111", "password": pass}
	b, _ := json.Marshal(reg)
	res, err := http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST /sellers/register failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("POST /sellers/register expected 201, got %d", res.StatusCode)
	}

	// Duplicate → 409
	res, _ = http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))
	if res.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate /sellers/register expected 409, got %d", res.StatusCode)
	}

	// Login
	login := map[string]string{"email": email, "password": pass}
	b, _ = json.Marshal(login)
	res, err = http.Post(baseURL+"/sellers/login", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST /sellers/login failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("POST /sellers/login expected 200, got %d", res.StatusCode)
	}
	var sl struct{ Token string `json:"token"` }
	mustDecode(t, res, &sl)
	if sl.Token == "" {
		t.Fatal("seller login: empty token")
	}

	// List all
	res, err = http.Get(baseURL + "/sellers")
	if err != nil {
		t.Fatalf("GET /sellers failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET /sellers expected 200, got %d", res.StatusCode)
	}
	var sellers []struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Verified bool   `json:"verified"`
	}
	mustDecode(t, res, &sellers)
	found := false
	for _, s := range sellers {
		if s.Email == email {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("GET /sellers did not include %s", email)
	}

	// Get by ID
	var sellerID string
	for _, s := range sellers {
		if s.Email == email {
			sellerID = s.ID
		}
	}
	res, err = http.Get(baseURL + "/sellers/" + sellerID)
	if err != nil {
		t.Fatalf("GET /sellers/%s failed: %v", sellerID, err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET /sellers/%s expected 200, got %d", sellerID, res.StatusCode)
	}
	var one struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	mustDecode(t, res, &one)
	if one.Email != email {
		t.Fatalf("GET /sellers/%s email=%s; want %s", sellerID, one.Email, email)
	}
}

// 3. Listings
func TestListingCRUD(t *testing.T) {
	// First register a seller to get a valid sellerId
	email := "listingseller+" + time.Now().Format("150507") + "@ex.com"
	pass := "pw"
	reg := map[string]string{"name": "LStar", "email": email, "phone": "+1000", "password": pass}
	b, _ := json.Marshal(reg)
	http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))
	
	// Get the sellerId
	res, _ := http.Get(baseURL + "/sellers")
	var sellers []struct{ ID, Email string }
	mustDecode(t, res, &sellers)
	var sellerID string
	for _, s := range sellers {
		if s.Email == email {
			sellerID = s.ID
		}
	}

	// Create
	create := map[string]interface{}{
		"sellerId":    sellerID,
		"title":       "Fresh Apples",
		"description": "Crisp and sweet",
		"price":       2.99,
		"available":   true,
	}
	b, _ = json.Marshal(create)
	res, err := http.Post(baseURL+"/listings", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST /listings failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("POST /listings expected 201, got %d", res.StatusCode)
	}
	var lr struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	}
	mustDecode(t, res, &lr)
	if lr.ID == "" {
		t.Fatal("listing create: empty id")
	}

	// Get by ID
	res, err = http.Get(baseURL + "/listings/" + lr.ID)
	if err != nil {
		t.Fatalf("GET /listings/%s failed: %v", lr.ID, err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET /listings/%s expected 200, got %d", lr.ID, res.StatusCode)
	}
	var linfo struct {
		ID        string  `json:"id"`
		SellerID  string  `json:"sellerId"`
		Title     string  `json:"title"`
		Available bool    `json:"available"`
		Price     float64 `json:"price"`
	}
	mustDecode(t, res, &linfo)
	if linfo.SellerID != sellerID {
		t.Fatalf("listing sellerId=%s; want %s", linfo.SellerID, sellerID)
	}

	// List all
	res, err = http.Get(baseURL + "/listings")
	if err != nil {
		t.Fatalf("GET /listings failed: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET /listings expected 200, got %d", res.StatusCode)
	}
	var all []linfo
	mustDecode(t, res, &all)
	found := false
	for _, x := range all {
		if x.ID == lr.ID {
			found = true
		}
	}
	if !found {
		t.Fatalf("GET /listings did not include %s", lr.ID)
	}

	// Update
	upp := map[string]interface{}{"price": 3.49, "available": false}
	b, _ = json.Marshal(upp)
	req, _ := http.NewRequest("PUT", baseURL+"/listings/"+lr.ID, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /listings/%s failed: %v", lr.ID, err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("PUT /listings/%s expected 200, got %d", lr.ID, res.StatusCode)
	}

	// Delete
	req, _ = http.NewRequest("DELETE", baseURL+"/listings/"+lr.ID, nil)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /listings/%s failed: %v", lr.ID, err)
	}
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE /listings/%s expected 204, got %d", lr.ID, res.StatusCode)
	}
}

// 4. Orders
func TestOrderLifecycle(t *testing.T) {
	// Register & login user
	email := "orderuser+" + time.Now().Format("150508") + "@ex.com"
	pass := "pw"
	reg := map[string]string{"name": "OUser", "email": email, "password": pass}
	b, _ := json.Marshal(reg)
	http.Post(baseURL+"/users/register", "application/json", bytes.NewReader(b))
	b, _ = json.Marshal(map[string]string{"email": email, "password": pass})
	res, _ := http.Post(baseURL+"/users/login", "application/json", bytes.NewReader(b))
	var lresp struct{ Token string `json:"token"` }
	mustDecode(t, res, &lresp)

	// Register seller and listing
	sellerEmail := "orderseller+" + time.Now().Format("150509") + "@ex.com"
	sreg := map[string]string{"name": "OSeller", "email": sellerEmail, "phone": "+2000", "password": "pw"}
	b, _ = json.Marshal(sreg)
	http.Post(baseURL+"/sellers/register", "application/json", bytes.NewReader(b))

	// get sellerId
	res, _ = http.Get(baseURL + "/sellers")
	var sellers []struct{ ID, Email string }
	mustDecode(t, res, &sellers)
	var sid string
	for _, s := range sellers {
		if s.Email == sellerEmail {
			sid = s.ID
		}
	}

	// create listing
	lreq := map[string]interface{}{
		"sellerId":    sid,
		"title":       "OrderItem",
		"description": "desc",
		"price":       15.0,
		"available":   true,
	}
	b, _ = json.Marshal(lreq)
	res, _ = http.Post(baseURL+"/listings", "application/json", bytes.NewReader(b))
	var lr struct{ ID string `json:"id"` }
	mustDecode(t, res, &lr)

	// place order
	oreq := map[string]interface{}{
		"listingIds": []string{lr.ID},
		"sellerId":   sid,
		"total":      15.0,
	}
	b, _ = json.Marshal(oreq)
	req, _ := http.NewRequest("POST", baseURL+"/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+lresp.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST /orders failed: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("POST /orders expected 201, got %d", res.StatusCode)
	}
	var or struct {
		ID        string   `json:"id"`
		Status    string   `json:"status"`
		ListingIDs []string `json:"listingIds"`
		SellerID  string   `json:"sellerId"`
	}
	mustDecode(t, res, &or)
	if or.Status != "pending" {
		t.Fatalf("new order status = %q; want pending", or.Status)
	}

	// list my orders
	req, _ = http.NewRequest("GET", baseURL+"/orders", nil)
	req.Header.Set("Authorization", "Bearer "+lresp.Token)
	res, _ = http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET /orders expected 200, got %d", res.StatusCode)
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
		t.Fatalf("GET /orders did not include %s", or.ID)
	}

	// accept order
	req, _ = http.NewRequest("PATCH", baseURL+"/orders/"+or.ID+"/accept", nil)
	req.Header.Set("Authorization", "Bearer "+lresp.Token)
	res, _ = http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("PATCH /orders/%s/accept expected 200, got %d", or.ID, res.StatusCode)
	}

	// verify accepted
	req, _ = http.NewRequest("GET", baseURL+"/orders/"+or.ID, nil)
	req.Header.Set("Authorization", "Bearer "+lresp.Token)
	res, _ = http.DefaultClient.Do(req)
	mustDecode(t, res, &or)
	if or.Status != "accepted" {
		t.Fatalf("order status after accept = %q; want accepted", or.Status)
	}

	// complete order
	req, _ = http.NewRequest("PATCH", baseURL+"/orders/"+or.ID+"/complete", nil)
	req.Header.Set("Authorization", "Bearer "+lresp.Token)
	res, _ = http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("PATCH /orders/%s/complete expected 200, got %d", or.ID, res.StatusCode)
	}

	// verify completed
	req, _ = http.NewRequest("GET", baseURL+"/orders/"+or.ID, nil)
	req.Header.Set("Authorization", "Bearer "+lresp.Token)
	res, _ = http.DefaultClient.Do(req)
	mustDecode(t, res, &or)
	if or.Status != "completed"
