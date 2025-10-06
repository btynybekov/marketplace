package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// store должен быть тем же, что используется в Login
var store = sessions.NewCookieStore([]byte("supersecretkey"))

// AuthMiddleware проверяет, авторизован ли пользователь
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")

		userID, ok := session.Values["userID"]
		if !ok || userID == nil {
			// Если пользователь не авторизован — редирект на страницу логина
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Пользователь авторизован — передаём управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
