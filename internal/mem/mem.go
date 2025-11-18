package mem

import (
	"fmt"
	"runtime"
)

// Package mem demonstrates Go's memory management and garbage collection.
// Package này minh họa quản lý bộ nhớ và garbage collection của Go.
//
// HOMEWORK YÊU CẦU 3: Minh họa các ý thuật cấp phát bộ nhớ và giải phóng bộ nhớ
//
// Nội dung demo:
// - Cấp phát bộ nhớ trên heap (Heap allocation)
// - Thống kê bộ nhớ (runtime.MemStats)
// - Kích hoạt garbage collection thủ công
// - Thu hồi bộ nhớ sau khi giải phóng tham chiếu
//
// So sánh với C/C++:
// - C/C++: Cần malloc()/free() thủ công -> Dễ bị memory leak
// - Go: Garbage Collector tự động thu hồi bộ nhớ không còn dùng

// RunMemoryDemo - Demo cấp phát bộ nhớ và garbage collection
//
// Demo này minh họa:
// 1. Cấp phát nhiều object lớn trên heap
// 2. Đo lường bộ nhớ sử dụng qua runtime.MemStats
// 3. Giải phóng tham chiếu (set = nil)
// 4. Kích hoạt GC và xem bộ nhớ được thu hồi
func RunMemoryDemo() {
	fmt.Println("This demo shows Go's automatic memory management.")
	fmt.Println("Unlike C/C++ where you manually malloc/free, Go uses GC.")
	fmt.Println()
	
	// Buộc GC chạy trước để có baseline sạch
	runtime.GC()
	runtime.GC() // Gọi 2 lần để đảm bảo cleanup hoàn toàn
	
	// In thống kê bộ nhớ ban đầu
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Println("=== Initial Memory Stats ===")
	printMemStats(&m1)
	
	// Cấp phát nhiều object lớn
	fmt.Println("\n\n=== Allocating Memory ===")
	const numObjects = 1000      // Số lượng object
	const objectSize = 1024       // Kích thước mỗi array (1024 integers)
	
	fmt.Printf("Allocating %d objects of %d integers each...\n", numObjects, objectSize)
	
	// Slice để giữ tham chiếu đến các object đã cấp phát
	// Quan trọng: Phải giữ tham chiếu thì object mới không bị GC thu hồi
	objects := make([]*[1024]int, numObjects)
	
	for i := 0; i < numObjects; i++ {
		// Cấp phát trên heap (dùng new hoặc make)
		// Mỗi array có 1024 integers = 1024 * 8 bytes = 8KB
		objects[i] = new([1024]int)
		
		// Khởi tạo data để đảm bảo bộ nhớ thực sự được dùng
		for j := 0; j < objectSize; j++ {
			objects[i][j] = i + j
		}
	}
	
	fmt.Printf("Allocated %d objects (approximately %.2f MB)\n", 
		numObjects, float64(numObjects*objectSize*8)/(1024*1024))
	
	// Print memory stats after allocation
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)
	fmt.Println("\n\n=== Memory Stats After Allocation ===")
	printMemStats(&m2)
	
	fmt.Printf("\nMemory increase:\n")
	fmt.Printf("  Alloc:      +%.2f MB\n", float64(m2.Alloc-m1.Alloc)/(1024*1024))
	fmt.Printf("  TotalAlloc: +%.2f MB\n", float64(m2.TotalAlloc-m1.TotalAlloc)/(1024*1024))
	fmt.Printf("  HeapAlloc:  +%.2f MB\n", float64(m2.HeapAlloc-m1.HeapAlloc)/(1024*1024))
	
	// Giải phóng tham chiếu (làm cho objects không còn truy cập được)
	fmt.Println("\n\n=== Releasing References ===")
	fmt.Println("Setting slice to nil to make objects unreachable...")
	
	// QUAN TRỌNG: Set objects = nil để loại bỏ tham chiếu
	// Khi không còn tham chiếu nào đến object, GC sẽ thu hồi bộ nhớ
	// Trong C/C++: Đây sẽ gây memory leak nếu không free() trước!
	objects = nil
	
	// Kích hoạt garbage collection thủ công
	fmt.Println("\nManually triggering garbage collection...")
	fmt.Println("(In practice, Go's GC runs automatically based on heap size)")
	runtime.GC()
	
	// Chờ GC hoàn tất
	runtime.GC()
	
	// Print memory stats after GC
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

// printMemStats - In các thống kê bộ nhớ quan trọng
// 
// Các số liệu quan trọng:
// - Alloc: Bộ nhớ đang được sử dụng
// - TotalAlloc: Tổng bộ nhớ đã cấp phát từ đầu chương trình
// - Sys: Bộ nhớ lấy từ hệ điều hành
// - HeapAlloc: Bộ nhớ trên heap đang dùng
// - NumGC: Số lần GC đã chạy
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

