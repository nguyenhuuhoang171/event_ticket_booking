package usecase

import (
	"context"
	"os"
	"regexp"
	"time"

	"event_ticket_booking/constant"
	refreshTokenEntity "event_ticket_booking/infrastructure/db/refresh_token/entity"
	userEntity "event_ticket_booking/infrastructure/db/user/entity"
	"event_ticket_booking/internal/domain/auth/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func verifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (u Usecase) generateAccessToken(user *userEntity.Entity) (string, error) {
	accessSecret := os.Getenv("ACCESS_SECRET")

	claims := model.Claims{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(u.cfg.Authentication.AccessTokenExpirationMinutes) * time.Minute)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return accessToken.SignedString([]byte(accessSecret))
}

func generateRefreshToken() string {
	return uuid.New().String()
}

func generateExpireAtOfRefreshToken(refreshTokenExpirationDays int) time.Time {
	return time.Now().Add(time.Duration(refreshTokenExpirationDays) * 24 * time.Hour)
}

func isValidEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(emailPattern)
	return regex.MatchString(email)
}

func isValidPassword(password string) bool {
	// Check minimum length
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case !('a' <= char && char <= 'z') && !('A' <= char && char <= 'Z') && !('0' <= char && char <= '9'):
			hasSpecial = true
		}
	}

	return hasUpper && hasNumber && hasSpecial
}

func (u Usecase) saveRefreshToken(ctx context.Context, userId uint64) (string, error) {
	refreshToken := generateRefreshToken()
	expireAt := generateExpireAtOfRefreshToken(u.cfg.Authentication.RefreshTokenExpirationDays)
	activeStatus := constant.REFRESH_TOKEN_STATUS_ACTIVE

	_, err := u.refreshTokenRepo.Create(ctx, &refreshTokenEntity.Entity{
		UserId:   userId,
		Token:    refreshToken,
		ExpireAt: expireAt,
		Status:   &activeStatus,
	})
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
