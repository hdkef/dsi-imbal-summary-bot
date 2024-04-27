package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hdkef/dsi-imbal-summary-bot/domain/entity"
	"github.com/hdkef/dsi-imbal-summary-bot/domain/service"
)

type pendanaan struct {
	ID uint32 `json:"id"`
}

type getAllPendanaanResponse struct {
	Pendanaan []pendanaan `json:"pendanaan"`
}

type detailImbal struct {
	ID               uint32  `json:"id"`
	TanggalPayOut    string  `json:"tanggal_payout"`
	Total            float64 `json:"total"`
	KeteranganPayout string  `json:"keterangan_payout"`
}

type detailImbalResponse struct {
	Data []detailImbal `json:"data"`
}

type DSISvc struct {
}

// GetAllPendanaanId implements service.DSIService.
func (d *DSISvc) GetAllPendanaanId(ctx context.Context, dto service.GetAllPendanaanIdDto) ([]service.Pendanaan, error) {

	req, err := http.NewRequest(http.MethodPost, "https://api.danasyariah.id/newProfile/homeLogin", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+dto.GetToken())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Request-Access", dto.GetToken())
	req.Header.Set("Host", "api.danasyariah.id")
	req.Header.Set("Origin", "https://lender.danasyariah.id")
	req.Header.Set("Referer", "https://lender.danasyariah.id/")

	client := &http.Client{}

	fmt.Println("about to hit get all pendanaan..")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	result := getAllPendanaanResponse{}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error(), string(bodyBytes))
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		fmt.Println(string(bodyBytes))
		return nil, err
	}

	pendanaan := []service.Pendanaan{}

	for _, v := range result.Pendanaan {
		p := service.Pendanaan{}
		p.SetID(v.ID)
		pendanaan = append(pendanaan, p)
	}

	fmt.Println("get all pendanaan succeed..")

	return pendanaan, nil
}

// GetImbalDetail implements service.DSIService.
func (d *DSISvc) GetImbalDetail(ctx context.Context, dto service.GetImbalDetailDto, retriesCount int) ([]entity.Imbal, error) {

	payload := struct {
		PendanaanID uint32 `json:"pendanaan_id"`
	}{
		PendanaanID: uint32(dto.GetID()),
	}

	payloadJSON, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, "https://api.danasyariah.id/v1/lender/imbal_hasil_pendanaan", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+dto.GetToken())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Request-Access", dto.GetToken())
	req.Header.Set("Host", "api.danasyariah.id")
	req.Header.Set("Origin", "https://lender.danasyariah.id")
	req.Header.Set("Referer", "https://lender.danasyariah.id/")

	client := &http.Client{}

	fmt.Printf("about to hit imbal detail with id : %d...\n", dto.GetID())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	maxRetries := 3
	secondBeforeRetries := 30
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			if retriesCount < maxRetries {
				fmt.Printf("too many request, will be retried in %d seconds\n", secondBeforeRetries)
				time.Sleep(time.Duration(secondBeforeRetries) * time.Second)
				return d.GetImbalDetail(ctx, dto, retriesCount+1)
			} else {
				return nil, fmt.Errorf("imbal detail with id : %d failed with status code %d", dto.GetID(), resp.StatusCode)
			}
		}
		return nil, fmt.Errorf("imbal detail with id : %d failed with status code %d", dto.GetID(), resp.StatusCode)
	}

	if resp.Body == nil {
		return nil, nil
	}

	result := detailImbalResponse{}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error(), string(bodyBytes))
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		fmt.Println(string(bodyBytes))
		return nil, err
	}

	imbal := []entity.Imbal{}

	for _, v := range result.Data {
		// append other than dana pokok
		if !strings.Contains(v.KeteranganPayout, "Dana Pokok") {
			i := entity.Imbal{}
			i.SetID(v.ID)

			date, err := time.Parse("02-01-2006", v.TanggalPayOut)
			if err != nil {
				return nil, err
			}

			// i.SetDate(v.TanggalPayOut)
			i.SetAmount(v.Total)
			i.SetDate(date)
			i.SetID(v.ID)
			imbal = append(imbal, i)
		}
	}

	fmt.Printf("imbal detail with id : %d succeed...\n", dto.GetID())

	return imbal, nil
}

func NewDSIService() service.DSIService {
	return &DSISvc{}
}
