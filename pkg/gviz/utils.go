package gviz

type At struct {
	m map[string]string
}

func NewAt() *At { return &At{m: map[string]string{}} }

func (a *At) Add(k, v string) *At    { a.m[k] = v; return a }
func (a *At) Map() map[string]string { return a.m }

func (a *At) Align(v string) *At   { return a.Add("align", v) }
func (a *At) BGColor(v string) *At { return a.Add("bgcolor", v) }
