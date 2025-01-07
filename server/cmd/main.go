package main

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sulavpanthi/BloomFilterPasswordChecker/internal/controller/http"
	"github.com/sulavpanthi/BloomFilterPasswordChecker/internal/usecase"
	"github.com/sulavpanthi/BloomFilterPasswordChecker/pkg/appcontext"
)

var context *appcontext.AppContext

func main() {

	if err := appcontext.Initialize(); err != nil {
		panic("Failed to initialize application context: " + err.Error())
	}

	context = appcontext.Get()

	bloomFilterUseCase := usecase.InitBloomFilterUseCase()
	bloomFilterHandler := controller.NewBloomFilterHandler(bloomFilterUseCase)

	app := gin.Default()

	app.POST("/add", bloomFilterHandler.AddPassword)

	app.POST("/check", bloomFilterHandler.CheckPassword)

	app.GET("/bloom-filter", bloomFilterHandler.GetBloomFilter)

	app.Run(":8000")
	context.Logger.Info().Msg("Server is running on port 8000...")
}
