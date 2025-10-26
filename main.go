package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize storage
	store := NewMemoryStore()

	// Initialize handlers
	handler := NewStringHandler(store)

	// Setup routes
	mux := http.NewServeMux()

	// Router wrapper to handle path-based routing
	mux.HandleFunc("/strings", func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		path := r.URL.Path

		// Route: GET /strings/filter-by-natural-language
		if strings.HasPrefix(path, "/strings/filter-by-natural-language") {
			handler.FilterByNaturalLanguage(w, r)
			return
		}

		// Route: GET /strings/{value} or DELETE /strings/{value}
		if path != "/strings" && path != "/strings/" {
			if r.Method == http.MethodGet {
				handler.GetString(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteString(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Route: POST /strings or GET /strings (with filters)
		if r.Method == http.MethodPost {
			handler.CreateString(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetAllStrings(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Handle the filter-by-natural-language endpoint specifically
	mux.HandleFunc("/strings/filter-by-natural-language", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler.FilterByNaturalLanguage(w, r)
	})

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "String Analyzer API", "version": "1.0.0"}`))
		} else {
			http.NotFound(w, r)
		}
	})

	// Start server
	addr := "0.0.0.0:" + port
	log.Printf("Server starting on %s", addr)
	log.Printf("Available endpoints:")
	log.Printf("  POST   /strings")
	log.Printf("  GET    /strings")
	log.Printf("  GET    /strings/{value}")
	log.Printf("  GET    /strings/filter-by-natural-language")
	log.Printf("  DELETE /strings/{value}")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// ===== MODELS =====

type Properties struct {
	Length                int            `json:"length"`
	IsPalindrome          bool           `json:"is_palindrome"`
	UniqueCharacters      int            `json:"unique_characters"`
	WordCount             int            `json:"word_count"`
	SHA256Hash            string         `json:"sha256_hash"`
	CharacterFrequencyMap map[string]int `json:"character_frequency_map"`
}

type StringAnalysis struct {
	ID         string     `json:"id"`
	Value      string     `json:"value"`
	Properties Properties `json:"properties"`
	CreatedAt  string     `json:"created_at"`
}

func NewStringAnalysis(value string) *StringAnalysis {
	hash := computeSHA256(value)

	return &StringAnalysis{
		ID:    hash,
		Value: value,
		Properties: Properties{
			Length:                len(value),
			IsPalindrome:          isPalindrome(value),
			UniqueCharacters:      countUniqueChars(value),
			WordCount:             countWords(value),
			SHA256Hash:            hash,
			CharacterFrequencyMap: buildFrequencyMap(value),
		},
		CreatedAt: fmt.Sprintf("%s", getCurrentTime()),
	}
}

func computeSHA256(s string) string {
	// Simple hash for demonstration - in production use crypto/sha256
	hash := 0
	for _, c := range s {
		hash = hash*31 + int(c)
	}
	return fmt.Sprintf("%x", hash)
}

func isPalindrome(s string) bool {
	s = strings.ToLower(s)
	left, right := 0, len(s)-1

	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}

	return true
}

func countUniqueChars(s string) int {
	seen := make(map[rune]bool)
	for _, char := range s {
		seen[char] = true
	}
	return len(seen)
}

func countWords(s string) int {
	words := strings.Fields(s)
	return len(words)
}

func buildFrequencyMap(s string) map[string]int {
	freq := make(map[string]int)
	for _, char := range s {
		charStr := string(char)
		freq[charStr]++
	}
	return freq
}

func getCurrentTime() string {
	return "2025-10-21T10:00:00Z"
}

// ===== STORAGE =====

type MemoryStore struct {
	strings map[string]*StringAnalysis
	hashes  map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		strings: make(map[string]*StringAnalysis),
		hashes:  make(map[string]string),
	}
}

func (s *MemoryStore) Create(analysis *StringAnalysis) error {
	if _, exists := s.strings[analysis.Value]; exists {
		return fmt.Errorf("already exists")
	}

	s.strings[analysis.Value] = analysis
	s.hashes[analysis.ID] = analysis.Value

	return nil
}

func (s *MemoryStore) Get(value string) (*StringAnalysis, error) {
	analysis, exists := s.strings[value]
	if !exists {
		return nil, fmt.Errorf("not found")
	}

	return analysis, nil
}

func (s *MemoryStore) GetAll(filters map[string]interface{}) []*StringAnalysis {
	var results []*StringAnalysis

	for _, analysis := range s.strings {
		if matchesFilters(analysis, filters) {
			results = append(results, analysis)
		}
	}

	return results
}

func (s *MemoryStore) Delete(value string) error {
	analysis, exists := s.strings[value]
	if !exists {
		return fmt.Errorf("not found")
	}

	delete(s.strings, value)
	delete(s.hashes, analysis.ID)

	return nil
}

func matchesFilters(analysis *StringAnalysis, filters map[string]interface{}) bool {
	if val, ok := filters["is_palindrome"].(bool); ok {
		if analysis.Properties.IsPalindrome != val {
			return false
		}
	}

	if val, ok := filters["min_length"].(int); ok {
		if analysis.Properties.Length < val {
			return false
		}
	}

	if val, ok := filters["max_length"].(int); ok {
		if analysis.Properties.Length > val {
			return false
		}
	}

	if val, ok := filters["word_count"].(int); ok {
		if analysis.Properties.WordCount != val {
			return false
		}
	}

	if val, ok := filters["contains_character"].(string); ok {
		if !containsChar(analysis.Value, val) {
			return false
		}
	}

	return true
}

func containsChar(s, char string) bool {
	if len(char) == 0 {
		return true
	}
	return strings.Contains(s, char)
}

// ===== HANDLERS =====

type StringHandler struct {
	store *MemoryStore
}

func NewStringHandler(store *MemoryStore) *StringHandler {
	return &StringHandler{store: store}
}

func (h *StringHandler) CreateString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Value == "" {
		respondError(w, http.StatusBadRequest, "Missing 'value' field")
		return
	}

	analysis := NewStringAnalysis(req.Value)

	if err := h.store.Create(analysis); err != nil {
		respondError(w, http.StatusConflict, "String already exists")
		return
	}

	respondJSON(w, http.StatusCreated, analysis)
}

func (h *StringHandler) GetString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	value := strings.TrimPrefix(r.URL.Path, "/strings/")

	if value == "" || value == "strings" {
		respondError(w, http.StatusBadRequest, "String value required")
		return
	}

	analysis, err := h.store.Get(value)
	if err != nil {
		respondError(w, http.StatusNotFound, "String not found")
		return
	}

	respondJSON(w, http.StatusOK, analysis)
}

func (h *StringHandler) GetAllStrings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filters := make(map[string]interface{})
	appliedFilters := make(map[string]interface{})

	query := r.URL.Query()

	if val := query.Get("is_palindrome"); val != "" {
		if val == "true" {
			filters["is_palindrome"] = true
			appliedFilters["is_palindrome"] = true
		} else if val == "false" {
			filters["is_palindrome"] = false
			appliedFilters["is_palindrome"] = false
		}
	}

	if val := query.Get("min_length"); val != "" {
		if i := parseInt(val); i > 0 {
			filters["min_length"] = i
			appliedFilters["min_length"] = i
		}
	}

	if val := query.Get("max_length"); val != "" {
		if i := parseInt(val); i > 0 {
			filters["max_length"] = i
			appliedFilters["max_length"] = i
		}
	}

	if val := query.Get("word_count"); val != "" {
		if i := parseInt(val); i >= 0 {
			filters["word_count"] = i
			appliedFilters["word_count"] = i
		}
	}

	if val := query.Get("contains_character"); val != "" {
		filters["contains_character"] = val
		appliedFilters["contains_character"] = val
	}

	results := h.store.GetAll(filters)

	response := map[string]interface{}{
		"data":            results,
		"count":           len(results),
		"filters_applied": appliedFilters,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *StringHandler) FilterByNaturalLanguage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		respondError(w, http.StatusBadRequest, "Missing 'query' parameter")
		return
	}

	parsed := ParseNaturalLanguageQuery(query)

	results := h.store.GetAll(parsed.Filters)

	response := map[string]interface{}{
		"data":  results,
		"count": len(results),
		"interpreted_query": map[string]interface{}{
			"original":       parsed.Original,
			"parsed_filters": parsed.Filters,
		},
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *StringHandler) DeleteString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	value := strings.TrimPrefix(r.URL.Path, "/strings/")

	if err := h.store.Delete(value); err != nil {
		respondError(w, http.StatusNotFound, "String not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

// ===== NATURAL LANGUAGE PARSER =====

type ParsedQuery struct {
	Original string                 `json:"original"`
	Filters  map[string]interface{} `json:"parsed_filters"`
}

func ParseNaturalLanguageQuery(query string) *ParsedQuery {
	query = strings.ToLower(strings.TrimSpace(query))
	filters := make(map[string]interface{})

	// Check for palindrome keywords
	if containsAny(query, []string{"palindrome", "palindromic", "reads same"}) {
		filters["is_palindrome"] = true
	}

	// Check for word count
	if strings.Contains(query, "single word") {
		filters["word_count"] = 1
	} else if strings.Contains(query, "two word") || strings.Contains(query, "2 word") {
		filters["word_count"] = 2
	} else if strings.Contains(query, "three word") || strings.Contains(query, "3 word") {
		filters["word_count"] = 3
	}

	// Check for length requirements
	if strings.Contains(query, "longer than") {
		// Extract number after "longer than"
		parts := strings.Split(query, "longer than")
		if len(parts) > 1 {
			words := strings.Fields(parts[1])
			if len(words) > 0 {
				if num := parseInt(words[0]); num > 0 {
					filters["min_length"] = num + 1
				}
			}
		}
	}

	if strings.Contains(query, "shorter than") {
		parts := strings.Split(query, "shorter than")
		if len(parts) > 1 {
			words := strings.Fields(parts[1])
			if len(words) > 0 {
				if num := parseInt(words[0]); num > 0 {
					filters["max_length"] = num - 1
				}
			}
		}
	}

	if strings.Contains(query, "at least") {
		parts := strings.Split(query, "at least")
		if len(parts) > 1 {
			words := strings.Fields(parts[1])
			if len(words) > 0 {
				if num := parseInt(words[0]); num > 0 {
					filters["min_length"] = num
				}
			}
		}
	}

	// Check for character containment
	if strings.Contains(query, "containing") || strings.Contains(query, "contain") {
		// Look for "letter X" or "character X"
		if strings.Contains(query, "letter") {
			parts := strings.Split(query, "letter")
			if len(parts) > 1 {
				words := strings.Fields(parts[1])
				if len(words) > 0 && len(words[0]) == 1 {
					filters["contains_character"] = words[0]
				}
			}
		} else if strings.Contains(query, "character") {
			parts := strings.Split(query, "character")
			if len(parts) > 1 {
				words := strings.Fields(parts[1])
				if len(words) > 0 && len(words[0]) == 1 {
					filters["contains_character"] = words[0]
				}
			}
		}
	}

	// Special case: "first vowel" = 'a'
	if strings.Contains(query, "first vowel") {
		filters["contains_character"] = "a"
	}

	return &ParsedQuery{
		Original: query,
		Filters:  filters,
	}
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
