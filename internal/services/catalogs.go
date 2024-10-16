package services

import (
	"context"
	m "lkserver/internal/models"
	"lkserver/internal/models/catalogs"
	"lkserver/internal/models/types"
	"lkserver/internal/repository"
)

type CatalogsService struct {
	provider *repository.Repo
}

func NewCatalogsService(r *repository.Repo) *CatalogsService {
	return &CatalogsService{provider: r}
}

type Result struct {
	Data any               `json:"data"`
	Len  int               `json:"len"`
	Rows int64             `json:"rows"`
	Meta map[string]m.META `json:"meta"`
}

func (c *CatalogsService) GetCato(ctx context.Context, ref m.JSONByte) (*Result, error) {
	data, err := c.provider.Catalogs.Cato.Get(ctx, ref)
	if err != nil {
		return nil, m.HandleError(err, "CatalogsService.GetCato")
	}
	return &Result{
		Data: data,
		Len:  1,
		Rows: c.provider.Catalogs.Cato.Count(ctx),
		Meta: map[string]m.META{types.Cato: catalogs.CatoMETA},
	}, nil
}

func (c *CatalogsService) CatoList(ctx context.Context, parent m.JSONByte, search string, limits ...int64) (*Result, error) {
	var result []*catalogs.Cato
	var err error
	if search != "" {
		result, err = c.provider.Catalogs.Cato.Find(ctx, search, limits...)
	} else {
		result, err = c.provider.Catalogs.Cato.List(ctx, parent, limits...)
	}
	if err != nil {
		return nil, err
	}
	return &Result{
		Data: result,
		Len:  len(result),
		Rows: c.provider.Catalogs.Cato.Count(ctx),
		Meta: map[string]m.META{types.Cato: catalogs.CatoMETA},
	}, nil
}

func (c *CatalogsService) GetVus(ctx context.Context, ref m.JSONByte) (*Result, error) {
	data, err := c.provider.Catalogs.Vus.Get(ctx, ref)
	if err != nil {
		return nil, m.HandleError(err, "CatalogsService.GetVus")
	}
	return &Result{
		Data: data,
		Len:  1,
		Rows: c.provider.Catalogs.Vus.Count(ctx),
		Meta: map[string]m.META{types.Vus: catalogs.CatoMETA},
	}, nil
}

func (c *CatalogsService) VusList(ctx context.Context, search string, limits ...int64) (*Result, error) {

	var result []*catalogs.Vus
	var err error
	if search != "" {
		result, err = c.provider.Catalogs.Vus.Find(ctx, search, limits...)
	} else {
		result, err = c.provider.Catalogs.Vus.List(ctx, limits...)
	}
	if err != nil {
		return nil, m.HandleError(err, "CatalogService.VusList")
	}

	return &Result{
		Data: result,
		Len:  len(result),
		Rows: c.provider.Catalogs.Vus.Count(ctx),
		Meta: map[string]m.META{types.Vus: catalogs.VusMETA},
	}, nil
}
