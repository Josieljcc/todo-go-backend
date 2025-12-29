package middleware

import (
	"strconv"
	"strings"
	"todo-go-backend/internal/config"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates a CORS middleware based on the provided configuration
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Parse allowed origins
	allowedOrigins := parseStringList(cfg.CORSAllowedOrigins)
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"} // Default: allow all
	}

	// Check if we have wildcard (*) configured
	hasWildcard := false
	for i, origin := range allowedOrigins {
		if strings.TrimSpace(origin) == "*" {
			hasWildcard = true
			allowedOrigins[i] = "*" // Normalize to exactly "*"
			break
		}
	}

	// Parse allowed methods
	allowedMethods := parseStringList(cfg.CORSAllowedMethods)
	if len(allowedMethods) == 0 {
		allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	}

	// Parse allowed headers
	allowedHeaders := parseStringList(cfg.CORSAllowedHeaders)
	if len(allowedHeaders) == 0 {
		allowedHeaders = []string{"Content-Type", "Authorization", "Accept", "Origin"}
	}

	// Parse exposed headers
	exposedHeaders := parseStringList(cfg.CORSExposedHeaders)

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		originAllowed := false
		var allowedOriginValue string

		// Debug: log CORS configuration (only in development)
		// Remove or comment out in production
		if gin.Mode() == gin.DebugMode {
			gin.DefaultWriter.Write([]byte(
				"[CORS] Request Origin: " + origin + "\n" +
					"[CORS] Allowed Origins: " + strings.Join(allowedOrigins, ", ") + "\n" +
					"[CORS] Allow Credentials: " + strconv.FormatBool(cfg.CORSAllowCredentials) + "\n",
			))
		}

		// Check if we have wildcard (*) configured
		if hasWildcard || (len(allowedOrigins) == 1 && allowedOrigins[0] == "*") {
			// Allow all origins
			originAllowed = true

			// If credentials are allowed, we CANNOT use "*" - must use specific origin
			// This is a CORS specification requirement
			if cfg.CORSAllowCredentials {
				// When credentials are allowed, we must return the specific origin
				if origin != "" {
					allowedOriginValue = origin
				} else {
					// No origin header means same-origin request
					// Construct origin from request
					scheme := "http"
					if c.Request.TLS != nil {
						scheme = "https"
					}
					allowedOriginValue = scheme + "://" + c.Request.Host
				}
			} else {
				// When credentials are NOT allowed, we can use "*"
				allowedOriginValue = "*"
			}
		} else {
			// If no origin header (same-origin request), allow it
			if origin == "" {
				originAllowed = true
				if cfg.CORSAllowCredentials {
					// For same-origin with credentials, use the request origin
					scheme := "http"
					if c.Request.TLS != nil {
						scheme = "https"
					}
					allowedOriginValue = scheme + "://" + c.Request.Host
				} else {
					allowedOriginValue = "*"
				}
			} else {
				// Check if origin is in allowed list
				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin {
						originAllowed = true
						allowedOriginValue = origin
						break
					}
				}
			}
		}

		if originAllowed {
			c.Header("Access-Control-Allow-Origin", allowedOriginValue)
		}

		// Set CORS headers
		if len(allowedMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		}

		if len(allowedHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
		}

		if len(exposedHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(exposedHeaders, ", "))
		}

		if cfg.CORSAllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if cfg.CORSMaxAge > 0 {
			c.Header("Access-Control-Max-Age", strconv.Itoa(cfg.CORSMaxAge))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// parseStringList parses a comma-separated string into a slice of trimmed strings
func parseStringList(s string) []string {
	if s == "" {
		return []string{}
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
