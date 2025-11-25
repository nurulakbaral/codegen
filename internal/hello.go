package internal

import domain "github.com/nurulakbaral/codegen/gen_domain"

func Hello() string {
	p, _ := domain.CreatePath("user")
	return p
}
