package chat

import (
	"ChatRoom/chat/internal/domain"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type mysqlRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) domain.ChatRepository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) Create(ctx context.Context, m *domain.Message) error {
	m.ID = uuid.NewString()
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO messages (id, username, content) VALUES (?, ?, ?)",
		m.Username, m.Content)
	return err
}

func (r *mysqlRepository) List(ctx context.Context) ([]*domain.Message, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, username, content FROM messages ORDER BY id DESC LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*domain.Message
	for rows.Next() {
		m := &domain.Message{}
		if err := rows.Scan(&m.ID, &m.Username, &m.Content); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
