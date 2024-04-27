package usecase

import (
	"context"

	"github.com/hdkef/dsi-imbal-summary-bot/domain/entity"
)

type ImbalSummaryUC interface {
	GetImbalHasil(ctx context.Context, dto entity.ImbalSummaryDto) (*entity.ImbalSummary, error)
}
