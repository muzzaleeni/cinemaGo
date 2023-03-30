package response

import "math"

type Page struct {
	Info PaginationInfo `json:"pagination"`
	Data any            `json:"data"`
}

type PaginationInfo struct {
	Number     uint64 `json:"page_number"`
	Size       uint64 `json:"page_size"`
	TotalPages uint64 `json:"page_count"`
	DataLen    int    `json:"total"`
}

func NewPage(number, size, totalRecords uint64, dataLen int, data any) *Page {
	totalPages := math.Ceil(float64(totalRecords) / float64(size))
	return &Page{
		Info: PaginationInfo{
			Number:     number,
			Size:       size,
			TotalPages: uint64(totalPages),
			DataLen:    dataLen,
		},
		Data: data,
	}
}
