package utils

import (
	"errors"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	BeforeError = errors.New("прохождения теста ещё не началось")
	AfterError  = errors.New("прохождение теста завершено")
)

func CheckDateLimit(startDate, endDate time.Time) error {
	today := time.Now().UTC().Add(3 * time.Hour)
	endDate = endDate.AddDate(0, 0, 1)

	log.Info().Any("startDate", startDate).Any("endDate", endDate).Any("now", today).Send()

	if today.Equal(startDate) || today.Equal(endDate) {
		return nil
	}

	if today.After(endDate) {
		return AfterError
	}

	if today.Before(startDate) {
		return BeforeError
	}

	return nil
}
