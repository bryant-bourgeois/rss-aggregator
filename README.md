## RSS Feed Aggregator

### Getting Started

---

#### Database Setup

- Run a postgres container:
```
 docker run -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
```
- Run database migrations (from ./sql/schema directry):
```
# Up Migrations: 
goose postgres postgres://postgres:postgres@localhost:5432/blogator up

# Down Migrations:
goose postgres postgres://postgres:postgres@localhost:5432/blogator down
```
- If making any database query changes to files in ./sql/queries, generate Go code in the internal/database package with ```sqlc generate``` (from project root).

**Note:** make sure not to edit any of the code generated by sqlc in the ```internal/database``` directory. For conventions like adding json tags to structs generated by sqlc, you can make conversion functions to copy the database structs to your own copies that have appropriate json tags.

#### Run the project (project root directory):
```
go run .
```