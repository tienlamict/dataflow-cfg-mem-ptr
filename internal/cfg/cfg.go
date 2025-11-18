package cfg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Package cfg implements Control Flow Graph construction.
// This relates to the CFG slides, demonstrating:
// - Basic block identification using leader rules
// - CFG construction from pseudo-code with labels and jumps
// - Graph representation with nodes and edges

// Block represents a basic block in the CFG
type Block struct {
	ID    int
	Label int      // Line number of first statement
	Stmts []string // Statements in this block
	Succs []int    // Successor block IDs
}

// CFG represents a Control Flow Graph
type CFG struct {
	Blocks []*Block
}

// RunCFGDemo demonstrates CFG construction from pseudo-code
func RunCFGDemo() {
	// Example pseudo-code with labels and conditional jumps
	// This represents a simple loop structure
	pseudoCode := []string{
		"1: x = 0",
		"2: y = 10",
		"3: if x >= y goto 8",
		"4: sum = 0",
		"5: sum = sum + x",
		"6: x = x + 1",
		"7: goto 3",
		"8: print sum",
		"9: return",
	}
	
	fmt.Println("Input Pseudo-Code:")
	for _, stmt := range pseudoCode {
		fmt.Printf("  %s\n", stmt)
	}
	
	// Build CFG
	cfg := buildCFG(pseudoCode)
	
	// Display results
	fmt.Println("\n\n=== Basic Block Identification ===")
	fmt.Println("Leader rules:")
	fmt.Println("  1. First statement is a leader")
	fmt.Println("  2. Target of a jump is a leader")
	fmt.Println("  3. Statement following a jump is a leader")
	
	fmt.Println("\n\n=== Control Flow Graph ===")
	for _, block := range cfg.Blocks {
		fmt.Printf("\nBlock %d (starts at line %d):\n", block.ID, block.Label)
		fmt.Println("  Statements:")
		for _, stmt := range block.Stmts {
			fmt.Printf("    %s\n", stmt)
		}
		if len(block.Succs) > 0 {
			fmt.Printf("  Successors: %v\n", block.Succs)
		} else {
			fmt.Println("  Successors: [EXIT]")
		}
	}
	
	// Display graph structure
	fmt.Println("\n\n=== Graph Edges ===")
	for _, block := range cfg.Blocks {
		if len(block.Succs) > 0 {
			for _, succ := range block.Succs {
				fmt.Printf("Block %d -> Block %d\n", block.ID, succ)
			}
		} else {
			fmt.Printf("Block %d -> EXIT\n", block.ID)
		}
	}
}

func buildCFG(pseudoCode []string) *CFG {
	// Step 1: Identify leaders
	leaders := make(map[int]bool)
	
	// Rule 1: First statement is a leader
	if len(pseudoCode) > 0 {
		label := extractLabel(pseudoCode[0])
		leaders[label] = true
	}
	
	// Rules 2 & 3: Process jumps
	for i, stmt := range pseudoCode {
		// Check if this is a jump statement
		if isJump(stmt) {
			// Rule 2: Target of jump is a leader
			target := extractJumpTarget(stmt)
			if target > 0 {
				leaders[target] = true
			}
			
			// Rule 3: Statement following jump is a leader
			if i+1 < len(pseudoCode) {
				nextLabel := extractLabel(pseudoCode[i+1])
				leaders[nextLabel] = true
			}
		}
	}
	
	// Step 2: Build basic blocks
	blocks := make([]*Block, 0)
	currentBlock := &Block{ID: 1, Stmts: make([]string, 0)}
	blockID := 1
	labelToBlockID := make(map[int]int)
	
	for _, stmt := range pseudoCode {
		label := extractLabel(stmt)
		
		if leaders[label] && len(currentBlock.Stmts) > 0 {
			// Start a new block
			blocks = append(blocks, currentBlock)
			blockID++
			currentBlock = &Block{ID: blockID, Stmts: make([]string, 0)}
		}
		
		if len(currentBlock.Stmts) == 0 {
			currentBlock.Label = label
			labelToBlockID[label] = blockID
		}
		
		currentBlock.Stmts = append(currentBlock.Stmts, stmt)
	}
	
	// Add last block
	if len(currentBlock.Stmts) > 0 {
		blocks = append(blocks, currentBlock)
	}
	
	// Step 3: Add edges (successors)
	for i, block := range blocks {
		lastStmt := block.Stmts[len(block.Stmts)-1]
		
		if isUnconditionalJump(lastStmt) {
			// Unconditional jump: only one successor
			target := extractJumpTarget(lastStmt)
			if targetBlockID, ok := labelToBlockID[target]; ok {
				block.Succs = append(block.Succs, targetBlockID)
			}
		} else if isConditionalJump(lastStmt) {
			// Conditional jump: two successors (fall-through and jump target)
			target := extractJumpTarget(lastStmt)
			if targetBlockID, ok := labelToBlockID[target]; ok {
				block.Succs = append(block.Succs, targetBlockID)
			}
			// Fall-through to next block
			if i+1 < len(blocks) {
				block.Succs = append(block.Succs, blocks[i+1].ID)
			}
		} else if !isReturn(lastStmt) {
			// Not a jump or return: fall through to next block
			if i+1 < len(blocks) {
				block.Succs = append(block.Succs, blocks[i+1].ID)
			}
		}
		// If it's a return statement, no successors (implicit exit)
	}
	
	return &CFG{Blocks: blocks}
}

// Helper functions

func extractLabel(stmt string) int {
	re := regexp.MustCompile(`^(\d+):`)
	matches := re.FindStringSubmatch(stmt)
	if len(matches) > 1 {
		label, _ := strconv.Atoi(matches[1])
		return label
	}
	return 0
}

func isJump(stmt string) bool {
	return strings.Contains(stmt, "goto") || 
	       (strings.Contains(stmt, "if") && strings.Contains(stmt, "goto"))
}

func isConditionalJump(stmt string) bool {
	return strings.Contains(stmt, "if") && strings.Contains(stmt, "goto")
}

func isUnconditionalJump(stmt string) bool {
	return strings.HasPrefix(strings.TrimSpace(strings.SplitN(stmt, ":", 2)[1]), "goto") &&
	       !strings.Contains(stmt, "if")
}

func isReturn(stmt string) bool {
	return strings.Contains(stmt, "return")
}

func extractJumpTarget(stmt string) int {
	re := regexp.MustCompile(`goto\s+(\d+)`)
	matches := re.FindStringSubmatch(stmt)
	if len(matches) > 1 {
		target, _ := strconv.Atoi(matches[1])
		return target
	}
	return 0
}

