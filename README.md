# 🛡️ e2e decentralized communication protocol

a secure, end-to-end encrypted communication protocol with full authentication and authorization — enabling private, decentralized messaging without centralized control.

---

## 🚀 getting started

### 1. clone the repository

```bash
git https://github.com/bok1c4/E2E-Decentralized-Communication-Protocol.git
cd into
```

### 2. set up environment variables

create a `.env` file in the root directory.

refer to:

- `config/config.go` for database settings (e.g. `db_host`, `db_port`, `db_user`, etc.)
- `session/session.go` for session-related variables (e.g. session key/secret)

example `.env` snippet:

```env
db_host=localhost
db_port=5432
db_user=postgres
db_pw=yourpassword
db_name=yourdb
db_sslmode=disable
session_key=your-secret-key
```

---

### 3. install go

make sure you have go installed and added to your system `path`.

> [download go](https://go.dev/dl/)

---

### 4. install `templ`

templ is required to generate html components.

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

---

### 5. generate components and run the app

```bash
templ generate
go run main.go
```

---

## 🧩 tech stack

- **go** — backend & session management
- **chi** — router and middleware
- **templ** — html component rendering
- **postgresql** — database
- **pgp** — public key validation and secure identity
- **gorilla sessions** — secure user session handling

---

## 📂 project structure

```bash
.
├── config/         # configuration loading
├── components/     # ui components built with templ
├── db/             # database connection and logic
├── handlers/       # route handlers
├── middleware/     # middleware (auth, headers, etc.)
├── session/        # session management
├── routes/         # route grouping and setup (optional)
├── main.go         # app entry point
└── .env            # environment variables (not committed)
```

---

## 🛡️ license

this project is open source and available under the [mit license](license).
