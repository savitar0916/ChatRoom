package domain

import "context"

type ChatUsecase interface {
	PostMessage(ctx context.Context, username, content string) error
	GetMessages(ctx context.Context) ([]*Message, error)
}
