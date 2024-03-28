package service

import (
	"App/internal/models"
	"errors"
	"github.com/rs/zerolog/log"
)

var (
	GroupCreateError = errors.New("ошибка при создании группы")
	GroupGetError    = errors.New("ошибка при получении группы")
	GroupUpdateError = errors.New("ошибка при обновлении группы")
	GroupNotFound    = errors.New("группа не найден")
	GroupDeleteError = errors.New("не удалось удалить группу")
)

type GroupRepository interface {
	Create(title string, userID int) (int, error)
	GetAll(userID int) ([]models.GroupOut, error)
	Get(groupID, userID int) (models.GroupOut, error)
	UpdateTitle(groupID, userID int, title string) error
	Delete(userID, groupID int) error
}

type GroupService struct {
	GroupRepository
}

func NewGroupService(gr GroupRepository) *GroupService {
	return &GroupService{gr}
}

func (s *GroupService) Create(in models.GroupIn, userID int) (int, error) {
	id, err := s.GroupRepository.Create(in.Title, userID)

	if err != nil {
		log.Err(err).Send()
		return 0, GroupCreateError
	}

	return id, nil
}

func (s *GroupService) GetAll(userID int) ([]models.GroupOut, error) {
	groups, err := s.GroupRepository.GetAll(userID)

	if err != nil {
		log.Err(err).Send()
		return nil, GroupGetError
	}

	return groups, nil
}

func (s *GroupService) Get(groupID, userID int) (models.GroupOut, error) {
	group, err := s.GroupRepository.Get(groupID, userID)

	if err != nil {
		log.Err(err).Send()
		return models.GroupOut{}, GroupGetError
	}

	return group, nil
}

func (s *GroupService) UpdateTitle(userID, groupID int, title string) error {
	err := s.GroupRepository.UpdateTitle(groupID, userID, title)

	if err != nil {
		log.Err(err).Send()
		return GroupUpdateError
	}

	return nil
}

func (s *GroupService) Delete(userID, groupID int) error {
	err := s.GroupRepository.Delete(userID, groupID)

	if err != nil {
		log.Err(err).Send()
		return GroupDeleteError
	}

	return nil
}
