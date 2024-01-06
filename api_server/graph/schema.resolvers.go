package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"fmt"
	"go-rustdesk-server/api_server/graph/graph_model"
	"go-rustdesk-server/data_server"
	"go-rustdesk-server/model"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*graph_model.Info, error) {
	db, err := data_server.GetDataSever()
	if err != nil {
		return nil, err
	}
	user, err := db.GetUserByName(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not exist")
	}
	if user.Password != password {
		return nil, fmt.Errorf("password error")
	}
	token, err := db.GenToken(username)
	if err != nil {
		return nil, err
	}
	return &graph_model.Info{
		UUID:     user.Uid,
		Username: username,
		IsAdmin:  user.IsAdmin,
		Token:    token,
	}, nil
}

// SelfInfo is the resolver for the selfInfo field.
func (r *queryResolver) SelfInfo(ctx context.Context) (*graph_model.Info, error) {
	db, err := data_server.GetDataSever()
	if err != nil {
		return nil, err
	}
	user := ctx.Value("user").(*model.User)
	user, err = db.GetUserByName(user.Name)
	if err != nil {
		return nil, err
	}
	return &graph_model.Info{
		UUID:     user.Uid,
		Username: user.Name,
		IsAdmin:  user.IsAdmin,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }