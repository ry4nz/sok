package types

type Namespaces struct {
	Namespaces []Namespace `json:"namespaces"`
	Active     string      `json:"active"`
}

type Namespace struct {
	Name string `json:"name"`
}
