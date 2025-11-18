package dataflow

import (
	"fmt"
	"sort"
	"strings"
)

// Package dataflow implements reaching definitions analysis.
// This relates to the Dataflow Analysis slides, demonstrating:
// - Forward dataflow analysis
// - GEN/KILL sets for each basic block
// - Fixpoint iteration using the transfer function
// - IN/OUT sets computation

// BasicBlock represents a node in the CFG with dataflow information
type BasicBlock struct {
	ID    int
	Stmts []string
	// GEN: definitions generated in this block
	GEN map[string]bool
	// KILL: definitions killed in this block
	KILL map[string]bool
	// IN: definitions reaching the entry of this block
	IN map[string]bool
	// OUT: definitions reaching the exit of this block
	OUT map[string]bool
	// Successors in the CFG
	Succs []*BasicBlock
	// Predecessors in the CFG
	Preds []*BasicBlock
}

// RunReachingDefinitions demonstrates reaching definitions analysis
// on a simple 3-block CFG
func RunReachingDefinitions() {
	// Example CFG representing:
	// Block 1: d1: x = 5
	//          d2: y = 3
	// Block 2: d3: x = x + y
	// Block 3: d4: y = x - 1
	
	// Initialize blocks
	block1 := &BasicBlock{
		ID:    1,
		Stmts: []string{"d1: x = 5", "d2: y = 3"},
		GEN:   make(map[string]bool),
		KILL:  make(map[string]bool),
		IN:    make(map[string]bool),
		OUT:   make(map[string]bool),
	}
	
	block2 := &BasicBlock{
		ID:    2,
		Stmts: []string{"d3: x = x + y"},
		GEN:   make(map[string]bool),
		KILL:  make(map[string]bool),
		IN:    make(map[string]bool),
		OUT:   make(map[string]bool),
	}
	
	block3 := &BasicBlock{
		ID:    3,
		Stmts: []string{"d4: y = x - 1"},
		GEN:   make(map[string]bool),
		KILL:  make(map[string]bool),
		IN:    make(map[string]bool),
		OUT:   make(map[string]bool),
	}
	
	// Set up CFG edges: 1 -> 2 -> 3, and 3 -> 2 (loop)
	block1.Succs = []*BasicBlock{block2}
	block2.Preds = []*BasicBlock{block1, block3}
	block2.Succs = []*BasicBlock{block3}
	block3.Preds = []*BasicBlock{block2}
	block3.Succs = []*BasicBlock{block2}
	
	// Compute GEN and KILL sets
	// Block 1: generates d1 (x=5) and d2 (y=3)
	block1.GEN["d1"] = true
	block1.GEN["d2"] = true
	// Block 1 kills any other definitions of x and y (none in this case for simplicity)
	
	// Block 2: generates d3 (x=x+y), kills d1
	block2.GEN["d3"] = true
	block2.KILL["d1"] = true
	
	// Block 3: generates d4 (y=x-1), kills d2
	block3.GEN["d4"] = true
	block3.KILL["d2"] = true
	
	blocks := []*BasicBlock{block1, block2, block3}
	
	fmt.Println("Initial GEN/KILL sets:")
	for _, b := range blocks {
		fmt.Printf("\nBlock %d:\n", b.ID)
		fmt.Printf("  Statements: %v\n", b.Stmts)
		fmt.Printf("  GEN  = {%s}\n", setToString(b.GEN))
		fmt.Printf("  KILL = {%s}\n", setToString(b.KILL))
	}
	
	// Perform fixpoint iteration
	fmt.Println("\n\n=== Fixpoint Iteration ===")
	fmt.Println("Using transfer functions:")
	fmt.Println("  IN[n]  = ∪ OUT[pred] for all pred of n")
	fmt.Println("  OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])")
	
	iteration := 0
	changed := true
	
	for changed {
		iteration++
		changed = false
		fmt.Printf("\n--- Iteration %d ---\n", iteration)
		
		for _, b := range blocks {
			// Compute IN[n] = ∪ OUT[pred]
			newIN := make(map[string]bool)
			for _, pred := range b.Preds {
				for def := range pred.OUT {
					newIN[def] = true
				}
			}
			
			// Compute OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])
			newOUT := make(map[string]bool)
			// Add GEN[n]
			for def := range b.GEN {
				newOUT[def] = true
			}
			// Add (IN[n] - KILL[n])
			for def := range newIN {
				if !b.KILL[def] {
					newOUT[def] = true
				}
			}
			
			// Check if changed
			if !setsEqual(b.IN, newIN) || !setsEqual(b.OUT, newOUT) {
				changed = true
			}
			
			b.IN = newIN
			b.OUT = newOUT
			
			fmt.Printf("Block %d: IN={%s}, OUT={%s}\n", 
				b.ID, setToString(b.IN), setToString(b.OUT))
		}
	}
	
	fmt.Printf("\nFixpoint reached after %d iterations.\n", iteration)
	
	fmt.Println("\n\n=== Final Results ===")
	for _, b := range blocks {
		fmt.Printf("\nBlock %d:\n", b.ID)
		fmt.Printf("  Statements: %v\n", b.Stmts)
		fmt.Printf("  IN  = {%s}\n", setToString(b.IN))
		fmt.Printf("  OUT = {%s}\n", setToString(b.OUT))
	}
}

// Helper functions

func setToString(set map[string]bool) string {
	if len(set) == 0 {
		return ""
	}
	elements := make([]string, 0, len(set))
	for k := range set {
		elements = append(elements, k)
	}
	sort.Strings(elements)
	return strings.Join(elements, ", ")
}

func setsEqual(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

