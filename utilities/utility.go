package utilities

import (
	"encoding/json"
	"strconv"
)

type PaginationParam struct {
	Page    string `json:"page"`
	PerPage string `json:"per_page"`
	Sort    string `json:"sort"`
	Search  string `json:"search"`
	Filter  string `json:"filter"`
}

type Sort struct {
	ComlumnName string `json:"column_name"`
	Value       string `json:"value"`
}

type Filter struct {
	ComlumnName string `json:"column_name"`
	Operator    string `json:"operator"`
	Value       string `json:"value"`
}

type FilterParam struct {
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	Sort    Sort     `json:"sort"`
	Search  string   `json:"search"`
	Filter  []Filter `json:"filter"`
}

func ExtractPagination(param PaginationParam) (FilterParam, error) {
	page, err := strconv.Atoi(param.Page)
	if page <= 0 || err != nil {
		page = 1
	}

	per_page, err := strconv.Atoi(param.PerPage)
	if per_page <= 0 || err != nil {
		per_page = 10
	}

	var sort Sort
	if param.Sort == "" {
		sort.ComlumnName = "created_at"
		sort.Value = "asc"
	} else {
		err := json.Unmarshal([]byte(param.Sort), &sort)
		if err != nil {
			return FilterParam{}, err
		}
	}

	var filter []Filter
	if param.Filter != "" {
		err := json.Unmarshal([]byte(param.Filter), &filter)
		if err != nil {
			return FilterParam{}, err
		}
	}

	return FilterParam{
		Page:    page,
		PerPage: per_page,
		Search:  param.Search,
		Sort:    sort,
		Filter:  filter,
	}, nil
}
