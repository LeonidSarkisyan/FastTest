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

	UserTokenError         = errors.New("ошибка при генерации токена")
	UserNotCorrectPassword = errors.New("невозможно захешировать пароль")
)

type UserRepository interface {
	Create(userIn models.UserIn) error
	GetByEmail(email string) (models.User, error)
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
		return models.User{}, err
	}

	return userExist, nil
}
