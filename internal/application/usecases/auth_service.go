package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/albkvv/student-job-finder-back/internal/domain/repositories"
	"github.com/albkvv/student-job-finder-back/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	otpLength     = 6
	otpExpiration = 5 * time.Minute
	maxOTPAttempts = 5
)

type AuthService struct {
	UserRepo repositories.UserRepository
	CodeRepo repositories.VerificationCodeRepository
}

func NewAuthService(u repositories.UserRepository, c repositories.VerificationCodeRepository) *AuthService {
	return &AuthService{UserRepo: u, CodeRepo: c}
}
func (a *AuthService) RegisterPassword(ctx context.Context, email, phone, password, role string) (*entities.User, string, error) {
	if email == "" && phone == "" {
		return nil, "", errors.New("email or phone is required")
	}
	
	if err := utils.ValidatePassword(password); err != nil {
		return nil, "", err
	}
	
	if role == "" {
		role = "student"
	}
	if err := utils.ValidateRole(role); err != nil {
		return nil, "", err
	}
	
	if email != "" {
		if err := utils.ValidateEmail(email); err != nil {
			return nil, "", err
		}
		existing, _ := a.UserRepo.FindByEmail(ctx, email)
		if existing != nil {
			return nil, "", errors.New("email already registered")
		}
	}
	
	if phone != "" {
		if err := utils.ValidatePhone(phone); err != nil {
			return nil, "", err
		}
		existing, _ := a.UserRepo.FindByPhone(ctx, phone)
		if existing != nil {
			return nil, "", errors.New("phone already registered")
		}
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	
	newUser := &entities.User{
		ID:           utils.GenerateID(),
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsVerified:   false,
		CreatedAt:    time.Now(),
	}
	
	if err := a.UserRepo.Create(ctx, newUser); err != nil {
		return nil, "", err
	}
	
	token, err := utils.GenerateJWT(newUser.ID)
	if err != nil {
		return nil, "", err
	}
	
	return newUser, token, nil
}

func (a *AuthService) LoginPassword(ctx context.Context, identifier, password string) (*entities.User, string, error) {
	if identifier == "" || password == "" {
		return nil, "", errors.New("identifier and password are required")
	}
	
	var foundUser *entities.User
	var err error
	
	if utils.ValidateEmail(identifier) == nil {
		foundUser, err = a.UserRepo.FindByEmail(ctx, identifier)
	} else if utils.ValidatePhone(identifier) == nil {
		foundUser, err = a.UserRepo.FindByPhone(ctx, identifier)
	} else {
		return nil, "", errors.New("invalid identifier format")
	}
	
	if err != nil || foundUser == nil {
		return nil, "", errors.New("invalid credentials")
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	
	token, err := utils.GenerateJWT(foundUser.ID)
	if err != nil {
		return nil, "", err
	}
	
	return foundUser, token, nil
}

func (a *AuthService) RequestPhoneCode(ctx context.Context, phone, role string) error {
	if err := utils.ValidatePhone(phone); err != nil {
		return err
	}
	
	if role == "" {
		role = "student"
	}
	if err := utils.ValidateRole(role); err != nil {
		return err
	}
	
	user, _ := a.UserRepo.FindByPhone(ctx, phone)
	if user == nil {
		newUser := &entities.User{
			ID:         utils.GenerateID(),
			Phone:      phone,
			Role:       role,
			IsVerified: false,
			CreatedAt:  time.Now(),
		}
		a.UserRepo.Create(ctx, newUser)
	}
	
	code := utils.GenerateNumericCode(otpLength)
	exp := time.Now().Add(otpExpiration).Unix()
	
	if err := a.CodeRepo.SetCode(ctx, phone, code, "phone", exp); err != nil {
		return err
	}
	
	utils.LogSMSToConsole(phone, code)
	return nil
}

func (a *AuthService) RequestEmailCode(ctx context.Context, email string) error {
	if err := utils.ValidateEmail(email); err != nil {
		return err
	}
	
	user, _ := a.UserRepo.FindByEmail(ctx, email)
	if user == nil {
		newUser := &entities.User{
			ID:         utils.GenerateID(),
			Email:      email,
			Role:       "student",
			IsVerified: false,
			CreatedAt:  time.Now(),
		}
		a.UserRepo.Create(ctx, newUser)
	}
	
	code := utils.GenerateNumericCode(otpLength)
	exp := time.Now().Add(otpExpiration).Unix()
	
	if err := a.CodeRepo.SetCode(ctx, email, code, "email", exp); err != nil {
		return err
	}
	
	utils.LogSMSToConsole(email, code)
	return nil
}
func (a *AuthService) VerifyPhoneCode(ctx context.Context, phone, code string) (*entities.User, string, error) {
	if err := utils.ValidatePhone(phone); err != nil {
		return nil, "", err
	}
	
	if len(code) != otpLength {
		return nil, "", errors.New("invalid code format")
	}
	
	vc, err := a.CodeRepo.GetCode(ctx, phone, "phone")
	if err != nil || vc == nil {
		return nil, "", errors.New("code not found or expired")
	}
	
	if time.Now().Unix() > vc.ExpiresAt {
		return nil, "", errors.New("code expired")
	}
	
	if vc.Attempts >= maxOTPAttempts {
		return nil, "", errors.New("maximum attempts exceeded")
	}
	
	if vc.Code != code {
		a.CodeRepo.IncrementAttempts(ctx, phone, "phone")
		return nil, "", errors.New("invalid code")
	}
	
	user, err := a.UserRepo.FindByPhone(ctx, phone)
	if err != nil || user == nil {
		return nil, "", errors.New("user not found")
	}
	
	user.IsVerified = true
	a.UserRepo.Update(ctx, user)
	a.CodeRepo.DeleteCode(ctx, phone, "phone")
	
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}
	
	return user, token, nil
}

func (a *AuthService) VerifyEmailCode(ctx context.Context, email, code string) (*entities.User, string, error) {
	if err := utils.ValidateEmail(email); err != nil {
		return nil, "", err
	}
	
	if len(code) != otpLength {
		return nil, "", errors.New("invalid code format")
	}
	
	vc, err := a.CodeRepo.GetCode(ctx, email, "email")
	if err != nil || vc == nil {
		return nil, "", errors.New("code not found or expired")
	}
	
	if time.Now().Unix() > vc.ExpiresAt {
		return nil, "", errors.New("code expired")
	}
	
	if vc.Attempts >= maxOTPAttempts {
		return nil, "", errors.New("maximum attempts exceeded")
	}
	
	if vc.Code != code {
		a.CodeRepo.IncrementAttempts(ctx, email, "email")
		return nil, "", errors.New("invalid code")
	}
	
	user, err := a.UserRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, "", errors.New("user not found")
	}
	
	user.IsVerified = true
	a.UserRepo.Update(ctx, user)
	a.CodeRepo.DeleteCode(ctx, email, "email")
	
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}
	
	return user, token, nil
}

