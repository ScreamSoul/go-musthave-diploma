package restapi

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type AccuralAPIRepository struct {
	actualSystemAddress string
}

func NewAccuralAPIRepository(actualSystemAddress string) *AccuralAPIRepository {
	return &AccuralAPIRepository{actualSystemAddress: actualSystemAddress}
}

func (r *AccuralAPIRepository) request(
	ctx context.Context,
	method,
	path string,
	resultStruct interface{},
) (*resty.Response, error) {

	req := resty.New().R().SetContext(ctx)

	if resultStruct != nil {
		req = req.SetResult(resultStruct)
	}

	req.Method = method
	req.URL = r.actualSystemAddress + path

	resp, err := req.Send()
	return resp, err
}
