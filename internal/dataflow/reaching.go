package dataflow

import (
	"fmt"
	"sort"
	"strings"
)

// Package dataflow implements reaching definitions analysis.
// Gói này triển khai phân tích reaching definitions (định nghĩa đạt được).
// 
// Liên hệ với slide Dataflow Analysis, minh họa:
// - Forward dataflow analysis (phân tích luồng dữ liệu tiến)
// - Tập GEN/KILL cho mỗi basic block
// - Lặp fixpoint sử dụng hàm chuyển (transfer function)
// - Tính toán tập IN/OUT
//
// Reaching Definitions: Phân tích xem định nghĩa (gán giá trị) nào có thể
// "đạt được" (reach) một điểm trong chương trình mà không bị gán đè lại.

// BasicBlock - Đại diện cho một node trong CFG với thông tin dataflow
type BasicBlock struct {
	ID    int      // ID của block
	Stmts []string // Các câu lệnh trong block
	
	// GEN: Tập các định nghĩa được tạo ra (generated) trong block này
	// Ví dụ: "x = 5" tạo ra định nghĩa d1 cho biến x
	GEN map[string]bool
	
	// KILL: Tập các định nghĩa bị loại bỏ (killed) bởi block này
	// Ví dụ: Nếu block có "x = 10" thì nó kill tất cả định nghĩa cũ của x
	KILL map[string]bool
	
	// IN: Tập các định nghĩa đạt được điểm vào (entry) của block
	// IN[n] = ∪ OUT[pred] (hợp của OUT từ tất cả predecessor)
	IN map[string]bool
	
	// OUT: Tập các định nghĩa đạt được điểm ra (exit) của block
	// OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])
	OUT map[string]bool
	
	// Succs: Các block kế tiếp (successors) trong CFG
	Succs []*BasicBlock
	
	// Preds: Các block đi trước (predecessors) trong CFG
	Preds []*BasicBlock
}

// RunReachingDefinitions - Demo phân tích reaching definitions
// trên CFG đơn giản gồm 3 blocks
func RunReachingDefinitions() {
	// Ví dụ CFG đại diện cho:
	// Block 1: d1: x = 5  (định nghĩa d1 cho biến x)
	//          d2: y = 3  (định nghĩa d2 cho biến y)
	// Block 2: d3: x = x + y  (định nghĩa d3 cho x, kill d1)
	// Block 3: d4: y = x - 1  (định nghĩa d4 cho y, kill d2)
	//
	// CFG edges: 1 -> 2 -> 3 -> 2 (có vòng lặp)
	
	// Khởi tạo các basic blocks
	// Block 1: khởi tạo x và y
	block1 := &BasicBlock{
		ID:    1,
		Stmts: []string{"d1: x = 5", "d2: y = 3"},
		GEN:   make(map[string]bool),
		KILL:  make(map[string]bool),
		IN:    make(map[string]bool),
		OUT:   make(map[string]bool),
	}
	
	// Block 2: cập nhật x (sử dụng x và y)
	block2 := &BasicBlock{
		ID:    2,
		Stmts: []string{"d3: x = x + y"},
		GEN:   make(map[string]bool),
		KILL:  make(map[string]bool),
		IN:    make(map[string]bool),
		OUT:   make(map[string]bool),
	}
	
	// Block 3: cập nhật y (sử dụng x)
	block3 := &BasicBlock{
		ID:    3,
		Stmts: []string{"d4: y = x - 1"},
		GEN:   make(map[string]bool),
		KILL:  make(map[string]bool),
		IN:    make(map[string]bool),
		OUT:   make(map[string]bool),
	}
	
	// Thiết lập cạnh của CFG: 1 -> 2 -> 3 -> 2 (có vòng lặp back edge)
	// Block 1 có successor là Block 2
	block1.Succs = []*BasicBlock{block2}
	// Block 2 có predecessors là Block 1 và Block 3 (vòng lặp)
	block2.Preds = []*BasicBlock{block1, block3}
	// Block 2 có successor là Block 3
	block2.Succs = []*BasicBlock{block3}
	// Block 3 có predecessor là Block 2
	block3.Preds = []*BasicBlock{block2}
	// Block 3 có successor là Block 2 (tạo vòng lặp)
	block3.Succs = []*BasicBlock{block2}
	
	// Tính toán tập GEN và KILL cho mỗi block
	// 
	// Block 1: sinh ra định nghĩa d1 (x=5) và d2 (y=3)
	block1.GEN["d1"] = true  // định nghĩa d1 cho biến x
	block1.GEN["d2"] = true  // định nghĩa d2 cho biến y
	// Block 1 không kill định nghĩa nào (vì là block đầu tiên)
	
	// Block 2: sinh ra định nghĩa d3 (x=x+y), kill định nghĩa d1 của x
	block2.GEN["d3"] = true   // định nghĩa mới d3 cho x
	block2.KILL["d1"] = true  // kill định nghĩa cũ d1 của x
	
	// Block 3: sinh ra định nghĩa d4 (y=x-1), kill định nghĩa d2 của y
	block3.GEN["d4"] = true   // định nghĩa mới d4 cho y
	block3.KILL["d2"] = true  // kill định nghĩa cũ d2 của y
	
	blocks := []*BasicBlock{block1, block2, block3}
	
	fmt.Println("Initial GEN/KILL sets:")
	for _, b := range blocks {
		fmt.Printf("\nBlock %d:\n", b.ID)
		fmt.Printf("  Statements: %v\n", b.Stmts)
		fmt.Printf("  GEN  = {%s}\n", setToString(b.GEN))
		fmt.Printf("  KILL = {%s}\n", setToString(b.KILL))
	}
	
	// Thực hiện lặp đến điểm bất động (Fixpoint Iteration)
	// Đây là trái tim của thuật toán dataflow analysis
	fmt.Println("\n\n=== Fixpoint Iteration ===")
	fmt.Println("Using transfer functions:")
	fmt.Println("  IN[n]  = ∪ OUT[pred] for all pred of n")
	fmt.Println("  OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])")
	
	iteration := 0
	changed := true  // Cờ để kiểm tra có thay đổi không
	
	// Lặp cho đến khi không còn thay đổi nào (đạt fixpoint)
	for changed {
		iteration++
		changed = false
		fmt.Printf("\n--- Iteration %d ---\n", iteration)
		
		// Với mỗi block, tính toán IN và OUT mới
		for _, b := range blocks {
			// Bước 1: Tính IN[n] = ∪ OUT[pred]
			// IN của block = hợp (union) của OUT từ tất cả predecessor
			newIN := make(map[string]bool)
			for _, pred := range b.Preds {
				for def := range pred.OUT {
					newIN[def] = true
				}
			}
			
			// Bước 2: Tính OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])
			// OUT = những định nghĩa được sinh ra + những định nghĩa đến nhưng không bị kill
			newOUT := make(map[string]bool)
			
			// Thêm tất cả định nghĩa trong GEN[n]
			for def := range b.GEN {
				newOUT[def] = true
			}
			
			// Thêm các định nghĩa từ IN[n] mà không bị KILL[n]
			for def := range newIN {
				if !b.KILL[def] {
					newOUT[def] = true
				}
			}
			
			// Kiểm tra xem có thay đổi so với iteration trước không
			if !setsEqual(b.IN, newIN) || !setsEqual(b.OUT, newOUT) {
				changed = true
			}
			
			// Cập nhật IN và OUT
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

// ===== Các hàm hỗ trợ (Helper functions) =====

// setToString - Chuyển map[string]bool thành chuỗi để in ra
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

// setsEqual - Kiểm tra 2 tập có bằng nhau không
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

