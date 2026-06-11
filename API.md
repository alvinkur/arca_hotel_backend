# Hotel Arca — Frontend API Reference

> Base URL: `http://localhost:8080`  
> The gateway proxies to 8 microservices behind a single port. CORS is open (`*`).

## Authentication

### Login — `POST /api/login` *(public)*

```json
// Request
{ "email": "user@example.com", "password": "secret", "role": "customer" }
// role ∈ "customer" | "owner" | "staff"

// Response 200
{ "token": "eyJ...", "user": { "id": 1, "name": "Alice", "email": "alice@example.com", "role": "customer" } }

// Response 400 — validation error
{ "error": "Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag" }

// Response 401 — wrong credentials
{ "error": "Email atau password salah" }
```

After login, pass the token on every request:  
`Authorization: Bearer <token>`

Token expires in **24 hours**. Claims contain `user_id`, `email`, `role`, `name`.

---

## All endpoints require `Authorization` header except `POST /api/login` and `POST /ai-recommend`.

---

## Customers

```
GET    /api/customers          → list
POST   /api/customers          → create
PUT    /api/customers/:id      → update
DELETE /api/customers/:id      → delete
```

```json
// POST/PUT body
{ "name": "Alice", "email": "alice@example.com", "password": "secret", "phone_number": "08123456789" }

// Response (password excluded from GET)
{ "id_customer": 1, "name": "Alice", "email": "alice@example.com", "phone_number": "08123456789" }

// DELETE response
{ "message": "Customer berhasil dihapus" }
```

---

## Owners

```
GET    /api/owners             → list
POST   /api/owners             → create
PUT    /api/owners/:id         → update
DELETE /api/owners/:id         → delete
```

```json
{ "name": "Bob", "email": "bob@example.com", "password": "secret" }
// Response: { "id_owner": 1, "name": "Bob", "email": "bob@example.com" }
```

---

## Staffs

```
GET    /api/staffs             → list
POST   /api/staffs             → create
PUT    /api/staffs/:id         → update
DELETE /api/staffs/:id         → delete
```

```json
{ "name": "Charlie", "email": "charlie@arca.com", "password": "secret" }
// Response: { "id_staff": 1, "name": "Charlie", "email": "charlie@arca.com" }
```

---

## Room Types

```
GET    /api/room-types         → list
POST   /api/room-types         → create
PUT    /api/room-types/:id     → update
DELETE /api/room-types/:id     → delete
```

```json
{ "name": "Deluxe", "price": 750000, "description": "King bed, city view, bathtub" }
// Response: { "id_room_type": 1, "name": "Deluxe", "price": 750000, "description": "..." }
```

---

## Rooms

```
GET    /api/rooms              → list
POST   /api/rooms              → create
PUT    /api/rooms/:id          → update (also used to toggle availability)
DELETE /api/rooms/:id          → delete
```

```json
// POST body
{ "room_number": "301", "id_room_type": 1, "availability": true }

// GET response includes nested room_type
{
  "id_room": 1,
  "room_number": "301",
  "id_room_type": 1,
  "availability": true,
  "room_type": { "id_room_type": 1, "name": "Deluxe", "price": 750000, "description": "..." }
}

// Toggle availability only — PUT with partial body
{ "availability": false }
```

---

## Bookings

```
GET    /api/bookings           → list
POST   /api/bookings           → create
PUT    /api/bookings/:id       → update
DELETE /api/bookings/:id       → delete
```

```json
// POST body
{
  "id_customer": 1,
  "id_room": 2,
  "date_in":  "2026-06-15T14:00:00Z",
  "date_out": "2026-06-17T12:00:00Z",
  "total_payment": 1500000,
  "status_payment": "pending"
}

// Response
{ "id_booking": 1, "id_customer": 1, "id_room": 2, "date_in": "...", "date_out": "...", "total_payment": 1500000, "status_payment": "pending" }
```

**Side effects on create:**
- Validates customer exists (calls auth-service)
- Validates room exists and is available (calls room-service)
- Sets room `availability` to `false` on success

**On delete:** releases the room (sets `availability` back to `true`)

**Status values:** `pending` | `paid` | `cancelled` (default: `pending`)

---

## Payments

```
GET    /api/payments           → list
POST   /api/payments           → create
PUT    /api/payments/:id       → update
DELETE /api/payments/:id       → delete
```

```json
// POST body
{
  "id_booking": 1,
  "total_payment": 1500000,
  "method": "transfer",
  "date": "2026-06-11T10:00:00Z",
  "status": "paid"
}

// Response
{ "id_payment": 1, "id_booking": 1, "total_payment": 1500000, "method": "transfer", "date": "...", "status": "paid" }
```

**Side effect on create:** updates the booking's `status_payment` to `"paid"`.

**Method values:** free-text string (e.g. `transfer`, `cash`, `credit_card`)

---

## Chats

```
GET    /api/chats              → list
POST   /api/chats              → create
PUT    /api/chats/:id          → update
DELETE /api/chats/:id          → delete (cascade-deletes messages)
```

```json
// POST body
{ "id_customer": 1, "id_staff": 1, "date": "2026-06-11T09:00:00Z" }
// Response: { "id_chat": 1, "id_customer": 1, "id_staff": 1, "date": "..." }
```

---

## Chat Messages

```
GET    /api/chats/:id/messages         → list messages for a chat
POST   /api/chats/:id/messages         → send a message
PUT    /api/chats/:id/messages/:msgId  → edit a message
DELETE /api/chats/:id/messages/:msgId  → delete a message
```

```json
// POST body
{ "sender_type": "customer", "message": "Halo, apakah kamar tersedia?" }

// Response
{ "id_chat_message": 1, "id_chat": 1, "sender_type": "customer", "message": "...", "date": "..." }
```

`sender_type` ∈ `"customer"` | `"staff"`

---

## Reviews

```
GET    /api/reviews            → list
POST   /api/reviews            → create
PUT    /api/reviews/:id        → update
DELETE /api/reviews/:id        → delete
```

```json
// POST body
{ "id_customer": 1, "id_room": 2, "rating": 5, "comment": "Kamar bersih, pelayanan ramah!" }

// Response
{ "id_review": 1, "id_customer": 1, "id_room": 2, "rating": 5, "comment": "..." }
```

- `rating`: **1–5** (validated)
- Validates customer and room exist on create

---

## Revenue Reports

```
GET    /api/revenue_reports    → list
POST   /api/revenue_reports    → create
PUT    /api/revenue_reports/:id → update
DELETE /api/revenue_reports/:id → delete
```

```json
// POST body
{ "period": "2026-06", "total_revenue": 15000000, "total_booking": 5, "total_review": 3, "detail_income": "Transfer: 10jt, Cash: 5jt" }

// Response
{ "id_revenue": 1, "period": "2026-06", "total_revenue": 15000000, "total_booking": 5, "total_review": 3, "detail_income": "Transfer: 10jt, Cash: 5jt" }
```

---

## AI Room Recommendation

### `POST /ai-recommend` *(public, no auth)*

```json
// Request
{ "message": "Saya cari kamar murah untuk keluarga" }

// Response
{ "reply": "Berdasarkan preferensi Anda, kami merekomendasikan kamar **Standard** dengan harga Rp350000/malam. ..." }
```

Fallback: if no `AI_API_KEY` is set, the service uses keyword matching (Indonesian). Works offline.

---

## Error shapes

All errors follow this format:

```json
{ "error": "descriptive message in Indonesian" }
```

| HTTP | Meaning |
|---|---|
| 400 | Bad request / validation |
| 401 | Missing or invalid token |
| 404 | Resource not found |
| 409 | Conflict (e.g. room unavailable) |
| 500 | Server error |

---

## Quick reference — ID field names

| Entity | JSON id field |
|---|---|
| Customer | `id_customer` |
| Owner | `id_owner` |
| Staff | `id_staff` |
| Room | `id_room` |
| Room Type | `id_room_type` |
| Booking | `id_booking` |
| Payment | `id_payment` |
| Chat | `id_chat` |
| Chat Message | `id_chat_message` |
| Review | `id_review` |
| Revenue Report | `id_revenue` |

All IDs are unsigned integers (uint).

---

## Minimal integration script

```js
const API = "http://localhost:8080";

let token = null;

async function login(email, password, role) {
  const res = await fetch(`${API}/api/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, role }),
  });
  if (!res.ok) throw new Error((await res.json()).error);
  const data = await res.json();
  token = data.token;
  return data.user;
}

function authHeaders() {
  return { "Content-Type": "application/json", Authorization: `Bearer ${token}` };
}

// CRUD helpers
const api = {
  list:   (r)     => fetch(`${API}/api/${r}`, { headers: authHeaders() }).then(r => r.json()),
  get:    (r, id) => fetch(`${API}/api/${r}/${id}`, { headers: authHeaders() }).then(r => r.json()),
  create: (r, b)  => fetch(`${API}/api/${r}`, { method: "POST", headers: authHeaders(), body: JSON.stringify(b) }).then(r => r.json()),
  update: (r, id, b) => fetch(`${API}/api/${r}/${id}`, { method: "PUT", headers: authHeaders(), body: JSON.stringify(b) }).then(r => r.json()),
  del:    (r, id) => fetch(`${API}/api/${r}/${id}`, { method: "DELETE", headers: authHeaders() }).then(r => r.json()),
};

// Usage
await login("alice@example.com", "secret", "customer");
const rooms = await api.list("rooms");
const booking = await api.create("bookings", { id_customer: 1, id_room: 2, date_in: "2026-06-15T14:00:00Z", date_out: "2026-06-17T12:00:00Z", total_payment: 500000 });
```
