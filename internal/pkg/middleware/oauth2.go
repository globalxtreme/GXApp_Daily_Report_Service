package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"service/internal/pkg/config"
)

// GXEmployee adalah representasi data karyawan dari GX OAuth2 server.
type GXEmployee struct {
	ID                string `json:"id"`
	EmployeeNo        string `json:"employeeNo"`
	FullName          string `json:"fullName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	Email             string `json:"email"`
	Gender            string `json:"gender"`
	CompanyOfficeID   int    `json:"companyOfficeId"`
	CompanyOfficeName string `json:"companyOfficeName"`
	DepartmentID      int    `json:"departmentId"`
	DepartmentName    string `json:"departmentName"`
	DivisionID        int    `json:"divisionId"`
	DivisionName      string `json:"divisionName"`
	SectionID         int    `json:"sectionId"`
	SectionName       string `json:"sectionName"`
	JobPositionID     int    `json:"jobPositionId"`
	JobPositionName   string `json:"jobPositionName"`
	JobLevelID        int    `json:"jobLevelId"`
	JobLevelName      string `json:"jobLevelName"`
}

// GXAccessToken adalah response dari endpoint /oauth2/token.
type GXAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type contextKey string

const (
	cookieState       = "gx_oauth2_state"
	cookieAccessToken = "gx_access_token"
	employeeCtxKey    contextKey = "employee"
)

// Redirect membuat state random, menyimpannya ke cookie, lalu redirect ke
// authorization URL GX OAuth2 server.
func Redirect(w http.ResponseWriter, r *http.Request) {
	state := generateState()

	http.SetCookie(w, &http.Cookie{
		Name:     cookieState,
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
	})

	authURL := fmt.Sprintf(
		"%s/oauth2/authorization?state=%s&client_id=%s",
		config.OAuth2.BaseURL,
		url.QueryEscape(state),
		url.QueryEscape(config.OAuth2.ClientID),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

// Callback memvalidasi state, menukar code → token, mengambil data employee,
// menyimpan access token ke cookie, lalu redirect ke frontend.
func Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	// Validasi state untuk mencegah CSRF
	cookie, err := r.Cookie(cookieState)
	if err != nil || state == "" || state != cookie.Value {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid or missing OAuth2 state",
		})
		return
	}

	// Hapus state cookie
	http.SetCookie(w, &http.Cookie{
		Name:    cookieState,
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
	})

	// Tukar authorization code → access token
	token, err := exchangeToken(code)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": fmt.Sprintf("token exchange failed: %v", err),
		})
		return
	}

	// Simpan access token ke cookie (dibaca oleh frontend untuk Auth header)
	http.SetCookie(w, &http.Cookie{
		Name:    cookieAccessToken,
		Value:   token.AccessToken,
		Path:    "/",
		Expires: time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	})

	http.Redirect(w, r, config.OAuth2.FrontendURL, http.StatusFound)
}

// AuthMiddleware memvalidasi Bearer token di header Authorization (atau cookie),
// mengambil data employee dari GX server, dan menyimpannya ke request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "unauthorized: missing token",
			})
			return
		}

		employee, err := fetchEmployee(token)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "unauthorized: invalid token",
			})
			return
		}

		ctx := context.WithValue(r.Context(), employeeCtxKey, employee)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetEmployee mengambil data GXEmployee yang sudah diinjek oleh AuthMiddleware
// dari request context.
func GetEmployee(r *http.Request) *GXEmployee {
	emp, _ := r.Context().Value(employeeCtxKey).(*GXEmployee)
	return emp
}

// --- internal helpers ---

func generateState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	// fallback: baca dari cookie (berguna saat frontend belum set header)
	if cookie, err := r.Cookie(cookieAccessToken); err == nil {
		return cookie.Value
	}
	return ""
}

func exchangeToken(code string) (*GXAccessToken, error) {
	resp, err := http.PostForm(
		config.OAuth2.BaseURL+"/oauth2/token",
		url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {config.OAuth2.ClientID},
			"client_secret": {config.OAuth2.ClientSecret},
			"code":          {code},
			"redirect_uri":  {config.OAuth2.RedirectURL},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	var token GXAccessToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decode token: %w", err)
	}

	return &token, nil
}

func fetchEmployee(accessToken string) (*GXEmployee, error) {
	req, err := http.NewRequest(http.MethodGet, config.OAuth2.BaseURL+"/oauth2/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d from /oauth2/user", resp.StatusCode)
	}

	var emp GXEmployee
	if err := json.NewDecoder(resp.Body).Decode(&emp); err != nil {
		return nil, fmt.Errorf("decode employee: %w", err)
	}

	return &emp, nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
