package core

type Section interface {
	IsArray() bool
	HasChildren() bool
	Navigate(path PathElement, create bool) Section

	Value(path PathElement) (value interface{}, ok bool)
	Set(path PathElement, value interface{}) (prev interface{})

	ValueDeep(path []PathElement) (value interface{}, ok bool)
	SetDeep(path []PathElement, createParents bool, value interface{}) (prev interface{})
}

type MapSection map[string]interface{}

func (_ *MapSection) IsArray() bool {
	return false
}

func (n *MapSection) HasChildren() bool {
	return len(*n) > 0
}

func (n *MapSection) Navigate(path PathElement, create bool) Section {
	if _, ok := path.(*ArrayPathElement); ok {
		return nil
	}

	self := *n
	key := string(*(path.(*ObjectPathElement)))

	v, ok := self[key]
	if !ok {
		goto end
	}

	if m, ok := v.(map[string]interface{}); ok {
		nm := MapSection(m)
		return &nm
	}

end:
	if create {
		newMap := make(map[string]interface{})
		self[key] = newMap

		nm := MapSection(newMap)
		return &nm
	}

	return nil
}

func (n *MapSection) Value(path PathElement) (value interface{}, ok bool) {
	if _, pok := path.(*ArrayPathElement); pok {
		return
	}

	self := *n
	key := string(*(path.(*ObjectPathElement)))

	value, ok = self[key]
	return
}

func (n *MapSection) Set(path PathElement, value interface{}) (prev interface{}) {
	if _, pok := path.(*ArrayPathElement); pok {
		return
	}

	self := *n
	key := string(*(path.(*ObjectPathElement)))

	prev = self[key]
	self[key] = value

	return
}

func (n *MapSection) ValueDeep(path []PathElement) (value interface{}, ok bool) {
	current := path[0]
	if len(path) > 1 {
		// we need to dig deeper
		nav := n.Navigate(current, false)
		if nav == nil {
			return
		}
		return nav.ValueDeep(path[1:])
	}

	return n.Value(current)
}

func (n *MapSection) SetDeep(path []PathElement, createParents bool, value interface{}) (prev interface{}) {
	current := path[0]
	if len(path) > 1 {
		// we need to dig deeper
		nav := n.Navigate(current, createParents)
		if nav == nil {
			return
		}

		return nav.SetDeep(path[1:], createParents, value)
	}

	return n.Set(current, value)
}

// TODO: implement
type ArraySection interface{}
