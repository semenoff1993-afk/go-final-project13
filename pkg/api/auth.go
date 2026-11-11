package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super_secret_key_todo_app")

// getPassHash создаёт хэш пароля для хранения в токене
func getPassHash(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return hex.EncodeToString(h.Sum(nil))
}

// signinHandler обрабатывает POST-запрос на /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {
	expectedPass := os.Getenv("TODO_PASSWORD")
	if expectedPass == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Аутентификация не настроена"})
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка запроса"})
		return
	}

	if req.Password != expectedPass {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": getPassHash(expectedPass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка генерации токена"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

// authMiddleware проверяет наличие и валидность токена при каждом запросе
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtStr string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtStr = cookie.Value
			}

			if jwtStr == "" {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			hash, ok := claims["hash"].(string)
			if !ok || hash != getPassHash(pass) {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}