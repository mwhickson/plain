package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index byte
}

type SymbolTable struct {
	Parent         *SymbolTable
	Store          map[string]Symbol
	NumDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Store: make(map[string]Symbol),
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Parent = outer
	return s
}

func (s *SymbolTable) Define(name string, index byte) Symbol {
	scope := GlobalScope
	if s.Parent != nil {
		scope = LocalScope
	}

	symbol := Symbol{Name: name, Scope: scope, Index: index}
	s.Store[name] = symbol
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	if symbol, ok := s.Store[name]; ok {
		return symbol, true
	}
	if s.Parent != nil {
		return s.Parent.Resolve(name)
	}
	return Symbol{}, false
}
