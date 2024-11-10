package database

import (
	"RelationalMDBGql/graph/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString string = "mongodb://127.0.0.1:27017"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		fmt.Println("Error:", err)
	}

	_ = client.Ping(ctx, readpref.Primary())

	return &DB{client: client}
}

func (db *DB) GetAllAuthors() []*model.Author {
	AuthorCollections := db.client.Database("RelationalMDBGQL").Collection("Authors")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Step 1: Lookup Books related to Author's bookIds
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Books"},
			{Key: "localField", Value: "bookIDs"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "books"},
		}}},
	}

	// Execute the aggregation
	cur, err := AuthorCollections.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("Failed to execute aggregation:", err)
		return nil
	}
	defer cur.Close(ctx)

	// Decode the results
	var authors []*model.Author
	for cur.Next(ctx) {
		var author model.Author
		if err := cur.Decode(&author); err != nil {
			fmt.Println("Error decoding author:", err)
			continue
		}
		fmt.Printf("Author: %+v\n", author) // Print each author for verification
		authors = append(authors, &author)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Cursor error:", err)
	}

	return authors
}

func (db *DB) GetAllPublishers() []*model.Publisher {
	PublisherCollections := db.client.Database("RelationalMDBGQL").Collection("Publishers")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Step 1: Lookup Books related to Publisher's bookIds
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Books"},
			{Key: "localField", Value: "bookIDs"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "books"},
		}}},
	}

	cur, err := PublisherCollections.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("Failed to execute aggregation:", err)
		return nil
	}
	defer cur.Close(ctx)

	// Decode the results
	var publishers []*model.Publisher
	for cur.Next(ctx) {
		var publisher model.Publisher
		if err := cur.Decode(&publisher); err != nil {
			fmt.Println("Error decoding publisher:", err)
			continue
		}
		publishers = append(publishers, &publisher)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Cursor error:", err)
	}

	return publishers
}

func (db *DB) GetAllBooks() []*model.Book {
	BookCollections := db.client.Database("RelationalMDBGQL").Collection("Books")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//    // Define the aggregation pipeline
	pipeline := mongo.Pipeline{

		// bson.D{
		// 	{Key: "$match", Value: bson.D{}},
		// },
		// Stage 1: Lookup authors
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Authors"},
				{Key: "localField", Value: "authorIds"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "authors"},
			}},
		},
		// Stage 2: Lookup publisher
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Publishers"},
				{Key: "localField", Value: "publisherId"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "publisher"},
			}},
		},
		// Stage 3: Unwind the publisher array
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$publisher"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
	}

	cur, err := BookCollections.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("failed to execute aggregation:", err)
	}
	defer cur.Close(ctx)

	var books []*model.Book
	for cur.Next(ctx) {
		var book model.Book
		if err := cur.Decode(&book); err != nil {
			fmt.Println("Error decoding book:", err)
			continue
		}
		books = append(books, &book)
	}

	fmt.Println("Books output:")
	for i, book := range books {
		fmt.Printf("Book %d: %+v\n", i+1, *book) // Prints book details to verify population
	}

	if err := cur.Err(); err != nil {
		fmt.Println("cursor error:", err)
	}
	return books
}

func (db *DB) GetAuthor(id string) *model.Author {
	AuthorCollections := db.client.Database("RelationalMDBGQL").Collection("Authors")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match the specific book ID
		bson.D{
			{Key: "$match", Value: bson.D{{Key: "_id", Value: _id}}},
		},
		// Step 2: Lookup Books related to Author's bookIds
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Books"},
			{Key: "localField", Value: "bookIDs"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "books"},
		}}},
	}

	// Execute the aggregation
	cur, err := AuthorCollections.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("Failed to execute aggregation:", err)
		return nil
	}
	defer cur.Close(ctx)

	// Decode the result (expecting a single document for this ID)
	var author model.Author
	if cur.Next(ctx) {
		if err := cur.Decode(&author); err != nil {
			fmt.Println("Error decoding author:", err)
		}
	} else {
		fmt.Println("Author not found for this ID:", id)
	}
	// Check for any cursor error
	if err := cur.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		return nil
	}

	return &author
}

func (db *DB) GetPublisher(id string) *model.Publisher {
	PublisherCollections := db.client.Database("RelationalMDBGQL").Collection("Publishers")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match the specific book ID
		bson.D{
			{Key: "$match", Value: bson.D{{Key: "_id", Value: _id}}},
		},

		// Step 2: Lookup Books related to Publisher's bookIds
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Books"},
			{Key: "localField", Value: "bookIDs"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "books"},
		}}},
	}

	cur, err := PublisherCollections.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("Failed to execute aggregation:", err)
		return nil
	}
	defer cur.Close(ctx)

	// Decode the result (expecting a single document for this ID)
	var publisher model.Publisher
	if cur.Next(ctx) {
		if err := cur.Decode(&publisher); err != nil {
			fmt.Println("Error decoding publisher:", err)
			return nil
		}
	} else {
		fmt.Println("Publisher not found for ID:", id)
		return nil
	}

	// Check for any cursor error
	if err := cur.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		return nil
	}

	return &publisher
}

func (db *DB) GetBook(id string) *model.Book {
	BookCollections := db.client.Database("RelationalMDBGQL").Collection("Books")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match the specific book ID
		bson.D{
			{Key: "$match", Value: bson.D{{Key: "_id", Value: _id}}},
		},
		// Stage 2: Lookup authors
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Authors"},
				{Key: "localField", Value: "authorIds"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "authors"},
			}},
		},
		// Stage 3: Lookup publisher
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "Publishers"},
				{Key: "localField", Value: "publisherId"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "publisher"},
			}},
		},
		// Stage 4: Unwind the publisher array
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$publisher"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
	}

	cur, err := BookCollections.Aggregate(ctx, pipeline)
	if err != nil {
		fmt.Println("Failed to execute aggregation:", err)
	}
	defer cur.Close(ctx)

	var Book model.Book
	if cur.Next(ctx) {
		if err := cur.Decode(&Book); err != nil {
			fmt.Println("Error decoding book:", err)
		}
	} else {
		fmt.Println("No book found for this ID")
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Cursor error:", err)
	}

	fmt.Printf("Book details: %+v\n", Book)
	return &Book
}

func (db *DB) CreateAuthor(input model.AuthorInput) *model.Author {
	AuthorCollections := db.client.Database("RelationalMDBGQL").Collection("Authors")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	inserted, err := AuthorCollections.InsertOne(ctx, bson.M{
		"name":  input.Name,
		"email": input.Email,
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	InsertedId := inserted.InsertedID.(primitive.ObjectID).Hex()

	returnAuthor := model.Author{
		ID:    InsertedId,
		Name:  input.Name,
		Email: input.Email,
	}

	return &returnAuthor
}

func (db *DB) CreatePublisher(input model.PublisherInput) *model.Publisher {
	PublisherCollections := db.client.Database("RelationalMDBGQL").Collection("Publishers")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	inserted, err := PublisherCollections.InsertOne(ctx, bson.M{
		"name":     input.Name,
		"location": input.Location,
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	InsertedId := inserted.InsertedID.(primitive.ObjectID).Hex()

	retuenPublisher := model.Publisher{
		ID:       InsertedId,
		Name:     input.Name,
		Location: input.Location,
	}

	return &retuenPublisher
}

func (db *DB) CreateBook(input model.BookInput) *model.Book {
	PublisherCollections := db.client.Database("RelationalMDBGQL").Collection("Publishers")
	AuthorCollections := db.client.Database("RelationalMDBGQL").Collection("Authors")
	BookCollections := db.client.Database("RelationalMDBGQL").Collection("Books")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Convert publisher ID to ObjectID
	publisherObjectId, _ := primitive.ObjectIDFromHex(input.PublisherID)

	// Convert author IDs to ObjectIDs
	authorObjectIds := []primitive.ObjectID{}
	for _, authorID := range input.AuthorIds {
		authorObjectId, _ := primitive.ObjectIDFromHex(authorID)
		authorObjectIds = append(authorObjectIds, authorObjectId)
	}

	// Insert the book into the Books collection
	inserted, err := BookCollections.InsertOne(ctx, bson.M{
		"title":       input.Title,
		"isbn":        input.Isbn,
		"authorIds":   authorObjectIds,
		"publisherId": publisherObjectId,
	})
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Get the inserted Book's ID
	InsertedId := inserted.InsertedID.(primitive.ObjectID)

	// Update each author to include the new book's ID in their Books array
	for _, authorObjectId := range authorObjectIds {
		_, err := AuthorCollections.UpdateOne(
			ctx,
			bson.M{"_id": authorObjectId},
			bson.M{"$push": bson.M{"bookIDs": InsertedId}},
		)
		if err != nil {
			fmt.Println("Error updating author with book ID:", err)
		}
	}

	// Update the publisher to include the new book's ID in their Books array
	_, err = PublisherCollections.UpdateOne(
		ctx,
		bson.M{"_id": publisherObjectId},
		bson.M{"$push": bson.M{"bookIDs": InsertedId}},
	)
	if err != nil {
		fmt.Println("Error updating publisher with book ID:", err)
	}

	// Fetch authors and publisher for the return response
	var authors []*model.Author
	for _, authorObjectId := range authorObjectIds {
		var author model.Author
		err := AuthorCollections.FindOne(ctx, bson.M{"_id": authorObjectId}).Decode(&author)
		if err != nil {
			fmt.Println("Error fetching Author(s):", err)
			continue
		}

		// fetch the books for this author
		var books []*model.Book
		for _, bookID := range author.BookIDs {
			var book model.Book
			err = BookCollections.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
			if err != nil {
				fmt.Println("Error fetching Book(s):", err)
				continue
			}
			books = append(books, &book)
		}
		author.Books = books
		authors = append(authors, &author)
	}

	var publisher model.Publisher
	err = PublisherCollections.FindOne(ctx, bson.M{"_id": publisherObjectId}).Decode(&publisher)
	if err != nil {
		fmt.Println("Error fetching Publisher:", err)
	}

	var books []*model.Book
	// bookIds := publisher.Books
	for _, bookID := range publisher.BookIDs {
		var book model.Book
		err = BookCollections.FindOne(ctx, bson.M{"_id": bookID}).Decode(&book)
		if err != nil {
			fmt.Println("Error fetch Book(s):", err)
			continue
		}
		books = append(books, &book)
	}
	publisher.Books = books

	returnBook := model.Book{
		ID:        InsertedId.Hex(),
		Title:     input.Title,
		Isbn:      input.Isbn,
		Authors:   authors,
		Publisher: &publisher,
	}

	return &returnBook
}
