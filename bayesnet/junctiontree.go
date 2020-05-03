package bayesnet

/* (V)----|S|----(W)
 *
 * S.table == ∑ {W\S} W.table == ∑ {V\S} V.table  --> global consistency
 *
 *
 */

// Clique aka Vertex of Junction Tree
type Clique struct {
	Variables  []*Node // Nodes that form the cluster
	Neighbours []*Clique
	Separators []*Separator
	Table      [][]float64
}

// Separator -> intersection of nodes between two Cliques
type Separator struct {
	Variables []*Node // Node intersection between two cluster
	Table     [][]float64
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
