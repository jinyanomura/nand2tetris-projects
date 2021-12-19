package symboltable

type Symbol struct {
	Type  string
	Kind  string
	Index int
}

type VarCount struct {
	Static, Field, Arg, Var int
}

type Table struct {
	Global map[string]Symbol
	Local  map[string]Symbol
	Count  VarCount
}

// Define adds given symbols to either Global or Local symbol table according to its kind.
func (t *Table) Define(name, varType, kind string) {
	entry := Symbol{
		Type:  varType,
		Kind:  kind,
		Index: 0,
	}

	switch kind {
	case "field":
		entry.Index = t.Count.Field
		t.Global[name] = entry
		t.Count.Field++
	case "static":
		entry.Index = t.Count.Static
		t.Global[name] = entry
		t.Count.Static++
	case "argument":
		entry.Index = t.Count.Arg
		t.Local[name] = entry
		t.Count.Arg++
	case "local":
		entry.Index = t.Count.Var
		t.Local[name] = entry
		t.Count.Var++
	}
}

// Reset resets local symbol tables and variable counts.
func (t *Table) Reset() {
	t.Local = make(map[string]Symbol)
	t.Count.Arg = 0
	t.Count.Var = 0
}

// KindOf returns the kind of given variable, or an empty string if not found.
func (t *Table) KindOf(name string) string {
	if symbol, ok := t.Local[name]; ok {
		return symbol.Kind
	} else if symbol, ok = t.Global[name]; ok {
		return symbol.Kind
	}
	return ""
}

// TypeOf returns the type of given variable, or an empty string if not found.
func (t *Table) TypeOf(name string) string {
	if symbol, ok := t.Local[name]; ok {
		return symbol.Type
	} else if symbol, ok = t.Global[name]; ok {
		return symbol.Type
	}
	return ""
}

// IndexOf returns the index of given variable, or -1 if not found.
func (t *Table) IndexOf(name string) int {
	if symbol, ok := t.Local[name]; ok {
		return symbol.Index
	} else if symbol, ok = t.Global[name]; ok {
		return symbol.Index
	}
	return -1
}

// New initializes a Table instance and returns its pointer address.
func New() *Table {
	return &Table{
		Global: make(map[string]Symbol),
		Local:  make(map[string]Symbol),
		Count: VarCount{
			Static: 0,
			Field:  0,
			Arg:    0,
			Var:    0,
		},
	}
}
