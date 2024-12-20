package services

import (
	"errors"
	"fmt"
	"net/http"
	"reson8-learning-api/managers"
	requests "reson8-learning-api/models/dtos/requests/auth"
	"reson8-learning-api/models/entities"
	"reson8-learning-api/util"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService struct {
	db        *gorm.DB
	hasher    managers.HashingManager
	logger    *util.Logger
	secretKey string
}

func NewAuthService(db *gorm.DB, hasher managers.HashingManager, logger *util.Logger, secretKey string) *AuthService {
	return &AuthService{db: db, hasher: hasher, logger: logger, secretKey: secretKey}
}

func (s *AuthService) Login(request requests.LoginRequest) (string, error) {
	var user entities.User
	s.logger.Debug().Msg("AuthService_Login has started")
	res := s.db.Preload("Password").Where("email = ?", request.Username).First(&user)
	if res.Error != nil {
		s.logger.Error().Msg(fmt.Sprintf("AuthService_Login had errors when looking for given username: %v", request.Username))
		return "", res.Error
	}

	check, err := s.hasher.Verify(user.Password.Value, request.Password)
	if !check {
		if err != nil {
			s.logger.Error().Err(err).Msg("AuthService_Login had failed while trying to match passwords via hashinManager")
			return "", err
		}
		return "", errors.New("passwordIncorrect")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(time.Hour * 8).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		s.logger.Error().Msg("AuthService_Login had errors when trying to get signed token string")
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) keyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Error().Err(http.ErrAbortHandler).Msg("AuthService_AuthMiddleware had errors when get key while trying to Parse token")
			return nil, http.ErrAbortHandler
		}
		return []byte(s.secretKey), nil
	}
}

func (s *AuthService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		token, err := jwt.Parse(tokenString, s.keyFunc())
		if err != nil {
			s.logger.Error().Err(err).Msg("AuthService_AuthMiddleware had an error when parsing jwt token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Token could not be parsed"})
			c.Abort()
			return
		}
		if !token.Valid {
			s.logger.Debug().Msg(fmt.Sprintf("AuthService_AuthMiddleware token is not valid: %v", tokenString))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Token is not valid"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			s.logger.Debug().Msg(fmt.Sprintf("AuthService_AuthMiddleware claims are not valid: %v", tokenString))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - Claims are not valid"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *AuthService) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(jwt.MapClaims)
		role := uint(claims["role"].(float64))

		if role != uint(entities.Admin) {
			s.logger.Debug().Msg(fmt.Sprintf("AuthService_AdminMiddleware role is not admin: %v", role))
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *AuthService) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *AuthService) GetUserIDFromToken(tokenString string) (string, error) {
	claims, err := s.ExtractClaims(tokenString)
	if err != nil {
		return "", err
	}
	id, ok := claims["userID"].(string)
	if !ok {
		return "", errors.New("tokenErrorClaimsNotValid")
	}
	return id, nil
}

func (s *AuthService) ExtractClaims(tokenString string) (map[string]any, error) {
	res := make(map[string]any)
	token, err := jwt.Parse(tokenString, s.keyFunc())
	if err != nil {
		s.logger.Error().Err(err).Msg("AuthService_AuthMiddleware had an error when parsing jwt token")
		return res, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return res, errors.New("tokenErrorClaimsNotValid")
	}
	return s.mapClaimsToDict(claims, res), nil
}

func (s *AuthService) mapClaimsToDict(claims jwt.MapClaims, res map[string]any) map[string]any {
	for key := range claims {
		res[key] = claims[key]
	}
	return res
}
