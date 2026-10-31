package main

import (
	"container/list"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type pstackItem struct {
	grammarSym string
	stateSym   string
}

func (item pstackItem) String() string {
	return item.grammarSym + item.stateSym
}

type parseStack struct {
	stack *list.List
}

func newParseStack() parseStack {
	return parseStack{stack: list.New()}
}

func (stack parseStack) push(item pstackItem) {
	stack.stack.PushFront(item)
}

func (stack parseStack) top() pstackItem {
	if element := stack.stack.Front(); element != nil {
		return element.Value.(pstackItem)
	}
	return pstackItem{}
}

func (stack parseStack) pop() pstackItem {
	if element := stack.stack.Front(); element != nil {
		stack.stack.Remove(element)
		return element.Value.(pstackItem)
	}
	return pstackItem{}
}

func (stack parseStack) popNum(count int) {
	for range count {
		stack.pop()
	}
}

func (stack parseStack) String() string {
	var result strings.Builder
	for element := stack.stack.Back(); element != nil; element = element.Prev() {
		result.WriteString(element.Value.(pstackItem).String())
	}
	return result.String()
}

type treeNode struct {
	symbol   string
	terminal bool
	children []*treeNode
}

func terminalNode(symbol string) *treeNode {
	return &treeNode{symbol: symbol, terminal: true}
}

func nonterminalNode(symbol string, children ...*treeNode) *treeNode {
	return &treeNode{symbol: symbol, children: children}
}

func (node *treeNode) compact() string {
	if node.terminal {
		return node.symbol
	}
	parts := make([]string, len(node.children))
	for i, child := range node.children {
		parts[i] = child.compact()
	}
	return "[" + node.symbol + " " + strings.Join(parts, " ") + "]"
}

func (node *treeNode) printIndented(indent int) {
	fmt.Printf("%s%s\n", strings.Repeat("    ", indent), node.symbol)
	for _, child := range node.children {
		child.printIndented(indent + 1)
	}
}

type treeStack struct {
	items []*treeNode
}

func (stack *treeStack) push(node *treeNode) {
	stack.items = append(stack.items, node)
}

func (stack *treeStack) pop() *treeNode {
	if len(stack.items) == 0 {
		return nil
	}
	last := len(stack.items) - 1
	node := stack.items[last]
	stack.items = stack.items[:last]
	return node
}

func (stack treeStack) String() string {
	parts := make([]string, len(stack.items))
	for i, node := range stack.items {
		parts[i] = node.compact()
	}
	return strings.Join(parts, " ")
}

type grammarRule struct {
	lhs string
	rhs []string
}

var grammar = []grammarRule{
	{lhs: "E", rhs: []string{"E", "+", "T"}},
	{lhs: "E", rhs: []string{"T"}},
	{lhs: "T", rhs: []string{"T", "*", "F"}},
	{lhs: "T", rhs: []string{"F"}},
	{lhs: "F", rhs: []string{"(", "E", ")"}},
	{lhs: "F", rhs: []string{"id"}},
}

var actionTable = [12][6]string{
	{"S5", "", "", "S4", "", ""},
	{"", "S6", "", "", "", "accept"},
	{"", "R2", "S7", "", "R2", "R2"},
	{"", "R4", "R4", "", "R4", "R4"},
	{"S5", "", "", "S4", "", ""},
	{"", "R6", "R6", "", "R6", "R6"},
	{"S5", "", "", "S4", "", ""},
	{"S5", "", "", "S4", "", ""},
	{"", "S6", "", "", "S11", ""},
	{"", "R1", "S7", "", "R1", "R1"},
	{"", "R3", "R3", "", "R3", "R3"},
	{"", "R5", "R5", "", "R5", "R5"},
}

var gotoTable = [12][3]string{
	{"1", "2", "3"},
	{"", "", ""},
	{"", "", ""},
	{"", "", ""},
	{"8", "2", "3"},
	{"", "", ""},
	{"", "9", "3"},
	{"", "", "10"},
	{"", "", ""},
	{"", "", ""},
	{"", "", ""},
	{"", "", ""},
}

var terminals = map[string]int{"id": 0, "+": 1, "*": 2, "(": 3, ")": 4, "$": 5}
var nonterminals = map[string]int{"E": 0, "T": 1, "F": 2}

type parser struct {
	input         []string
	position      int
	stack         parseStack
	trees         treeStack
	accepted      bool
	ungrammatical bool
	step          int
}

func newParser(input []string) *parser {
	tokens := append(append([]string{}, input...), "$")
	stack := newParseStack()
	stack.push(pstackItem{grammarSym: "", stateSym: "0"})
	return &parser{input: tokens, stack: stack}
}

func (p *parser) currentToken() string {
	return p.input[p.position]
}

func (p *parser) inputString() string {
	return strings.Join(p.input[p.position:], "")
}

func (p *parser) parse() {
	fmt.Println("step | stack | input | lookup | action | details | parse tree stack")
	fmt.Println("-----+-------+-------+--------+--------+---------+-----------------")
	for !p.accepted && !p.ungrammatical {
		p.parse1Step()
	}
	if p.ungrammatical {
		fmt.Println("\nResult: input is not grammatical.")
		return
	}
	fmt.Println("\nResult: input is grammatical.")
	fmt.Println("\nParse tree:")
	if len(p.trees.items) == 1 {
		p.trees.items[0].printIndented(0)
	}
}

func (p *parser) parse1Step() {
	p.step++
	beforeStack := p.stack.String()
	token := p.currentToken()
	state, err := strconv.Atoi(p.stack.top().stateSym)
	if err != nil {
		p.ungrammatical = true
		return
	}
	tokenColumn, ok := terminals[token]
	if !ok {
		p.ungrammatical = true
		return
	}
	lookup := fmt.Sprintf("[%d,%s]", state, token)
	action := actionTable[state][tokenColumn]
	if action == "" {
		p.ungrammatical = true
		fmt.Printf("%4d | %-5s | %-5s | %-6s | error  | invalid action | %s\n", p.step, beforeStack, p.inputString(), lookup, p.trees)
		return
	}

	details := ""
	switch {
	case action == "accept":
		p.accepted = true
	case strings.HasPrefix(action, "S"):
		nextState := strings.TrimPrefix(action, "S")
		p.stack.push(pstackItem{grammarSym: token, stateSym: nextState})
		p.trees.push(terminalNode(token))
		p.position++
		details = "push " + token + nextState
	case strings.HasPrefix(action, "R"):
		ruleNumber, _ := strconv.Atoi(strings.TrimPrefix(action, "R"))
		rule := grammar[ruleNumber-1]
		p.stack.popNum(len(rule.rhs))
		previousState, _ := strconv.Atoi(p.stack.top().stateSym)
		nextState := gotoTable[previousState][nonterminals[rule.lhs]]
		if nextState == "" {
			p.ungrammatical = true
			return
		}
		p.stack.push(pstackItem{grammarSym: rule.lhs, stateSym: nextState})
		p.reduceTree(rule)
		details = fmt.Sprintf("R%d: %s -> %s; goto %s", ruleNumber, rule.lhs, strings.Join(rule.rhs, " "), nextState)
	}
	fmt.Printf("%4d | %-5s | %-5s | %-6s | %-6s | %-40s | %s\n", p.step, beforeStack, p.inputString(), lookup, action, details, p.trees)
}

func (p *parser) reduceTree(rule grammarRule) {
	nodes := make([]*treeNode, len(rule.rhs))
	for i := len(nodes) - 1; i >= 0; i-- {
		nodes[i] = p.trees.pop()
	}
	if len(nodes) == 1 {
		p.trees.push(nonterminalNode(rule.lhs, nodes[0]))
		return
	}
	children := make([]*treeNode, len(nodes))
	for i, symbol := range rule.rhs {
		if symbol == "E" || symbol == "T" || symbol == "F" {
			children[i] = nodes[i]
		} else {
			children[i] = terminalNode(symbol)
		}
	}
	p.trees.push(nonterminalNode(rule.lhs, children...))
}

func main() {
	inputFlag := flag.String("input", "id+id*id", "tokens to parse (id, +, *, (, and ))")
	flag.Parse()
	tokens := tokenize(*inputFlag)
	if len(tokens) == 0 {
		fmt.Println("error: input must contain at least one token")
		return
	}
	newParser(tokens).parse()
}

func tokenize(input string) []string {
	input = strings.ReplaceAll(input, " ", "")
	result := make([]string, 0)
	for len(input) > 0 {
		if strings.HasPrefix(input, "id") {
			result = append(result, "id")
			input = input[2:]
			continue
		}
		if strings.ContainsRune("+-*()", rune(input[0])) {
			result = append(result, string(input[0]))
			input = input[1:]
			continue
		}
		return nil
	}
	return result
}
