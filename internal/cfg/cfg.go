package cfg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Package cfg implements Control Flow Graph construction.
// Package này minh họa xây dựng đồ thị luồng điều khiển (CFG).
//
// HOMEWORK YÊU CẦU 2: Minh họa ý tưởng phân tích mã nguồn chương trình sử dụng kỹ thuật CFG
//
// Nội dung demo:
// - Xác định basic block bằng quy tắc leaders
// - Xây dựng CFG từ pseudo-code có nhãn và lệnh nhảy
// - Biểu diễn đồ thị với các node (blocks) và cạnh (edges)
//
// Quy tắc Leaders (để xác định đầu của basic block):
// 1. Câu lệnh đầu tiên là một leader
// 2. Đích của lệnh nhảy (goto) là một leader  
// 3. Câu lệnh ngay sau lệnh nhảy là một leader

// Block - Biểu diễn một basic block trong CFG
// Basic block là một chuỗi câu lệnh liên tiếp mà:
// - Chỉ có thể vào qua câu lệnh đầu tiên
// - Chỉ có thể ra qua câu lệnh cuối cùng
// - Không có nhảy vào giữa block
type Block struct {
	ID    int        // ID của block
	Label int        // Số dòng của câu lệnh đầu tiên
	Stmts []string   // Danh sách câu lệnh trong block
	Succs []int      // Danh sách ID của các block kế tiếp
}

// CFG - Biểu diễn đồ thị luồng điều khiển
type CFG struct {
	Blocks []*Block  // Danh sách tất cả các blocks
}

// RunCFGDemo - Demo xây dựng CFG từ pseudo-code
//
// Demo này minh họa cách xây dựng Control Flow Graph từ mã nguồn có:
// - Nhãn (labels): 1:, 2:, 3:, ...
// - Lệnh nhảy có điều kiện: if x >= y goto 8
// - Lệnh nhảy vô điều kiện: goto 3
// - Lệnh return
//
// Pseudo-code dưới đây biểu diễn một vòng lặp đơn giản:
// - Khởi tạo x=0, y=10
// - Lặp cho đến khi x >= y
// - Mỗi lần lặp: cộng x vào sum, tăng x lên 1
func RunCFGDemo() {
	// Mã pseudo-code mẫu
	pseudoCode := []string{
		"1: x = 0",             // Khởi tạo x
		"2: y = 10",            // Khởi tạo y
		"3: if x >= y goto 8",  // Điều kiện thoát vòng lặp
		"4: sum = 0",           // Khởi tạo sum
		"5: sum = sum + x",     // Tính tổng
		"6: x = x + 1",         // Tăng x
		"7: goto 3",            // Quay lại kiểm tra điều kiện
		"8: print sum",         // In kết quả
		"9: return",            // Kết thúc
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

// buildCFG - Xây dựng CFG từ danh sách pseudo-code
func buildCFG(pseudoCode []string) *CFG {
	// Bước 1: Xác định các leaders (câu lệnh đầu của mỗi basic block)
	leaders := make(map[int]bool)
	
	// Quy tắc 1: Câu lệnh đầu tiên là một leader
	if len(pseudoCode) > 0 {
		label := extractLabel(pseudoCode[0])
		leaders[label] = true
	}
	
	// Quy tắc 2 & 3: Xử lý các lệnh nhảy
	for i, stmt := range pseudoCode {
		// Kiểm tra xem có phải lệnh nhảy không
		if isJump(stmt) {
			// Quy tắc 2: Đích của lệnh nhảy là một leader
			target := extractJumpTarget(stmt)
			if target > 0 {
				leaders[target] = true
			}
			
			// Quy tắc 3: Câu lệnh ngay sau lệnh nhảy là một leader
			if i+1 < len(pseudoCode) {
				nextLabel := extractLabel(pseudoCode[i+1])
				leaders[nextLabel] = true
			}
		}
	}
	
	// Bước 2: Xây dựng các basic blocks
	// Mỗi block bắt đầu tại một leader và kết thúc trước leader tiếp theo
	blocks := make([]*Block, 0)
	currentBlock := &Block{ID: 1, Stmts: make([]string, 0)}
	blockID := 1
	labelToBlockID := make(map[int]int)  // Map từ label sang block ID
	
	for _, stmt := range pseudoCode {
		label := extractLabel(stmt)
		
		// Nếu gặp leader mới và block hiện tại đã có câu lệnh
		if leaders[label] && len(currentBlock.Stmts) > 0 {
			// Lưu block hiện tại và bắt đầu block mới
			blocks = append(blocks, currentBlock)
			blockID++
			currentBlock = &Block{ID: blockID, Stmts: make([]string, 0)}
		}
		
		// Nếu là câu lệnh đầu tiên của block, lưu label
		if len(currentBlock.Stmts) == 0 {
			currentBlock.Label = label
			labelToBlockID[label] = blockID
		}
		
		// Thêm câu lệnh vào block hiện tại
		currentBlock.Stmts = append(currentBlock.Stmts, stmt)
	}
	
	// Thêm block cuối cùng
	if len(currentBlock.Stmts) > 0 {
		blocks = append(blocks, currentBlock)
	}
	
	// Bước 3: Thêm các cạnh (edges) giữa các blocks
	// Xác định successor của mỗi block dựa vào câu lệnh cuối cùng
	for i, block := range blocks {
		lastStmt := block.Stmts[len(block.Stmts)-1]
		
		if isUnconditionalJump(lastStmt) {
			// Lệnh nhảy vô điều kiện (goto): chỉ có 1 successor
			target := extractJumpTarget(lastStmt)
			if targetBlockID, ok := labelToBlockID[target]; ok {
				block.Succs = append(block.Succs, targetBlockID)
			}
		} else if isConditionalJump(lastStmt) {
			// Lệnh nhảy có điều kiện (if...goto): có 2 successors
			// 1. Nhảy đến đích (nếu điều kiện đúng)
			target := extractJumpTarget(lastStmt)
			if targetBlockID, ok := labelToBlockID[target]; ok {
				block.Succs = append(block.Succs, targetBlockID)
			}
			// 2. Rơi xuống block tiếp theo (nếu điều kiện sai)
			if i+1 < len(blocks) {
				block.Succs = append(block.Succs, blocks[i+1].ID)
			}
		} else if !isReturn(lastStmt) {
			// Không phải nhảy hay return: rơi xuống block tiếp theo
			if i+1 < len(blocks) {
				block.Succs = append(block.Succs, blocks[i+1].ID)
			}
		}
		// Nếu là lệnh return: không có successor (thoát chương trình)
	}
	
	return &CFG{Blocks: blocks}
}

// ===== Các hàm hỗ trợ (Helper functions) =====

// extractLabel - Trích xuất số label từ câu lệnh (vd: "3: x = 5" -> 3)
func extractLabel(stmt string) int {
	re := regexp.MustCompile(`^(\d+):`)
	matches := re.FindStringSubmatch(stmt)
	if len(matches) > 1 {
		label, _ := strconv.Atoi(matches[1])
		return label
	}
	return 0
}

// isJump - Kiểm tra xem có phải lệnh nhảy không (goto hoặc if...goto)
func isJump(stmt string) bool {
	return strings.Contains(stmt, "goto") || 
	       (strings.Contains(stmt, "if") && strings.Contains(stmt, "goto"))
}

// isConditionalJump - Kiểm tra xem có phải lệnh nhảy có điều kiện không
func isConditionalJump(stmt string) bool {
	return strings.Contains(stmt, "if") && strings.Contains(stmt, "goto")
}

// isUnconditionalJump - Kiểm tra xem có phải lệnh nhảy vô điều kiện không
func isUnconditionalJump(stmt string) bool {
	return strings.HasPrefix(strings.TrimSpace(strings.SplitN(stmt, ":", 2)[1]), "goto") &&
	       !strings.Contains(stmt, "if")
}

// isReturn - Kiểm tra xem có phải lệnh return không
func isReturn(stmt string) bool {
	return strings.Contains(stmt, "return")
}

// extractJumpTarget - Trích xuất đích nhảy từ lệnh goto (vd: "goto 8" -> 8)
func extractJumpTarget(stmt string) int {
	re := regexp.MustCompile(`goto\s+(\d+)`)
	matches := re.FindStringSubmatch(stmt)
	if len(matches) > 1 {
		target, _ := strconv.Atoi(matches[1])
		return target
	}
	return 0
}

