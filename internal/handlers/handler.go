package handlers

import "merch_shop/internal/service"

type Controller struct {
	service service.MerchShopService
}

func New(s service.MerchShopService) *Controller {
	return &Controller{
		service: s,
	}
}
