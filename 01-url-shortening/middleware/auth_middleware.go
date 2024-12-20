package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
)

type OIDCConfig struct {
	Issuer   string
	ClientID string
}

// AuthMiddleware returns a Gin middleware that validates JWT access tokens from Keycloak using go-oidc.
func AuthMiddleware(cfg OIDCConfig) (gin.HandlerFunc, error) {
	// Set up the provider using the given issuer URL (Keycloak realm URL).
	provider, err := oidc.NewProvider(context.Background(), cfg.Issuer)
	if err != nil {
		return nil, err
	}

	// Configure an OIDC verifier.
	// If you do not have a Client Secret (Public Client), you can omit some fields.
	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return func(c *gin.Context) {
		// Extract the token from the Authorization header: "Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHORIZED", "error_message": "Missing Authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": "INVALID_TOKEN_FORMAT", "error_message": "Invalid Authorization header format"})
			return
		}

		rawToken := parts[1]

		// Verify the token
		idToken, err := verifier.Verify(context.Background(), rawToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": "INVALID_TOKEN", "error_message": "Token verification failed"})
			log.Println("Verification error: ", err.Error())
			return
		}

		// Extract custom claims if needed
		var claims struct {
			Email             string   `json:"email"`
			PreferredUsername string   `json:"preferred_username"`
			Roles             []string `json:"realm_access.roles"` // Example for Keycloak
		}
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": "INVALID_CLAIMS", "error_message": "Failed to parse token claims"})
			return
		}

		// Store claims or user info in context
		c.Set("email", claims.Email)
		c.Set("username", claims.PreferredUsername)
		c.Set("roles", claims.Roles)

		c.Next()
	}, nil
}
