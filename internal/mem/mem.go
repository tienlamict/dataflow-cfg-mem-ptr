package mem

import (
	"fmt"
	"runtime"
)

func RunMemoryDemo() {
	fmt.Println("This demo shows Go's automatic memory management.")
	fmt.Println("Unlike C/C++ where you manually malloc/free, Go uses GC.")
	fmt.Println()
	
	runtime.GC()
	runtime.GC()
	
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Println("=== Initial Memory Stats ===")
	printMemStats(&m1)
	
	fmt.Println("\n\n=== Allocating Memory ===")
	const numObjects = 1000
	const objectSize = 1024
	
	fmt.Printf("Allocating %d objects of %d integers each...\n", numObjects, objectSize)
	
	objects := make([]*[1024]int, numObjects)
	
	for i := 0; i < numObjects; i++ {
		objects[i] = new([1024]int)
		
		for j := 0; j < objectSize; j++ {
			objects[i][j] = i + j
		}
	}
	
	fmt.Printf("Allocated %d objects (approximately %.2f MB)\n", 
		numObjects, float64(numObjects*objectSize*8)/(1024*1024))
	
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Println("\n\n=== Memory Stats After Allocation ===")
	printMemStats(&m2)
	
	fmt.Printf("\nMemory increase:\n")
	fmt.Printf("  Alloc:      +%.2f MB\n", float64(m2.Alloc-m1.Alloc)/(1024*1024))
	fmt.Printf("  TotalAlloc: +%.2f MB\n", float64(m2.TotalAlloc-m1.TotalAlloc)/(1024*1024))
	fmt.Printf("  HeapAlloc:  +%.2f MB\n", float64(m2.HeapAlloc-m1.HeapAlloc)/(1024*1024))
	
	fmt.Println("\n\n=== Releasing References ===")
	fmt.Println("Setting slice to nil to make objects unreachable...")
	
	objects = nil
	
	fmt.Println("\nManually triggering garbage collection...")
	fmt.Println("(In practice, Go's GC runs automatically based on heap size)")
	runtime.GC()
	
	runtime.GC()
	
	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)
	fmt.Println("\n\n=== Memory Stats After GC ===")
	printMemStats(&m3)
	
	fmt.Printf("\nMemory after GC (vs after allocation):\n")
	fmt.Printf("  Alloc:     %.2f MB (was %.2f MB, freed %.2f MB)\n", 
		float64(m3.Alloc)/(1024*1024),
		float64(m2.Alloc)/(1024*1024),
		float64(m2.Alloc-m3.Alloc)/(1024*1024))
	fmt.Printf("  HeapAlloc: %.2f MB (was %.2f MB, freed %.2f MB)\n",
		float64(m3.HeapAlloc)/(1024*1024),
		float64(m2.HeapAlloc)/(1024*1024),
		float64(m2.HeapAlloc-m3.HeapAlloc)/(1024*1024))
	
	fmt.Println("\n\n=== Summary ===")
	fmt.Println("Go's garbage collector automatically reclaimed the memory")
	fmt.Println("once we removed all references to the allocated objects.")
	fmt.Println("No manual free() calls needed like in C!")
}

func printMemStats(m *runtime.MemStats) {
	fmt.Printf("  Alloc:       %10.2f MB (bytes allocated and still in use)\n", 
		float64(m.Alloc)/(1024*1024))
	fmt.Printf("  TotalAlloc:  %10.2f MB (total bytes allocated over time)\n", 
		float64(m.TotalAlloc)/(1024*1024))
	fmt.Printf("  Sys:         %10.2f MB (bytes obtained from OS)\n", 
		float64(m.Sys)/(1024*1024))
	fmt.Printf("  HeapAlloc:   %10.2f MB (bytes allocated on heap)\n", 
		float64(m.HeapAlloc)/(1024*1024))
	fmt.Printf("  HeapSys:     %10.2f MB (bytes obtained from OS for heap)\n", 
		float64(m.HeapSys)/(1024*1024))
	fmt.Printf("  NumGC:       %10d (number of completed GC cycles)\n", m.NumGC)
}
