package main

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"log"

	gosp "github.com/SparkPost/gosparkpost"

	"fmt"

	"github.com/joho/godotenv"
)

func DotEnv(key string) string {
	env := make(chan string, 1)

	if os.Getenv("GO_ENV") != "production" {
		godotenv.Load(".env")
		env <- os.Getenv(key)
	} else {
		 env <- os.Getenv(key)
	}

	return <-env
}

func main(){
	app := fiber.New()

	//cors
	app.Use(cors.New())

	//config sparkpost
	cfg := &gosp.Config{ApiKey: DotEnv("SPARKPOST_API_KEY")}
	var sp gosp.Client 
	sp.Init(cfg)
	
	//sendmail route
	app.Post("/sendmail", func(c *fiber.Ctx) error {
		payload := struct {
			Name string `json:"name"`
			Email string `json:"email"`
			Subject string `json:"subject"`
			Message string `json:"message"`
			
		}{}

		if err := c.BodyParser(&payload); err != nil {
			log.Fatalln(err)
			return err
		}

		//setup mail data
		html := fmt.Sprintf(`You got an email from <br/>%s`, payload.Name)
		content := gosp.Content{
			From: DotEnv("EMAIL"),
			Subject: payload.Subject,
			HTML: html,

		}

		//send mail
		sender := &gosp.Transmission{
			Content: content,
			Recipients: []string{"sayeedmondal1412@gmail.com"},
		}

		id, _, err := sp.Send(sender)

		if err != nil {
			log.Fatal(err)
		}




		return c.JSON(id)
	})

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "3001"
	}
  	app.Listen(":" + PORT)
}