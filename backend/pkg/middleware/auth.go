package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"gymlog-backend/internal/services"
)

// Clerk JWT Claims struct
type ClerkClaims struct {
	Sub string `json:"sub"` // Clerk user ID
	jwt.RegisteredClaims
}

// JWKS structures for Clerk public keys
type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var (
	jwksCache      *JWKS
	jwksCacheMutex sync.RWMutex
	jwksCacheTime  time.Time
	jwksCacheTTL   = 1 * time.Hour
)

// fetchClerkJWKS fetches the Clerk JWKS (JSON Web Key Set)
func fetchClerkJWKS() (*JWKS, error) {
	jwksCacheMutex.RLock()
	if jwksCache != nil && time.Since(jwksCacheTime) < jwksCacheTTL {
		jwksCacheMutex.RUnlock()
		return jwksCache, nil
	}
	jwksCacheMutex.RUnlock()

	// Get Clerk publishable key to determine the instance
	clerkPublishableKey := os.Getenv("EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY")
	if clerkPublishableKey == "" {
		return nil, fmt.Errorf("EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY environment variable not set")
	}

	// Extract instance from publishable key (format: pk_test_xxx or pk_live_xxx)
	// The JWKS URL is: https://clerk.{instance}.com/.well-known/jwks.json
	// For most cases, we can use the standard Clerk JWKS endpoint
	jwksURL := "https://clerk.com/.well-known/jwks.json"

	// If using a custom domain, you might need to adjust this
	// For now, we'll use the standard endpoint which works for most Clerk instances

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	jwksCacheMutex.Lock()
	jwksCache = &jwks
	jwksCacheTime = time.Now()
	jwksCacheMutex.Unlock()

	return &jwks, nil
}

// getPublicKey retrieves the RSA public key for a given key ID
func getPublicKey(kid string) (*rsa.PublicKey, error) {
	jwks, err := fetchClerkJWKS()
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return jwkToPublicKey(key)
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

// jwkToPublicKey converts a JWK to an RSA public key
func jwkToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)

	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

// validateClerkJWT validates a Clerk JWT token and returns the claims
func validateClerkJWT(tokenString string) (*ClerkClaims, error) {
	// Parse the token to get the header (which contains the kid)
	token, err := jwt.ParseWithClaims(tokenString, &ClerkClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the kid from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		// Get the public key for this kid
		return getPublicKey(kid)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*ClerkClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	return claims, nil
}

// AuthInterceptor provides authentication for gRPC services using Clerk JWT tokens
func AuthInterceptor(userService *services.UserService) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip auth for certain methods
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract token from metadata
		token, err := extractToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
		}

		// Validate Clerk JWT token
		claims, err := validateClerkJWT(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid or expired token: %v", err))
		}

		// Add Clerk user ID to context
		// The 'sub' claim contains the Clerk user ID
		ctx = context.WithValue(ctx, "clerk_user_id", claims.Sub)

		return handler(ctx, req)
	}
}

func isPublicMethod(method string) bool {
	// With Clerk, all methods require authentication
	// User creation/management is handled by Clerk
	publicMethods := []string{
		// Add any truly public methods here (e.g., health checks)
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}
	return false
}

func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	token := authHeader[0]
	if !strings.HasPrefix(token, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	return strings.TrimPrefix(token, "Bearer "), nil
}
