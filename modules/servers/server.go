package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/SurrealGod/jelly-shop-code/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
}
type server struct {
	app *fiber.App
	cfg config.IConfig
	db *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	 
	return &server{
		cfg: cfg,
		db: db,
		app: fiber.New(fiber.Config{
			AppName: cfg.App().Name(),
			BodyLimit: cfg.App().BodyLimit(),
			ReadTimeout: cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder: json.Marshal,
		    JSONDecoder: json.Unmarshal,   
		}),
	}
}

func (s *server) Start() { 
     //Middlewares
    middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())
	

	 //modules
	 //http:localhost:3000/
	 v1 := s.app.Group("v1")

     modules := InitModule(v1, s, middlewares)

	 modules.MonitorModule()

	 s.app.Use(middlewares.RouterCheck())

    //Geaceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_= <-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()

	}()
	
	//listen to host
	log.Printf("server is starting on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}