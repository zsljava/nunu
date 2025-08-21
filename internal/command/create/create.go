package create

import (
	"bytes"
	"fmt"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/go-nunu/nunu/tpl"
	"github.com/spf13/cobra"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Create struct {
	ProjectName          string
	CreateType           string
	FilePath             string
	FileName             string
	BasePkgName          string
	PkgName              string
	StructName           string
	StructNameLowerFirst string
	StructNameFirstChar  string
	StructNameSnakeCase  string
	IsFull               bool
}

func NewCreate() *Create {
	return &Create{}
}

var CmdCreate = &cobra.Command{
	Use:     "create [type] [handler-name]",
	Short:   "Create a new handler/service/repository/model",
	Example: "nunu create handler user",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var (
	tplPath string
)

func init() {
	CmdCreateHandler.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateService.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateRepository.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateModel.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")
	CmdCreateAll.Flags().StringVarP(&tplPath, "tpl-path", "t", tplPath, "template path")

}

var CmdCreateHandler = &cobra.Command{
	Use:     "handler",
	Short:   "Create a new handler",
	Example: "nunu create handler user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CmdCreateService = &cobra.Command{
	Use:     "service",
	Short:   "Create a new service",
	Example: "nunu create service user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CmdCreateRepository = &cobra.Command{
	Use:     "repository",
	Short:   "Create a new repository",
	Example: "nunu create repository user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CmdCreateModel = &cobra.Command{
	Use:     "model",
	Short:   "Create a new model",
	Example: "nunu create model user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CmdCreateAll = &cobra.Command{
	Use:     "all",
	Short:   "Create a new handler & service & repository & model",
	Example: "nunu create all user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	c := NewCreate()
	c.ProjectName = helper.GetProjectName(".")
	c.CreateType = cmd.Use
	c.FilePath, c.StructName = filepath.Split(args[0])
	split := strings.Split(c.StructName, ".")
	if len(split) == 1 {
		c.BasePkgName = split[0]
		c.FileName = strings.ReplaceAll(split[0], ".go", "")
	} else {
		c.BasePkgName = strings.Join(split[:len(split)-1], "/")
		c.FileName = strings.ReplaceAll(split[len(split)-1], ".go", "")
	}
	//c.FileName = strings.ReplaceAll(c.StructName, ".go", "")
	c.StructName = strutil.UpperFirst(strutil.CamelCase(c.FileName))
	c.StructNameLowerFirst = strutil.LowerFirst(c.StructName)
	c.StructNameFirstChar = string(c.StructNameLowerFirst[0])
	c.StructNameSnakeCase = strutil.SnakeCase(c.StructName)

	switch c.CreateType {
	case "handler", "service", "repository", "model":
		c.genFile()
	case "all":

		c.CreateType = "handler"
		c.PkgName = c.BasePkgName + "/" + c.CreateType
		c.genFile()

		c.CreateType = "service"
		c.PkgName = c.BasePkgName + "/" + c.CreateType
		c.genFile()

		c.CreateType = "repository"
		c.PkgName = c.BasePkgName + "/" + c.CreateType
		c.genFile()

		c.CreateType = "model"
		c.PkgName = c.BasePkgName + "/" + c.CreateType
		c.genFile()

		c.CreateType = "request"
		c.PkgName = c.BasePkgName + "/model/" + c.CreateType
		c.genFile()

		c.CreateType = "response"
		c.PkgName = c.BasePkgName + "/model/" + c.CreateType
		c.genFile()

		c.CreateType = "router"
		c.PkgName = c.BasePkgName + "/" + c.CreateType
		c.genFile()

		c.CreateType = "domain"
		c.PkgName = c.BasePkgName + "/" + c.CreateType
		c.genFile()

		c.CreateType = "listener"
		c.FilePath = "listener"
		c.PkgName = c.BasePkgName + "/listener"
		c.genFile()

		_ = c.rewrite()
	default:
		log.Fatalf("Invalid handler type: %s", c.PkgName)
	}
}

// rewrite 是重构后的顶层函数，负责协调整个重写流程
func (c *Create) rewrite() error {
	if err := c.rewriteRouterFile(); err != nil {
		return fmt.Errorf("failed to rewrite router file: %w", err)
	}
	if err := c.rewriteWireFile(); err != nil {
		return fmt.Errorf("failed to rewrite wire file: %w", err)
	}
	return nil
}

// rewriteWireFile 负责处理 wire.go 文件的重写逻辑
func (c *Create) rewriteWireFile() error {
	filePath := "cmd/server/wire/wire.go"

	// 定义需要添加的 providers
	providers := []struct {
		varName  string // wire.NewSet 所在的变量名
		pkgAlias string // import 的别名
		pkgPath  string // import 的路径
		param    string // 要添加到 NewSet 的参数
	}{
		{
			"repositorySet",
			c.StructNameLowerFirst + "Repository",
			fmt.Sprintf("%s/internal/%s/repository", c.ProjectName, c.BasePkgName),
			fmt.Sprintf("New%sRepository", c.StructName),
		},
		{
			"serviceSet",
			c.StructNameLowerFirst + "Service",
			fmt.Sprintf("%s/internal/%s/service", c.ProjectName, c.BasePkgName),
			fmt.Sprintf("New%sService", c.StructName),
		},
		{
			"handlerSet",
			c.StructNameLowerFirst + "Handler",
			fmt.Sprintf("%s/internal/%s/handler", c.ProjectName, c.BasePkgName),
			fmt.Sprintf("New%sHandler", c.StructName),
		},
		{
			"listenerSet",
			c.StructNameLowerFirst + "Listener",
			fmt.Sprintf("%s/listener/%s", c.ProjectName, c.BasePkgName),
			fmt.Sprintf("New%sListener", c.StructName),
		},
	}

	// 使用通用的文件处理函数
	return processSourceFile(filePath, func(fset *token.FileSet, file *ast.File) (bool, error) {
		var modified bool
		for _, p := range providers {
			// 添加 import
			imported, err := addNamedImport(file, p.pkgAlias, p.pkgPath)
			if err != nil {
				return false, err
			}

			// 添加 wire provider
			// 调用简化后的函数
			paramToAdd := fmt.Sprintf("%s.%s", p.pkgAlias, p.param)
			added, err := addWireSetProvider(fset, file, p.varName, paramToAdd) // <-- 不再需要 fset
			if err != nil {
				return false, err
			}

			if imported || added {
				modified = true
			}
		}
		return modified, nil
	})
}

// rewriteRouterFile 负责处理 router.go 文件的重写逻辑
func (c *Create) rewriteRouterFile() error {
	filePath := "common/base/router/router.go"
	structName := "Routers"
	handlerPkgPath := fmt.Sprintf(`%s/internal/%s/handler`, c.ProjectName, c.BasePkgName)

	return processSourceFile(filePath, func(fset *token.FileSet, file *ast.File) (bool, error) {
		// 1. 添加 import
		imported, err := addImport(file, handlerPkgPath)
		if err != nil {
			return false, err
		}

		// 2. 添加结构体字段
		newFieldName := c.StructName + "Handler"
		newFieldType := "*handler." + c.StructName + "Handler"
		added, err := addStructField(file, structName, newFieldName, newFieldType)
		if err != nil {
			return false, err
		}

		return imported || added, nil
	})
}

// addImport 检查并添加一个新的 import 路径
func addImport(file *ast.File, importPath string) (bool, error) {
	quotedPath := `"` + importPath + `"`
	for _, imp := range file.Imports {
		if imp.Path.Value == quotedPath {
			fmt.Printf("Import already exists: %s\n", quotedPath)
			return false, nil
		}
	}

	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: quotedPath,
		},
	}

	// 找到 import 声明块并添加
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			genDecl.Specs = append(genDecl.Specs, newImport)
			fmt.Printf("Added import: %s\n", quotedPath)
			return true, nil
		}
	}

	// 如果没有 import 块，则创建一个新的
	importDecl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{newImport},
	}
	file.Decls = append([]ast.Decl{importDecl}, file.Decls...)
	fmt.Printf("Created import block and added: %s\n", quotedPath)
	return true, nil
}

// addNamedImport 检查并添加一个带别名的 import
func addNamedImport(file *ast.File, name, importPath string) (bool, error) {
	quotedPath := `"` + importPath + `"`
	for _, imp := range file.Imports {
		if imp.Path.Value == quotedPath {
			// 如果路径已存在，检查别名是否匹配
			if (imp.Name == nil && name == "") || (imp.Name != nil && imp.Name.Name == name) {
				fmt.Printf("Named import already exists: %s %s\n", name, quotedPath)
				return false, nil
			}
		}
	}

	newImport := &ast.ImportSpec{
		Name: ast.NewIdent(name),
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: quotedPath,
		},
	}

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			genDecl.Specs = append(genDecl.Specs, newImport)
			fmt.Printf("Added named import: %s %s\n", name, quotedPath)
			return true, nil
		}
	}

	importDecl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{newImport},
	}
	file.Decls = append([]ast.Decl{importDecl}, file.Decls...)
	fmt.Printf("Created import block and added named import: %s %s\n", name, quotedPath)
	return true, nil
}

// addStructField 向指定的结构体中添加一个新的字段
func addStructField(file *ast.File, structName, fieldName, fieldType string) (bool, error) {
	var modified bool
	ast.Inspect(file, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok || ts.Name.Name != structName {
			return true
		}

		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return false
		}

		// 检查字段是否已存在
		for _, f := range st.Fields.List {
			for _, name := range f.Names {
				if name.Name == fieldName {
					fmt.Printf("Field already exists in struct %s: %s\n", structName, fieldName)
					modified = false
					return false
				}
			}
		}

		// 添加新字段
		// 注意：这里的 fieldType 直接作为 ast.Ident 是简化的处理方式
		// 一个更完整的实现需要解析 fieldType 字符串来构建复杂的类型表达式（如指针、选择器等）
		// 但对于 "*handler.UserHandler" 这种形式，直接创建 Ident 也能工作
		newField := &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(fieldName)},
			Type:  ast.NewIdent(fieldType),
		}
		st.Fields.List = append(st.Fields.List, newField)
		fmt.Printf("Added field to struct %s: %s %s\n", structName, fieldName, fieldType)
		modified = true
		return false
	})
	return modified, nil
}

// addWireSetProvider 向指定的 wire.NewSet 调用中添加一个新的 provider。
// 它不再需要 fset，也不再关心格式问题。
func addWireSetProvider(fset *token.FileSet, file *ast.File, varName, param string) (bool, error) {
	var modified bool
	parts := strings.Split(param, ".")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid param format for wire provider, expected 'package.Function': %s", param)
	}
	pkgName := parts[0]
	funcName := parts[1]

	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if name.Name != varName {
					continue
				}

				callExpr, ok := valueSpec.Values[i].(*ast.CallExpr)
				if !ok {
					return false
				}

				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok || selExpr.Sel.Name != "NewSet" {
					return false
				}

				// 检查参数是否已存在
				for _, arg := range callExpr.Args {
					if sel, ok := arg.(*ast.SelectorExpr); ok {
						if x, ok := sel.X.(*ast.Ident); ok {
							if x.Name == pkgName && sel.Sel.Name == funcName {
								fmt.Printf("Provider already exists in %s: %s\n", varName, param)
								modified = false
								return false
							}
						}
					}
				}

				// *** 核心手术逻辑 ***
				// 1. 获取括号的位置信息
				lparenPos := fset.Position(callExpr.Lparen)
				rparenPos := fset.Position(callExpr.Rparen)

				// 2. 如果它们在同一行，说明是单行格式，需要强制转换
				if lparenPos.Line == rparenPos.Line {
					// 3. 强制转换的技巧：将右括号的位置设置为紧跟在左括号之后。
					// 这会告诉 printer，括号之间没有任何内容可以放在同一行，
					// 迫使它在打印参数时进行换行。
					callExpr.Rparen = callExpr.Lparen + 1
					fmt.Printf("Detected single-line format in %s. Forcing multi-line.\n", varName)
				}

				// 直接添加新参数，无需任何格式化技巧
				newArg := &ast.SelectorExpr{
					X:   ast.NewIdent("\n\t" + pkgName),
					Sel: ast.NewIdent(funcName + ",\n"),
				}
				callExpr.Args = append(callExpr.Args, newArg)
				fmt.Printf("Added provider to %s: %s\n", varName, param)
				modified = true
				return false
			}
		}
		return true
	})
	return modified, nil
}

// processSourceFile 是一个通用的文件处理函数，封装了读、改、格式化、写的完整流程
func processSourceFile(filePath string, modifier func(*token.FileSet, *ast.File) (bool, error)) error {
	fset := token.NewFileSet()
	// 注意：这里的模式需要是 parser.ParseComments，以便 format 包能更好地处理注释
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// 调用我们自己的修改逻辑
	modified, err := modifier(fset, file)
	if err != nil {
		return fmt.Errorf("failed to modify AST for %s: %w", filePath, err)
	}

	// 如果我们的逻辑没有做出任何修改，就直接返回
	if !modified {
		fmt.Printf("No changes needed for %s, skipping write.\n", filePath)
		return nil
	}

	// *** 核心改动：使用 go/format ***
	// 1. 创建一个内存缓冲区
	var buf bytes.Buffer

	// 2. 将修改后的 AST 打印到缓冲区
	if err := printer.Fprint(&buf, fset, file); err != nil {
		return fmt.Errorf("failed to print ast to buffer for %s: %w", filePath, err)
	}

	// 3. 使用 go/format 包格式化缓冲区中的 Go 代码
	formattedSource, err := format.Source(buf.Bytes())
	if err != nil {
		// 如果格式化失败，通常说明我们的 AST 修改引入了语法错误
		return fmt.Errorf("failed to format source code for %s: %w", filePath, err)
	}

	// 4. 将格式化后的代码写回文件
	if err := os.WriteFile(filePath, formattedSource, 0644); err != nil {
		return fmt.Errorf("failed to write formatted code to %s: %w", filePath, err)
	}

	fmt.Printf("Successfully updated and formatted %s\n", filePath)
	return nil
}

func (c *Create) genFile() {
	filePath := c.FilePath
	if filePath == "" {
		filePath = fmt.Sprintf("internal/%s/", c.PkgName)
	}
	f := createFile(filePath, strings.ToLower(c.FileName)+".go")
	if f == nil {
		log.Printf("warn: file %s%s %s", filePath, strings.ToLower(c.FileName)+".go", "already exists.")
		return
	}
	defer f.Close()
	var t *template.Template
	var err error
	if tplPath == "" {
		t, err = template.ParseFS(tpl.CreateTemplateFS, fmt.Sprintf("create/%s.tpl", c.CreateType))
	} else {
		t, err = template.ParseFiles(path.Join(tplPath, fmt.Sprintf("%s.tpl", c.CreateType)))
	}
	if err != nil {
		log.Fatalf("create %s error: %s", c.CreateType, err.Error())
	}
	err = t.Execute(f, c)
	if err != nil {
		log.Fatalf("create %s error: %s", c.CreateType, err.Error())
	}
	log.Printf("Created new %s: %s", c.CreateType, filePath+strings.ToLower(c.FileName)+".go")

}
func createFile(dirPath string, filename string) *os.File {
	filePath := filepath.Join(dirPath, filename)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create dir %s: %v", dirPath, err)
	}
	stat, _ := os.Stat(filePath)
	if stat != nil {
		return nil
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}

	return file
}
