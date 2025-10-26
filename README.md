# String Analyzer API

A RESTful API service that analyzes strings and computes their properties including length, palindrome detection, SHA-256 hash, character frequency, and more.

## Features

- ✨ Analyze strings and compute multiple properties
- 🔐 SHA-256 hash-based unique identification  
- 🔍 Advanced filtering capabilities
- 🤖 Natural language query support
- 🔒 Thread-safe in-memory storage
- 🚀 Zero external dependencies

## Tech Stack

- **Language:** Go 1.21+
- **HTTP:** Standard library (net/http)
- **Storage:** In-memory with sync.RWMutex

## Project Structure

```
string-analyzer/
├── main.go          # All-in-one implementation
├── go.mod           # Go module file
├── README.md        # This file
└── .env.example     # Environment variables template
```

## Local Setup

### Prerequisites
- Go 1.21 or higher installed

### Installation

1. **Clone the repository:**
```bash
git clone https://github.com/yourusername/string-analyzer.git
cd string-analyzer
```

2. **Initialize Go module:**
```bash
go mod init github.com/yourusername/string-analyzer
go mod tidy
```

3. **Run the application:**
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Build Binary

To create a standalone executable:
```bash
go build -o string-analyzer
./string-analyzer
```

## Environment Variables

- `PORT`: Server port (default: 8080)

Create a `.env` file (optional):
```
PORT=8080
```

## API Endpoints

### 1. Create/Analyze String

**Endpoint:** `POST /strings`

**Request:**
```json
{
  "value": "hello world"
}
```

**Response (201 Created):**
```json
{
  "id": "abc123...",
  "value": "hello world",
  "properties": {
    "length": 11,
    "is_palindrome": false,
    "unique_characters": 8,
    "word_count": 2,
    "sha256_hash": "abc123...",
    "character_frequency_map": {
      "h": 1,
      "e": 1,
      "l": 3,
      "o": 2,
      " ": 1,
      "w": 1,
      "r": 1,
      "d": 1
    }
  },
  "created_at": "2025-10-21T10:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or missing "value" field
- `409 Conflict`: String already exists
- `422 Unprocessable Entity`: Invalid data type

---

### 2. Get Specific String

**Endpoint:** `GET /strings/{string_value}`

**Example:**
```bash
GET /strings/hello%20world
```

**Response (200 OK):**
```json
{
  "id": "abc123...",
  "value": "hello world",
  "properties": { ... },
  "created_at": "2025-10-21T10:00:00Z"
}
```

**Error Response:**
- `404 Not Found`: String does not exist

---

### 3. Get All Strings with Filters

**Endpoint:** `GET /strings`

**Query Parameters:**
- `is_palindrome`: boolean (true/false)
- `min_length`: integer (minimum string length)
- `max_length`: integer (maximum string length)
- `word_count`: integer (exact word count)
- `contains_character`: string (single character)

**Examples:**
```bash
GET /strings?is_palindrome=true
GET /strings?min_length=5&max_length=20
GET /strings?word_count=2&contains_character=a
GET /strings?is_palindrome=true&min_length=5
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "hash1",
      "value": "racecar",
      "properties": { ... },
      "created_at": "2025-10-21T10:00:00Z"
    }
  ],
  "count": 1,
  "filters_applied": {
    "is_palindrome": true,
    "min_length": 5
  }
}
```

---

### 4. Natural Language Filtering

**Endpoint:** `GET /strings/filter-by-natural-language`

**Query Parameter:**
- `query`: Natural language string describing filters

**Examples:**
```bash
GET /strings/filter-by-natural-language?query=single%20word%20palindromes
GET /strings/filter-by-natural-language?query=strings%20longer%20than%2010%20characters
GET /strings/filter-by-natural-language?query=containing%20letter%20z
GET /strings/filter-by-natural-language?query=palindromic%20strings%20with%20first%20vowel
```

**Supported Query Patterns:**
- "single word palindromes" → `word_count=1, is_palindrome=true`
- "strings longer than 10 characters" → `min_length=11`
- "strings shorter than 20 characters" → `max_length=19`
- "containing letter z" → `contains_character=z`
- "first vowel" → `contains_character=a`

**Response (200 OK):**
```json
{
  "data": [ ... ],
  "count": 3,
  "interpreted_query": {
    "original": "single word palindromes",
    "parsed_filters": {
      "word_count": 1,
      "is_palindrome": true
    }
  }
}
```

**Error Responses:**
- `400 Bad Request`: Missing or invalid query parameter

---

### 5. Delete String

**Endpoint:** `DELETE /strings/{string_value}`

**Example:**
```bash
DELETE /strings/hello%20world
```

**Response:** `204 No Content` (empty body)

**Error Response:**
- `404 Not Found`: String does not exist

---

## Testing Examples

### Using cURL

**Create a palindrome:**
```bash
curl -X POST http://localhost:8080/strings \
  -H "Content-Type: application/json" \
  -d '{"value": "racecar"}'
```

**Create multiple strings:**
```bash
curl -X POST http://localhost:8080/strings -H "Content-Type: application/json" -d '{"value": "A man a plan a canal Panama"}'
curl -X POST http://localhost:8080/strings -H "Content-Type: application/json" -d '{"value": "hello world"}'
curl -X POST http://localhost:8080/strings -H "Content-Type: application/json" -d '{"value": "noon"}'
```

**Get all palindromes:**
```bash
curl "http://localhost:8080/strings?is_palindrome=true"
```

**Get single-word strings:**
```bash
curl "http://localhost:8080/strings?word_count=1"
```

**Natural language query:**
```bash
curl "http://localhost:8080/strings/filter-by-natural-language?query=single%20word%20palindromes"
```

**Get specific string:**
```bash
curl http://localhost:8080/strings/racecar
```

**Delete string:**
```bash
curl -X DELETE http://localhost:8080/strings/racecar
```

### Using Postman

1. Import the following as a collection
2. Set `{{base_url}}` to `http://localhost:8080`

**Collection:**
- POST `{{base_url}}/strings` with JSON body
- GET `{{base_url}}/strings`
- GET `{{base_url}}/strings/racecar`
- GET `{{base_url}}/strings?is_palindrome=true`
- DELETE `{{base_url}}/strings/racecar`

---

## Deploy to Railway

### Quick Deploy

1. **Push to GitHub:**
```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin YOUR_REPO_URL
git push -u origin main
```

2. **Deploy on Railway:**
   - Go to [railway.app](https://railway.app)
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository
   - Railway auto-detects Go and deploys!

3. **Access your API:**
   - Railway provides a URL like: `https://your-project.up.railway.app`
   - Test: `https://your-project.up.railway.app/health`

### Railway Configuration (Optional)

Create `railway.json`:
```json
{
  "build": {
    "builder": "nixpacks"
  },
  "deploy": {
    "startCommand": "./main",
    "restartPolicyType": "always"
  }
}
```

Or create `Procfile`:
```
web: ./string-analyzer
```

### Environment Variables on Railway

Railway automatically sets the `PORT` variable. No manual configuration needed!

---

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful GET request
- `201 Created`: String created successfully
- `204 No Content`: String deleted successfully
- `400 Bad Request`: Invalid request body or query parameters
- `404 Not Found`: String doesn't exist
- `409 Conflict`: String already exists
- `422 Unprocessable Entity`: Invalid data type
- `500 Internal Server Error`: Server error

Error response format:
```json
{
  "error": "Error message description"
}
```

---

## Implementation Details

### String Properties Computed

1. **length**: Total number of characters
2. **is_palindrome**: Case-insensitive palindrome check
3. **unique_characters**: Count of distinct characters
4. **word_count**: Number of whitespace-separated words
5. **sha256_hash**: Unique identifier (simplified hash in this implementation)
6. **character_frequency_map**: Character occurrence counts

### Storage

- **In-memory storage**: Data persists only during server runtime
- **Thread-safe**: Uses mutexes for concurrent access
- **Key-based lookup**: Fast O(1) retrieval by string value

### Natural Language Processing

Simple pattern matching for common query patterns:
- Keyword detection (palindrome, single word, etc.)
- Number extraction (longer than X, at least Y)
- Character detection (containing letter Z)

---

## Testing Checklist

Before submission, verify:

- ✅ All 5 endpoints work correctly
- ✅ Error responses return correct status codes
- ✅ Palindrome detection is case-insensitive
- ✅ Character frequency map is accurate
- ✅ Natural language queries parse correctly
- ✅ Filters combine properly
- ✅ Duplicate strings return 409 Conflict
- ✅ Deleting returns 204 No Content
- ✅ API is deployed and accessible

---

## Troubleshooting

**Server won't start:**
- Check if port 8080 is already in use
- Try `PORT=3000 go run main.go`

**Build fails:**
- Ensure Go 1.21+ is installed: `go version`
- Run `go mod tidy`

**Railway deployment fails:**
- Check Railway logs
- Ensure `main.go` is in root directory
- Verify go.mod exists

---

## Future Enhancements

Potential improvements:
- Persistent storage (PostgreSQL, MongoDB)
- Pagination for GET /strings
- Rate limiting
- API authentication
- Caching layer
- Metrics and logging
- Unit tests
- Docker containerization

---

## License

MIT License - feel free to use this project however you'd like!

---

## Support

For issues or questions:
- Create an issue on GitHub
- Check Railway documentation: https://docs.railway.app
- Go documentation: https://go.dev/doc

---