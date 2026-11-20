package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/john22n/gator/internal/config"
	"github.com/john22n/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	Args := os.Args
	if len(Args) < 2 {
		log.Fatal("usage: cli <command> [args...]")
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	dbConfig := database.New(db)

	programState := &state{cfg: &cfg, db: dbConfig}

	cmd := command{name: Args[1], args: Args[2:]}
	cmds := commands{registeredCommands: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetAllFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerAddFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if err := cmds.run(programState, cmd); err != nil {
		log.Fatal(err)
	}
}
