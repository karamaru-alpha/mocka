package internal

import (
	"fmt"
	"strings"

	"github.com/scylladb/go-set/strset"
	"golang.org/x/tools/go/packages"

	"github.com/karamaru-alpha/mocka/internal/config"
	"github.com/karamaru-alpha/mocka/internal/generator"
	"github.com/karamaru-alpha/mocka/internal/visitor"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	packagePatternSet := strset.NewWithSize(len(cfg.Packages))
	packageConfigMap := make(map[string]*config.Package, len(cfg.Packages))
	for _, packages := range cfg.Packages {
		for pkg, packageConfig := range packages {
			packageConfigMap[pkg] = packageConfig
			pattern := pkg
			if packageConfig.Recursive {
				pattern = fmt.Sprintf("%s/...", pkg)
			}
			packagePatternSet.Add(pattern)
		}
	}
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
	}, packagePatternSet.List()...)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		if strings.Contains(pkg.PkgPath, "/mocks/") {
			continue
		}
		if len(pkg.Errors) > 0 {
			return pkg.Errors[0]
		}

		fileMap := visitor.Visit(pkg, packageConfigMap[pkg.PkgPath])
		if len(fileMap) == 0 {
			continue
		}
		if err := generator.Generate(pkg, fileMap); err != nil {
			return err
		}
	}
	return nil
}
