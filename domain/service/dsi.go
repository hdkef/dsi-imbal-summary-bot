package service

import (
	"context"
	"strings"

	"github.com/hdkef/dsi-imbal-summary-bot/domain/entity"
)

type GetAllPendanaanIdDto struct {
	token string
}

type GetImbalDetailDto struct {
	id    uint32
	token string
}

type Pendanaan struct {
	id uint32
}

func (p *Pendanaan) SetID(id uint32) {
	p.id = id
}

func (p *Pendanaan) GetID() uint32 {
	return p.id
}

func (g *GetAllPendanaanIdDto) SetToken(token string) {
	g.token = strings.ReplaceAll(token, "\n", "")
}

func (g *GetAllPendanaanIdDto) GetToken() string {
	return g.token
}

func (g *GetImbalDetailDto) SetToken(token string) {
	g.token = strings.ReplaceAll(token, "\n", "")
}

func (g *GetImbalDetailDto) GetToken() string {
	return g.token
}

func (g *GetImbalDetailDto) SetID(id uint32) {
	g.id = id
}

func (g *GetImbalDetailDto) GetID() uint32 {
	return g.id
}

type DSIService interface {
	GetAllPendanaanId(ctx context.Context, dto GetAllPendanaanIdDto) ([]Pendanaan, error)
	GetImbalDetail(ctx context.Context, dto GetImbalDetailDto, retriesCount int) ([]entity.Imbal, error)
}
