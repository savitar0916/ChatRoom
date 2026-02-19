package domain

import "context"

type ChatRepository interface {
	Create(ctx context.Context, msg *Message) error
	List(ctx context.Context) ([]*Message, error)
}
