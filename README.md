# API Endpoint Documentation

This document lists all API endpoints, their parameters, and example requests/responses.

---

## Authentication Header

All *protected* endpoints require the following HTTP header:

```
Authorization: Bearer <JWT_TOKEN>
```

---

## 1. Users

### 1.1 Register

* **Endpoint**: `POST /users/register`

* **Description**: Create a new user account.

* **Request Body (JSON)**:

  | Field    | Type   | Required | Description          |
  | -------- | ------ | -------- | -------------------- |
  | name     | string | yes      | Full name            |
  | email    | string | yes      | Unique email address |
  | password | string | yes      | Plain-text password  |

* **Example Request**:

  ```http
  POST /users/register HTTP/1.1
  Content-Type: application/json

  {
    "name": "Alice Example",
    "email": "alice@example.com",
    "password": "p@ssw0rd"
  }
  ```

* **Example Response** (`201 Created`):

  ```json
  { "message": "registered" }
  ```

### 1.2 Login

* **Endpoint**: `POST /users/login`

* **Description**: Authenticate and receive a JWT.

* **Request Body (JSON)**:

  | Field    | Type   | Required | Description         |
  | -------- | ------ | -------- | ------------------- |
  | email    | string | yes      | Registered email    |
  | password | string | yes      | Plain-text password |

* **Example Request**:

  ```http
  POST /users/login HTTP/1.1
  Content-Type: application/json

  {
    "email": "alice@example.com",
    "password": "p@ssw0rd"
  }
  ```

* **Example Response** (`200 OK`):

  ```json
  { "token": "eyJhbGciOiJI..." }
  ```

### 1.3 Get Profile (Protected)

* **Endpoint**: `GET /users/profile`

* **Description**: Retrieve the profile of the authenticated user.

* **Headers**:

  | Header        | Value           |
  | ------------- | --------------- |
  | Authorization | Bearer \<token> |

* **Example Request**:

  ```http
  GET /users/profile HTTP/1.1
  Authorization: Bearer eyJhbGciOiJI...
  ```

* **Example Response** (`200 OK`):

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

- **Endpoint**: `POST /sellers/register`
- **Description**: Create a new seller account.

#### Request Body

| Field    | Type   | Required | Description                     |
| -------- | ------ | -------- | ------------------------------- |
| name     | string | yes      | Seller’s display name           |
| email    | string | yes      | Contact email (must be unique)  |
| phone    | string | yes      | Phone number                    |
| password | string | yes      | Plain-text password (min length 8) |

#### Example Request

```http
POST /sellers/register HTTP/1.1
Content-Type: application/json

{
  "name": "Bob’s Burgers",
  "email": "bob@burgers.com",
  "phone": "+1-555-1234",
  "password": "hunter2!"
}
````

#### Responses

**201 Created**

```json
{
  "message": "seller registered"
}
```

**409 Conflict** (email already exists)

```json
{
  "error": "seller already exists"
}
```

**400 Bad Request** (validation error)

```json
{
  "error": "invalid request payload"
}
```

---

### 2.2 Seller Login

* **Endpoint**: `POST /sellers/login`
* **Description**: Authenticates a seller and returns a JWT token for future authenticated requests.

#### Request Body

| Field    | Type   | Required | Description          |
| -------- | ------ | -------- | -------------------- |
| email    | string | yes      | Seller email address |
| password | string | yes      | Plain-text password  |

#### Example Request

```http
POST /sellers/login HTTP/1.1
Content-Type: application/json

{
  "email": "bob@burgers.com",
  "password": "hunter2!"
}
```

#### Responses

**200 OK**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}
```

**401 Unauthorized** (invalid credentials)

```json
{
  "error": "invalid credentials"
}
```

**400 Bad Request** (invalid request format)

```json
{
  "error": "invalid request payload"
}
```

---

### 2.3 Get Seller by ID

* **Endpoint**: `GET /sellers/{id}`
* **Description**: Fetch a seller’s public details by their unique ID.

#### Path Parameters

| Name | Type   | Description |
| ---- | ------ | ----------- |
| id   | string | Seller UUID |

#### Example Request

```http
GET /sellers/123e4567-e89b-12d3-a456-426614174000 HTTP/1.1
```

#### Responses

**200 OK**

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Bob’s Burgers",
  "email": "bob@burgers.com",
  "phone": "+1-555-1234",
  "verified": false
}
```

**404 Not Found**

```json
{
  "error": "not found"
}
```

---

### 2.4 List All Sellers

* **Endpoint**: `GET /sellers`
* **Description**: Retrieve a list of all sellers, sorted by ID.

#### Example Request

```http
GET /sellers HTTP/1.1
```

#### Response (200 OK)

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

## 3. Listings

### 3.1 Create Listing

* **Endpoint**: `POST /listings`

* **Description**: Create a new listing. The server generates a unique listing ID.

* **Request Body (JSON)**:

  | Field       | Type    | Required | Description          |
  | ----------- | ------- | -------- | -------------------- |
  | sellerId    | string  | yes      | Existing Seller UUID |
  | title       | string  | yes      | Listing title        |
  | description | string  | yes      | Detailed description |
  | price       | float   | yes      | Price in USD         |
  | available   | boolean | yes      | Availability flag    |

> **Note:** Do **not** include an `id` field in the request body. The ID is generated by the server.

* **Example Request**:

  ```http
  POST /listings HTTP/1.1
  Content-Type: application/json

  {
    "sellerId": "seller-uuid",
    "title": "Fresh Apples",
    "description": "Crisp and sweet",
    "price": 2.99,
    "available": true
  }```

* **Example Response** (`201 Created`):

  ```json
  {
    "message": "listing created",
    "id": "generated-uuid"
  }
  ```

---

### 3.2 Get Listing by ID

* **Endpoint**: `GET /listings/{id}`

* **Description**: Fetch a single listing by its ID.

* **Path Parameters**:

  | Name | Type   | Description  |
  | ---- | ------ | ------------ |
  | id   | string | Listing UUID |

* **Example Request**:

  ```http
  GET /listings/abc123-def456 HTTP/1.1
  ```

* **Example Response** (`200 OK`):

  ```json
  {
    "id": "abc123-def456",
    "sellerId": "seller-uuid",
    "title": "Fresh Apples",
    "description": "Crisp and sweet",
    "price": 2.99,
    "available": true
  }
  ```

---

### 3.3 List Listings (Optional Filter)

* **Endpoint**: `GET /listings`

* **Description**: Get all listings, or only those for a specific seller.

* **Query Parameters**:

  | Name     | Type   | Description            |
  | -------- | ------ | ---------------------- |
  | sellerId | string | (optional) Seller UUID |

* **Example Requests**:

  ```http
  GET /listings HTTP/1.1
  ```

  ```http
  GET /listings?sellerId=seller-uuid HTTP/1.1
  ```

* **Example Response** (`200 OK`):

  ```json
  [
    {
      "id": "abc123",
      "sellerId": "seller-uuid",
      "title": "Fresh Apples",
      "description": "Crisp and sweet",
      "price": 2.99,
      "available": true
    },
    {
      "id": "def456",
      "sellerId": "seller-uuid",
      "title": "Organic Honey",
      "description": "Local wildflower",
      "price": 9.50,
      "available": false
    }
  ]
``` ```

---

### 3.4 Update Listing

* **Endpoint**: `PUT /listings/{id}`

* **Description**: Update one or more fields of an existing listing.

* **Path Parameters**:

  | Name | Type   | Description  |
  | ---- | ------ | ------------ |
  | id   | string | Listing UUID |

* **Request Body (JSON)**: include any subset of updatable fields. The `id` in the URL is authoritative.

  | Field       | Type    | Required | Description                |
  | ----------- | ------- | -------- | -------------------------- |
  | sellerId    | string  | no       | Change seller (if allowed) |
  | title       | string  | no       | New listing title          |
  | description | string  | no       | Updated description        |
  | price       | float   | no       | New price in USD           |
  | available   | boolean | no       | New availability flag      |

* **Example Request**:

  ```http
  PUT /listings/abc123-def456 HTTP/1.1
  Content-Type: application/json

  {
    "price": 3.49,
    "available": false
  }
  ```

* **Example Response** (`200 OK`):

  ```json
  {
    "message": "listing updated"
  }
  ```

---

### 3.5 Delete Listing

* **Endpoint**: `DELETE /listings/{id}`

* **Description**: Remove a listing by its ID.

* **Path Parameters**:

  | Name | Type   | Description  |
  | ---- | ------ | ------------ |
  | id   | string | Listing UUID |

* **Example Request**:

  ```http
  DELETE /listings/abc123-def456 HTTP/1.1
  ```

* **Example Response** (`204 No Content`):
  No body.

```
---

## 4. Orders (Protected)

All endpoints below require the `Authorization` header.

### 4.1 Create Order

* **Endpoint**: `POST /orders`

* **Description**: Place a new order.

* **Request Body (JSON)**:

  | Field      | Type      | Required | Description                           |
  | ---------- | --------- | -------- | ------------------------------------- |
  | id         | string    | no       | Client-supplied Order UUID (optional) |
  | listingIds | string\[] | yes      | Array of Listing UUIDs                |
  | sellerId   | string    | yes      | Seller UUID                           |
  | total      | float     | yes      | Order total in USD                    |

* **Example Request**:

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

* **Example Response** (`201 Created`):

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

### 4.2 Get Order by ID

* **Endpoint**: `GET /orders/{id}`

* **Description**: Retrieve an order (only the owner).

* **Path Parameters**:

  | Name | Type   | Description |
  | ---- | ------ | ----------- |
  | id   | string | Order UUID  |

* **Example Request**:

  ```http
  GET /orders/order-uuid HTTP/1.1
  Authorization: Bearer eyJhbGci...
  ```

* **Example Response** (`200 OK`):

  ```json
  { /* Order object */ }
  ```

### 4.3 List My Orders

* **Endpoint**: `GET /orders`

* **Description**: List all orders placed by the authenticated user.

* **Example Request**:

  ```http
  GET /orders HTTP/1.1
  Authorization: Bearer eyJhbGci...
  ```

* **Example Response** (`200 OK`):

  ```json
  [
    { /* order 1 */ },
    { /* order 2 */ }
  ]
  ```

### 4.4 Accept Order

* **Endpoint**: `PATCH /orders/{id}/accept`

* **Description**: Mark an order as accepted (by seller or owner).

* **Path Parameters**:

  | Name | Type   | Description |
  | ---- | ------ | ----------- |
  | id   | string | Order UUID  |

* **Example Request**:

  ```http
  PATCH /orders/order-uuid/accept HTTP/1.1
  Authorization: Bearer eyJhbGci...
  ```

* **Example Response** (`200 OK`):

  ```json
  { "message": "order accepted" }
  ```

### 4.5 Complete Order

* **Endpoint**: `PATCH /orders/{id}/complete`

* **Description**: Mark an order as completed.

* **Path Parameters**:

  | Name | Type   | Description |
  | ---- | ------ | ----------- |
  | id   | string | Order UUID  |

* **Example Request**:

  ```http
  PATCH /orders/order-uuid/complete HTTP/1.1
  Authorization: Bearer eyJhbGci...
  ```

* **Example Response** (`200 OK`):

  ```json
  { "message": "order completed" }
  ```
