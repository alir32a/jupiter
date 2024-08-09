package handler

import "github.com/alir32a/jupiter/internal/model"

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
}

func toCtrlPagination(req model.Pagination) Pagination {
	return Pagination{
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
		Total:       req.Total,
		TotalPages:  req.TotalPages,
	}
}
