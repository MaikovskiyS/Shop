package api

import "myproject/internal/domain"

type GetAllResponse struct {
	Result []*domain.Order `json:"orders"`
}
