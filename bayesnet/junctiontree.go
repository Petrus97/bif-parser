package bayesnet

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
	Table      [][]float64
}

// Separator -> intersection of nodes between two Cliques
type Separator struct {
	Variables   []*Node // Node intersection between two cluster
	Table       [][]float64
	LeftClique  *Clique
	RightClique *Clique
}

// JunctionTree structure
type JunctionTree struct {
	Cliques    []*Clique
	Separators []*Separator
}

func (jt *JunctionTree) addCliques(cliques ...*Clique) {
	for _, clique := range cliques {
		jt.Cliques = append(jt.Cliques, clique)
	}
}

func (jt *JunctionTree) addSeparators(separators ...*Separator) {
	for _, sep := range separators {
		jt.Separators = append(jt.Separators, sep)
	}
}

func (c *Clique) AddVariables(nodes ...*Node) {
	for _, node := range nodes {
		c.Variables = append(c.Variables, node)
	}
}

// InitIcyRoadJT initialize junctiontree relative to icy_roads.bif
func InitIcyRoadJT(bn *BN) (jt *JunctionTree) {
	iw := Clique{
		Name:       "IW",
		Variables:  make([]*Node, 0),
		Neighbours: make([]*Clique, 0),
		Separators: make([]*Separator, 0),
		Table:      make([][]float64, 0),
	}
	iw.AddVariables(bn.GetNode("Icy"), bn.GetNode("Watson"))
	ih := Clique{
		Name:       "IH",
		Variables:  make([]*Node, 0),
		Neighbours: make([]*Clique, 0),
		Separators: make([]*Separator, 0),
		Table:      make([][]float64, 0),
	}
	ih.AddVariables(bn.GetNode("Icy"), bn.GetNode("Holmes"))
	var junc JunctionTree
	return &junc
}
