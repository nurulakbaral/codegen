package domain

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
)

type DomainTestSuite struct {
	suite.Suite
	DefaultConfig Config
}

type MockTemplateData struct {
	ModuleName       string
	PascalDomainName string
	LowerDomainName  string
}

func (s *DomainTestSuite) SetupTest() {
	// @Notes If you want to one level up, use `../app`
	root := "app"
	arg := "user"

	config := Config{
		Root:   root,
		Domain: arg,
		Force:  false,
		Dirs: []Dir{
			Dir{Input: "../examples/domain/entity", Output: filepath.Join(root, arg)},
			Dir{Input: "../examples/domain/handler", Output: filepath.Join(root, arg)},
			Dir{Input: "../examples/domain/repository", Output: filepath.Join(root, arg)},
			Dir{Input: "../examples/domain/service", Output: filepath.Join(root, arg)},
		},
		Data: MockTemplateData{
			ModuleName:       "github.com/nurulakbaral/codegen",
			PascalDomainName: lo.PascalCase(arg),
			LowerDomainName:  arg,
		},
	}

	s.DefaultConfig = config
}

func (s *DomainTestSuite) TestDomain() {
	s.Run("Check domain config", func() {
		s.Equal(s.DefaultConfig, s.DefaultConfig)
	})

	s.Run("Copy templates to root dir ", func() {
		s.DefaultConfig.Domain = "user"
		gen := New(s.DefaultConfig)
		err := gen.Generate()

		if !errors.Is(err, ErrDomainExists) && err != nil {
			s.T().Fatalf("Error codegen.Generate(): %v", err)
		}

		if err != nil {
			s.ErrorIs(err, ErrDomainExists)
		}

		expected := map[string][]string{
			"dirs": []string{"entity", "handler", "repository", "service"},
			"files": []string{
				s.DefaultConfig.Domain + "_entity.go",
				s.DefaultConfig.Domain + "_seeds.go",
				s.DefaultConfig.Domain + "_dto.go",
				s.DefaultConfig.Domain + "_handler.go",
				s.DefaultConfig.Domain + "_routes.go",
				"postgre_" + s.DefaultConfig.Domain + "_repository.go",
				s.DefaultConfig.Domain + "_helper.go",
				s.DefaultConfig.Domain + "_service.go",
			},
		}

		domainPath, err := CreatePath("app/user")

		if err != nil {
			s.T().Fatalf("Error helper.CreatePath(): %v", err)
		}

		actual := WalkDirs(domainPath)

		s.Equal(expected, actual)
	})

	s.Run("Copy nested templates to root dir ", func() {
		// @Notes if you change the domain name, order will be change.
		s.DefaultConfig.Domain = "commodity"

		s.DefaultConfig.Dirs = []Dir{
			Dir{Input: "../examples/domain_nested/entity", Output: filepath.Join(s.DefaultConfig.Root, s.DefaultConfig.Domain)},
			Dir{Input: "../examples/domain_nested/handler", Output: filepath.Join(s.DefaultConfig.Root, s.DefaultConfig.Domain)},
			Dir{Input: "../examples/domain_nested/repository", Output: filepath.Join(s.DefaultConfig.Root, s.DefaultConfig.Domain)},
			Dir{Input: "../examples/domain_nested/service", Output: filepath.Join(s.DefaultConfig.Root, s.DefaultConfig.Domain)},
		}

		gen := New(s.DefaultConfig)
		err := gen.Generate()

		if !errors.Is(err, ErrDomainExists) && err != nil {
			s.T().Fatalf("Error codegen.Generate(): %v", err)
		}

		domainPath, err := CreatePath("app", s.DefaultConfig.Domain)

		if err != nil {
			s.T().Fatalf("Error helper.CreatePath(): %v", err)
		}

		expected := map[string][]string{
			"dirs": []string{"entity", "handler", "http", "dto", "rabbitmq", "rtsp", "repository", "service"},
			"files": []string{
				s.DefaultConfig.Domain + "_entity.go",
				s.DefaultConfig.Domain + "_seeds.go",
				s.DefaultConfig.Domain + "_handler.go",
				s.DefaultConfig.Domain + "_routes.go",
				s.DefaultConfig.Domain + "_dto.go",
				".keep",
				".keep",
				"postgre_" + s.DefaultConfig.Domain + "_repository.go",
				s.DefaultConfig.Domain + "_helper.go",
				s.DefaultConfig.Domain + "_service.go",
			},
		}
		actual := WalkDirs(domainPath)

		// @Notes the order is matters.
		// @Refactor Just assert only item, not order.
		s.Equal(expected, actual)

	})

	s.Run("Directory domain should be exist (Domain Auth)", func() {
		s.DefaultConfig.Domain = "auth"
		gen := New(s.DefaultConfig)
		err := gen.Generate()

		s.Error(err)
	})
}

func (s *DomainTestSuite) TestDomainHelper() {
	s.Run("Directory should be not exist (Domain User)", func() {
		s.DefaultConfig.Domain = "heksaldasbas"
		completePath, _ := CreatePath(s.DefaultConfig.Root, s.DefaultConfig.Domain)
		currPath, err := CheckDir(s.DefaultConfig)

		s.T().Logf("currPath %s", currPath)
		s.Equal(currPath, completePath)
		s.Equal(err, nil)
	})

	s.Run("Directory should be exist (Domain Auth)", func() {
		s.DefaultConfig.Domain = "auth"
		completePath, _ := CreatePath(s.DefaultConfig.Root, s.DefaultConfig.Domain)
		currPath, err := CheckDir(s.DefaultConfig)

		s.Equal(currPath, completePath)
		s.Error(err)
	})

	s.Run("Create pair paths should be right", func() {
		s.DefaultConfig.Domain = "user"
		input, err := CreatePath("../examples/domain/entity")
		if err != nil {
			s.T().Fatalf("Error helper.CreatePath(): %v", err)
		}

		output, err := CreatePath("app/user")
		if err != nil {
			s.T().Fatalf("Error helper.CreatePath(): %v", err)
		}

		res, err := CreatePairFilePath(s.DefaultConfig.Domain, input, output)

		if err != nil {
			s.T().Fatalf("Error helper.CreatePairFilePath(): %v", err)
		}

		expected := [][]string{
			[]string{filepath.Join(input, "domain_entity.tmpl"), filepath.Join(output, "entity", s.DefaultConfig.Domain+"_entity.go")},
			[]string{filepath.Join(input, "domain_seeds.tmpl"), filepath.Join(output, "entity", s.DefaultConfig.Domain+"_seeds.go")},
		}

		s.Equal(expected, res)
	})
}

func TestDomainTestSuite(t *testing.T) {
	suite.Run(t, new(DomainTestSuite))
}
