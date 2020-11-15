package fcpartitioninglib

import (
	lsplib "github.com/Rakiiii/goBipartitonLocalSearch"
)

type solutionTreeNode struct {
	lNode  *solutionTreeNode
	rNode  *solutionTreeNode
	vector []bool
}

func (s *solutionTreeNode) constructLeftNode() {
	s.lNode = &solutionTreeNode{lNode: nil, rNode: nil, vector: append([]bool{false}, s.vector...)}
}

func (s *solutionTreeNode) constructRightNode() {
	s.rNode = &solutionTreeNode{lNode: nil, rNode: nil, vector: append([]bool{true}, s.vector...)}
}

func (s *solutionTreeNode) isFinnalNode(g lsplib.IGraph) bool {
	return len(s.vector) == g.AmountOfVertex()-g.GetAmountOfIndependent()
}

func reverseBool(s []bool) []bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

type FCPartitioner struct {
	solutionTreeRoot *solutionTreeNode
}

func NewFCPartitioner() *FCPartitioner {
	return &FCPartitioner{solutionTreeRoot: &solutionTreeNode{vector: nil, lNode: &solutionTreeNode{vector: []bool{false}, lNode: nil, rNode: nil}, rNode: &solutionTreeNode{vector: []bool{true}, lNode: nil, rNode: nil}}}
}

func boolVectorToInt(b []bool) []int {
	i := make([]int, len(b))
	for pos, val := range b {
		if val {
			i[pos] = 1
		} else {
			i[pos] = 0
		}
	}
	return i
}

func checkNode(s *solutionTreeNode, g lsplib.IGraph, baseSoulution *FCPartitionSolution, groupSize int) {
	// fmt.Println("checking solution:", boolVectorToInt(s.vector), " best value:", baseSoulution.Value)
	if len(s.vector) > 4 {
		if s.isFinnalNode(g) {
			newSol := new(lsplib.Solution)

			newSol.Init(g)
			newSol.Vector = append(make([]bool, g.GetAmountOfIndependent()), reverseBool(s.vector)...)
			//newSol.CountMark()
			mark := newSol.CountMark()
			// fmt.Println("full solution:", boolVectorToInt(newSol.Vector), " mark is:", mark)
			if baseSoulution.Value == -1 {
				if flag := newSol.PartIndependent(groupSize); flag {
					baseSoulution.Solution = *newSol
					baseSoulution.CountParameter()
					return
				} else {
					return
				}
			}
			if mark < baseSoulution.Value {
				if flag := newSol.PartIndependent(groupSize); flag {
					if newSol.CountParameter() < baseSoulution.Value {
						baseSoulution.Solution = *newSol
						return
					}
				}
			}
			return
		}

		if baseSoulution.Value != -1 {
			markList := baseSoulution.matchVectorWithAllMarks(s.vector)
			for _, mark := range markList {
				if mark >= baseSoulution.Value {
					// fmt.Println("solution droped:", boolVectorToInt(s.vector), " with mark:", mark, " and base value:", baseSoulution.Value)
					return
				}
			}
		}
		s.constructLeftNode()
		s.constructRightNode()
		checkNode(s.lNode, g, baseSoulution, groupSize)
		checkNode(s.rNode, g, baseSoulution, groupSize)

	} else {
		s.constructLeftNode()
		s.constructRightNode()
		checkNode(s.lNode, g, baseSoulution, groupSize)
		checkNode(s.rNode, g, baseSoulution, groupSize)
	}
}

func (f *FCPartitioner) Partition(g lsplib.IGraph, baseSoulution *FCPartitionSolution, groupSize int) (*FCPartitionSolution, error) {
	if g.GetAmountOfIndependent() <= 4 {
		solution := lsplib.LSPartiotionAlgorithmNonRecFast(g, &baseSoulution.Solution, groupSize)
		baseSoulution.Solution = *solution
		return baseSoulution, nil
	}

	if err := baseSoulution.constructMarkMap(); err != nil {
		return nil, err
	}

	// fmt.Println("mark map output")
	// for key, value := range baseSoulution.markMap {
	// 	fmt.Println("verteces:", key, " mark:", value)
	// }

	checkNode(f.solutionTreeRoot.lNode, g, baseSoulution, groupSize)
	checkNode(f.solutionTreeRoot.rNode, g, baseSoulution, groupSize)

	// for i := 4; i < g.GetAmountOfIndependent(); i++{
	// 	baseSoulution.fcvector = make([]int, i)

	// 	for vertexPosition := 0 ; vertexPosition < i; vertexPosition ++{
	// 		baseSoulution.fcvector[vertexPosition] = vertexPosition
	// 	}

	// 	for vertexPosition := 0 ; vertexPosition < i; vertexPosition ++{
	// 		for vertex := vertexPosition; vertex < g.GetAmountOfIndependent(); i ++{
	// 			baseSoulution.fcvector[vertexPosition] = vertex
	// 			subPartiotion := make([]bool,i)
	// 			for k := 0; k < int(math.Pow(2,float64(i))); k++ {

	// 			}
	// 		}
	// 	}
	// }

	return baseSoulution, nil
}

func setDependentAsBinnary(num int64, vector []bool) {
	for i := len(vector) - 1; i >= 0; i-- {
		if num%2 == 1 {
			vector[i] = true
		} else {
			vector[i] = false
		}
		num /= 2
	}
}
