package ptr

import (
	"fmt"
	"runtime"
)

// Package ptr demonstrates pointer usage and automatic garbage collection.
// Package này minh họa sử dụng con trỏ và garbage collection tự động.
//
// HOMEWORK YÊU CẦU 4: Minh họa ý tưởng sử dụng và giải phóng pointer
//
// Nội dung demo:
// - Cấu trúc dữ liệu động (linked list)
// - Cấp phát bộ nhớ dựa trên con trỏ
// - Quản lý bộ nhớ tự động vs thủ công (C's malloc/free)
// - Garbage collection cho các object không còn truy cập được
//
// So sánh Go vs C:
// - C: Cần malloc() để cấp phát, free() để giải phóng
// - Go: Chỉ cần tạo object (&Node{}), GC tự động giải phóng

// Node - Biểu diễn một node trong linked list
// Mỗi node chứa:
// - Value: giá trị của node
// - Next: con trỏ đến node tiếp theo
type Node struct {
	Value int   // Giá trị lưu trong node
	Next  *Node // Con trỏ đến node tiếp theo (nil nếu là node cuối)
}

// RunPointerDemo - Demo tạo linked list và garbage collection
//
// Demo này minh họa:
// 1. Tạo cấu trúc dữ liệu động (linked list) bằng con trỏ
// 2. Hiển thị các node và địa chỉ bộ nhớ
// 3. Giải phóng tham chiếu (head = nil)
// 4. GC tự động thu hồi tất cả các node
func RunPointerDemo() {
	fmt.Println("This demo shows pointer usage and garbage collection in Go.")
	fmt.Println()
	fmt.Println("=== Comparison with C ===")
	fmt.Println("In C:")
	fmt.Println("  struct Node* node = (struct Node*)malloc(sizeof(struct Node));")
	fmt.Println("  // ... use node ...")
	fmt.Println("  free(node);  // Manual memory management!")
	fmt.Println()
	fmt.Println("In Go:")
	fmt.Println("  node := &Node{Value: 42}")
	fmt.Println("  // ... use node ...")
	fmt.Println("  // No free() needed - GC handles it automatically!")
	fmt.Println()

	// Tạo linked list với nhiều nodes
	const numNodes = 10
	fmt.Printf("=== Creating Linked List (%d nodes) ===\n", numNodes)

	// Tạo linked list: head -> node1 -> node2 -> ... -> nodeN -> nil
	// Trong C: Mỗi node cần malloc(), sau đó phải free() từng node
	// Trong Go: Chỉ cần &Node{}, GC tự động quản lý
	head := createLinkedList(numNodes)

	// Hiển thị một vài node đầu tiên và địa chỉ bộ nhớ
	fmt.Println("\nLinked list created. First 5 nodes:")
	printList(head, 5)

	// Get memory stats before GC
	runtime.GC()
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	fmt.Println("\n\n=== Memory Stats (with linked list) ===")
	fmt.Printf("  Alloc:     %.2f KB\n", float64(m1.Alloc)/1024)
	fmt.Printf("  HeapAlloc: %.2f KB\n", float64(m1.HeapAlloc)/1024)
	fmt.Printf("  NumGC:     %d\n", m1.NumGC)

	// Count nodes (to verify they're still accessible)
	count := countNodes(head)
	fmt.Printf("\nTotal nodes in list: %d\n", count)

	// Giải phóng tham chiếu đến head
	fmt.Println("\n\n=== Dropping Reference to List ===")
	fmt.Println("Setting head = nil makes all nodes unreachable...")
	fmt.Println("(In C, this would be a memory leak without manual free!)")

	// QUAN TRỌNG: Khi set head = nil:
	// - Không còn cách nào truy cập vào bất kỳ node nào
	// - Trong C: Đây là memory leak nghiêm trọng!
	// - Trong Go: GC sẽ tự động phát hiện và thu hồi tất cả các node
	head = nil

	// Kích hoạt GC để thu hồi bộ nhớ
	fmt.Println("\nTriggering garbage collection...")
	runtime.GC()
	runtime.GC()

	// Get memory stats after GC
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	fmt.Println("\n\n=== Memory Stats (after GC) ===")
	fmt.Printf("  Alloc:     %.2f KB\n", float64(m2.Alloc)/1024)
	fmt.Printf("  HeapAlloc: %.2f KB\n", float64(m2.HeapAlloc)/1024)
	fmt.Printf("  NumGC:     %d\n", m2.NumGC)

	if m1.HeapAlloc > m2.HeapAlloc {
		fmt.Printf("\nMemory reclaimed: %.2f KB\n",
			float64(m1.HeapAlloc-m2.HeapAlloc)/1024)
	} else {
		fmt.Printf("\nMemory change: %.2f KB (GC completed, runtime overhead may vary)\n",
			float64(m2.HeapAlloc-m1.HeapAlloc)/1024)
	}

	fmt.Println("\n\n=== Summary ===")
	fmt.Println("✓ Created linked list with dynamic allocation (like malloc in C)")
	fmt.Println("✓ Dropped all references to the list (head = nil)")
	fmt.Println("✓ GC automatically freed all unreachable nodes")
	fmt.Println("✓ No manual free() needed - Go manages memory automatically!")
	fmt.Println()
	fmt.Println("Key advantage: No memory leaks, no use-after-free bugs!")
}

// ===== Các hàm hỗ trợ (Helper functions) =====

// createLinkedList - Tạo một linked list với n nodes
// Mỗi node được cấp phát động bằng &Node{...}
// Trong C tương đương: malloc(sizeof(struct Node))
func createLinkedList(n int) *Node {
	if n <= 0 {
		return nil
	}

	// Tạo node đầu tiên (head)
	// &Node{Value: 1} tương đương malloc(sizeof(Node)) trong C
	head := &Node{Value: 1}
	current := head

	// Tạo các node còn lại và liên kết với nhau
	for i := 2; i <= n; i++ {
		newNode := &Node{Value: i}
		current.Next = newNode // Liên kết node hiện tại với node mới
		current = newNode      // Di chuyển con trỏ current
	}

	return head
}

// printList - In n node đầu tiên của linked list
// Hiển thị giá trị, địa chỉ của node, và địa chỉ của Next
func printList(head *Node, n int) {
	current := head
	count := 0

	for current != nil && count < n {
		fmt.Printf("  Node %d: Value = %d, Addr = %p, Next = %p\n",
			count+1, current.Value, current, current.Next)
		current = current.Next
		count++
	}

	if current != nil {
		fmt.Printf("  ... (%d more nodes)\n", countNodes(current))
	}
}

// countNodes - Đếm số lượng node trong linked list
func countNodes(head *Node) int {
	count := 0
	current := head

	for current != nil {
		count++
		current = current.Next
	}

	return count
}
