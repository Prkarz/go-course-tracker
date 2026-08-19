package apilayer

import (
	"database/sql"

	"google.golang.org/genai"
)

// APIServer struct represents the API server with a database connection.
// It contains a single field, DB, which is a pointer to an sql.DB instance.
// The APIServer struct is used to handle API requests and interact with the database.
type APIServer struct {
	DB       *sql.DB
	AIClient *genai.Client
}
