package internal_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/albus-droid/Capstone-Project-Backend/internal/listing"
	"github.com/albus-droid/Capstone-Project-Backend/internal/order"
	"github.com/albus-droid/Capstone-Project-Backend/internal/seller"
	"github.com/albus-droid/Capstone-Project-Backend/internal/user"
)

// setupRouter wires all module routes into one Gin engine.
func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// User
	usvc := user.NewInMemoryService()
	user.RegisterRoutes(r, usvc)

	// Seller
	ssvc := seller.NewInMemoryService()
	seller.RegisterRoutes(r, ssvc)

	// Listing
	lsvc := listing.NewInMemoryService()
	listing.RegisterRoutes(r, lsvc)

	// Order
	osvc := order.NewInMemoryService()
	order.RegisterRoutes(r, osvc)

	return r
}

func TestUserModule(t *testing.T) {
	r := setupRouter()

	// Register
	reqBody := map[string]string{"id":"u1","name":"Test","email":"test@ex.com","password":"pw"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("User register expected 201, got %d", w.Code)
	}

	// Login
	creds := map[string]string{"email":"test@ex.com","password":"pw"}
	body, _ = json.Marshal(creds)
	req = httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("User login expected 200, got %d", w.Code)
	}
}

func TestSellerModule(t *testing.T) {
	r := setupRouter()

	// Register Seller
	sellerBody := map[string]interface{}{"id":"s1","name":"Bob","email":"bob@ex.com","phone":"1234","verified":false}
	body, _ := json.Marshal(sellerBody)
	req := httptest.NewRequest("POST", "/sellers/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Seller register expected 201, got %d", w.Code)
	}

	// Get by ID
	req = httptest.NewRequest("GET", "/sellers/s1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Seller get expected 200, got %d", w.Code)
	}

	// List all
	req = httptest.NewRequest("GET", "/sellers", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Seller list expected 200, got %d", w.Code)
	}
}

func TestListingModule(t *testing.T) {
	r := setupRouter()

	// Create Listing
	l := listing.Listing{ID:"l1",SellerID:"s1",Title:"Item",Description:"desc",Price:9.9,Available:true}
	body, _ := json.Marshal(l)
	req := httptest.NewRequest("POST", "/listings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Listing create expected 201, got %d", w.Code)
	}

	// Get by ID
	req = httptest.NewRequest("GET", "/listings/l1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Listing get expected 200, got %d", w.Code)
	}

	// List all
	req = httptest.NewRequest("GET", "/listings", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Listing list expected 200, got %d", w.Code)
	}
}

func TestOrderModule(t *testing.T) {
	r := setupRouter()

	// Create Order
	o := order.Order{ID:"o1",UserID:"u1",ListingIDs:[]string{"l1"},Total:9.9}
	body, _ := json.Marshal(o)
	req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Order create expected 201, got %d", w.Code)
	}

	// Get by ID
	req = httptest.NewRequest("GET", "/orders/o1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Order get expected 200, got %d", w.Code)
	}

	// List by user
	req = httptest.NewRequest("GET", "/orders?userId=u1", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Order list expected 200, got %d", w.Code)
	}
}

