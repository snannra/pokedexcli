# Pokedex CLI 🧭

A command-line Pokedex application written in **Go** that allows users to explore the Pokémon world using the **PokeAPI**.  
The project demonstrates clean Go architecture, API integration, caching, concurrency, and stateful CLI design.

---

## 🚀 Features

- **Explore Pokémon locations**
  - Paginated exploration of Pokémon *location areas* using the PokeAPI
  - Navigate forward (`map`) and backward (`mapb`) through the world
- **Inspect Pokémon** (if applicable in your project)
  - Query detailed Pokémon information by name
- **Built-in HTTP response caching**
  - Reduces redundant network requests
  - Improves performance and responsiveness
- **Interactive REPL-style CLI**
  - Persistent state across commands
  - Clean command parsing and dispatch
- **Graceful error handling**
  - Network, JSON parsing, and user input errors handled explicitly

---

## 🛠️ Technologies Used

- **Go (Golang)**
  - Strong typing, explicit error handling, and concurrency primitives
- **PokeAPI**
  - RESTful API used to retrieve Pokémon world data
- **Standard Go Libraries**
  - `net/http` – HTTP client
  - `encoding/json` – JSON decoding
  - `sync` – thread-safe caching
  - `time` – cache expiration and background cleanup
- **Go Modules**
  - Proper module management and dependency isolation

---

### Key Design Decisions

- **Internal Packages (`internal/`)**
  - Prevents external consumers from depending on implementation details
- **Client-based API access**
  - A `pokeapi.Client` encapsulates HTTP logic and caching
- **Stateful CLI config**
  - Pagination state (`next` / `previous`) shared across commands
- **Thread-safe cache**
  - Uses `sync.Mutex` and a background goroutine to evict stale entries

---

## 🧩 Caching Strategy

- Each API URL is used as a **cache key**
- Cached responses store:
  - Raw response bytes
  - Creation timestamp
- A background **reaper goroutine** periodically removes expired entries
- Cache access is protected with a mutex to ensure thread safety

This design minimizes network calls while remaining safe for concurrent access.

---

## 📖 Example Usage

```text
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex > map
sunyshore-city-area
...

Pokedex > mapb
canalave-city-area
...

Pokedex > explore canalave-city-area
