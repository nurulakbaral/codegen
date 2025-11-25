package domain

import "fmt"

type Domain interface {
	Generate() error
}

type domain struct {
	Config
}

var _ Domain = (*domain)(nil)

func New(c ...Config) *domain {
	gen := &domain{}

	if len(c) > 0 {
		gen.Config = c[0]
	}

	return gen
}

func (c *Config) Generate() error {
	_, err := CheckDir(*c)

	if err != nil {
		return fmt.Errorf("Error helper.CheckDir(): %v. Message: %w", err, ErrDomainExists)
	}

	for _, dir := range c.Dirs {
		inputpath, err := CreatePath(dir.Input)
		if err != nil {
			return fmt.Errorf("Error helper.CreatePath(): Input %v", err)
		}

		outputPath, err := CreatePath(dir.Output)
		if err != nil {
			return fmt.Errorf("Error helper.CreatePath(): Output %v", err)
		}

		pairPaths, err := CreatePairFilePath(c.Domain, inputpath, outputPath)

		if err != nil {
			return fmt.Errorf("Error helper.CreatePairFilePath(): %v", err)
		}

		for _, pairPath := range pairPaths {
			input := pairPath[0]
			output := pairPath[1]

			if err := GenerateTemplateFiles(c.Data, input, output); err != nil {
				return fmt.Errorf("Error GenerateTemplateFiles(): %v", err)
			}
		}

	}

	return nil
}
