package service

import (
	"context"
	"time"
	"github.com/DamikaAlwis-Gif/shorten-url-app/database"
	"github.com/DamikaAlwis-Gif/shorten-url-app/repository"
)

type ClickLogService struct {
	broker repository.MessageBroker
}

func NewClickLogService(broker repository.MessageBroker) *ClickLogService {
  return &ClickLogService{broker: broker}
}

func (c *ClickLogService) LogClick(ctx context.Context, shortURL, ipAddress, userAgent string) error {
	clickLog := database.Click{
		ShortURL: shortURL,
		IPAddress: ipAddress,
    UserAgent: userAgent,
		Timestamp: time.Now(),

	}
	return c.broker.Publish(ctx, "click_logs", clickLog)

}