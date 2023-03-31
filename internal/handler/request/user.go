package request

import "strings"

type UserCreate struct {
	Name string `json:"name"`
}

func (r *UserCreate) Normalise() {
	r.Name = strings.TrimSpace(r.Name)
}

func (r *UserCreate) Validate() error {
	if r.Name == "" {
		return BlankFieldError{"name"}
	}

	return nil
}

type UserUpdate struct {
	Name string `json:"name"`
}

func (r *UserUpdate) Normalise() {
	r.Name = strings.TrimSpace(r.Name)
}

func (r *UserUpdate) Validate() error {
	if r.Name == "" {
		return BlankFieldError{"name"}
	}

	return nil
}
