# API Endpoint Documentation

This document lists all API endpoints, their parameters, and example requests/responses.

---

## Authentication Header

All *protected* endpoints require:

```
Authorization: Bearer <JWT_TOKEN>
```

---

## 1. Users

### 1.1 Register

* **Endpoint:** `POST /users/register`
* **Description:** Create a new user account.

**Request Body (JSON):**

| Field    | Type   | Required | Description          |
| -------- | ------ | -------- | -------------------- |
| name     | string | yes      | Full name            |
| email    | string | yes      | Unique email address |
| password | string | yes      | Plain-text password  |

**Example Request:**

```http
POST /users/register HTTP/1.1
Content-Type: application/json

{
  "name": "Alice Example",
  "email": "alice@example.com",
  "password": "p@ssw0rd"
}
```

**Example Response** (`201 Created`):

```json
{ "message": "registered" }
```

---

### 1.2 Login

* **Endpoint:** `POST /users/login`
* **Description:** Authenticate and receive a JWT.

**Request Body (JSON):**

| Field    | Type   | Required | Description         |
| -------- | ------ | -------- | ------------------- |
| email    | string | yes      | Registered email    |
| password | string | yes      | Plain-text password |

**Example Request:**

```http
POST /users/login HTTP/1.1
Content-Type: application/json

{
  "email": "alice@example.com",
  "password": "p@ssw0rd"
}
```

**Example Response** (`200 OK`):

```json
{ "token": "eyJhbGciOiJI..." }
```

---

### 1.3 Get Profile (Protected)

* **Endpoint:** `GET /users/profile`
* **Description:** Retrieve the profile of the authenticated user.

**Headers:**

| Header        | Value          |
| ------------- | -------------- |
| Authorization | Bearer <token> |

**Example Request:**

```http
GET /users/profile HTTP/1.1
Authorization: Bearer eyJhbGciOiJI...
```

**Example Response** (`200 OK`):

```json
{
  "id": "user-uuid",
  "name": "Alice Example",
  "email": "alice@example.com"
}
```

---

## 2. Sellers

### 2.1 Register Seller

* **Endpoint:** `POST /sellers/register`
* **Description:** Create a new seller account.

**Request Body:**

| Field    | Type   | Required | Description                 |
| -------- | ------ | -------- | --------------------------- |
| name     | string | yes      | Seller’s display name       |
| email    | string | yes      | Contact email (unique)      |
| phone    | string | yes      | Phone number                |
| password | string | yes      | Plain-text password (min 8) |

**Example Request:**

```http
POST /sellers/register HTTP/1.1
Content-Type: application/json

{
  "name": "Bob’s Burgers",
  "email": "bob@burgers.com",
  "phone": "+1-555-1234",
  "password": "hunter2!"
}
```

**201 Created:**

```json
{ "message": "seller registered" }
```

**409 Conflict:**

```json
{ "error": "seller already exists" }
```

**400 Bad Request:**

```json
{ "error": "invalid request payload" }
```

---

### 2.2 Seller Login

* **Endpoint:** `POST /sellers/login`
* **Description:** Authenticates a seller and returns a JWT.

**Request Body:**

| Field    | Type   | Required | Description          |
| -------- | ------ | -------- | -------------------- |
| email    | string | yes      | Seller email address |
| password | string | yes      | Plain-text password  |

**Example Request:**

```http
POST /sellers/login HTTP/1.1
Content-Type: application/json

{
  "email": "bob@burgers.com",
  "password": "hunter2!"
}
```

**200 OK:**

```json
{ "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..." }
```

**401 Unauthorized:**

```json
{ "error": "invalid credentials" }
```

**400 Bad Request:**

```json
{ "error": "invalid request payload" }
```

---

### 2.3 Get Seller by ID

* **Endpoint:** `GET /sellers/{id}`
* **Description:** Fetch seller’s public details by UUID.

**Path Parameters:**

| Name | Type   | Description |
| ---- | ------ | ----------- |
| id   | string | Seller UUID |

**Example Request:**

```http
GET /sellers/123e4567-e89b-12d3-a456-426614174000 HTTP/1.1
```

**200 OK:**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Bob’s Burgers",
  "email": "bob@burgers.com",
  "phone": "+1-555-1234",
  "verified": false
}
```

**404 Not Found:**

```json
{ "error": "not found" }
```

---

### 2.4 List All Sellers

* **Endpoint:** `GET /sellers`
* **Description:** Retrieve all sellers.

**Example Request:**

```http
GET /sellers HTTP/1.1
```

**Response (200 OK):**

```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Bob’s Burgers",
    "email": "bob@burgers.com",
    "phone": "+1-555-1234",
    "verified": false
  },
  {
    "id": "223e4567-e89b-12d3-a456-426614174001",
    "name": "Alice’s Empanadas",
    "email": "alice@empanadas.com",
    "phone": "+1-555-5678",
    "verified": true
  }
]
```

---

## 3. Listings (Protected)

> All endpoints below require:
>
> ```
> Authorization: Bearer <SELLER_JWT_TOKEN>
> ```

### 3.1 Create Listing

* **Endpoint:** `POST /listings`
* **Description:** Create a new listing. The server generates a unique listing ID.

**Headers:**

| Name          | Value                 |
| ------------- | --------------------- |
| Authorization | Bearer `<SELLER_JWT>` |
| Content-Type  | application/json      |

**Request Body:**

| Field       | Type    | Required | Description             |
| ----------- | ------- | -------- | ----------------------- |
| sellerId    | string  | yes      | Existing Seller UUID    |
| title       | string  | yes      | Listing title           |
| description | string  | yes      | Detailed description    |
| price       | float   | yes      | Price in USD            |
| available   | boolean | yes      | Availability flag       |
| portionSize | int     | yes      | Size of each portion    |
| leftSize    | int     | yes      | Number of portions left |

**Example Request:**

```http
POST /listings HTTP/1.1
Authorization: Bearer eyJhbGciOiJI…
Content-Type: application/json

{
  "sellerId": "seller-uuid",
  "title": "Fresh Apples",
  "description": "Crisp and sweet",
  "price": 2.99,
  "available": true,
  "portionSize": 1,
  "leftSize": 10
}
```

**201 Created:**

```json
{
  "message": "listing created",
  "id": "generated-uuid"
}
```

---

### 3.2 Get Listing by ID

* **Endpoint:** `GET /listings/{id}`
* **Description:** Retrieve a listing by its ID.

**Headers:**

| Name          | Value                 |
| ------------- | --------------------- |
| Authorization | Bearer `<SELLER_JWT>` |

**Example Request:**

```http
GET /listings/abc123-def456 HTTP/1.1
Authorization: Bearer eyJhbGciOiJI…
```

**200 OK:**

```json
{
  "id": "abc123-def456",
  "sellerId": "seller-uuid",
  "title": "Fresh Apples",
  "description": "Crisp and sweet",
  "price": 2.99,
  "available": true,
  "portionSize": 1,
  "leftSize": 10
}
```

---

### 3.3 List Listings (Optional Filter)

* **Endpoint:** `GET /listings`
* **Description:** Get all listings or filter by seller.

**Headers:**

| Name          | Value                 |
| ------------- | --------------------- |
| Authorization | Bearer `<SELLER_JWT>` |

**Query Parameters:**

| Name     | Type   | Description            |
| -------- | ------ | ---------------------- |
| sellerId | string | (optional) Seller UUID |

**Example Request:**

```http
GET /listings?sellerId=seller-uuid HTTP/1.1
Authorization: Bearer eyJhbGciOiJI…
```

**200 OK:**

```json
[
  {
    "id": "abc123-def456",
    "sellerId": "seller-uuid",
    "title": "Fresh Apples",
    "description": "Crisp and sweet",
    "price": 2.99,
    "available": true,
    "portionSize": 1,
    "leftSize": 10
  }
]
```

---

### 3.4 Update Listing

* **Endpoint:** `PUT /listings/{id}`
* **Description:** Update fields of a listing. `{id}` is the listing UUID.

**Headers:**

| Name          | Value                 |
| ------------- | --------------------- |
| Authorization | Bearer `<SELLER_JWT>` |
| Content-Type  | application/json      |

**Request Body:** (any subset)

| Field       | Type    | Description                |
| ----------- | ------- | -------------------------- |
| sellerId    | string  | Change seller (if allowed) |
| title       | string  | New listing title          |
| description | string  | Updated description        |
| price       | float   | New price in USD           |
| available   | boolean | New availability flag      |
| portionSize | int     | Portion size (optional)    |
| leftSize    | int     | Number of portions left    |

**Example Request:**

```http
PUT /listings/abc123-def456 HTTP/1.1
Authorization: Bearer eyJhbGciOiJI…
Content-Type: application/json

{
  "price": 3.49,
  "available": false
}
```

**200 OK:**

```json
{ "message": "listing updated" }
```

---

### 3.5 Delete Listing

* **Endpoint:** `DELETE /listings/{id}`
* **Description:** Remove a listing by its ID.

**Headers:**

| Name          | Value                 |
| ------------- | --------------------- |
| Authorization | Bearer `<SELLER_JWT>` |

**Example Request:**

```http
DELETE /listings/abc123-def456 HTTP/1.1
Authorization: Bearer eyJhbGciOiJI…
```

**Response:** (`204 No Content`)

---


## **3.6 Upload Listing Image (Protected)**

* **Endpoint:** `POST /listings/{id}/image`

* **Description:** Upload an image for a listing. (Saves image URL in listing.)

* **Headers:**

  * `Authorization: Bearer <SELLER_JWT>`
  * `Content-Type: multipart/form-data`

* **Path Parameters:**

  | Name | Type   | Description  |
  | ---- | ------ | ------------ |
  | id   | string | Listing UUID |

* **Form Fields:**

  | Field | Type | Required | Description                                    |
  | ----- | ---- | -------- | ---------------------------------------------- |
  | file  | file | yes      | The image file to upload (e.g. `.jpg`, `.png`) |

**Example Request (using `curl`):**

```bash
curl -X POST \
  -H "Authorization: Bearer <SELLER_JWT>" \
  -F "file=@apples.jpg" \
  http://localhost:8000/listings/abc123-def456/image
```

**Response (`200 OK`):**

```json
{
  "image_url": "/listings/abc123-def456/image/apples.jpg"
}
```

*If your API also updates the Listing’s image field, future GETs of this listing will include the image URL:*

**Example Listing with Image:**

```json
{
  "id": "abc123-def456",
  "sellerId": "seller-uuid",
  "title": "Fresh Apples",
  "description": "Crisp and sweet",
  "price": 2.99,
  "available": true,
  "portionSize": 1,
  "leftSize": 10,
  "image": "/listings/abc123-def456/image/apples.jpg"
}
```

---

## **3.7 Get Signed Image URL (Protected)**

* **Endpoint:** `GET /listings/{id}/image/{filename}`

* **Description:** Get a signed URL for viewing a listing image (valid for 1 hour).

* **Headers:**

  * `Authorization: Bearer <SELLER_JWT>`

* **Path Parameters:**

  | Name     | Type   | Description        |
  | -------- | ------ | ------------------ |
  | id       | string | Listing UUID       |
  | filename | string | Name of image file |

**Example Request:**

```http
GET /listings/abc123-def456/image/apples.jpg HTTP/1.1
Authorization: Bearer eyJhbGciOiJI…
```

**Response (`200 OK`):**

```json
{
  "signed_url": "http://localhost:9000/listing-images/listings/abc123-def456/apples.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=..."
}
```

* Use this signed URL as the `src` in an `<img>` tag in your frontend.

---


## 4. Orders (Protected)

All endpoints below require the `Authorization` header.

### 4.1 Create Order

* **Endpoint:** `POST /orders`
* **Description:** Place a new order.

**Request Body (JSON):**

| Field      | Type      | Required | Description                           |
| ---------- | --------- | -------- | ------------------------------------- |
| id         | string    | no       | Client-supplied Order UUID (optional) |
| listingIds | string\[] | yes      | Array of Listing UUIDs                |
| sellerId   | string    | yes      | Seller UUID                           |
| total      | float     | yes      | Order total in USD                    |

**Example Request:**

```http
POST /orders HTTP/1.1
Authorization: Bearer eyJhbGci...
Content-Type: application/json

{
  "listingIds": ["l1","l2"],
  "sellerId": "seller-uuid",
  "total": 19.98
}
```

**201 Created:**

```json
{
  "id": "order-uuid",
  "user_email": "alice@example.com",
  "sellerId": "seller-uuid",
  "listingIds": ["l1","l2"],
  "total": 19.98,
  "createdAt": 1620000000,
  "status": "pending"
}
```

---

### 4.2 Get Order by ID

* **Endpoint:** `GET /orders/{id}`
* **Description:** Retrieve an order (only the owner).

**Example Request:**

```http
GET /orders/order-uuid HTTP/1.1
Authorization: Bearer eyJhbGci...
```

**Example Response** (`200 OK`):

```json
{
  "id": "order-uuid",
  "user_email": "alice@example.com",
  "sellerId": "seller-uuid",
  "listingIds": ["l1","l2"],
  "total": 19.98,
  "createdAt": 1620000000,
  "status": "accepted"
}
```

---

### 4.3 List My Orders

* **Endpoint:** `GET /orders`
* **Description:** List all orders placed by the authenticated user.

**Example Request:**

```http
GET /orders HTTP/1.1
Authorization: Bearer eyJhbGci...
```

**Example Response** (`200 OK`):

```json
[
  { /* order 1 object */ },
  { /* order 2 object */ }
]
```

---

### 4.4 Accept Order

* **Endpoint:** `PATCH /orders/{id}/accept`
* **Description:** Mark an order as accepted (by seller or owner).

**Example Request:**

```http
PATCH /orders/order-uuid/accept HTTP/1.1
Authorization: Bearer eyJhbGci...
```

**Example Response** (`200 OK`):

```json
{ "message": "order accepted" }
```

---

### 4.5 Complete Order

* **Endpoint:** `PATCH /orders/{id}/complete`
* **Description:** Mark an order as completed.

**Example Request:**

```http
PATCH /orders/order-uuid/complete HTTP/1.1
Authorization: Bearer eyJhbGci...
```

**Example Response** (`200 OK`):

```json
{ "message": "order completed" }
```

---