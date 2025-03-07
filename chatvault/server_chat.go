package chatvault

import (
	"context"
	"database/sql"

	. "github.com/ZolaraProject/chat-vault-service/chatvaultrpc"
	grpctoken "github.com/ZolaraProject/library/grpctoken"
	logger "github.com/ZolaraProject/library/logger"
	_ "github.com/lib/pq"
)

func (*server) CreateConversation(ctx context.Context, conversation *ConversationRequest) (*ConversationResponse, error) {
	grpcToken := grpctoken.GetToken(ctx)

	db, err := sql.Open("postgres", DbUrl())
	if err != nil {
		logger.Err(grpcToken, "Open error : %v", err)
		return nil, err
	}
	defer db.Close()

	query := "INSERT INTO user_conversations (user_id, title) VALUES ($1, $2) RETURNING id"

	var convId sql.NullInt64
	err = db.QueryRow(query, conversation.UserId, "temp").Scan(&convId)
	if err != nil {
		logger.Err(grpcToken, "failed to create conversation: %s", err)
		return nil, err
	}

	return &ConversationResponse{
		Id: convId.Int64,
	}, nil
}

func (*server) GetConversationList(ctx context.Context, conversationList *ConversationListRequest) (*ConversationListResponse, error) {
	grpcToken := grpctoken.GetToken(ctx)

	db, err := sql.Open("postgres", DbUrl())
	if err != nil {
		logger.Err(grpcToken, "Open error : %v", err)
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, title FROM user_conversations WHERE user_id = $1"

	rows, err := db.Query(query, conversationList.UserId)
	if err != nil {
		logger.Err(grpcToken, "failed to get conversation list: %s", err)
		return nil, err
	}
	defer rows.Close()

	var conversations []*ConversationResponse
	for rows.Next() {
		var conv ConversationResponse
		err = rows.Scan(&conv.Id, &conv.Title)
		if err != nil {
			logger.Err(grpcToken, "failed to scan conversation list: %s", err)
			return nil, err
		}

		conversations = append(conversations, &conv)
	}

	return &ConversationListResponse{
		Conversations: conversations,
	}, nil
}

func (*server) SaveChatMessage(ctx context.Context, message *ChatMessageRequest) (*SaveChatMessageResponse, error) {
	grpcToken := grpctoken.GetToken(ctx)

	db, err := sql.Open("postgres", DbUrl())
	if err != nil {
		logger.Err(grpcToken, "Open error : %v", err)
		return nil, err
	}
	defer db.Close()

	query := "INSERT INTO conversations_messages (user_conversations_id, content, message_type) VALUES ($1, $2, $3)"

	_, err = db.Exec(query, message.ConversationId, message.Content, MessageTypes_name[int32(message.MessageType)])
	if err != nil {
		logger.Err(grpcToken, "failed to save chat message: %s", err)
		return nil, err
	}

	return &SaveChatMessageResponse{
		Message: message.Content,
	}, nil
}

func (*server) GetChatMessageList(ctx context.Context, messageList *ChatMessageListRequest) (*ChatMessageListResponse, error) {
	grpcToken := grpctoken.GetToken(ctx)

	db, err := sql.Open("postgres", DbUrl())
	if err != nil {
		logger.Err(grpcToken, "Open error : %v", err)
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, content, message_type, created_at FROM conversations_messages WHERE user_conversations_id = $1"

	rows, err := db.Query(query, messageList.ConversationId)
	if err != nil {
		logger.Err(grpcToken, "failed to get chat message list: %s", err)
		return nil, err
	}
	defer rows.Close()

	var messages []*ChatMessageStrResponse
	for rows.Next() {
		var message ChatMessageStrResponse
		err = rows.Scan(&message.Id, &message.Message, &message.MessageType, &message.CreatedAt)
		if err != nil {
			logger.Err(grpcToken, "failed to scan chat message list: %s", err)
			return nil, err
		}

		messages = append(messages, &message)
	}

	return &ChatMessageListResponse{
		Messages: messages,
	}, nil
}
