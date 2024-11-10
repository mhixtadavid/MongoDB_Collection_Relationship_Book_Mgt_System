# Book Management System API

A backend API for managing books, authors, and publishers in a MongoDB database using Go and GraphQL. This project allows querying and managing a book library, where each book is associated with authors and publishers.

## Table of Contents
- Overview
- Features
- Tech Stack
- Getting Started
  - Prerequisites
  - Installation
  - Running the Application
- API Endpoints
- Database Schema
- Contributing
- License

---

## Overview
This project is designed to manage books, authors, and publishers with MongoDB as the database backend. Each book references one publisher and can have multiple authors. The API supports retrieving full book details with associated authors and publishers, thanks to MongoDB's aggregation capabilities.

## Features
- **Books Collection**: Stores book details with references to authors and publishers.
- **Authors Collection**: Manages author details with associated book references.
- **Publishers Collection**: Manages publisher details with associated book references.
- **Aggregated Queries**: Uses MongoDB `$lookup` aggregation to retrieve nested data (e.g., books with authors and publisher details).

## Tech Stack
- **Programming Language**: Go
- **Database**: MongoDB
- **GraphQL**: gqlgen
- **API Testing**: GraphQL Playground

## Getting Started

### Prerequisites
Ensure you have the following installed:
- [Go](https://golang.org/doc/install) (version 1.18+ recommended)
- [MongoDB](https://www.mongodb.com/try/download/community)
- [gqlgen](https://github.com/99designs/gqlgen) - A Go library for building GraphQL servers

### Installation
1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/book-management-system.git
   cd book-management-system
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Set Up MongoDB**
   Ensure your MongoDB instance is running and accessible. Update the MongoDB connection URI in the code if necessary.

4. **Configure Environment Variables**
   Create a `.env` file with your configuration settings, including MongoDB URI:
   ```env
   MONGODB_URI="your_mongodb_connection_string"
   ```

### Running the Application
```bash
go run server.go
```
This will start the server, and you can access the GraphQL Playground at `http://localhost:8080` to test the API.

---

## API Endpoints

The API provides the following GraphQL queries and mutations:

### Queries
- `GetAllBooks`: Fetches all books with associated authors and publisher details.
- `GetAllAuthors`: Retrieves all authors with their associated books.
- `GetAllPublishers`: Retrieves all publishers with their associated books.
- `GetPublisher(id: String!)`: Retrieves a specific publisher by ID, along with its books.

### Mutations
- `CreateBook`: Adds a new book with references to authors and a publisher.
- `CreateAuthor`: Adds a new author.
- `CreatePublisher`: Adds a new publisher.

## Database Schema
The project contains three main collections in MongoDB:

- **Books**:
  - `title`: String
  - `authorIds`: Array of ObjectID (References Authors)
  - `publisherId`: ObjectID (References Publisher)
  - Other fields related to book information
  
- **Authors**:
  - `_id`: ObjectID
  - `name`: String
  - `bookIds`: Array of ObjectID (References Books)

- **Publishers**:
  - `_id`: ObjectID
  - `name`: String
  - `bookIds`: Array of ObjectID (References Books)

---

## Contributing
Contributions are welcome! To contribute:
1. Fork the project.
2. Create a new branch for your feature or bugfix.
3. Commit your changes.
4. Submit a pull request.

---

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Acknowledgements
Thanks to MongoDB, gqlgen, and the Go community for their excellent resources and tools.
```







mutation CreateBook {
  CreateBook(input: {
    title: "Hello world!",
    isbn: "1234"
    authorIds: ["67301ee96f4e76b7de66f693"]
    publisherId: "67301efe6f4e76b7de66f694"
  })
    {
      id
      isbn
    }
}

mutation CreateAuthor {
  CreateAuthor (input: {
    name: "Esther O"
    email: "dav@dav.com"
  })
  {
    id
    name
    email
  }
}

mutation CreatePub {
  CreatePublisher (input: {
    name: "Hello"
    location: "Earth"
  }){
    id
    name
    location
  }
}

query GetAuthor {
  GetAllAuthors {
    name
    id
    books{
      id
      title
    }
  }
}

query GetPub {
  GetAllPublishers {
    name
    id
    books {
      id
      title
    }
  }
}

query GetBooks {
  GetAllBooks {
    authors {
      id
    }
    id
    publisher {
      id
    }
  }
}

query GetOneAuthor {
  GetAuthor (id: "67301ee96f4e76b7de66f693") {
    id
    name
    books {
      id
      title
    }
  }
}

query GetOnePub {
  GetPublisher (id: "67301efe6f4e76b7de66f694") {
    id
    name
    books {
      id
      title
    }
  }
}

query GetOneBook {
  GetBook (id: "67301f106f4e76b7de66f695") {
    id
    title
    authors {
      id
      name
    }
    publisher {
      id
      name
    }
  }
}