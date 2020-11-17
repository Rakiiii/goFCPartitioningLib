package fcpartitioninglib

import (
	lsplib "github.com/Rakiiii/goBipartitonLocalSearch"
)

type solutionTreeNode struct {
	parentNode *solutionTreeNode
	lNode      *solutionTreeNode
	rNode      *solutionTreeNode
	vector     []bool
	isChecked  bool
}

func (s *solutionTreeNode) constructLeftNode() {
	s.lNode = &solutionTreeNode{lNode: nil, rNode: nil, vector: append([]bool{false}, s.vector...), parentNode: s, isChecked: false}
}

func (s *solutionTreeNode) constructRightNode() {
	s.rNode = &solutionTreeNode{lNode: nil, rNode: nil, vector: append([]bool{true}, s.vector...), parentNode: s, isChecked: false}
}

func (s *solutionTreeNode) isFinnalNode(g lsplib.IGraph) bool {
	return len(s.vector) == g.AmountOfVertex()-g.GetAmountOfIndependent()
}

func (s *solutionTreeNode) isAnyNodeCheckable() bool {
	return s.lNode == nil || s.rNode == nil || !s.lNode.isChecked || !s.rNode.isChecked || s.parentNode != nil
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
	partitioner := &FCPartitioner{solutionTreeRoot: &solutionTreeNode{vector: []bool{}, lNode: nil, rNode: nil, isChecked: false, parentNode: nil}}
	partitioner.solutionTreeRoot.constructLeftNode()
	partitioner.solutionTreeRoot.constructRightNode()
	return partitioner
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
	s.isChecked = true
	if len(s.vector) > 4 {
		if s.isFinnalNode(g) {
			newSol := new(lsplib.Solution)

			newSol.Init(g)
			newSol.Vector = append(make([]bool, g.GetAmountOfIndependent()), reverseBool(s.vector)...)
			mark := newSol.CountMark()
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

	checkNode(f.solutionTreeRoot.lNode, g, baseSoulution, groupSize)
	checkNode(f.solutionTreeRoot.rNode, g, baseSoulution, groupSize)

	return baseSoulution, nil
}

func (f *FCPartitioner) PartitionNonRec(g lsplib.IGraph, baseSoulution *FCPartitionSolution, groupSize int) (*FCPartitionSolution, error) {
	if g.GetAmountOfIndependent() <= 4 {
		solution := lsplib.LSPartiotionAlgorithmNonRecFast(g, &baseSoulution.Solution, groupSize)
		baseSoulution.Solution = *solution
		return baseSoulution, nil
	}

	if err := baseSoulution.constructMarkMap(); err != nil {
		return nil, err
	}

	checkNode := f.solutionTreeRoot
	for checkNode.isAnyNodeCheckable() {
		if checkNode.lNode == nil && checkNode.rNode == nil {
			if checkNode.isFinnalNode(g) {
				newSol := new(lsplib.Solution)

				newSol.Init(g)
				newSol.Vector = append(make([]bool, g.GetAmountOfIndependent()), reverseBool(checkNode.vector)...)
				mark := newSol.CountMark()
				if baseSoulution.Value == -1 {
					if flag := newSol.PartIndependent(groupSize); flag {
						baseSoulution.Solution = *newSol
						baseSoulution.CountParameter()
						checkNode.isChecked = true
						checkNode = checkNode.parentNode
						continue
					} else {
						checkNode.isChecked = true
						checkNode = checkNode.parentNode
						continue
					}
				}
				if mark < baseSoulution.Value {
					if flag := newSol.PartIndependent(groupSize); flag {
						if newSol.CountParameter() < baseSoulution.Value {
							baseSoulution.Solution = *newSol
							checkNode.isChecked = true
							checkNode = checkNode.parentNode
							continue
						}
					}
				}
				checkNode.isChecked = true
				checkNode = checkNode.parentNode
				continue
			} else {
				if baseSoulution.Value != -1 && len(checkNode.vector) >= 4 {
					markList := baseSoulution.matchVectorWithAllMarks(checkNode.vector)
					continueFlag := false
					for _, mark := range markList {
						if mark >= baseSoulution.Value {
							continueFlag = true
							break
						}
					}
					if continueFlag {
						checkNode.isChecked = true
						checkNode = checkNode.parentNode
						continue
					}
				}
				checkNode.constructLeftNode()
				checkNode.constructRightNode()
				checkNode = checkNode.lNode
				continue
			}
		} else {
			if checkNode.lNode.isChecked && checkNode.rNode.isChecked && checkNode.parentNode != nil {
				checkNode.isChecked = true
				checkNode = checkNode.parentNode
			} else {
				if !checkNode.lNode.isChecked {
					checkNode = checkNode.lNode
				} else {
					if !checkNode.rNode.isChecked {
						checkNode = checkNode.rNode
					}
				}
			}
		}
	}
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
