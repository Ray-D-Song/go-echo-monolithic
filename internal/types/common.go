package types

type PageRequestBox struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type PageResponseBox[T any] struct {
	List  []*T  `json:"list"`
	Total int64 `json:"total"`
}
