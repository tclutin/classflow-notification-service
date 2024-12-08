package service

import (
	"context"
	"fmt"
	"github.com/tclutin/classflow-notification-service/internal/model"
	"github.com/tclutin/classflow-notification-service/internal/repository"
	"github.com/tclutin/classflow-notification-service/pkg/client/telegram"
	"log/slog"
	"sync"
	"time"
)

type NotificationService struct {
	logger     *slog.Logger
	tgClient   telegram.Client
	repository *repository.ScheduleRepository
}

func NewNotificationService(logger *slog.Logger, tgClient telegram.Client, repository *repository.ScheduleRepository) *NotificationService {
	return &NotificationService{
		logger:     logger,
		tgClient:   tgClient,
		repository: repository,
	}
}

func (n *NotificationService) Start(ctx context.Context) {
	go n.SendTelegram(ctx)
}

// SendTelegram TODO: worker pool
func (n *NotificationService) SendTelegram(ctx context.Context) {
	var wg sync.WaitGroup

	for {
		now := time.Now()

		n.logger.Debug("Starting schedule check",
			"time", now.Format("15:04:05"),
		)

		nextMinute := now.Truncate(time.Minute).Add(time.Minute * 1)

		durationUntilNextMinute := nextMinute.Sub(now)

		n.logger.Debug("Waiting until next minute",
			"next_minute", nextMinute.Format("15:04:05"),
		)

		time.Sleep(durationUntilNextMinute)

		weekDay := int(now.Weekday())
		if weekDay == 0 {
			weekDay = 7
		}

		_, week := now.ISOWeek()

		isEven := week%2 != 0

		n.logger.Info("Schedule check",
			"week_day", weekDay,
			"is_even_week", isEven,
		)

		schedules, err := n.repository.FindUpcomingSchedule(ctx, weekDay, isEven)

		if err != nil {
			n.logger.Error("Failed to get schedule",
				"error", err,
				"day_of_week", weekDay,
				"is_even_week", isEven)
			continue
		}
		n.logger.Debug("found schedules",
			"count", len(schedules),
		)

		for _, schedule := range schedules {
			startProcessing := time.Now()
			n.logger.Info("Start processing schedule",
				"start_time", startProcessing.Format("15:04:05"),
			)

			startTime := schedule.StartTime.Format("15:04")
			endTime := schedule.EndTime.Format("15:04")

			message := fmt.Sprintf(
				"⏰ Через %d минут\n"+
					"📚 %s\n"+
					"🏫 %s\n"+
					"👨‍🏫 %s\n"+
					"🕒 %s - %s",
				schedule.NotificationDelay,
				schedule.SubjectName,
				schedule.Room,
				schedule.Teacher,
				startTime,
				endTime,
			)

			n.logger.Info("Sending message",
				"telegram_chat", schedule.TelegramChat,
				"message", message,
			)

			wg.Add(1)
			go func(notification model.Notification) {
				defer wg.Done()
				err = n.tgClient.SendMessage(schedule.TelegramChat, message)
				if err != nil {
					n.logger.Error("Error sending message to Telegram chat",
						"error", err,
						"telegram_chat", notification.TelegramChat,
					)
					return
				}
				n.logger.Info("Message sent successfully",
					"telegram_chat", notification.TelegramChat,
				)
			}(schedule)

			n.logger.Info("Finished processing schedule",
				"duration", time.Since(startProcessing).String(),
			)
		}
		wg.Wait()
		n.logger.Info("Finished schedule check",
			"time", time.Now().Format("15:04:05"),
		)
	}
}
