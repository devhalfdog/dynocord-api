package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/devhalfdog/dynocord-api/handlers"
	"github.com/devhalfdog/dynocord-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	/* Handlers Initial */
	err := handlers.Initial()
	if err != nil {
		log.Fatal(err)
	}

	/* Fiber Instance Create */
	app := fiber.New(fiber.Config{})

	/* Middleware */
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(func(ctx *fiber.Ctx) error {
		ip := ctx.IP()
		if xff := ctx.Get(fiber.HeaderXForwardedFor); xff != "" {
			ips := strings.Split(xff, ", ")
			if len(ips) > 0 {
				ip = ips[0]
			}

			log.Printf("%s - %s : %s", ip, ctx.Method(), ctx.Path())
		}

		return ctx.Next()
	})

	/* Route */
	// 스트리머 채팅 최근 10개를 불러옴.
	app.Get("/chat/:streamer", handlers.GetChat)
	// 스트리머 채팅을 before 기준으로 최근 10개를 가져옴.
	app.Get("/chat/:streamer/before::before<int>", handlers.GetChat)
	// 스트리머 채팅을 전송받음.
	app.Post("/chat", handlers.SaveChat)
	// 스트리머 이름을 전송받아, 해당 스트리머의 방송 화면을 전송함.
	app.Get("/screen/:streamer", handlers.GetStreamCapture)
	// 정적 이미지 파일을 제공함.
	app.Static("/static/image", "./static")

	/* Server Start */
	err = app.Listen(fmt.Sprintf(":%s", utils.Environment("API_PORT")))

	if err != nil {
		log.Fatal(err)
	}
}
