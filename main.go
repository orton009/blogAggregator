package main

import (
	"blogAggregator/internal/database"
	"database/sql"
	"fmt"
	"os"

	"internal/config"

	_ "github.com/lib/pq"
)

func exitWithError(e error) {
	fmt.Println(e)
	fmt.Println("Exiting....")
	os.Exit(1)
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		exitWithError(fmt.Errorf("error reading configuration file, make sure that gatorconfig.json exists in home directory %w", err))
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		exitWithError(fmt.Errorf("failed to open db connection\n %w", err))
	}
	dbQueries := database.New(db)

	s := state{config: &cfg, db: dbQueries}
	commands := commands{
		CommandMap: map[string](func(*state, command) error){},
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerList)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerCreateFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", handlerFollowing)
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	args := os.Args
	if len(args) < 2 {
		exitWithError(fmt.Errorf("expected command argument"))
	}

	command := command{
		Name: args[1],
		Args: args[2:],
	}

	if err := commands.run(&s, command); err != nil {
		exitWithError(fmt.Errorf("error executing command: %s \n %w", command.Name, err))
	}
}
