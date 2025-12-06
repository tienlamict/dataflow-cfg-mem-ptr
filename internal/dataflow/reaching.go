package dataflow

import (
	"fmt"
	"sort"
	"strings"
)

type BasicBlock struct {
	ID    int
	Stmts []string
	
	GEN map[string]bool
	
	KILL map[string]bool
	
	IN map[string]bool
	
	OUT map[string]bool
	
	Succs []*BasicBlock
	
	Preds []*BasicBlock
}

func RunReachingDefinitions() {
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
	
	block1.Succs = []*BasicBlock{block2}
	block2.Preds = []*BasicBlock{block1, block3}
	block2.Succs = []*BasicBlock{block3}
	block3.Preds = []*BasicBlock{block2}
	block3.Succs = []*BasicBlock{block2}
	
	block1.GEN["d1"] = true
	block1.GEN["d2"] = true
	
	block2.GEN["d3"] = true
	block2.KILL["d1"] = true
	
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
			newIN := make(map[string]bool)
			for _, pred := range b.Preds {
				for def := range pred.OUT {
					newIN[def] = true
				}
			}
			
			newOUT := make(map[string]bool)
			
			for def := range b.GEN {
				newOUT[def] = true
			}
			
			for def := range newIN {
				if !b.KILL[def] {
					newOUT[def] = true
				}
			}
			
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
