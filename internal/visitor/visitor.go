package visitor

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/scylladb/go-set/strset"
	"golang.org/x/tools/go/packages"

	"github.com/karamaru-alpha/mocka/internal/config"
	"github.com/karamaru-alpha/mocka/internal/types"
)

type visitor struct {
	pkg     *packages.Package
	config  *config.Package
	fileMap map[string]*types.File
}

func Visit(pkg *packages.Package, cfg *config.Package) map[string]*types.File {
	v := &visitor{
		fileMap: make(map[string]*types.File),
		pkg:     pkg,
		config:  cfg,
	}
	for _, syntax := range pkg.Syntax {
		ast.Walk(v, syntax)
	}
	return v.fileMap
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	typeSpec, ok := node.(*ast.TypeSpec)
	if !ok {
		return v
	}
	interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok || len(interfaceType.Methods.List) == 0 {
		return v
	}
	if !v.config.IsTarget(typeSpec.Name.Name) {
		return v
	}
	interfaceInfo := v.visit(typeSpec.Name.Name, interfaceType)
	fileName := filepath.Base(v.pkg.Fset.Position(node.Pos()).Filename)
	if _, ok := v.fileMap[fileName]; !ok {
		v.fileMap[fileName] = &types.File{
			ImportPathSet: interfaceInfo.ImportPathSet,
			Interfaces:    []*types.Interface{interfaceInfo},
		}
		return v
	}
	v.fileMap[fileName].ImportPathSet.Merge(interfaceInfo.ImportPathSet)
	v.fileMap[fileName].Interfaces = append(v.fileMap[fileName].Interfaces, interfaceInfo)
	return v
}

func (v *visitor) visit(interfaceName string, interfaceType *ast.InterfaceType) *types.Interface {
	importPathSet := strset.New()
	methods := make([]*types.Method, 0, len(interfaceType.Methods.List))
	for _, method := range interfaceType.Methods.List {
		switch typ := method.Type.(type) {
		case *ast.InterfaceType:
			// TODO: implement me (embed interface)
		case *ast.FuncType:
			args := make([]*types.Tuple, 0, len(typ.Params.List))
			for index, param := range typ.Params.List {
				if len(param.Names) == 0 {
					typ, importPaths := v.extractType(param.Type)
					importPathSet.Add(importPaths...)
					args = append(args, &types.Tuple{
						Name: fmt.Sprintf("arg%d", index),
						Type: typ,
					})
					continue
				}
				for _, name := range param.Names {
					typ, importPaths := v.extractType(param.Type)
					importPathSet.Add(importPaths...)
					args = append(args, &types.Tuple{
						Name: name.Name,
						Type: typ,
					})
				}
			}
			if typ.Results == nil {
				methods = append(methods, &types.Method{
					Name: method.Names[0].Name,
					Args: args,
				})
				continue
			}
			results := make([]*types.Tuple, 0, len(typ.Results.List))
			for index, result := range typ.Results.List {
				if len(result.Names) == 0 {
					typeName, importPaths := v.extractType(result.Type)
					importPathSet.Add(importPaths...)
					results = append(results, &types.Tuple{
						Name: fmt.Sprintf("result%d", index),
						Type: typeName,
					})
					continue
				}
				for _, name := range result.Names {
					typ, importPaths := v.extractType(result.Type)
					importPathSet.Add(importPaths...)
					results = append(results, &types.Tuple{
						Name: name.Name,
						Type: typ,
					})
				}
			}
			methods = append(methods, &types.Method{
				Name:    method.Names[0].Name,
				Args:    args,
				Results: results,
			})
		}
	}

	return &types.Interface{
		Name:          interfaceName,
		Methods:       methods,
		ImportPathSet: importPathSet,
	}
}

func (v *visitor) extractType(expr ast.Expr) (name string, importPaths []string) {
	switch t := expr.(type) {
	case *ast.Ident:
		obj := v.pkg.TypesInfo.ObjectOf(t)
		name := t.Name
		if obj != nil && obj.Pkg() != nil {
			if obj.Pkg().Path() == v.pkg.PkgPath {
				name = fmt.Sprintf("%s.%s", v.pkg.Name, name)
			}
			return name, []string{obj.Pkg().Path()}
		}
		return name, []string{}
	case *ast.SelectorExpr:
		ident := t.X.(*ast.Ident)
		obj := v.pkg.TypesInfo.ObjectOf(ident)
		name := fmt.Sprintf("%s.%s", ident.Name, t.Sel.Name)
		if obj != nil && obj.Pkg() != nil {
			return name, []string{obj.Pkg().Path()}
		}
		return name, []string{}
	case *ast.StarExpr:
		name, importPaths := v.extractType(t.X)
		return fmt.Sprintf("*%s", name), importPaths
	case *ast.ArrayType:
		name, importPaths := v.extractType(t.Elt)
		return fmt.Sprintf("[]%s", name), importPaths
	case *ast.MapType:
		keyName, keyImportPaths := v.extractType(t.Key)
		valueName, valueImportPaths := v.extractType(t.Value)
		return fmt.Sprintf("map[%s]%s", keyName, valueName), append(keyImportPaths, valueImportPaths...)
	case *ast.FuncType:
		importPathSet := strset.New()
		var builder strings.Builder
		builder.WriteString("func(")
		for i, param := range t.Params.List {
			if i > 0 {
				builder.WriteString(", ")
			}
			name, importPaths := v.extractType(param.Type)
			builder.WriteString(name)
			importPathSet.Add(importPaths...)
		}
		builder.WriteString(")")
		if t.Results != nil {
			builder.WriteString(" (")
			for i, result := range t.Results.List {
				if i > 0 {
					builder.WriteString(", ")
				}
				name, importPaths := v.extractType(result.Type)
				builder.WriteString(name)
				importPathSet.Add(importPaths...)
			}
			builder.WriteString(")")
		}
		return builder.String(), importPathSet.List()
	default:
		return "", []string{}
	}
}
