package controllers

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (c *Controllers) RegisterRouters() {
	fiberApp := fiber.New()
	fiberApp.Use(recover.New())

	conf := cors.ConfigDefault
	conf.AllowOrigins = "*"
	conf.AllowCredentials = true
	conf.AllowMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	conf.AllowHeaders = "Content-Type, Accept, Authorization, X-CSRF-Token"
	fiberApp.Use(cors.New(conf))

	fiberApp.Get("/getStocks", c.GetStocks)
	fiberApp.Post("/addStocks", c.CreateStocks)
	fiberApp.Post("/addOrEditStock", c.CreateOrEditStock)

	fiberApp.Post("/daily/saveHistoricalStocks", c.SaveHistoricalStocks)
	fiberApp.Post("/daily/saveHistoricalStock", c.CreateOrEditStock)
	// fiberApp.Get("/daily/getDailyStockBySymbol", c.CreateOrEditStock)
	fiberApp.Get("/daily/getDailyStockByDurations", c.CreateOrEditStock)

	log.Println("Starting application on port", c.Port)
	fiberApp.Listen(fmt.Sprintf(":%s", c.Port))
}
