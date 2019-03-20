package middleware

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/yaoice/gocele/pkg/log"
	"net/http"
	"sync"
	"time"
)

var (
	authMiddleware *jwt.GinJWTMiddleware
	err            error
	once           sync.Once
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(http.StatusOK, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func GetAuthMiddleware() *jwt.GinJWTMiddleware {
	once.Do(func() {
		// the jwt middleware
		authMiddleware = &jwt.GinJWTMiddleware{
			Realm:      "test zone",
			Key:        []byte("secret key"),
			Timeout:    time.Hour,
			MaxRefresh: time.Hour,
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				if v, ok := data.(*User); ok {
					return jwt.MapClaims{
						"id": v.UserName,
					}
				}
				return jwt.MapClaims{}
			},
			Authenticator: func(c *gin.Context) (interface{}, error) {
				var loginVals login
				if err := c.Bind(&loginVals); err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				userID := loginVals.Username
				password := loginVals.Password

				if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
					return &User{
						UserName:  userID,
						LastName:  "Bo-Yi",
						FirstName: "Wu",
					}, nil
				}

				return nil, jwt.ErrFailedAuthentication
			},
			Authorizator: func(data interface{}, c *gin.Context) bool {
				if v, ok := data.(string); ok && v == "admin" {
					return true
				}

				return false
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.JSON(code, gin.H{
					"code":    code,
					"message": message,
				})
			},
			// TokenLookup is a string in the form of "<source>:<name>" that is used
			// to extract token from the request.
			// Optional. Default value "header:Authorization".
			// Possible values:
			// - "header:<name>"
			// - "query:<name>"
			// - "cookie:<name>"
			TokenLookup: "header: Authorization, query: token, cookie: jwt",
			// TokenLookup: "query:token",
			// TokenLookup: "cookie:token",

			// TokenHeadName is a string in the header. Default value is "Bearer"
			TokenHeadName: "Bearer",

			// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
			TimeFunc: time.Now,
		}
		if err != nil {
			log.Fatal("JWT Error:" + err.Error())
		}
	})
	return authMiddleware
}
