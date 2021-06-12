package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/jwt_api"
	"go-site/session"
	"net/http"
)

func CreateDoLogout(_ *pgx.Conn, rdb *redis.Client) {

	doLogout := func(writer http.ResponseWriter, request *http.Request) {
		var payload *jwt_api.Payload

		if request.Method != "POST" {
			return
		}

		{ // check user authed
			JWTToken, _ := request.Cookie("JWT")
			payload, _ = jwt_api.VerifyToken(JWTToken.Value)
		}

		{ // delete session
			sessionId, _ := request.Cookie("SessionId")
			_ = session.DeleteSession(writer, rdb, sessionId.Value)
			_ = jwt_api.DeleteJWTToken(writer, rdb, *payload)
			_ = jwt_api.DeleteRefreshToken(writer, rdb, *payload)
		}
	}

	http.HandleFunc("/do/logout", doLogout)
}
