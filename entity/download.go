package entity

type DlDir struct {
	Username string `json:"username" validate:"nonzero,nonnil"`
	Dir      string `json:"dir" validate:"nonzero,nonnil"`
}

type DlFiles struct {
	Username string   `json:"username" validate:"nonzero,nonnil"`
	Files    []string `json:"files" validate:"nonzero,nonnil"`
}
