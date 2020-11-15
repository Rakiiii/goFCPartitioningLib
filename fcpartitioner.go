package fcpartitioninglib

import (
	lsplib "github.com/Rakiiii/goBipartitonLocalSearch"
	cglib "github.com/Rakiiii/goCoarseningGraphLib"
	gopair "github.com/Rakiiii/goPair"
	pmlib "github.com/Rakiiii/goPerfectMathcingLib"
)

type vertexConfiguration struct {
	FirstSubGraph  gopair.IntPair
	SecondSubGraph gopair.IntPair
}

const randomNumber int64 = 12389712371246

func (v *vertexConfiguration) toFixedVertexesSlice() []gopair.IntPair {
	return []gopair.IntPair{v.FirstSubGraph, v.SecondSubGraph}
}

func (v *vertexConfiguration) maxVertexNumber() int {
	maxInt := -1
	for _, pair := range v.toFixedVertexesSlice() {
		if pair.First > maxInt {
			maxInt = pair.First
		}
		if pair.Second > maxInt {
			maxInt = pair.Second
		}
	}
	return maxInt
}

func (v *vertexConfiguration) minVertexNumber() int {
	//Граф на более чем 10000 вершин в любом случае точно разбить не получиться
	minInt := 10000
	for _, pair := range v.toFixedVertexesSlice() {
		if pair.First < minInt {
			minInt = pair.First
		}
		if pair.Second < minInt {
			minInt = pair.Second
		}
	}
	return minInt + 1
}

func (v *vertexConfiguration) String() string {
	return "firstSubGraph-> first:" + string(v.FirstSubGraph.First) + " second:" + string(v.FirstSubGraph.Second) + " secondSubGraph-> first:" + string(v.SecondSubGraph.First) + " second:" + string(v.SecondSubGraph.Second)
}

func newVertexConfiguration(fv, fe, sv, se int) *vertexConfiguration {
	return &vertexConfiguration{FirstSubGraph: gopair.IntPair{First: fv, Second: fe}, SecondSubGraph: gopair.IntPair{First: sv, Second: se}}
}

type FCPartitionSolution struct {
	lsplib.Solution
	markMap  map[vertexConfiguration]int64
	fcvector []int
	banList  []vertexConfiguration
}

func NewFCPartitionSolution(g lsplib.IGraph) *FCPartitionSolution {
	baseSolution := lsplib.Solution{}
	baseSolution.Init(g)
	return &FCPartitionSolution{Solution: baseSolution, markMap: map[vertexConfiguration]int64{}, fcvector: make([]int, 0), banList: make([]vertexConfiguration, 0)}
}

func (c *FCPartitionSolution) constructMarkMap() error {
	matcher := pmlib.NewRandomMathcerWithNilFixedVertexes()
	c.constructMarkableSet()
	contractableGraph := cglib.NewGraph(c.Solution.Gr)
	for set := range c.markMap {
		matcher.SetFixedVertexes(set.toFixedVertexesSlice()) //= pmlib.NewRandomMathcerWithFixedVertexes(set.toFixedVertexesSlice())
		// matcher = &pmlib.RandomMathcerWithFixedVertexes{FixedVertexes: set.toFixedVertexesSlice(), RandomMatcher: pmlib.RandomMatcher{Rnd: rand.New(rand.NewSource(randomNumber))}}
		if matcher.IsPerfectMatchingExist(contractableGraph.IGraph) {
			contractedGraph, contr, err := contractableGraph.GetPerfectlyContractedGraph(matcher, matcher)
			if err == pmlib.NoPerfectMatching {
				delete(c.markMap, set)
			} else if err != nil {
				return err
			} else {
				contractedGraph.HungryNumIndependent()
				subSolution := *lsplib.LSPartiotionAlgorithmNonRecFast(contractedGraph.IGraph, nil, contractedGraph.AmountOfVertex()/2)
				value := cglib.UncontractedGraphBipartition(contr, subSolution.Vector)
				solution := lsplib.Solution{}
				solution.Init(contractableGraph)
				solution.Vector = value
				c.markMap[set] = solution.CountParameter()
			}
		} else {
			delete(c.markMap, set)
		}
	}
	return nil
}

func (c *FCPartitionSolution) matchVectorWithAllMarks(vector []bool) []int64 {
	ver := c.Gr.AmountOfVertex()
	markList := make([]int64, 0)
	for markConfig, mark := range c.markMap {
		len := len(vector)
		if ver-markConfig.minVertexNumber() < len {
			if vector[len-(ver-markConfig.FirstSubGraph.First)] == vector[len-(ver-markConfig.FirstSubGraph.Second)] && vector[len-(ver-markConfig.SecondSubGraph.First)] == vector[len-(ver-markConfig.SecondSubGraph.Second)] {
				if vector[len-(ver-markConfig.FirstSubGraph.First)] != vector[len-(ver-markConfig.SecondSubGraph.First)] {
					markList = append(markList, mark)
				}
			}
		}
	}
	return markList
}

func (c *FCPartitionSolution) constructMarkableSet() {
	indep := c.Solution.Gr.GetAmountOfIndependent()
	for vertex := indep; vertex < c.Solution.Gr.AmountOfVertex(); vertex++ {
		for _, edge := range c.Solution.Gr.GetEdges(vertex) {
			if edge > indep {
				for secVertex := vertex + 1; secVertex < c.Solution.Gr.AmountOfVertex(); secVertex++ {
					if secVertex != vertex && secVertex != edge {
						for _, secEdge := range c.Solution.Gr.GetEdges(secVertex) {
							if secEdge != vertex && secEdge != edge && secEdge > indep {
								c.markMap[*newVertexConfiguration(vertex, edge, secVertex, secEdge)] = -1
							}
						}
					}
				}
			}
		}
	}
}
