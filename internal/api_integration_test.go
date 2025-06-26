// mvp_tests.go – Full coverage for User, Seller, Listing, Order, and Event flow
// -----------------------------------------------------------------------------
// Run: go test ./...
// -----------------------------------------------------------------------------
package internal_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
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
)

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

// ensure we use the same secret everywhere
func init() { _ = os.Setenv("JWT_SECRET", "replace-with-secure-secret") }

// generateTestToken signs a JWT valid for 1 h
func generateTestToken(email string) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	s, _ := tok.SignedString(auth.Secret())
	return s
}

// build a full router with in-memory services wired together
func newRouter() (*gin.Engine, struct {
	usvc user.Service
	ssvc seller.Service
	lsvc listing.Service
	osvc order.Service
}) {
	gin.SetMode(gin.TestMode)

	svcs := struct {
		usvc user.Service
		ssvc seller.Service
		lsvc listing.Service
		osvc order.Service
	}{
		usvc: user.NewInMemoryService(),
		ssvc: seller.NewInMemoryService(),
		lsvc: listing.NewInMemoryService(),
		osvc: order.NewInMemoryService(),
	}

	r := gin.New()
	r.Use(gin.Recovery())

	user.RegisterRoutes(r, svcs.usvc)
	seller.RegisterRoutes(r, svcs.ssvc)
	listing.RegisterRoutes(r, svcs.lsvc)
	order.RegisterRoutes(r, svcs.osvc)

	return r, svcs
}

// -----------------------------------------------------------------------------
// USER MODULE
// -----------------------------------------------------------------------------
func TestUser_RegisterAndLogin(t *testing.T) {
	r, _ := newRouter()

	reg := map[string]string{"id": "u1", "name": "Tom", "email": "tom@ex.com", "password": "pw"}
	body, _ := json.Marshal(reg)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	login := map[string]string{"email": "tom@ex.com", "password": "pw"}
	body, _ = json.Marshal(login)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// -----------------------------------------------------------------------------
// SELLER MODULE
// -----------------------------------------------------------------------------
func TestSeller_CRUD(t *testing.T) {
	r, _ := newRouter()

	s := map[string]any{"id": "s1", "name": "Bob", "email": "bob@ex.com"}
	body, _ := json.Marshal(s)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/sellers/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// duplicate → 409
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/sellers/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}

	// get by ID
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/sellers/s1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// list all
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/sellers", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "s1") {
		t.Fatalf("list all failed: status %d body %s", w.Code, w.Body.String())
	}
}

// -----------------------------------------------------------------------------
// LISTING MODULE
// -----------------------------------------------------------------------------
func TestListing_CRUD(t *testing.T) {
	r, _ := newRouter()

	l := listing.Listing{ID: "l1", SellerID: "s1", Title: "Item", Description: "desc", Price: 9.9, Available: true}
	body, _ := json.Marshal(l)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/listings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	// get by ID
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/listings/l1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// list all
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/listings", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

// -----------------------------------------------------------------------------
// ORDER MODULE + EVENT FLOW
// -----------------------------------------------------------------------------
func TestOrder_CreateAndEventFlow(t *testing.T) {
	// isolated event bus
	bus := make(chan event.Event, 2)
	event.Bus = bus

	r, svcs := newRouter()

	// seed seller & listing
	svcs.ssvc.Register(seller.Seller{ID: "s1", Name: "Bob", Email: "bob@ex.com"})
	svcs.lsvc.Create(listing.Listing{ID: "l1", SellerID: "s1", Title: "Prod", Price: 10.0})

	// register & login user
	svcs.usvc.Register(user.User{ID: "u1", Name: "Tom", Email: "tom@ex.com", Password: "pw"})
	token := generateTestToken("tom@ex.com")

	// create order
	{
		payload := map[string]any{
			"id":         "o1",
			"listingIds": []string{"l1"},
			"sellerId":   "s1",
			"total":      10.0,
		}
		body, _ := json.Marshal(payload)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d body %s", w.Code, w.Body.String())
		}

		select {
		case evt := <-bus:
			if evt.Type != "OrderPlaced" {
				t.Fatalf("expected OrderPlaced, got %s", evt.Type)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for OrderPlaced")
		}
	}

	// accept order
	{
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/orders/o1/accept", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}

		select {
		case evt := <-bus:
			if evt.Type != "OrderAccepted" {
				t.Fatalf("expected OrderAccepted, got %s", evt.Type)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout waiting for OrderAccepted")
		}
	}
}

// -----------------------------------------------------------------------------
// CONCURRENCY – Seller.Register race safety
// -----------------------------------------------------------------------------
func TestSeller_ConcurrentRegister(t *testing.T) {
	svc := seller.NewInMemoryService()
	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("s-%d", i)
			_ = svc.Register(seller.Seller{ID: id, Name: "X", Email: id + "@x.com"})
		}(i)
	}
	wg.Wait()
	if got := len(svc.ListAll()); got != n {
		t.Fatalf("expected %d sellers, got %d", n, got)
	}
}
