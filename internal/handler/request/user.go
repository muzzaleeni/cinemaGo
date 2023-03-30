package request

type UserCreate struct {
	Name string `json:"name"`
}

func (r UserCreate) Validate() error {
	if r.Name == "" {
		return BlankFieldError{"name"}
	}

	return nil
}

type UserUpdate struct {
	Name string `json:"name"`
}

func (r UserUpdate) Validate() error {
	if r.Name == "" {
		return BlankFieldError{"name"}
	}

	return nil
}
