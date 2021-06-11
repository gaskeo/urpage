package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/jwt"
	"go-site/session"
	"go-site/structs"
	"net/http"
)

func CreateDoLogout(_ *pgx.Conn, rdb *redis.Client) {

	doLogout := func(writer http.ResponseWriter, request *http.Request) {
		var (
			payload *structs.Payload
		)

		if request.Method != "POST" {
			return
		}

		{ // check user authed
			JWTToken, _ := request.Cookie("JWT")
			payload, _ = jwt.VerifyToken(JWTToken.Value)
		}

		{ // delete session
			sessionId, _ := request.Cookie("SessionId")
			_ = session.DeleteSession(writer, rdb, sessionId.Value)
			_ = jwt.DeleteJWTToken(writer, rdb, *payload)
			_ = jwt.DeleteRefreshToken(writer, rdb, *payload)
		}
	}

	http.HandleFunc("/do/logout", doLogout)
}
