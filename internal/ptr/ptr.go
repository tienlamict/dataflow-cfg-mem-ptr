package ptr

import (
	"fmt"
	"runtime"
)

type Node struct {
	Value int
	Next  *Node
}

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

	const numNodes = 10
	fmt.Printf("=== Creating Linked List (%d nodes) ===\n", numNodes)

	head := createLinkedList(numNodes)

	fmt.Println("\nLinked list created. First 5 nodes:")
	printList(head, 5)

	runtime.GC()
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	fmt.Println("\n\n=== Memory Stats (with linked list) ===")
	fmt.Printf("  Alloc:     %.2f KB\n", float64(m1.Alloc)/1024)
	fmt.Printf("  HeapAlloc: %.2f KB\n", float64(m1.HeapAlloc)/1024)
	fmt.Printf("  NumGC:     %d\n", m1.NumGC)

	count := countNodes(head)
	fmt.Printf("\nTotal nodes in list: %d\n", count)

	fmt.Println("\n\n=== Dropping Reference to List ===")
	fmt.Println("Setting head = nil makes all nodes unreachable...")
	fmt.Println("(In C, this would be a memory leak without manual free!)")

	head = nil

	fmt.Println("\nTriggering garbage collection...")
	runtime.GC()
	runtime.GC()

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

func createLinkedList(n int) *Node {
	if n <= 0 {
		return nil
	}

	head := &Node{Value: 1}
	current := head

	for i := 2; i <= n; i++ {
		newNode := &Node{Value: i}
		current.Next = newNode
		current = newNode
	}

	return head
}

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

func countNodes(head *Node) int {
	count := 0
	current := head

	for current != nil {
		count++
		current = current.Next
	}

	return count
}
