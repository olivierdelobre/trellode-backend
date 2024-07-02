package middlewares

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"trellode-go/internal/utils/logging"
	"trellode-go/internal/utils/messages"

	"net/http"

	toolbox_api "github.com/epfl-si/go-toolbox/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// don't really do authentication as it is delegated upstream in KrakenD, but this middleware aims to extract the userid from the Authorization header to use it as identity when creating/updating/deleting entities in the API
func AuthenticationMiddleware(db *gorm.DB, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Printf("---------- fullpath:%s\n", c.FullPath())
		// no control for some endpoints
		if strings.Contains(c.FullPath(), "/docs/") || strings.Contains(c.FullPath(), "/healthcheck") || strings.HasSuffix(c.FullPath(), "/liveness") || strings.Contains(c.FullPath(), "/users/register") || strings.Contains(c.FullPath(), "/users/authenticate") {
			c.Set("userId", "probe")
			c.Next()
			return
		}

		authorization := c.GetHeader("Authorization")
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		lang := c.Request.Header.Get("Content-Language")
		if lang == "" {
			lang = "fr"
		}
		c.Set("lang", lang)

		if reqMethod == "OPTIONS" {
			c.Next()
			return
		}

		bearerPattern, _ := regexp.Compile("Bearer (.*)")
		if bearerPattern.MatchString(authorization) {
			if os.Getenv("MODE") == "local" {
				// get Bearer value from Authorization header
				// Authorization: Bearer <token>
				matchBearer := regexp.MustCompile(`Bearer (.*)`)
				matches := matchBearer.FindStringSubmatch(c.Request.Header.Get("Authorization"))
				if len(matches) > 1 {
					c.Set("userId", matches[1])
					c.Next()
					return
				}
			}

			checkSuccessful, user, err := DecodeToken(log, c, authorization)
			if err != nil {
				logging.LogCustom(log, "error", reqMethod, reqUri, 500, "", "AuthenticationMiddleware: DecodeToken: "+err.Error())
				c.JSON(http.StatusForbidden, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(lang, "NotAuthorized"), err.Error(), "", nil))
				c.Abort()
				return
			}

			if !checkSuccessful {
				logging.LogCustom(log, "error", reqMethod, reqUri, 500, "", "AuthenticationMiddleware: passed Bearer authorization '"+authorization+"' could not be validated")
				c.JSON(http.StatusForbidden, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(lang, "NotAuthorized"), messages.GetMessage(lang, "InvalidToken"), "", nil))
				c.Abort()
				return
			}

			c.Set("userId", user.Uniqueid)
			c.Next()
			return
		} else {
			logging.LogCustom(log, "error", c.Request.Method, c.Request.RequestURI, 401, "", messages.GetMessage(lang, "NotAuthorized"))
			c.JSON(http.StatusForbidden, toolbox_api.MakeError(c, "", http.StatusBadRequest, messages.GetMessage(lang, "NotAuthorized"), messages.GetMessage(lang, "UnsupportedAuthorizationSchema"), "", nil))
			c.Abort()
			return
		}
	}
}

type TokenInfo struct {
	Uniqueid string `json:"uniqueid"`
	Email    string `json:"username"`
}

type UserClaims struct {
	Id      string `json:"id"`
	Email   string `json:"firstname"`
	Profile string `json:"profile"`
	jwt.StandardClaims
}

func DecodeToken(log *zap.Logger, c *gin.Context, tokenB64 string) (bool, TokenInfo, error) {
	tokenB64 = strings.Replace(tokenB64, "Bearer ", "", 1)

	// Parse the JWT
	parsedAccessToken, err := jwt.Parse(tokenB64, func(token *jwt.Token) (interface{}, error) {
		// Make sure the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return false, TokenInfo{}, err
	}

	tokenInfo := TokenInfo{}
	// Check if the token is valid
	if claims, ok := parsedAccessToken.Claims.(jwt.MapClaims); ok && parsedAccessToken.Valid {
		//fmt.Printf("Token valid, claims: %v\n", claims)
		if id, ok := claims["id"].(string); ok {
			tokenInfo.Uniqueid = id
			fmt.Printf("id: %s\n", id)
		}
		if email, ok := claims["email"].(string); ok {
			tokenInfo.Email = email
			fmt.Printf("email: %s\n", email)
		}
	} else {
		return false, TokenInfo{}, err
	}

	return true, tokenInfo, nil
}
