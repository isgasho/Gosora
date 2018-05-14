package tmpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template/parse"
)

// TODO: Turn this file into a library
var textOverlapList = make(map[string]int)

type VarItem struct {
	Name        string
	Destination string
	Type        string
}

type VarItemReflect struct {
	Name        string
	Destination string
	Value       reflect.Value
}

type CTemplateConfig struct {
	Minify        bool
	Debug         bool
	SuperDebug    bool
	SkipHandles   bool
	SkipInitBlock bool
	PackageName   string
}

// nolint
type CTemplateSet struct {
	templateList          map[string]*parse.Tree
	fileDir               string
	funcMap               map[string]interface{}
	importMap             map[string]string
	TemplateFragmentCount map[string]int
	Fragments             map[string]int
	fragmentCursor        map[string]int
	FragOut               string
	varList               map[string]VarItem
	localVars             map[string]map[string]VarItemReflect
	hasDispInt            bool
	localDispStructIndex  int
	langIndexToName       []string
	stats                 map[string]int
	previousNode          parse.NodeType
	currentNode           parse.NodeType
	nextNode              parse.NodeType
	//tempVars map[string]string
	config        CTemplateConfig
	baseImportMap map[string]string
	expectsInt    interface{}
}

func NewCTemplateSet() *CTemplateSet {
	return &CTemplateSet{
		config: CTemplateConfig{
			PackageName: "main",
		},
		baseImportMap: map[string]string{},
		funcMap: map[string]interface{}{
			"and":      "&&",
			"not":      "!",
			"or":       "||",
			"eq":       true,
			"ge":       true,
			"gt":       true,
			"le":       true,
			"lt":       true,
			"ne":       true,
			"add":      true,
			"subtract": true,
			"multiply": true,
			"divide":   true,
			"dock":     true,
			"lang":     true,
			"scope":    true,
		},
	}
}

func (c *CTemplateSet) SetConfig(config CTemplateConfig) {
	if config.PackageName == "" {
		config.PackageName = "main"
	}
	c.config = config
}

func (c *CTemplateSet) GetConfig() CTemplateConfig {
	return c.config
}

func (c *CTemplateSet) SetBaseImportMap(importMap map[string]string) {
	c.baseImportMap = importMap
}

func (c *CTemplateSet) Compile(name string, fileDir string, expects string, expectsInt interface{}, varList map[string]VarItem, imports ...string) (out string, err error) {
	if c.config.Debug {
		fmt.Println("Compiling template '" + name + "'")
	}
	c.importMap = map[string]string{}
	for index, item := range c.baseImportMap {
		c.importMap[index] = item
	}
	if len(imports) > 0 {
		for _, importItem := range imports {
			c.importMap[importItem] = importItem
		}
	}

	c.fileDir = fileDir
	c.varList = varList
	c.hasDispInt = false
	c.localDispStructIndex = 0
	c.stats = make(map[string]int)
	c.expectsInt = expectsInt
	holdreflect := reflect.ValueOf(expectsInt)

	res, err := ioutil.ReadFile(fileDir + "overrides/" + name)
	if err != nil {
		c.detail("override path: ", fileDir+"overrides/"+name)
		c.detail("override err: ", err)
		res, err = ioutil.ReadFile(fileDir + name)
		if err != nil {
			return "", err
		}
	}
	content := string(res)
	if c.config.Minify {
		content = minify(content)
	}

	tree := parse.New(name, c.funcMap)
	var treeSet = make(map[string]*parse.Tree)
	tree, err = tree.Parse(content, "{{", "}}", treeSet, c.funcMap)
	if err != nil {
		return "", err
	}
	c.detail(name)

	fname := strings.TrimSuffix(name, filepath.Ext(name))
	c.templateList = map[string]*parse.Tree{fname: tree}
	varholder := "tmpl_" + fname + "_vars"

	c.detail(c.templateList)
	c.localVars = make(map[string]map[string]VarItemReflect)
	c.localVars[fname] = make(map[string]VarItemReflect)
	c.localVars[fname]["."] = VarItemReflect{".", varholder, holdreflect}
	if c.Fragments == nil {
		c.Fragments = make(map[string]int)
	}
	c.fragmentCursor = map[string]int{fname: 0}
	c.langIndexToName = nil

	// TODO: Is this the first template loaded in? We really should have some sort of constructor for CTemplateSet
	if c.TemplateFragmentCount == nil {
		c.TemplateFragmentCount = make(map[string]int)
	}

	out += c.rootIterate(c.templateList[fname], varholder, holdreflect, fname)
	c.TemplateFragmentCount[fname] = c.fragmentCursor[fname] + 1

	var importList string
	for _, item := range c.importMap {
		importList += "import \"" + item + "\"\n"
	}

	var varString string
	for _, varItem := range c.varList {
		varString += "var " + varItem.Name + " " + varItem.Type + " = " + varItem.Destination + "\n"
	}

	fout := "// +build !no_templategen\n\n// Code generated by Gosora. More below:\n/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */\n"
	fout += "package " + c.config.PackageName + "\n" + importList + "\n"

	if !c.config.SkipInitBlock {
		if len(c.langIndexToName) > 0 {
			fout += "var " + fname + "_tmpl_phrase_id int\n\n"
		}
		fout += "// nolint\nfunc init() {\n"

		if !c.config.SkipHandles {
			fout += "\tcommon.Template_" + fname + "_handle = Template_" + fname + "\n"

			fout += "\tcommon.Ctemplates = append(common.Ctemplates,\"" + fname + "\")\n\tcommon.TmplPtrMap[\"" + fname + "\"] = &common.Template_" + fname + "_handle\n"
		}

		fout += "\tcommon.TmplPtrMap[\"o_" + fname + "\"] = Template_" + fname + "\n"
		if len(c.langIndexToName) > 0 {
			fout += "\t" + fname + "_tmpl_phrase_id = common.RegisterTmplPhraseNames([]string{\n"
			for _, name := range c.langIndexToName {
				fout += "\t\t" + `"` + name + `"` + ",\n"
			}
			fout += "\t})\n"
		}
		fout += "}\n\n"
	}

	fout += "// nolint\nfunc Template_" + fname + "(tmpl_" + fname + "_vars " + expects + ", w io.Writer) error {\n"
	if len(c.langIndexToName) > 0 {
		fout += "var phrases = common.GetTmplPhrasesBytes(" + fname + "_tmpl_phrase_id)\n"
	}
	fout += varString + out + "return nil\n}\n"

	fout = strings.Replace(fout, `))
w.Write([]byte(`, " + ", -1)
	fout = strings.Replace(fout, "` + `", "", -1)
	//spstr := "`([:space:]*)`"
	//whitespaceWrites := regexp.MustCompile(`(?s)w.Write\(\[\]byte\(`+spstr+`\)\)`)
	//fout = whitespaceWrites.ReplaceAllString(fout,"")

	if c.config.Debug {
		for index, count := range c.stats {
			fmt.Println(index+": ", strconv.Itoa(count))
		}
		fmt.Println(" ")
	}
	c.detail("Output!")
	c.detail(fout)
	return fout, nil
}

func (c *CTemplateSet) rootIterate(tree *parse.Tree, varholder string, holdreflect reflect.Value, fname string) (out string) {
	c.detail(tree.Root)
	treeLength := len(tree.Root.Nodes)
	for index, node := range tree.Root.Nodes {
		c.detail("Node:", node.String())
		c.previousNode = c.currentNode
		c.currentNode = node.Type()
		if treeLength != (index + 1) {
			c.nextNode = tree.Root.Nodes[index+1].Type()
		}
		out += c.compileSwitch(varholder, holdreflect, fname, node)
	}
	return out
}

func (c *CTemplateSet) compileSwitch(varholder string, holdreflect reflect.Value, templateName string, node parse.Node) (out string) {
	c.detail("in compileSwitch")
	switch node := node.(type) {
	case *parse.ActionNode:
		c.detail("Action Node")
		if node.Pipe == nil {
			break
		}
		for _, cmd := range node.Pipe.Cmds {
			out += c.compileSubswitch(varholder, holdreflect, templateName, cmd)
		}
	case *parse.IfNode:
		c.detail("If Node:")
		c.detail("node.Pipe", node.Pipe)
		var expr string
		for _, cmd := range node.Pipe.Cmds {
			c.detail("If Node Bit:", cmd)
			c.detail("Bit Type:", reflect.ValueOf(cmd).Type().Name())
			expr += c.compileVarswitch(varholder, holdreflect, templateName, cmd)
			c.detail("Expression Step:", c.compileVarswitch(varholder, holdreflect, templateName, cmd))
		}

		c.detail("Expression:", expr)
		c.previousNode = c.currentNode
		c.currentNode = parse.NodeList
		c.nextNode = -1
		out = "if " + expr + " {\n" + c.compileSwitch(varholder, holdreflect, templateName, node.List) + "}"
		if node.ElseList == nil {
			c.detail("Selected Branch 1")
			return out + "\n"
		}
		c.detail("Selected Branch 2")
		return out + " else {\n" + c.compileSwitch(varholder, holdreflect, templateName, node.ElseList) + "}\n"
	case *parse.ListNode:
		c.detail("List Node")
		for _, subnode := range node.Nodes {
			out += c.compileSwitch(varholder, holdreflect, templateName, subnode)
		}
	case *parse.RangeNode:
		return c.compileRangeNode(varholder, holdreflect, templateName, node)
	case *parse.TemplateNode:
		return c.compileSubtemplate(varholder, holdreflect, node)
	case *parse.TextNode:
		c.previousNode = c.currentNode
		c.currentNode = node.Type()
		c.nextNode = 0
		tmpText := bytes.TrimSpace(node.Text)
		if len(tmpText) == 0 {
			return ""
		}

		fragmentName := templateName + "_" + strconv.Itoa(c.fragmentCursor[templateName])
		fragmentPrefix := templateName + "_frags[" + strconv.Itoa(c.fragmentCursor[templateName]) + "]"
		_, ok := c.Fragments[fragmentName]
		if !ok {
			c.Fragments[fragmentName] = len(node.Text)
			c.FragOut += fragmentPrefix + " = []byte(`" + string(node.Text) + "`)\n"
		}
		c.fragmentCursor[templateName] = c.fragmentCursor[templateName] + 1
		return "w.Write(" + fragmentPrefix + ")\n"
	default:
		return c.unknownNode(node)
	}
	return out
}

func (c *CTemplateSet) compileRangeNode(varholder string, holdreflect reflect.Value, templateName string, node *parse.RangeNode) (out string) {
	c.detail("Range Node!")
	c.detail(node.Pipe)
	var outVal reflect.Value
	for _, cmd := range node.Pipe.Cmds {
		c.detail("Range Bit:", cmd)
		out, outVal = c.compileReflectSwitch(varholder, holdreflect, templateName, cmd)
	}
	c.detail("Returned:", out)
	c.detail("Range Kind Switch!")

	switch outVal.Kind() {
	case reflect.Map:
		var item reflect.Value
		for _, key := range outVal.MapKeys() {
			item = outVal.MapIndex(key)
		}

		c.detail("Range item:", item)
		if !item.IsValid() {
			panic("item" + "^\n" + "Invalid map. Maybe, it doesn't have any entries for the template engine to analyse?")
		}

		out = "if len(" + out + ") != 0 {\nfor _, item := range " + out + " {\n" + c.compileSwitch("item", item, templateName, node.List) + "}\n}"
		if node.ElseList != nil {
			out += " else {\n" + c.compileSwitch("item", item, templateName, node.ElseList) + "}\n"
		}
	case reflect.Slice:
		if outVal.Len() == 0 {
			panic("The sample data needs at-least one or more elements for the slices. We're looking into removing this requirement at some point!")
		}
		item := outVal.Index(0)
		out = "if len(" + out + ") != 0 {\nfor _, item := range " + out + " {\n" + c.compileSwitch("item", item, templateName, node.List) + "}\n}"
		if node.ElseList != nil {
			out += " else {\n" + c.compileSwitch(varholder, holdreflect, templateName, node.ElseList) + "}"
		}
	case reflect.Invalid:
		return ""
	}

	return out + "\n"
}

func (c *CTemplateSet) compileSubswitch(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string) {
	c.detail("in compileSubswitch")
	firstWord := node.Args[0]
	switch n := firstWord.(type) {
	case *parse.FieldNode:
		c.detail("Field Node:", n.Ident)
		/* Use reflect to determine if the field is for a method, otherwise assume it's a variable. Variable declarations are coming soon! */
		cur := holdreflect

		var varbit string
		if cur.Kind() == reflect.Interface {
			cur = cur.Elem()
			varbit += ".(" + cur.Type().Name() + ")"
		}

		// ! Might not work so well for non-struct pointers
		skipPointers := func(cur reflect.Value, id string) reflect.Value {
			if cur.Kind() == reflect.Ptr {
				c.detail("Looping over pointer")
				for cur.Kind() == reflect.Ptr {
					cur = cur.Elem()
				}
				c.detail("Data Kind:", cur.Kind().String())
				c.detail("Field Bit:", id)
			}
			return cur
		}

		var assLines string
		var multiline = false
		for _, id := range n.Ident {
			c.detail("Data Kind:", cur.Kind().String())
			c.detail("Field Bit:", id)
			cur = skipPointers(cur, id)

			if !cur.IsValid() {
				c.error("Debug Data:")
				c.error("Holdreflect:", holdreflect)
				c.error("Holdreflect.Kind():", holdreflect.Kind())
				if !c.config.SuperDebug {
					c.error("cur.Kind():", cur.Kind().String())
				}
				c.error("")
				if !multiline {
					panic(varholder + varbit + "^\n" + "Invalid value. Maybe, it doesn't exist?")
				}
				panic(varbit + "^\n" + "Invalid value. Maybe, it doesn't exist?")
			}

			c.detail("in-loop varbit: " + varbit)
			if cur.Kind() == reflect.Map {
				cur = cur.MapIndex(reflect.ValueOf(id))
				varbit += "[\"" + id + "\"]"
				cur = skipPointers(cur, id)

				if cur.Kind() == reflect.Struct || cur.Kind() == reflect.Interface {
					// TODO: Move the newVarByte declaration to the top level or to the if level, if a dispInt is only used in a particular if statement
					var dispStr, newVarByte string
					if cur.Kind() == reflect.Interface {
						dispStr = "Int"
						if !c.hasDispInt {
							newVarByte = ":"
							c.hasDispInt = true
						}
					}
					// TODO: De-dupe identical struct types rather than allocating a variable for each one
					if cur.Kind() == reflect.Struct {
						dispStr = "Struct" + strconv.Itoa(c.localDispStructIndex)
						newVarByte = ":"
						c.localDispStructIndex++
					}
					varholder = "disp" + dispStr
					varbit = varholder + " " + newVarByte + "= " + varholder + varbit + "\n"
					multiline = true
				} else {
					continue
				}
			}
			if cur.Kind() != reflect.Interface {
				cur = cur.FieldByName(id)
				varbit += "." + id
			}

			// TODO: Handle deeply nested pointers mixed with interfaces mixed with pointers better
			if cur.Kind() == reflect.Interface {
				cur = cur.Elem()
				varbit += ".("
				// TODO: Surely, there's a better way of doing this?
				if cur.Type().PkgPath() != "main" && cur.Type().PkgPath() != "" {
					c.importMap["html/template"] = "html/template"
					varbit += strings.TrimPrefix(cur.Type().PkgPath(), "html/") + "."
				}
				varbit += cur.Type().Name() + ")"
			}
			c.detail("End Cycle: ", varbit)
		}

		if multiline {
			assSplit := strings.Split(varbit, "\n")
			varbit = assSplit[len(assSplit)-1]
			assSplit = assSplit[:len(assSplit)-1]
			assLines = strings.Join(assSplit, "\n") + "\n"
		}
		varbit = varholder + varbit
		out = c.compileVarsub(varbit, cur, assLines)

		for _, varItem := range c.varList {
			if strings.HasPrefix(out, varItem.Destination) {
				out = strings.Replace(out, varItem.Destination, varItem.Name, 1)
			}
		}
		return out
	case *parse.DotNode:
		c.detail("Dot Node:", node.String())
		return c.compileVarsub(varholder, holdreflect, "")
	case *parse.NilNode:
		panic("Nil is not a command x.x")
	case *parse.VariableNode:
		c.detail("Variable Node:", n.String())
		c.detail(n.Ident)
		varname, reflectVal := c.compileIfVarsub(n.String(), varholder, templateName, holdreflect)
		return c.compileVarsub(varname, reflectVal, "")
	case *parse.StringNode:
		return n.Quoted
	case *parse.IdentifierNode:
		c.detail("Identifier Node:", node)
		c.detail("Identifier Node Args:", node.Args)
		out, outval, lit := c.compileIdentSwitch(varholder, holdreflect, templateName, node)
		if lit {
			return out
		}
		return c.compileVarsub(out, outval, "")
	default:
		return c.unknownNode(node)
	}
}

func (c *CTemplateSet) compileVarswitch(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string) {
	c.detail("in compileVarswitch")
	firstWord := node.Args[0]
	switch n := firstWord.(type) {
	case *parse.FieldNode:
		if c.config.SuperDebug {
			fmt.Println("Field Node:", n.Ident)
			for _, id := range n.Ident {
				fmt.Println("Field Bit:", id)
			}
		}

		/* Use reflect to determine if the field is for a method, otherwise assume it's a variable. Coming Soon. */
		return c.compileBoolsub(n.String(), varholder, templateName, holdreflect)
	case *parse.ChainNode:
		c.detail("Chain Node:", n.Node)
		c.detail("Node Args:", node.Args)
	case *parse.IdentifierNode:
		c.detail("Identifier Node:", node)
		c.detail("Node Args:", node.Args)
		return c.compileIdentSwitchN(varholder, holdreflect, templateName, node)
	case *parse.DotNode:
		return varholder
	case *parse.VariableNode:
		c.detail("Variable Node:", n.String())
		c.detail("Node Identifier:", n.Ident)
		out, _ = c.compileIfVarsub(n.String(), varholder, templateName, holdreflect)
	case *parse.NilNode:
		panic("Nil is not a command x.x")
	case *parse.PipeNode:
		c.detail("Pipe Node!")
		c.detail(n)
		c.detail("Node Args:", node.Args)
		out += c.compileIdentSwitchN(varholder, holdreflect, templateName, node)
		c.detail("Out:", out)
	default:
		return c.unknownNode(firstWord)
	}
	return out
}

func (c *CTemplateSet) unknownNode(node parse.Node) (out string) {
	fmt.Println("Unknown Kind:", reflect.ValueOf(node).Elem().Kind())
	fmt.Println("Unknown Type:", reflect.ValueOf(node).Elem().Type().Name())
	panic("I don't know what node this is! Grr...")
}

func (c *CTemplateSet) compileIdentSwitchN(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string) {
	c.detail("in compileIdentSwitchN")
	out, _, _ = c.compileIdentSwitch(varholder, holdreflect, templateName, node)
	return out
}

func (c *CTemplateSet) dumpSymbol(pos int, node *parse.CommandNode, symbol string) {
	c.detail("symbol: ", symbol)
	c.detail("node.Args[pos + 1]", node.Args[pos+1])
	c.detail("node.Args[pos + 2]", node.Args[pos+2])
}

func (c *CTemplateSet) compareFunc(varholder string, holdreflect reflect.Value, templateName string, pos int, node *parse.CommandNode, compare string) (out string) {
	c.dumpSymbol(pos, node, compare)
	return c.compileIfVarsubN(node.Args[pos+1].String(), varholder, templateName, holdreflect) + " " + compare + " " + c.compileIfVarsubN(node.Args[pos+2].String(), varholder, templateName, holdreflect)
}

func (c *CTemplateSet) simpleMath(varholder string, holdreflect reflect.Value, templateName string, pos int, node *parse.CommandNode, symbol string) (out string, val reflect.Value) {
	leftParam, val2 := c.compileIfVarsub(node.Args[pos+1].String(), varholder, templateName, holdreflect)
	rightParam, val3 := c.compileIfVarsub(node.Args[pos+2].String(), varholder, templateName, holdreflect)

	if val2.IsValid() {
		val = val2
	} else if val3.IsValid() {
		val = val3
	} else {
		// TODO: What does this do?
		numSample := 1
		val = reflect.ValueOf(numSample)
	}
	c.dumpSymbol(pos, node, symbol)
	return leftParam + " " + symbol + " " + rightParam, val
}

func (c *CTemplateSet) compareJoin(varholder string, holdreflect reflect.Value, templateName string, pos int, node *parse.CommandNode, symbol string) (pos2 int, out string) {
	c.detailf("Building %s function", symbol)
	if pos == 0 {
		fmt.Println("pos:", pos)
		panic(symbol + " is missing a left operand")
	}
	if len(node.Args) <= pos {
		fmt.Println("post pos:", pos)
		fmt.Println("len(node.Args):", len(node.Args))
		panic(symbol + " is missing a right operand")
	}

	left := c.compileBoolsub(node.Args[pos-1].String(), varholder, templateName, holdreflect)
	_, funcExists := c.funcMap[node.Args[pos+1].String()]

	var right string
	if !funcExists {
		right = c.compileBoolsub(node.Args[pos+1].String(), varholder, templateName, holdreflect)
	}
	out = left + " " + symbol + " " + right

	c.detail("Left operand:", node.Args[pos-1])
	c.detail("Right operand:", node.Args[pos+1])
	if !funcExists {
		pos++
	}
	c.detail("pos:", pos)
	c.detail("len(node.Args):", len(node.Args))

	return pos, out
}

func (c *CTemplateSet) compileIdentSwitch(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string, val reflect.Value, literal bool) {
	c.detail("in compileIdentSwitch")
ArgLoop:
	for pos := 0; pos < len(node.Args); pos++ {
		id := node.Args[pos]
		c.detail("pos:", pos)
		c.detail("ID:", id)
		switch id.String() {
		case "not":
			out += "!"
		case "or", "and":
			var rout string
			pos, rout = c.compareJoin(varholder, holdreflect, templateName, pos, node, c.funcMap[id.String()].(string)) // TODO: Test this
			out += rout
		case "le": // TODO: Can we condense these comparison cases down into one?
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, "<=")
			break ArgLoop
		case "lt":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, "<")
			break ArgLoop
		case "gt":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, ">")
			break ArgLoop
		case "ge":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, ">=")
			break ArgLoop
		case "eq":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, "==")
			break ArgLoop
		case "ne":
			out += c.compareFunc(varholder, holdreflect, templateName, pos, node, "!=")
			break ArgLoop
		case "add":
			rout, rval := c.simpleMath(varholder, holdreflect, templateName, pos, node, "+")
			out += rout
			val = rval
			break ArgLoop
		case "subtract":
			rout, rval := c.simpleMath(varholder, holdreflect, templateName, pos, node, "-")
			out += rout
			val = rval
			break ArgLoop
		case "divide":
			rout, rval := c.simpleMath(varholder, holdreflect, templateName, pos, node, "/")
			out += rout
			val = rval
			break ArgLoop
		case "multiply":
			rout, rval := c.simpleMath(varholder, holdreflect, templateName, pos, node, "*")
			out += rout
			val = rval
			break ArgLoop
		case "dock":
			var leftParam, rightParam string
			// TODO: Implement string literals properly
			leftOperand := node.Args[pos+1].String()
			rightOperand := node.Args[pos+2].String()

			if len(leftOperand) == 0 || len(rightOperand) == 0 {
				panic("The left or right operand for function dock cannot be left blank")
			}
			if leftOperand[0] == '"' {
				leftParam = leftOperand
			} else {
				leftParam, _ = c.compileIfVarsub(leftOperand, varholder, templateName, holdreflect)
			}
			if rightOperand[0] == '"' {
				panic("The right operand for function dock cannot be a string")
			}

			rightParam, val3 := c.compileIfVarsub(rightOperand, varholder, templateName, holdreflect)
			if !val3.IsValid() {
				panic("val3 is invalid")
			}
			val = val3

			// TODO: Refactor this
			out = "w.Write([]byte(common.BuildWidget(" + leftParam + "," + rightParam + ")))\n"
			literal = true
			break ArgLoop
		case "lang":
			var leftParam string
			// TODO: Implement string literals properly
			leftOperand := node.Args[pos+1].String()
			if len(leftOperand) == 0 {
				panic("The left operand for the language string cannot be left blank")
			}

			if leftOperand[0] == '"' {
				// ! Slightly crude but it does the job
				leftParam = strings.Replace(leftOperand, "\"", "", -1)
			} else {
				panic("Phrase names cannot be dynamic")
			}

			c.langIndexToName = append(c.langIndexToName, leftParam)
			out = "w.Write(phrases[" + strconv.Itoa(len(c.langIndexToName)-1) + "])\n"
			literal = true
			break ArgLoop
		case "scope":
			literal = true
			break ArgLoop
		default:
			c.detail("Variable!")
			if len(node.Args) > (pos + 1) {
				nextNode := node.Args[pos+1].String()
				if nextNode == "or" || nextNode == "and" {
					continue
				}
			}
			out += c.compileIfVarsubN(id.String(), varholder, templateName, holdreflect)
		}
	}
	return out, val, literal
}

func (c *CTemplateSet) compileReflectSwitch(varholder string, holdreflect reflect.Value, templateName string, node *parse.CommandNode) (out string, outVal reflect.Value) {
	c.detail("in compileReflectSwitch")
	firstWord := node.Args[0]
	switch n := firstWord.(type) {
	case *parse.FieldNode:
		if c.config.SuperDebug {
			fmt.Println("Field Node:", n.Ident)
			for _, id := range n.Ident {
				fmt.Println("Field Bit:", id)
			}
		}
		/* Use reflect to determine if the field is for a method, otherwise assume it's a variable. Coming Soon. */
		return c.compileIfVarsub(n.String(), varholder, templateName, holdreflect)
	case *parse.ChainNode:
		c.detail("Chain Node:", n.Node)
		c.detail("node.Args:", node.Args)
	case *parse.DotNode:
		return varholder, holdreflect
	case *parse.NilNode:
		panic("Nil is not a command x.x")
	default:
		//panic("I don't know what node this is")
	}
	return "", outVal
}

func (c *CTemplateSet) compileIfVarsubN(varname string, varholder string, templateName string, cur reflect.Value) (out string) {
	c.detail("in compileIfVarsubN")
	out, _ = c.compileIfVarsub(varname, varholder, templateName, cur)
	return out
}

func (c *CTemplateSet) compileIfVarsub(varname string, varholder string, templateName string, cur reflect.Value) (out string, val reflect.Value) {
	c.detail("in compileIfVarsub")
	if varname[0] != '.' && varname[0] != '$' {
		return varname, cur
	}

	bits := strings.Split(varname, ".")
	if varname[0] == '$' {
		var res VarItemReflect
		if varname[1] == '.' {
			res = c.localVars[templateName]["."]
		} else {
			res = c.localVars[templateName][strings.TrimPrefix(bits[0], "$")]
		}
		out += res.Destination
		cur = res.Value

		if cur.Kind() == reflect.Interface {
			cur = cur.Elem()
		}
	} else {
		out += varholder
		if cur.Kind() == reflect.Interface {
			cur = cur.Elem()
			out += ".(" + cur.Type().Name() + ")"
		}
	}
	bits[0] = strings.TrimPrefix(bits[0], "$")

	c.detail("Cur Kind:", cur.Kind())
	c.detail("Cur Type:", cur.Type().Name())
	for _, bit := range bits {
		c.detail("Variable Field:", bit)
		if bit == "" {
			continue
		}

		// TODO: Fix this up so that it works for regular pointers and not just struct pointers. Ditto for the other cur.Kind() == reflect.Ptr we have in this file
		if cur.Kind() == reflect.Ptr {
			c.detail("Looping over pointer")
			for cur.Kind() == reflect.Ptr {
				cur = cur.Elem()
			}
			c.detail("Data Kind:", cur.Kind().String())
			c.detail("Field Bit:", bit)
		}

		cur = cur.FieldByName(bit)
		out += "." + bit
		if cur.Kind() == reflect.Interface {
			cur = cur.Elem()
			out += ".(" + cur.Type().Name() + ")"
		}
		if !cur.IsValid() {
			panic(out + "^\n" + "Invalid value. Maybe, it doesn't exist?")
		}
		c.detail("Data Kind:", cur.Kind())
		c.detail("Data Type:", cur.Type().Name())
	}

	c.detail("Out Value:", out)
	c.detail("Out Kind:", cur.Kind())
	c.detail("Out Type:", cur.Type().Name())

	for _, varItem := range c.varList {
		if strings.HasPrefix(out, varItem.Destination) {
			out = strings.Replace(out, varItem.Destination, varItem.Name, 1)
		}
	}

	c.detail("Out Value:", out)
	c.detail("Out Kind:", cur.Kind())
	c.detail("Out Type:", cur.Type().Name())

	_, ok := c.stats[out]
	if ok {
		c.stats[out]++
	} else {
		c.stats[out] = 1
	}

	return out, cur
}

func (c *CTemplateSet) compileBoolsub(varname string, varholder string, templateName string, val reflect.Value) string {
	c.detail("in compileBoolsub")
	out, val := c.compileIfVarsub(varname, varholder, templateName, val)
	// TODO: What if it's a pointer or an interface? I *think* we've got pointers handled somewhere, but not interfaces which we don't know the types of at compile time
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		out += " > 0"
	case reflect.Bool: // Do nothing
	case reflect.String:
		out += " != \"\""
	case reflect.Slice, reflect.Map:
		out = "len(" + out + ") != 0"
	default:
		fmt.Println("Variable Name:", varname)
		fmt.Println("Variable Holder:", varholder)
		fmt.Println("Variable Kind:", val.Kind())
		panic("I don't know what this variable's type is o.o\n")
	}
	return out
}

func (c *CTemplateSet) compileVarsub(varname string, val reflect.Value, assLines string) (out string) {
	c.detail("in compileVarsub")

	// Is this a literal string?
	if len(varname) != 0 && varname[0] == '"' {
		return assLines + "w.Write([]byte(" + varname + "))\n"
	}

	for _, varItem := range c.varList {
		if strings.HasPrefix(varname, varItem.Destination) {
			varname = strings.Replace(varname, varItem.Destination, varItem.Name, 1)
		}
	}

	_, ok := c.stats[varname]
	if ok {
		c.stats[varname]++
	} else {
		c.stats[varname] = 1
	}

	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	c.detail("varname: ", varname)
	c.detail("assLines: ", assLines)
	switch val.Kind() {
	case reflect.Int:
		c.importMap["strconv"] = "strconv"
		out = "w.Write([]byte(strconv.Itoa(" + varname + ")))\n"
	case reflect.Bool:
		out = "if " + varname + " {\nw.Write([]byte(\"true\"))} else {\nw.Write([]byte(\"false\"))\n}\n"
	case reflect.String:
		if val.Type().Name() != "string" && !strings.HasPrefix(varname, "string(") {
			varname = "string(" + varname + ")"
		}
		out = "w.Write([]byte(" + varname + "))\n"
	case reflect.Int64:
		c.importMap["strconv"] = "strconv"
		out = "w.Write([]byte(strconv.FormatInt(" + varname + ", 10)))"
	default:
		if !val.IsValid() {
			panic(assLines + varname + "^\n" + "Invalid value. Maybe, it doesn't exist?")
		}
		fmt.Println("Unknown Variable Name:", varname)
		fmt.Println("Unknown Kind:", val.Kind())
		fmt.Println("Unknown Type:", val.Type().Name())
		panic("-- I don't know what this variable's type is o.o\n")
	}
	c.detail("out: ", out)
	return assLines + out
}

func (c *CTemplateSet) compileSubtemplate(pvarholder string, pholdreflect reflect.Value, node *parse.TemplateNode) (out string) {
	c.detail("in compileSubtemplate")
	c.detail("Template Node: ", node.Name)

	fname := strings.TrimSuffix(node.Name, filepath.Ext(node.Name))
	varholder := "tmpl_" + fname + "_vars"
	var holdreflect reflect.Value
	if node.Pipe != nil {
		for _, cmd := range node.Pipe.Cmds {
			firstWord := cmd.Args[0]
			switch firstWord.(type) {
			case *parse.DotNode:
				varholder = pvarholder
				holdreflect = pholdreflect
			case *parse.NilNode:
				panic("Nil is not a command x.x")
			default:
				c.detail("Unknown Node: ", firstWord)
				panic("")
			}
		}
	}

	// TODO: Cascade errors back up the tree to the caller?
	res, err := ioutil.ReadFile(c.fileDir + "overrides/" + node.Name)
	if err != nil {
		c.detail("override path: ", c.fileDir+"overrides/"+node.Name)
		c.detail("override err: ", err)
		res, err = ioutil.ReadFile(c.fileDir + node.Name)
		if err != nil {
			log.Fatal(err)
		}
	}
	content := string(res)
	if c.config.Minify {
		content = minify(content)
	}

	tree := parse.New(node.Name, c.funcMap)
	var treeSet = make(map[string]*parse.Tree)
	tree, err = tree.Parse(content, "{{", "}}", treeSet, c.funcMap)
	if err != nil {
		log.Fatal(err)
	}

	c.templateList[fname] = tree
	subtree := c.templateList[fname]
	c.detail("subtree.Root", subtree.Root)

	c.localVars[fname] = make(map[string]VarItemReflect)
	c.localVars[fname]["."] = VarItemReflect{".", varholder, holdreflect}
	c.fragmentCursor[fname] = 0

	out += c.rootIterate(subtree, varholder, holdreflect, fname)
	c.TemplateFragmentCount[fname] = c.fragmentCursor[fname] + 1
	return out
}

// TODO: Should we rethink the way the log methods work or their names?

func (c *CTemplateSet) detail(args ...interface{}) {
	if c.config.SuperDebug {
		fmt.Println(args...)
	}
}

func (c *CTemplateSet) detailf(left string, args ...interface{}) {
	if c.config.SuperDebug {
		fmt.Printf(left, args...)
	}
}

func (c *CTemplateSet) error(args ...interface{}) {
	if c.config.Debug {
		fmt.Println(args...)
	}
}
