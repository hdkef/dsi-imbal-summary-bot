package delivery

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/hdkef/dsi-imbal-summary-bot/domain/entity"
	"github.com/hdkef/dsi-imbal-summary-bot/domain/usecase"
)

type CLIHandler struct {
	imbalSummaryUC usecase.ImbalSummaryUC
}

func NewHandleCLI(imbalSummaryUC usecase.ImbalSummaryUC) *CLIHandler {
	return &CLIHandler{
		imbalSummaryUC: imbalSummaryUC,
	}
}

func (c *CLIHandler) Handle() {

	// prompt jwt token
	fmt.Println("Please insert your JWT token:")
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	dto := entity.ImbalSummaryDto{}
	dto.SetToken(token)

	// execute usecase
	result, err := c.imbalSummaryUC.GetImbalHasil(context.TODO(), dto)
	if err != nil {
		fmt.Println("Error executing usecase:", err.Error())
		return
	}

	// print result
	v, grandTotal := result.GetImbal()
	for _, v := range v {
		fmt.Println("###################")
		fmt.Printf("date : %s\n", v.Key)
		fmt.Printf("total : %f\n", v.GetTotal())
		fmt.Println("list : ")
		for _, k := range v.DailySummaryResult {
			fmt.Println("-------")
			fmt.Printf("date : %s\n", k.Key)
			fmt.Printf("total : %f\n", k.Total)
			fmt.Println("list : ")
			for i, l := range k.Imbals {
				fmt.Printf("%d. Amount : %f\n", i+1, l.GetAmount())
			}
			fmt.Println("-------")
		}
		fmt.Println("###################")
	}

	fmt.Printf("Grand Total : %f\n", grandTotal)
}
