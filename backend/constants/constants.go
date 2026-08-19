package constants

import "os"

// getEnv returns the value of the given environment variable, or fallback
// if it is unset/empty. Keeps all external configuration (Elasticsearch
// credentials, GCS bucket, OpenAI key, etc.) out of source control.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

var (
	USER_INDEX = "user"
	POST_INDEX = "post"

	ES_URL      = getEnv("ES_URL", "http://localhost:9200")
	ES_USERNAME = getEnv("ES_USERNAME", "elastic")
	ES_PASSWORD = getEnv("ES_PASSWORD", "changeme")

	GCS_BUCKET = getEnv("GCS_BUCKET", "your-gcs-bucket-name")

	// JWT_SECRET signs auth tokens. The fallback is for local dev only —
	// always set a real JWT_SECRET env var outside local development.
	JWT_SECRET = getEnv("JWT_SECRET", "dev-only-change-me")

	// OPENAI_API_KEY is used server-side only to call DALL-E 3. It never
	// reaches the frontend bundle.
	OPENAI_API_KEY = getEnv("OPENAI_API_KEY", "")
)
