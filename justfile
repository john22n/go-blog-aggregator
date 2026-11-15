DbURL := "postgres://johnnoriega:@localhost:5432/gator"

recipe-name:
    echo 'This is a recipe!'

# this is a comment
another-recipe:
    @echo 'This is another recipe.'

migration-up:
    goose -dir ./sql/schema up

migration-down:
    goose -dir ./sql/schema down

db-connect:
    psql {{ DbURL }}

gen:
    sqlc generate
