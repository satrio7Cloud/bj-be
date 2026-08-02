package qris

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"bj/internal/config"
)

type QrisGenerator interface {
	GenerateDynamicQris(amount float64, useUniqueCode bool) (*QrisResponseData, error)
}

type qrisGenerator struct {
	config *config.Config
}

func NewQrisGenerator(cfg *config.Config) QrisGenerator {
	return &qrisGenerator{
		config: cfg,
	}
}

type QrisResponseData struct {
	QrisPayload string  `json:"qrisPayload"`
	Amount      float64 `json:"amount"`
}

func (g *qrisGenerator) GenerateDynamicQris(amount float64, useUniqueCode bool) (*QrisResponseData, error) {
	if g.config.Qris.StaticPayload == "" {
		return nil, errors.New("QRIS_STATIC_PAYLOAD is not configured")
	}

	finalAmount := int64(amount)
	if useUniqueCode {
		finalAmount += int64(rand.Intn(99) + 1)
	}

	payload, err := convertStaticToDynamic(g.config.Qris.StaticPayload, finalAmount)
	if err != nil {
		return nil, err
	}

	return &QrisResponseData{
		QrisPayload: payload,
		Amount:      float64(finalAmount),
	}, nil
}

func convertStaticToDynamic(staticPayload string, amount int64) (string, error) {
	nodes, err := ParseTLV(strings.TrimSpace(staticPayload), true)
	if err != nil {
		return "", err
	}

	nodes = RemoveTag(nodes, "63")
	nodes = UpsertTag(nodes, "01", "12", map[string]bool{
		"26": true, "27": true, "28": true, "29": true, "30": true, "31": true,
	})
	nodes = UpsertTag(nodes, "54", strconv.FormatInt(amount, 10), map[string]bool{
		"58": true, "59": true, "60": true,
	})

	return AppendCRC(SerializeTLV(nodes)), nil
}
