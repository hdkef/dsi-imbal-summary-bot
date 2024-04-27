package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hdkef/dsi-imbal-summary-bot/domain/entity"
	"github.com/hdkef/dsi-imbal-summary-bot/domain/service"
	"github.com/hdkef/dsi-imbal-summary-bot/domain/usecase"
	"golang.org/x/sync/errgroup"
)

type ImbalSummaryUsecase struct {
	dsi service.DSIService
}

// GetImbalHasil implements usecase.ImbalSummaryUC.
func (i *ImbalSummaryUsecase) GetImbalHasil(ctx context.Context, dto entity.ImbalSummaryDto) (*entity.ImbalSummary, error) {

	dtoGetAll := service.GetAllPendanaanIdDto{}
	dtoGetAll.SetToken(dto.GetToken())

	// get all pendanaan ids
	pendanaan, err := i.dsi.GetAllPendanaanId(ctx, dtoGetAll)
	if err != nil {
		return nil, err
	}

	fmt.Printf("total pendanaan %d\n", len(pendanaan))

	// get imbal detail for each filtered id

	results := []entity.Imbal{}

	errGrp := &errgroup.Group{}
	mtx := &sync.Mutex{}

	chunkSize := 30
	for idx := 0; idx < len(pendanaan); idx += chunkSize {
		end := idx + chunkSize
		if end > len(pendanaan) {
			end = len(pendanaan)
		}
		chunk := pendanaan[idx:end]
		fmt.Println(chunk)
		for _, v := range chunk {

			id := v.GetID()

			errGrp.Go(func() error {
				dtoImbal := service.GetImbalDetailDto{}
				dtoImbal.SetID(id)
				dtoImbal.SetToken(dto.GetToken())
				imbals, err := i.dsi.GetImbalDetail(ctx, dtoImbal, 0)
				if err != nil {
					return err
				}

				for _, imbal := range imbals {
					// filter by month and year if exist
					if dto.GetStartMonth() != nil && dto.GetEndMonth() != nil {

						year := dto.GetYear()

						// if year not set, set to current year
						if year == nil {
							nowYear := time.Now().Year()
							year = &nowYear
						}

						// if year is match and month > start & month < end
						// appends filtered ids

						if imbal.GetDate().Year() == *year && imbal.GetDate().Month() >= time.Month(*dto.GetStartMonth()) && imbal.GetDate().Month() <= time.Month(*dto.GetEndMonth()) {
							mtx.Lock()
							results = append(results, imbal)
							mtx.Unlock()
						}

					} else {
						nowYear := time.Now().Year()
						if imbal.GetDate().Year() == nowYear {
							mtx.Lock()
							results = append(results, imbal)
							mtx.Unlock()
						}
					}
				}
				return nil
			})

		}

		err = errGrp.Wait()
		if err != nil {
			return nil, err
		}
	}

	sum := &entity.ImbalSummary{}

	for _, v := range results {
		sum.SetImbal(v)
	}

	// sort by date

	return sum, nil
}

func NewImbalSummaryUsecase(dsi service.DSIService) usecase.ImbalSummaryUC {
	return &ImbalSummaryUsecase{
		dsi: dsi,
	}
}
