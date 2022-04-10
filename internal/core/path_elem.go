package core

import "fmt"

type PathElement interface {
	String() string
}

type ObjectPathElement string

func (e *ObjectPathElement) String() string {
	return string(*e)
}

func NewObjectPathElement(elem string) PathElement {
	v := ObjectPathElement(elem)
	return &v
}

type ArrayPathElement uint

func NewArrayPathElement(elem uint) PathElement {
	v := ArrayPathElement(elem)
	return &v
}

func (e *ArrayPathElement) String() string {
	return fmt.Sprintf("%d", *e)
}

func ObjectPathElements(elements ...string) (p []PathElement) {
	for _, elem := range elements {
		p = append(p, NewObjectPathElement(elem))
	}
	return
}

func ArrayPathElements(elements ...uint) (p []PathElement) {
	for _, elem := range elements {
		p = append(p, NewArrayPathElement(elem))
	}
	return
}

func PathElements(elements ...interface{}) (p []PathElement) {
	for _, elem := range elements {
		switch t := elem.(type) {
		case string:
			p = append(p, NewObjectPathElement(t))
		case uint:
			p = append(p, NewArrayPathElement(t))
		case int:
			p = append(p, NewArrayPathElement(uint(t)))
		default:
			panic(fmt.Errorf("unsupported path element '%#v'", t))
		}
	}
	return
}
