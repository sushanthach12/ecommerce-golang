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
	return a.svc.GetById(ctx, id)
}
