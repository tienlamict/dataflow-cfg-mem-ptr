package cfg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Block struct {
	ID    int
	Label int
	Stmts []string
	Succs []int
}

type CFG struct {
	Blocks []*Block
}

func RunCFGDemo() {
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
	
	cfg := buildCFG(pseudoCode)
	
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
	leaders := make(map[int]bool)
	
	if len(pseudoCode) > 0 {
		label := extractLabel(pseudoCode[0])
		leaders[label] = true
	}
	
	for i, stmt := range pseudoCode {
		if isJump(stmt) {
			target := extractJumpTarget(stmt)
			if target > 0 {
				leaders[target] = true
			}
			
			if i+1 < len(pseudoCode) {
				nextLabel := extractLabel(pseudoCode[i+1])
				leaders[nextLabel] = true
			}
		}
	}
	
	blocks := make([]*Block, 0)
	currentBlock := &Block{ID: 1, Stmts: make([]string, 0)}
	blockID := 1
	labelToBlockID := make(map[int]int)
	
	for _, stmt := range pseudoCode {
		label := extractLabel(stmt)
		
		if leaders[label] && len(currentBlock.Stmts) > 0 {
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
	
	if len(currentBlock.Stmts) > 0 {
		blocks = append(blocks, currentBlock)
	}
	
	for i, block := range blocks {
		lastStmt := block.Stmts[len(block.Stmts)-1]
		
		if isUnconditionalJump(lastStmt) {
			target := extractJumpTarget(lastStmt)
			if targetBlockID, ok := labelToBlockID[target]; ok {
				block.Succs = append(block.Succs, targetBlockID)
			}
		} else if isConditionalJump(lastStmt) {
			target := extractJumpTarget(lastStmt)
			if targetBlockID, ok := labelToBlockID[target]; ok {
				block.Succs = append(block.Succs, targetBlockID)
			}
			if i+1 < len(blocks) {
				block.Succs = append(block.Succs, blocks[i+1].ID)
			}
		} else if !isReturn(lastStmt) {
			if i+1 < len(blocks) {
				block.Succs = append(block.Succs, blocks[i+1].ID)
			}
		}
	}
	
	return &CFG{Blocks: blocks}
}

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
