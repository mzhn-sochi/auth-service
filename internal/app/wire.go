//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/mzhn-sochi/auth-service/internal/config"
)

func Init() (*App, func(), error) {
	panic(
		wire.Build(
			newApp,
			wire.NewSet(config.New),
			//wire.NewSet(initDB),

			//wire.NewSet(pg.NewTicketStorage),
			//wire.NewSet(ticketservice.New),
			//
			//wire.Bind(new(server.TicketService), new(*ticketservice.TicketService)),
			//wire.Bind(new(ticketservice.TicketStorage), new(*pg.TicketStorage)),
			//
			//wire.NewSet(server.New),
		),
	)
}

//func initDB(cfg *config.Config) (*sqlx.DB, func(), error) {
//
//	host := cfg.DB.Host
//	port := cfg.DB.Port
//	user := cfg.DB.User
//	pass := cfg.DB.Pass
//	name := cfg.DB.Name
//
//	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, pass, host, port, name)
//
//	log.Printf("connecting to %s\n", cs)
//
//	db, err := sqlx.Open("postgres", cs)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	return db, func() { db.Close() }, nil
//}
