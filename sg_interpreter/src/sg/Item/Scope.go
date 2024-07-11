package Item

type Scope struct {
	Mp    map[string]Item
	outer *Scope
}

func NewScope() *Scope {
	m := make(map[string]Item)
	return &Scope{Mp: m, outer: nil}
}

func NewEnclosedScope(outer *Scope) *Scope {
	s := NewScope()
	s.outer = outer
	return s
}

func (scope *Scope) Get(key string) (Item, bool) {
	item, ok := scope.Mp[key]
	if !ok && scope.outer != nil {
		item, ok = scope.outer.Get(key)
	}
	return item, ok
}

func (scope *Scope) Set(key string, item Item) Item {
	scope.Mp[key] = item
	return item
}
