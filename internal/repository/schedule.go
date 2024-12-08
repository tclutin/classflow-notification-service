package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/classflow-notification-service/internal/model"
	"log/slog"
	"time"
)

type ScheduleRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewScheduleRepository(db *pgxpool.Pool, logger *slog.Logger) *ScheduleRepository {
	return &ScheduleRepository{
		db:     db,
		logger: logger,
	}
}

func (s *ScheduleRepository) FindUpcomingSchedule(ctx context.Context, DayOfWeek int, IsEven bool) ([]model.Notification, error) {
	currentTime := time.Now().Truncate(time.Minute)

	sql := `
		SELECT
			u.telegram_chat,
			u.notification_delay,
			s.subject_name,
			s.room,
			s.teacher,
			s.start_time,
			s.end_time
		FROM
			public.users AS u
		JOIN public.members AS m
		    ON m.user_id = u.user_id
	    JOIN public.groups AS g
	        ON g.group_id = m.group_id
		JOIN public.schedule AS s
		    ON g.group_id = s.group_id
		WHERE
			u.notifications_enabled = true
			AND u.notification_delay IS NOT NULL
			AND s.day_of_week = $1
			AND s.is_even = $2
			AND s.start_time > $3
		`

	rows, err := s.db.Query(ctx, sql, DayOfWeek, IsEven, currentTime)
	if err != nil {
		s.logger.Error("Failed to execute query",
			"error", err,
			"day_of_week", DayOfWeek,
			"is_even", IsEven,
		)
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		err = rows.Scan(
			&notification.TelegramChat,
			&notification.NotificationDelay,
			&notification.SubjectName,
			&notification.Room,
			&notification.Teacher,
			&notification.StartTime,
			&notification.EndTime,
		)

		if err != nil {
			s.logger.Error("Failed to scan row",
				"error", err,
			)
			return nil, err
		}

		notificationTime := notification.StartTime.Add(-time.Duration(notification.NotificationDelay) * time.Minute)

		if notificationTime.Format("15:04:05") == currentTime.Format("15:04:05") {
			notifications = append(notifications, notification)
		}

	}

	return notifications, nil
}
