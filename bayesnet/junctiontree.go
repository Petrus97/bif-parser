package bayesnet

import (
	"fmt"

	"github.com/fatih/color"
)

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

// AddCliques to junction tree
func (jt *JunctionTree) AddCliques(cliques ...*Clique) {
	for _, clique := range cliques {
		jt.Cliques = append(jt.Cliques, clique)
	}
}

// SetRoot set the first clique in junction tree as root
func (jt *JunctionTree) SetRoot() {
	jt.Root = jt.Cliques[0]
}

// AddSeparators to junction tree
func (jt *JunctionTree) AddSeparators(separators ...*Separator) {
	for _, sep := range separators {
		jt.Separators = append(jt.Separators, sep)
		sep.LeftClique.AddSeparator(sep)
		sep.RightClique.AddSeparator(sep)
	}
}

// NewSeparator create a new separator from factor table, assign right and left clique
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
	fmt.Println("my separator", sep)
	return &sep
}

func (c *Clique) getSeparatorBetween(other *Clique) *Separator {
	// fmt.Println("c", c.Separators)
	// fmt.Println("othe", other.Separators)
	for _, sep := range c.Separators {
		// fmt.Println("sep", sep.LeftClique.Name, sep.RightClique.Name, c.Name, other.Name)
		if (sep.LeftClique.Name == other.Name && sep.RightClique.Name == c.Name) || (sep.RightClique.Name == other.Name && sep.LeftClique.Name == c.Name) {
			return sep
		}
	}
	return nil
}

func (jt *JunctionTree) collectEvidence(prev *Clique, actual *Clique, sep *Separator) {
	color.Green("Collect evidence")
	actual.Visited = true
	for _, neighbour := range actual.Neighbours {
		fmt.Println(actual.Name, neighbour.Name)
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
	color.Green("Distribute evidence")
	c.Visited = true
	for _, neighbour := range c.Neighbours {
		var sep *Separator
		sep = neighbour.getSeparatorBetween(c)
		if sep != nil {
			if neighbour.Visited == false {
				jt.passMessage(c, neighbour, sep)
				jt.distributeEvidence(neighbour)
			}
		}
	}

}

func (jt *JunctionTree) passMessage(from *Clique, to *Clique, sep *Separator) {
	color.Green("Pass Message")
	fmt.Println("passing message from", from.Name, "to", to.Name)
	fmt.Println(from.Table)
	fmt.Println(to.Table)
	m := difference(from.Variables, sep.Variables)
	color.Red("Var to remove:")
	for i, v := range m {
		fmt.Print(i, ". ", v.Name, " ")
	}
	fmt.Print("\n")
	fmt.Println("from", from.Name, from.Table)
	new := from.Table.Marginalize(true, m...) // t*_s sep table
	color.Red("Var in new: ")
	for i, v := range new.Scope {
		fmt.Print(i, ". ", v.Name, " ")
	}
	fmt.Print("\n")
	fmt.Println("NEW", new)
	fmt.Println("SEP", sep.Table)
	// t*_w x (t*_s / t_S)
	// tmp := DivideFactor(new, sep.Table, true) // t*_s / t_s
	// MultiplyFactor(to.Table, tmp, false)      // t*_w x (t*_s / t_s)
	MultiplyFactor(to.Table, new, false) // t*_w x (t*_s / t_s)
	fmt.Println("to", to.Name, to.Table)
	DivideFactor(to.Table, sep.Table, false)
	fmt.Println("to", to.Name, to.Table)
	sep.Table = new
	fmt.Println("sep", sep.Table)
	color.HiMagenta("HRS", jt.Cliques[1].Table)

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

// EnterEvidence in in a tree
func (jt *JunctionTree) EnterEvidence(n *Node, values []float64) {
	color.Yellow("EnterEvidence")
	fmt.Println("Evidence on", n.Name)
	for _, c := range jt.Cliques {
		if ok := containsNodeV2(c.Table, n); ok == true {
			fmt.Println("Entering evidence in", c.Name)
			err := EnterEvidenceFactor(c.Table, n, values)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// EnterEvidenceFactor add evidence to clique tables (factor)
func EnterEvidenceFactor(old *FactorV2, n *Node, values []float64) error {
	color.Yellow("Enter evidence factor")
	ev := n
	if len(values) != n.Numvalues {
		return fmt.Errorf("Evidence list: %d does not match numvalues: %d", len(values), len(n.CPT))
	}
	// ev.CPT = values
	evfact := CreateFactorV2(ev)
	j := 0
	for i := 0; i < len(evfact.CPT); i++ {
		if j < len(values) {
			evfact.CPT[i] = values[j]
		} else {
			j = 0
			evfact.CPT[i] = values[j]
		}
		j++
	}

	fmt.Println("ev", evfact)
	fmt.Println("old", old)
	for i := 0; i < len(ev.CPT); i++ {
		old.CPT[i] *= ev.CPT[i]
	}
	// MultiplyFactor(old, evfact, false)

	color.Cyan("AFTER EV", old)

	return nil
}

// AddSeparator to clique
func (c *Clique) AddSeparator(sep *Separator) {
	fmt.Println("adding separator...", sep)
	c.Separators = append(c.Separators, sep)
	if c == sep.LeftClique {
		c.Neighbours = append(c.Neighbours, sep.RightClique)
	} else {
		c.Neighbours = append(c.Neighbours, sep.LeftClique)
	}
}

// Propagate perform Hugin algorithm
func (jt *JunctionTree) Propagate() {
	color.Green("Propagate")
	for _, c := range jt.Cliques {
		c.Visited = false
	}
	// Collect evidence
	jt.collectEvidence(nil, jt.Root, nil)
	for _, c := range jt.Cliques {
		c.Visited = false
	}
	color.HiMagenta("HRS", jt.Cliques[1].Table.CPT)
	// Distribute Evidence
	jt.distributeEvidence(jt.Root)
	color.HiMagenta("HRS", jt.Cliques[1].Table.CPT)
	// Normalize
	color.Green("Normalizing")
	for _, c := range jt.Cliques {
		c.Table.Normalize()
	}
	for _, s := range jt.Separators {
		s.Table.Normalize()
	}
}
