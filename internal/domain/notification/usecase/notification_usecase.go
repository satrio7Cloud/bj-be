package usecase

import (
	"bj/internal/domain/notification/dto"
	"bj/internal/domain/notification/repository"
)

type NotificationUsecase interface {
	GetUnreadNotifications(limit int) ([]dto.NotificationResponse, error)
	MarkNotificationAsRead(id string) error
	StartNotificationScheduler()
}

type notificationUsecase struct {
	repo repository.NotificationRepository
}

func NewNotificationUsecase(repo repository.NotificationRepository) NotificationUsecase {
	return &notificationUsecase{repo: repo}
}

func (u *notificationUsecase) GetUnreadNotifications(limit int) ([]dto.NotificationResponse, error) {
	notifs, err := u.repo.GetUnreadNotifs(limit)
	if err != nil {
		return nil, err
	}

	var res []dto.NotificationResponse
	for _, n := range notifs {
		res = append(res, dto.NotificationResponse{
			ID:          n.ID,
			Title:       n.Title,
			Message:     n.Message,
			Type:        n.Type,
			ReferenceID: n.ReferenceID,
			IsRead:      n.IsRead,
			CreatedAt:   n.CreatedAt,
		})
	}
	
	if res == nil {
		res = []dto.NotificationResponse{}
	}
	
	return res, nil
}

func (u *notificationUsecase) MarkNotificationAsRead(id string) error {
	return u.repo.MarkAsRead(id)
}
