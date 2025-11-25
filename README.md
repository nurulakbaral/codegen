# Introduction

## Install

```sh
go get github.com/nurulakbaral/codegen@<latest_commit_hash>
```

## Usage

```go
import (
	"github.com/samber/lo"
   codegen "github.com/nurulakbaral/codegen/gen_domain"
)

func main() {
	root := "app"
	arg := "user"

	config := Config{
		Root:   root,
		Domain: arg,
		Force:  false, // @Notes This feature is not yet implemented.
		Dirs: []Dir{
			Dir{Input: "/templates/domain/entity", Output: filepath.Join(root, arg)},
			Dir{Input: "/templates/domain/handler", Output: filepath.Join(root, arg)},
			Dir{Input: "/templates/domain/repository", Output: filepath.Join(root, arg)},
			Dir{Input: "/templates/domain/service", Output: filepath.Join(root, arg)},
		},
		Data: MockTemplateData{
			ModuleName:       "github.com/nurulakbaral/codegen",
			PascalDomainName: lo.PascalCase(arg),
			LowerDomainName:  arg,
		},
	}
	
	gen := codegen.New(config)
	err := gen.Generate()
}
```
