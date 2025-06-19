// mvp_tests.go – Full coverage for User, Seller, Listing, Order, and Event flow
// -----------------------------------------------------------------------------
// Run: go test ./...
// -----------------------------------------------------------------------------
package internal_test

import (
    "bytes"
    "encoding/json"
    "fmt"
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "sync"
    "testing"
    "time"

    "github.com/albus-droid/Capstone-Project-Backend/internal/event"
    "github.com/albus-droid/Capstone-Project-Backend/internal/listing"
    "github.com/albus-droid/Capstone-Project-Backend/internal/order"
    "github.com/albus-droid/Capstone-Project-Backend/internal/seller"
    "github.com/albus-droid/Capstone-Project-Backend/internal/user"
    "github.com/gin-gonic/gin"
)

// -----------------------------------------------------------------------------
// Helpers: build a complete Gin router with in‑memory services wired together
// -----------------------------------------------------------------------------
func newRouter() (*gin.Engine, struct {
    usvc user.Service
    ssvc seller.Service
    lsvc listing.Service
    osvc order.Service
}) {
    gin.SetMode(gin.TestMode)

    services := struct {
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

    user.RegisterRoutes(r, services.usvc)
    seller.RegisterRoutes(r, services.ssvc)
    listing.RegisterRoutes(r, services.lsvc)
    order.RegisterRoutes(r, services.osvc)

    return r, services
}

// -----------------------------------------------------------------------------
// USER MODULE
// -----------------------------------------------------------------------------
func TestUser_RegisterAndLogin(t *testing.T) {
    r, _ := newRouter()

    // Register
    reg := map[string]string{"id": "u1", "name": "Tom", "email": "tom@ex.com", "password": "pw"}
    b, _ := json.Marshal(reg)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }

    // Login
    login := map[string]string{"email": "tom@ex.com", "password": "pw"}
    b, _ = json.Marshal(login)
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(b))
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

    // Register
    s := map[string]interface{}{`id`: `s1`, `name`: `Bob`, `email`: `bob@ex.com`}
    b, _ := json.Marshal(s)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/sellers/register", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }

    // Duplicate → 409
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodPost, "/sellers/register", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusConflict {
        t.Fatalf("expected 409, got %d", w.Code)
    }

    // Get by ID
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodGet, "/sellers/s1", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }

    // List all
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodGet, "/sellers", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "s1") {
        t.Fatalf("list all failed, status %d, body %s", w.Code, w.Body.String())
    }
}

// -----------------------------------------------------------------------------
// LISTING MODULE
// -----------------------------------------------------------------------------
func TestListing_CRUD(t *testing.T) {
    r, _ := newRouter()

    l := listing.Listing{ID: "l1", SellerID: "s1", Title: "Item", Description: "desc", Price: 9.9, Available: true}
    b, _ := json.Marshal(l)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/listings", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d", w.Code)
    }

    // Get by ID
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodGet, "/listings/l1", nil)
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }

    // List all
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
    // Create an isolated event bus per test
    bus := make(chan event.Event, 2)
    event.Bus = bus

    r, services := newRouter()

    // Pre-seed a seller & listing
    services.ssvc.Register(seller.Seller{ID: "s1", Name: "Bob", Email: "bob@ex.com"})
    services.lsvc.Create(listing.Listing{ID: "l1", SellerID: "s1", Title: "Prod", Price: 10.0})

    // Create order → expect OrderPlaced
    ord := order.Order{ID: "o1", UserID: "u1", SellerID: "s1", ListingIDs: []string{"l1"}, Total: 10.0}
    b, _ := json.Marshal(ord)
    w := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    if w.Code != http.StatusCreated {
        t.Fatalf("expected 201, got %d, body %s", w.Code, w.Body.String())
    }

    select {
    case evt := <-bus:
        if evt.Type != "OrderPlaced" {
            t.Fatalf("expected OrderPlaced, got %s", evt.Type)
        }
    case <-time.After(time.Second):
        t.Fatal("timeout waiting for OrderPlaced event")
    }

    // Accept → expect OrderAccepted
    w = httptest.NewRecorder()
    req = httptest.NewRequest(http.MethodPatch, "/orders/o1/accept", nil)
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
        t.Fatal("timeout waiting for OrderAccepted event")
    }
}

// -----------------------------------------------------------------------------
// CONCURRENCY test for Seller.Register (race‑safety)
// -----------------------------------------------------------------------------
func TestSeller_ConcurrentRegister(t *testing.T) {
    svc := seller.NewInMemoryService()
    const n = 10 // reduced to avoid heavy goroutines
    wg := sync.WaitGroup{}
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            id := fmt.Sprintf("s-%d", i)
            _ = svc.Register(seller.Seller{ID: id, Name: "X", Email: id + "@x.com"})
        }(i)
    }
    wg.Wait()
    got := svc.ListAll()
    if len(got) != n {
        t.Fatalf("expected %d sellers, got %d", n, len(got))
    }
}
