package products

import "context"

type adapter struct {
	svc Service
}

func newAdapter(svc Service) Adapter {
	return &adapter{
		svc: svc,
	}
}

func (a *adapter) GetByID(ctx context.Context, id string) (Product, error) {
	response, err := a.svc.GetById(ctx, id)

	if err != nil {
		return Product{}, err
	}

	return response.Data, nil
}
