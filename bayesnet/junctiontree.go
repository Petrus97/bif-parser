package bayesnet

import "fmt"

/* (V)----|S|----(W)
 *
 * S.table == ∑ {W\S} W.table == ∑ {V\S} V.table  --> global consistency
 *
 *
 */

// Clique aka Vertex of Junction Tree
type Clique struct {
	Name       string
	Variables  []*Node // Nodes that form the cluster
	Neighbours []*Clique
	Separators []*Separator
	Table      *FactorV2
	Visited    bool
}

// Separator -> intersection of nodes between two Cliques
type Separator struct {
	Variables   []*Node // Node intersection between two cluster
	Table       *FactorV2
	LeftClique  *Clique
	RightClique *Clique
}

// JunctionTree structure
type JunctionTree struct {
	Root       *Clique
	Cliques    []*Clique
	Separators []*Separator
}

func (jt *JunctionTree) AddCliques(cliques ...*Clique) {
	for _, clique := range cliques {
		jt.Cliques = append(jt.Cliques, clique)
	}
}

func (jt *JunctionTree) SetRoot() {
	jt.Root = jt.Cliques[0]
}

func (jt *JunctionTree) AddSeparators(separators ...*Separator) {
	for _, sep := range separators {
		jt.Separators = append(jt.Separators, sep)
		sep.LeftClique.AddSeparator(sep)
		sep.RightClique.AddSeparator(sep)
	}
}

func NewSeparator(f *FactorV2, left *Clique, right *Clique) *Separator {
	sep := Separator{
		Table:       new(FactorV2),
		LeftClique:  left,
		RightClique: right,
	}
	vars := intersectV2(left.Table, right.Table)
	sep.Variables = vars
	factors := make([]*FactorV2, 0)
	for i, v := range vars {
		vfact := CreateFactorV2(v)
		if i == 0 {
			sep.Table = vfact
		} else {
			factors = append(factors, vfact)
		}
	}
	fmt.Println("TMP SEP", sep.Table)
	for _, fact := range factors {
		MultiplyFactor(sep.Table, fact, false)
	}
	return &sep
}

func (c *Clique) getSeparatorBetween(other *Clique) *Separator {
	// fmt.Println("c", c.Separators)
	// fmt.Println("othe", other.Separators)
	for _, sep := range c.Separators {
		if (sep.LeftClique == other && sep.RightClique == c) || (sep.RightClique == other && sep.LeftClique == c) {
			return sep
		}
	}
	return nil
}

func (jt *JunctionTree) collectEvidence(prev *Clique, actual *Clique, sep *Separator) {
	actual.Visited = true
	for _, neighbour := range actual.Neighbours {
		var sep *Separator
		sep = neighbour.getSeparatorBetween(actual)
		// fmt.Println(sep)
		if sep != nil {
			if neighbour.Visited == false {
				jt.collectEvidence(actual, neighbour, sep)
			}
		}
	}
	if actual != jt.Root {
		jt.passMessage(actual, prev, sep)
	}
}

func (jt *JunctionTree) distributeEvidence(c *Clique) {
	c.Visited = true
	for _, neighbour := range c.Neighbours {
		var sep *Separator
		sep = neighbour.getSeparatorBetween(c)
		// fmt.Println(sep)
		if sep != nil {
			if neighbour.Visited == false {
				jt.passMessage(c, neighbour, sep)
				jt.distributeEvidence(neighbour)
			}
		}
	}

}

func (jt *JunctionTree) passMessage(from *Clique, to *Clique, sep *Separator) {
	fmt.Println("passign message from", from.Name, "to", to.Name)
	m := difference(from.Variables, sep.Variables)
	for _, v := range m {
		fmt.Println(v.Name)
	}
	fmt.Println(from.Table)
	new := from.Table.Marginalize(true, m...)
	fmt.Println("NEW", new)
	MultiplyFactor(to.Table, new, false)
	fmt.Println(to.Table)
	DivideFactor(to.Table, sep.Table, false)
	fmt.Println(to.Table)
	sep.Table = new
	fmt.Println("sep", sep.Table)
	// for _, v := range from.Table.Scope {
	// 	fmt.Println(v.Name)
	// }
	// strides := make([]int, 0)
	// for _, v := range from.Table.Scope {
	// 	strides = append(strides, to.Table.Strides[v])
	// }
	// fmt.Println(to.Table)
	// fmt.Println(strides)
	// for _, s := range strides {
	// 	for i := range to.Table.CPT {
	// 		if i%s == 0 {
	// 			to.Table.CPT[i] *= from.Table.CPT[1]
	// 		} else {
	// 			to.Table.CPT[i] *= from.Table.CPT[0]
	// 		}
	// 	}
	// }
	// fmt.Println(to.Table)
	// fmt.Println(sep.Table)
	// Multiply to.Table with sep.Table using from.Table strides
	// Divide to.Table with sep.Table using sep.Table strides

}

// NewClique build a new clique from a factor
func NewClique(f *FactorV2, name string) *Clique {
	clique := Clique{
		Name:       name,
		Variables:  f.Scope,
		Neighbours: make([]*Clique, 0),
		Separators: make([]*Separator, 0),
		Table:      f,
	}
	return &clique
}

func (jt *JunctionTree) EnterEvidence(n *Node, values []float64) {
	for _, c := range jt.Cliques {
		if ok := containsNodeV2(c.Table, n); ok == true {
			EnterEvidenceFactor(c.Table, n, values)
		}
	}
}

func EnterEvidenceFactor(old *FactorV2, n *Node, values []float64) error {
	ev := n
	if len(values) != len(n.CPT) {
		return fmt.Errorf("Evidence list does not match len of CPT values")
	}
	ev.CPT = values
	evfact := CreateFactorV2(ev)
	fmt.Println(evfact)
	fmt.Println(old)
	MultiplyFactor(old, evfact, false)
	// stride := old.Strides[n]
	// fmt.Println(n.Name, stride)
	// for i := range old.CPT {
	// 	if i%stride == 0 {
	// 		old.CPT[i] *= ev.CPT[1]
	// 	} else {
	// 		old.CPT[i] *= ev.CPT[0]
	// 	}
	// }
	fmt.Println("AFTER EV", old)
	return nil
}

func (c *Clique) AddSeparator(sep *Separator) {
	fmt.Println("adding", sep)
	c.Separators = append(c.Separators, sep)
	if c == sep.LeftClique {
		c.Neighbours = append(c.Neighbours, sep.RightClique)
	} else {
		c.Neighbours = append(c.Neighbours, sep.RightClique)
	}
}

func (jt *JunctionTree) Propagate() {
	for _, c := range jt.Cliques {
		c.Visited = false
	}
	// Collect evidence
	jt.collectEvidence(nil, jt.Root, nil)
	for _, c := range jt.Cliques {
		c.Visited = false
	}
	// Distribute Evidence
	jt.distributeEvidence(jt.Root)
	// Normalize
	for _, c := range jt.Cliques {
		c.Table.Normalize()
	}
	for _, s := range jt.Separators {
		s.Table.Normalize()
	}
}

// func (c *Clique) AddVariables(nodes ...*Node) {
// 	for _, node := range nodes {
// 		c.Variables = append(c.Variables, node)
// 	}
// }

// InitIcyRoadJT initialize junctiontree relative to icy_roads.bif
// func InitIcyRoadJT(bn *BN) (jt *JunctionTree) {
// 	iw := Clique{
// 		Name:       "IW",
// 		Variables:  make([]*Node, 0),
// 		Neighbours: make([]*Clique, 0),
// 		Separators: make([]*Separator, 0),
// 		Table:      make([][]float64, 0),
// 	}
// 	iw.AddVariables(bn.GetNode("Icy"), bn.GetNode("Watson"))
// 	ih := Clique{
// 		Name:       "IH",
// 		Variables:  make([]*Node, 0),
// 		Neighbours: make([]*Clique, 0),
// 		Separators: make([]*Separator, 0),
// 		Table:      make([][]float64, 0),
// 	}
// 	ih.AddVariables(bn.GetNode("Icy"), bn.GetNode("Holmes"))
// 	var junc JunctionTree
// 	return &junc
// }

func (bn *BN) MatchDomains() {
	for _, n := range bn.Nodes {
		matches := map[interface{}]float64{}
		for i := 0; i < len(n.CPT); i++ {
			matches[n.Domain[i]] = n.CPT[i]
		}
	}
}
