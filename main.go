package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/alghiffari10/blog_aggregator/internal/commands"
	"github.com/alghiffari10/blog_aggregator/internal/config"
	"github.com/alghiffari10/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	defer db.Close()
	dbQueries := database.New(db)

	programState := &commands.State{
		Db:  dbQueries,
		Cfg: &cfg,
	}

	cmds := commands.Commands{
		RegisteredCommands: make(map[string]func(*commands.State, commands.Command) error),
	}

	// Register the commands
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAggregator)
	cmds.Register("feeds", commands.HandlerFeeds)
	cmds.Register("addfeed", commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	cmds.Register("follow", commands.MiddlewareLoggedIn(commands.HandlerFollows))
	cmds.Register("following", commands.MiddlewareLoggedIn(commands.HandlerFollowing))
	cmds.Register("unfollow", commands.MiddlewareLoggedIn(commands.HandlerUnfollow))
	cmds.Register("browse", commands.MiddlewareLoggedIn(commands.HandlerBrowse))
	cmds.Register("version", commands.HandlerVersion)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.Run(programState, commands.Command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}
