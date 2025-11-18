package ptr

import (
	"fmt"
	"runtime"
)

// Package ptr demonstrates pointer usage and automatic garbage collection.
// This relates to the Pointer slides, demonstrating:
// - Dynamic data structures (linked list)
// - Pointer-based memory allocation
// - Automatic memory management vs manual (C's malloc/free)
// - Garbage collection of unreachable objects

// Node represents a node in a linked list
type Node struct {
	Value int
	Next  *Node
}

// RunPointerDemo demonstrates linked list creation and GC
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
	
	// Create linked list
	const numNodes = 10
	fmt.Printf("=== Creating Linked List (%d nodes) ===\n", numNodes)
	
	head := createLinkedList(numNodes)
	
	// Display first few nodes
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
	
	// Drop reference to head
	fmt.Println("\n\n=== Dropping Reference to List ===")
	fmt.Println("Setting head = nil makes all nodes unreachable...")
	fmt.Println("(In C, this would be a memory leak without manual free!)")
	head = nil
	
	// Trigger GC
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

// createLinkedList creates a linked list with n nodes
func createLinkedList(n int) *Node {
	if n <= 0 {
		return nil
	}
	
	// Create head
	head := &Node{Value: 1}
	current := head
	
	// Create remaining nodes
	for i := 2; i <= n; i++ {
		newNode := &Node{Value: i}
		current.Next = newNode
		current = newNode
	}
	
	return head
}

// printList prints the first n nodes of the linked list
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

// countNodes counts the number of nodes in a linked list
func countNodes(head *Node) int {
	count := 0
	current := head
	
	for current != nil {
		count++
		current = current.Next
	}
	
	return count
}

