package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
	FreeScope   SymbolScope = "FREE"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index byte
}

type SymbolTable struct {
	Parent         *SymbolTable
	Store          map[string]Symbol
	FreeSymbols    []Symbol
	NumDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Store:       make(map[string]Symbol),
		FreeSymbols: []Symbol{},
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

	if s.Parent == nil {
		return Symbol{}, false
	}

	// Look up in outer scope
	symbol, ok := s.Parent.Resolve(name)
	if !ok {
		return Symbol{}, false
	}

	// Globals don't need to be captured as free variables
	if symbol.Scope == GlobalScope {
		return symbol, true
	}

	// It's an outer local or free variable—capture it in this table
	return s.defineFree(symbol), true
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	// If we've already captured this symbol, return the existing free symbol reference
	for _, sym := range s.FreeSymbols {
		if sym.Name == original.Name {
			return sym
		}
	}

	freeSymbol := Symbol{
		Name:  original.Name,
		Scope: FreeScope,
		Index: byte(len(s.FreeSymbols)),
	}
	s.FreeSymbols = append(s.FreeSymbols, original)
	return freeSymbol
}
