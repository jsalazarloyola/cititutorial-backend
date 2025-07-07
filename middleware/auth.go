package middleware

import (
	"errors"
	"go-template/models"
	"go-template/services"
	"log"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Authenticator struct {
	service *services.LoginService
}

func NewAuthenticator(service *services.LoginService) Authenticator {
	return Authenticator{
		service: service,
	}
}

func (auth Authenticator) Login(c *gin.Context) (any, error) {
	var loginValues models.Login
	// log.Println("holi", c.Request.Body)

	if err := c.BindJSON(&loginValues); err != nil {
		return nil, err
	}
	// log.Println(loginValues)
	response, err := auth.service.DoRequestLogin(loginValues.User, loginValues.Password)

	if err != nil {
		return nil, err
	}
	if len(response.Data) == 0 {
		return nil, errors.New("usuario o contraseña incorrectas")
	}

	// log.Println(response.Data)

	user := models.User{
		Rut:  response.Data["rut"].(string),
		User: response.Data["user"].(string),
	}

	return user, nil
}

// En caso de rellenar después
func (auth Authenticator) Authorize(data any, c *gin.Context) bool {
	if data == nil {
		log.Println("ño >:c")
		return false
	}
	userData := data.(map[string]any)
	log.Println(userData)

	return true
}

// La función IdentityHandler por defecto asume que la identificación viene
// dada por IdentityKey, pero acá no usamos eso, sino que user, así que esto.
// Una alternativa es usar el campo IdentityKey, que especifica cómo se llama
// el campo que obtiene el usuario

// func (auth Authenticator) IdentityHandler(c *gin.Context) any {
// 	jwtClaims := jwt.ExtractClaims(c)
// 	return jwtClaims["user"]
// }

// La carga con los datos que tendrá el token
func (auth Authenticator) Payload(data any) jwt.MapClaims {
	user := data.(models.User)

	claims := jwt.MapClaims{
		"user": map[string]any{
			"user": user.User,
			"rut":  user.Rut,
		},
	}

	return claims
}

// La cosa que manejará las sesiones al final
func LoadJWTAuth(service *services.LoginService) *jwt.GinJWTMiddleware {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		key = "asdf"
	}

	auth := NewAuthenticator(service)

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "Tutorial",
		Key:   []byte(key),
		// Tiempo de vencimiento del token
		Timeout:    time.Hour * 24, // el día
		MaxRefresh: time.Hour * 24, // máximo tiempo para refrescar

		Authenticator: auth.Login,
		Authorizator:  auth.Authorize,
		PayloadFunc:   auth.Payload,
		// Unauthorized:
		// LoginResponse:
		// LogoutResponse:
		// RefreshResponse:
		// IdentityHandler:
		IdentityKey:   "user",
		TokenLookup:   "header: Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		// HTTPStatusMessageFunc:
		// PrivKeyFile:          key,
		// PrivKeyBytes:         []byte{},
		// PubKeyFile:           key,
		// PrivateKeyPassphrase: key,
		// PubKeyBytes:          []byte{},
		// SendCookie:           false,
		// CookieMaxAge:         0,
		// SecureCookie:         false,
		// CookieHTTPOnly:       false,
		// CookieDomain:         "",
		// SendAuthorization:    false,
		// DisabledAbort:        false,
		// CookieName:           "",
		// CookieSameSite:       0,
		// ParseOptions:         []jwt.ParserOption{},
		// ExpField:             "",
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return authMiddleware
}
