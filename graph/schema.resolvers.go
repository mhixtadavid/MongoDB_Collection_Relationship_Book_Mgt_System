package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"RelationalMDBGql/database"
	"RelationalMDBGql/graph/model"
	"context"
)

var db = database.Connect()

// CreateBook is the resolver for the CreateBook field.
func (r *mutationResolver) CreateBook(ctx context.Context, input model.BookInput) (*model.Book, error) {
	return db.CreateBook(input), nil
}

// CreateAuthor is the resolver for the CreateAuthor field.
func (r *mutationResolver) CreateAuthor(ctx context.Context, input model.AuthorInput) (*model.Author, error) {
	return db.CreateAuthor(input), nil
}

// CreatePublisher is the resolver for the CreatePublisher field.
func (r *mutationResolver) CreatePublisher(ctx context.Context, input model.PublisherInput) (*model.Publisher, error) {
	return db.CreatePublisher(input), nil
}

// GetBook is the resolver for the GetBook field.
func (r *queryResolver) GetBook(ctx context.Context, id string) (*model.Book, error) {
	return db.GetBook(id), nil
}

// GetAuthor is the resolver for the GetAuthor field.
func (r *queryResolver) GetAuthor(ctx context.Context, id string) (*model.Author, error) {
	return db.GetAuthor(id), nil
}

// GetPublisher is the resolver for the GetPublisher field.
func (r *queryResolver) GetPublisher(ctx context.Context, id string) (*model.Publisher, error) {
	return db.GetPublisher(id), nil
}

// GetAllBooks is the resolver for the GetAllBooks field.
func (r *queryResolver) GetAllBooks(ctx context.Context) ([]*model.Book, error) {
	return db.GetAllBooks(), nil
}

// GetAllAuthors is the resolver for the GetAllAuthors field.
func (r *queryResolver) GetAllAuthors(ctx context.Context) ([]*model.Author, error) {
	return db.GetAllAuthors(), nil
}

// GetAllPublishers is the resolver for the GetAllPublishers field.
func (r *queryResolver) GetAllPublishers(ctx context.Context) ([]*model.Publisher, error) {
	return db.GetAllPublishers(), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}
*/
