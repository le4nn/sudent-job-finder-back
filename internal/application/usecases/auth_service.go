package usecases

import (
	"time"

	"github.com/albkvv/student-job-finder-back/internal/domain/entities"
	"github.com/albkvv/student-job-finder-back/internal/domain/repositories"
	"github.com/albkvv/student-job-finder-back/internal/utils"
)

type AuthService struct {
	UserRepo repositories.UserRepository
	CodeRepo repositories.VerificationCodeRepository
}

func NewAuthService(u repositories.UserRepository, c repositories.VerificationCodeRepository) *AuthService {
	return &AuthService{UserRepo: u, CodeRepo: c}
}
func (a *AuthService) RequestCode(phone, role string) error {
	code := utils.GenerateNumericCode(6)
	exp := time.Now().Add(5 * time.Minute).Unix()
	if err := a.CodeRepo.SetCode(phone, code, exp); err != nil {
		return err
	}
	user, err := a.UserRepo.FindByPhone(phone)
	if err != nil || user == nil {
		u := &entities.User{Phone: phone, Role: role}
		a.UserRepo.Create(u)
	}
	utils.LogSMSToConsole(phone, code)
	return nil
}
func (a *AuthService) VerifyCode(phone, code string) (string, *entities.User, int64, error) {
	vc, err := a.CodeRepo.GetCode(phone)
	if err != nil || vc == nil || vc.Code != code || time.Now().Unix() > vc.ExpiresAt {
		return "", nil, 0, err
	}
	user, err := a.UserRepo.FindByPhone(phone)
	if err != nil || user == nil {
		return "", nil, 0, err
	}
	token, exp, err := utils.GenerateJWT(user.ID, 30*24*time.Hour)
	if err != nil {
		return "", nil, 0, err
	}
	a.CodeRepo.DeleteCode(phone)
	return token, user, exp.Unix(), nil
}

