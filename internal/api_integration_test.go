// api_endpoints_test.go
// -----------------------------------------------------------------------------
// Run: go test ./internal
// -----------------------------------------------------------------------------
package internal_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/albus-droid/Capstone-Project-Backend/internal/auth"
	"github.com/albus-droid/Capstone-Project-Backend/internal/event"
	"github.com/albus-droid/Capstone-Project-Backend/internal/listing"
	"github.com/albus-droid/Capstone-Project-Backend/internal/order"
	"github.com/albus-droid/Capstone-Project-Backend/internal/seller"
	"github.com/albus-droid/Capstone-Project-Backend/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ensure we use the same secret everywhere
func init() {
	_ = os.Setenv("JWT_SECRET", "replace-with-secure-secret")
}

// helper to sign a token
func generateTestToken(email string) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	s, _ := t.SignedString(auth.Secret())
	return s
}

// newRouter spins up Gin with GORM+SQLite services
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// open in‑memory SQLite
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// auto‑migrate all tables
	if err := user.Migrate(db); err != nil {
		panic(err)
	}
	if err := seller.Migrate(db); err != nil {
		panic(err)
	}
	if err := listing.Migrate(db); err != nil {
		panic(err)
	}
	if err := order.Migrate(db); err != nil {
		panic(err)
	}

	// wire services
	usvc := user.NewPostgresService(db)
	ssvc := seller.NewPostgresService(db)
	lsvc := listing.NewPostgresService(db)
	osvc := order.NewPostgresService(db)

	r := gin.New()
	r.Use(gin.Recovery())

	user.RegisterRoutes(r, usvc)
	seller.RegisterRoutes(r, ssvc)
	listing.RegisterRoutes(r, lsvc)
	order.RegisterRoutes(r, osvc)

	return r
}

func TestUserEndpoints(t *testing.T) {
	r := newRouter()

	// 1. Register
	reg := map[string]string{
		"name":     "Alice",
		"email":    "alice@ex.com",
		"password": "pw1234",
	}
	b, _ := json.Marshal(reg)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/users/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Register: got %d, body=%s", w.Code, w.Body)
	}

	// 2. Login
	login := map[string]string{"email": "alice@ex.com", "password": "pw1234"}
	b, _ = json.Marshal(login)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Login: got %d, body=%s", w.Code, w.Body)
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	token, ok := resp["token"]
	if !ok || token == "" {
		t.Fatalf("Login: missing token, body=%s", w.Body)
	}

	// 3. Profile
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Profile: expected 200, got %d", w.Code)
	}
}

func TestSellerEndpoints(t *testing.T) {
	r := newRouter()

	// Register
	sellerPayload := map[string]string{
		"name":     "Bob's Burgers",
		"email":    "bob@ex.com",
		"phone":    "+15551234",
		"password": "passw0rd",
	}
	b, _ := json.Marshal(sellerPayload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/sellers/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Seller register: %d", w.Code)
	}

	// Duplicate → 409
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/sellers/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("Seller duplicate: expected 409, got %d", w.Code)
	}

	// Login
	login := map[string]string{"email": "bob@ex.com", "password": "passw0rd"}
	b, _ = json.Marshal(login)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/sellers/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Seller login: got %d", w.Code)
	}
	var tok map[string]string
	json.Unmarshal(w.Body.Bytes(), &tok)
	sellerToken := tok["token"]

	// List all
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/sellers", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("List sellers: %d", w.Code)
	}
}

func TestListingEndpoints(t *testing.T) {
	r := newRouter()

	// create a seller first
	_ = httptest.NewRecorder()
	payload := map[string]string{
		"name":     "X",
		"email":    "x@ex.com",
		"phone":    "000",
		"password": "pw",
	}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/sellers/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), req)

	// Create listing
	create := listing.Listing{
		SellerID:    "00000000-0000-0000-0000-000000000000", // any valid UUID
		Title:       "Item",
		Description: "desc",
		Price:       9.9,
		Available:   true,
	}
	b, _ = json.Marshal(create)
	w := httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/listings", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Create listing: got %d", w.Code)
	}
	var cr map[string]string
	json.Unmarshal(w.Body.Bytes(), &cr)
	listingID := cr["id"]

	// Get by ID
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/listings/"+listingID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Get listing: %d", w.Code)
	}

	// List all
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/listings", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("List listings: %d", w.Code)
	}

	// Update
	update := map[string]any{"price": 5.5}
	b, _ = json.Marshal(update)
	w = httptest.NewRecorder()
	req = httptest.NewRequest("PUT", "/listings/"+listingID, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Update listing: %d", w.Code)
	}

	// Delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest("DELETE", "/listings/"+listingID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("Delete listing: %d", w.Code)
	}
}

func TestOrderEndpoints(t *testing.T) {
	// prepare a fresh bus
	bus := make(chan event.Event, 2)
	event.Bus = bus

	r := newRouter()

	// seed seller & listing
	_ = user.NewPostgresService(nil) // no-op for order flow
	_ = seller.NewPostgresService(nil)
	_ = listing.NewPostgresService(nil)
	// register a user & get token
	usvc := user.NewPostgresService(nil)
	usvc.Register(user.User{Name: "A", Email: "a@ex.com", Password: "pw"})
	token := generateTestToken("a@ex.com")

	// create
	payload := map[string]any{
		"listingIds": []string{"l1"},
		"sellerId":   "s1",
		"total":      42.0,
	}
	b, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Create order: %d", w.Code)
	}

	// receive OrderPlaced
	select {
	case ev := <-bus:
		if ev.Type != "OrderPlaced" {
			t.Fatalf("expected OrderPlaced, got %s", ev.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for OrderPlaced")
	}

	// accept
	w = httptest.NewRecorder()
	req = httptest.NewRequest("PATCH", "/orders/"+ev.Data.(order.Order).ID+"/accept", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Accept order: %d", w.Code)
	}

	// receive OrderAccepted
	select {
	case ev2 := <-bus:
		if ev2.Type != "OrderAccepted" {
			t.Fatalf("expected OrderAccepted, got %s", ev2.Type)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for OrderAccepted")
	}
}
