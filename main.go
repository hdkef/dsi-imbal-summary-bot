package main

import (
	"github.com/hdkef/dsi-imbal-summary-bot/internal/delivery"
	"github.com/hdkef/dsi-imbal-summary-bot/internal/service"
	"github.com/hdkef/dsi-imbal-summary-bot/internal/usecase"
)

func main() {

	// init service
	dsiSvc := service.NewDSIService()

	// init usecase
	imbalSummaryUsecase := usecase.NewImbalSummaryUsecase(dsiSvc)

	// handler
	handlerCli := delivery.NewHandleCLI(imbalSummaryUsecase)
	handlerCli.Handle()
}
