package fcpartitioninglib

import (
	"fmt"
	"testing"

	lsplib "github.com/Rakiiii/goBipartitonLocalSearch"
)

var newTestDebugFlag bool = false

var Dir string = "Testing/"

var testGraphPerfectMatching string = Dir + "GetPerfectMatchingGraph"
var testGraphGetHungryContractedGraphNI = Dir + "GetHungryContractedGraphNI"
var benchGraph string = Dir + "graph_bench_1"

func TestConstructMarkableSet(t *testing.T) {
	if newTestDebugFlag {
		t.Skip()
	}

	fmt.Println("Start TestConstructMarkableSet...")
	graph := lsplib.NewGraph()
	graph.ParseGraph(testGraphGetHungryContractedGraphNI)

	graph.HungryNumIndependent()

	solution := NewFCPartitionSolution(graph)

	solution.constructMarkableSet()

	for key, value := range solution.markMap {
		fmt.Println("key:", key, " valur:", value)
	}

	fmt.Println("TestConstructMarkableSet=[ok]")
}

func TestConstructMarkMap(t *testing.T) {
	if newTestDebugFlag {
		t.Skip()
	}

	fmt.Println("Start TestConstructMarkMap....")
	graph := lsplib.NewGraph()
	graph.ParseGraph(benchGraph)

	graph.HungryNumIndependent()

	solution := NewFCPartitionSolution(graph)

	solution.constructMarkMap()

	for key, value := range solution.markMap {
		fmt.Println("key:", key, " valur:", value)
	}

	fmt.Println("TestConstructMarkMap=[ok]")
}

func TestFcPartiotioner(t *testing.T) {
	if newTestDebugFlag {
		t.Skip()
	}

	fmt.Println("Start TestFcPartiotioner...")

	graph := lsplib.NewGraph()
	if err := graph.ParseGraph(benchGraph); err != nil {
		fmt.Println(err)
		return
	}

	groupSize := graph.AmountOfVertex() / 2

	res := NewFCPartitionSolution(graph)
	graph.HungryNumIndependent()

	partitioner := NewFCPartitioner()

	newRes, err := partitioner.Partition(graph, res, groupSize)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(newRes.Value)
	correctSolution := []bool{false, false, true, true, true, false, true, false, true, false, false, true, false, false, false, false, false, true, false, true, true, true, true, false, true, true, false, true, true, false}
	amountOfTrues := 0
	for pos, i := range newRes.Vector {
		if correctSolution[pos] != i {
			t.Error("Wrong partition at:", pos, " expected:", correctSolution[pos], " found:", i)
		}
		if i {
			fmt.Print("1 ")
			amountOfTrues++
		} else {
			fmt.Print("0 ")
		}
	}
	fmt.Println()

	fmt.Println("first group size:", amountOfTrues, " second group size:", len(res.Vector)-amountOfTrues, " expected group size:", groupSize)

	if newRes.Value != 14 {
		t.Error("Wrong value in partition")
	}

	if amountOfTrues != 15 {
		t.Error("Wrong disbalance in partition")
	}

	fmt.Println("TestFcPartiotioner=[ok]")
}

func TestPartitionNonRec(t *testing.T) {
	if newTestDebugFlag {
		t.Skip()
	}

	fmt.Println("Start TestPartitionNonRec...")

	graph := lsplib.NewGraph()
	if err := graph.ParseGraph(benchGraph); err != nil {
		fmt.Println(err)
		return
	}

	groupSize := graph.AmountOfVertex() / 2

	res := NewFCPartitionSolution(graph)
	graph.HungryNumIndependent()

	partitioner := NewFCPartitioner()

	newRes, err := partitioner.PartitionNonRec(graph, res, groupSize)

	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println(newRes.Value)
	correctSolution := []bool{false, false, true, true, true, false, true, false, true, false, false, true, false, false, false, false, false, true, false, true, true, true, true, false, true, true, false, true, true, false}
	amountOfTrues := 0
	for pos, i := range newRes.Vector {
		if correctSolution[pos] != i {
			t.Error("Wrong partition at:", pos, " expected:", correctSolution[pos], " found:", i)
		}
		if i {
			fmt.Print("1 ")
			amountOfTrues++
		} else {
			fmt.Print("0 ")
		}
	}
	fmt.Println()

	fmt.Println("first group size:", amountOfTrues, " second group size:", len(res.Vector)-amountOfTrues, " expected group size:", groupSize)

	if newRes.Value != 14 {
		t.Error("Wrong value in partition")
	}

	if amountOfTrues != 15 {
		t.Error("Wrong disbalance in partition")
	}

	fmt.Println("TestPartitionNonRec=[ok]")
}
