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

func (uc *chatUsecase) PostMessage(ctx context.Context, username, content string) error {
	msg := &domain.Message{Username: username, Content: content}
	return uc.repo.Create(ctx, msg)
}

func (uc *chatUsecase) GetMessages(ctx context.Context) ([]*domain.Message, error) {
	return uc.repo.List(ctx)
}
