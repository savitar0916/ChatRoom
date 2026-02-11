package chat

import (
	domain "ChatRoom/chat/internal/domain"
	"context"
)

type chatUsecase struct {
	repo domain.ChatRepository
}

func NewChatUsecase(r domain.ChatRepository) domain.ChatUsecase {
	return &chatUsecase{repo: r}
}

func (uc *chatUsecase) PostMessage(ctx context.Context, username, content string) (*domain.Message, error) {
	msg := &domain.Message{Username: username, Content: content}
	if err := uc.repo.Create(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (uc *chatUsecase) GetMessages(ctx context.Context) ([]*domain.Message, error) {
	return uc.repo.List(ctx)
}
