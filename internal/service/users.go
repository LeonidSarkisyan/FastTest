package service

import (
	"App/internal/models"
	"App/pkg/utils"
	"errors"
	"github.com/rs/zerolog/log"
)

var (
	UserAlreadyExists      = errors.New("пользователь с таким email уже существует")
	UserBadLoginOrPassword = errors.New("неверный логин или пароль")

	UserTokenError          = errors.New("ошибка при генерации токена")
	UserNotCorrectPassword  = errors.New("невозможно захешировать пароль")
	UserGetError            = errors.New("ошибка при получении пользователя")
	UserChangePasswordError = errors.New("ошибка при обновлении пароля пользователя")
)

type UserRepository interface {
	Create(userIn models.UserIn) error
	GetByEmail(email string) (models.User, error)
	GetByID(userID int) (models.User, error)
	ChangePassword(userID int, newPassword models.NewPassword) error
}

type UserService struct {
	UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{r}
}

func (s *UserService) Register(in models.UserIn) error {
	userExist, err := s.UserRepository.GetByEmail(in.Email)

	if err != nil && userExist.ID == 0 {
		log.Err(err).Send()
		return UserAlreadyExists
	}

	in.Password, err = utils.HashPassword(in.Password)

	if err != nil {
		log.Err(err).Send()
		return UserNotCorrectPassword
	}

	err = s.UserRepository.Create(in)

	if err != nil {
		log.Err(err).Send()
		return UserAlreadyExists
	}

	return nil
}

func (s *UserService) Login(in models.UserIn) (string, error) {
	userExist, err := s.UserRepository.GetByEmail(in.Email)

	if err != nil || userExist.ID == 0 {
		log.Err(err).Send()
		return "", UserBadLoginOrPassword
	}

	if !utils.ComparePassword(in.Password, userExist.Password) {
		return "", UserBadLoginOrPassword
	}

	log.Info().Int("user_id", userExist.ID).Send()

	token, err := utils.GenerateJWTToken(userExist.ID)

	if err != nil {
		log.Err(err).Send()
		return "", UserTokenError
	}

	return token, nil
}

func (s *UserService) Exist(email string) bool {
	userExist, err := s.UserRepository.GetByEmail(email)

	if err != nil {
		log.Err(err).Send()
		return true
	}

	return userExist.ID != 0
}

func (s *UserService) GetByEmail(email string) (models.User, error) {
	userExist, err := s.UserRepository.GetByEmail(email)

	if err != nil {
		log.Err(err).Send()
		return models.User{}, UserGetError
	}

	return userExist, nil
}

func (s *UserService) GetByID(userID int) (models.User, error) {
	userExist, err := s.UserRepository.GetByID(userID)

	if err != nil {
		log.Err(err).Send()
		return models.User{}, UserGetError
	}

	return userExist, nil
}

func (s *UserService) ChangePassword(userID int, newPassword models.NewPassword) error {
	hashPassword, err := utils.HashPassword(newPassword.Password)

	if err != nil {
		return UserChangePasswordError
	}

	newPassword.Password = hashPassword

	err = s.UserRepository.ChangePassword(userID, newPassword)

	if err != nil {
		return UserChangePasswordError
	}

	return nil
}
