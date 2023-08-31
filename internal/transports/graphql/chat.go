package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.33

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skinnykaen/rpa_clone/graph"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/pkg/utils"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CreateChat is the resolver for the CreateChat field.
func (r *mutationResolver) CreateChat(ctx context.Context, userID string) (*models.ChatMutationResult, error) {
	user1ID := ctx.Value(consts.KeyId).(uint)
	user2ID, err := strconv.Atoi(userID)

	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrAtoi,
				},
			},
		}
	}

	chat, err := r.chatService.CreateChat(user1ID, uint(user2ID))

	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}

	return &models.ChatMutationResult{
		ID:      strconv.Itoa(int(chat.ID)),
		User1Id: strconv.Itoa(int(chat.User1ID)),
		User2Id: strconv.Itoa(int(chat.User2ID)),
	}, nil
}

// DeleteChat is the resolver for the DeleteChat field.
func (r *mutationResolver) DeleteChat(ctx context.Context, chatID string) (*models.Response, error) {
	userID := ctx.Value(consts.KeyId).(uint)

	id, err := strconv.Atoi(chatID)

	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": utils.ResponseError{
					Code:    http.StatusBadRequest,
					Message: consts.ErrAtoi,
				},
			},
		}
	}

	if err := r.chatService.DeleteChat(uint(id), userID); err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}

	return &models.Response{Ok: true}, nil
}

// Chats is the resolver for the Chats field.
func (r *queryResolver) Chats(ctx context.Context, page *int, pageSize *int) (*models.ChatsList, error) {
	userID := ctx.Value(consts.KeyId).(uint)

	chats, countRows, err := r.chatService.Chats(userID, page, pageSize)

	if err != nil {
		r.loggers.Err.Printf("%s", err.Error())
		return nil, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"err": err,
			},
		}
	}

	res := make([]*models.ChatHTTP, len(chats))

	for i, chat := range chats {
		var chatHttp models.ChatHTTP
		chatHttp.FromCore(chat)
		res[i] = &chatHttp
	}

	return &models.ChatsList{
		Chats:     res,
		CountRows: int(countRows),
	}, nil
}

// UserJoined is the resolver for the UserJoined field.
func (r *subscriptionResolver) UserJoined(ctx context.Context, userID string, chatID string) (<-chan *models.MessageHTTP, error) {
	panic(fmt.Errorf("not implemented: UserJoined - UserJoined"))
}

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
