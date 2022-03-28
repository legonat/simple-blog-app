package main

import (
	"awesomeProjectRentaTeam/internal/config"
	"awesomeProjectRentaTeam/internal/db"
	"awesomeProjectRentaTeam/internal/handler"
	"awesomeProjectRentaTeam/internal/logger"
	"awesomeProjectRentaTeam/internal/server"
	"awesomeProjectRentaTeam/internal/service"
	"awesomeProjectRentaTeam/pkg/erx"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.GetConfigInstance()
	log := logger.GetLogrusInstance()

	err := db.InitBlogDb(cfg.BlogDb.Path)
	if err != nil {
		log.Error(err)
	}

	function := flag.String("f", "default", "Specify one of the commands: reinitVmDb, startServer")
	flag.Parse()

	switch *function {
	case "initVmDb":
		err = db.InitBlogDb(cfg.BlogDb.Path)
		if err != nil {
			log.Error(err)
			return
		}
	case "reinitVmDb":
		err = db.ForceInitBlogDb(cfg.BlogDb.Path)
		if err != nil {
			log.Error(err)
			return
		}

	case "startServer":
		dbConn, err := db.NewSqliteDB("./data/blog.db")
		if err != nil {
			log.Error(err)
			return
		}
		defer func(dbConn *sql.DB) {
			err := dbConn.Close()
			if err != nil {
				log.Error(err)
			}
		}(dbConn)

		repository := db.NewRepository(dbConn)
		serv := service.NewService(repository)
		h := handler.NewHandler(serv)
		srv := new(server.Server)
		port := fmt.Sprintf("%v", cfg.Server.Port)
		go func() {
			if err := srv.Run(port, h.InitRoutes()); err != nil {
				log.Error(erx.New(err))
			}
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
		<-quit

		log.Info("Server Shutting Down")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error(erx.New(err))
		}

		if err := dbConn.Close(); err != nil {
			log.Error(erx.New(err))
		}

	default:
		println("Expected flag (-f)")
		flag.PrintDefaults()
	}

}
