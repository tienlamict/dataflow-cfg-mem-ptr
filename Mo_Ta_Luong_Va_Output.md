# MÔ TẢ LUỒNG VÀ OUTPUT CỦA PROJECT

## TỔNG QUAN PROJECT

Project "dataflow-cfg-mem-ptr" là một công cụ CLI (Command Line Interface) được viết bằng Go, bao gồm 4 module độc lập minh họa các kỹ thuật phân tích chương trình và quản lý bộ nhớ.

### Cấu trúc chương trình chính

```
main.go (hoặc cmd/tool/main.go)
  |
  +-- Nhận subcommand từ command line
  |
  +-- Phân phối đến module tương ứng:
      |
      +-- dataflow: internal/dataflow/reaching.go
      +-- cfg:      internal/cfg/cfg.go
      +-- mem:      internal/mem/mem.go
      +-- ptr:      internal/ptr/ptr.go
```

### Cách chạy chương trình

```bash
go run main.go <subcommand>
```

Trong đó `<subcommand>` có thể là:
- `dataflow` - Chạy phân tích reaching definitions
- `cfg` - Xây dựng và hiển thị Control Flow Graph
- `mem` - Demo quản lý bộ nhớ và garbage collection
- `ptr` - Demo con trỏ và GC với linked list

---

## MODULE 1: DATAFLOW ANALYSIS

### Luồng xử lý

**Bước 1: Khởi tạo CFG**
```
- Tạo 3 basic blocks với các predecessors và successors
- Block 1 -> Block 2 -> Block 3 -> Block 2 (có vòng lặp)
```

**Bước 2: Thiết lập các câu lệnh cho mỗi block**
```
Block 1: d1: x = 5, d2: y = 3
Block 2: d3: x = x + y
Block 3: d4: y = x - 1
```

**Bước 3: Tính toán tập GEN và KILL**
```
Block 1:
  - GEN = {d1, d2} (sinh ra định nghĩa d1 cho x và d2 cho y)
  - KILL = {} (không kill định nghĩa nào)

Block 2:
  - GEN = {d3} (sinh ra định nghĩa mới d3 cho x)
  - KILL = {d1} (kill định nghĩa cũ d1 của x)

Block 3:
  - GEN = {d4} (sinh ra định nghĩa mới d4 cho y)
  - KILL = {d2} (kill định nghĩa cũ d2 của y)
```

**Bước 4: Fixpoint Iteration**
```
Khởi tạo: IN = {}, OUT = {} cho tất cả blocks

Lặp lại cho đến khi không thay đổi:
  Với mỗi block:
    1. Tính IN[n] = hợp của OUT từ tất cả predecessors
    2. Tính OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])
    3. Kiểm tra có thay đổi không
```

**Bước 5: Đạt fixpoint và kết thúc**

### Output đầu ra

```
=== Running Dataflow Analysis (Reaching Definitions) ===

Initial GEN/KILL sets:

Block 1:
  Statements: [d1: x = 5 d2: y = 3]
  GEN  = {d1, d2}
  KILL = {}

Block 2:
  Statements: [d3: x = x + y]
  GEN  = {d3}
  KILL = {d1}

Block 3:
  Statements: [d4: y = x - 1]
  GEN  = {d4}
  KILL = {d2}


=== Fixpoint Iteration ===
Using transfer functions:
  IN[n]  = ∪ OUT[pred] for all pred of n
  OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])

--- Iteration 1 ---
Block 1: IN={}, OUT={d1, d2}
Block 2: IN={d1, d2}, OUT={d2, d3}
Block 3: IN={d2, d3}, OUT={d3, d4}

--- Iteration 2 ---
Block 1: IN={}, OUT={d1, d2}
Block 2: IN={d1, d2, d3, d4}, OUT={d2, d3, d4}
Block 3: IN={d2, d3, d4}, OUT={d3, d4}

Fixpoint reached after 2 iterations.


=== Final Results ===

Block 1:
  Statements: [d1: x = 5 d2: y = 3]
  IN  = {}
  OUT = {d1, d2}

Block 2:
  Statements: [d3: x = x + y]
  IN  = {d1, d2, d3, d4}
  OUT = {d2, d3, d4}

Block 3:
  Statements: [d4: y = x - 1]
  IN  = {d2, d3, d4}
  OUT = {d3, d4}
```

### Giải thích kết quả

**Iteration 1:**
- Block 1 không có predecessor nên IN rỗng, OUT = GEN = {d1, d2}
- Block 2 nhận IN từ Block 1 = {d1, d2}, OUT = {d3} ∪ ({d1, d2} - {d1}) = {d2, d3}
- Block 3 nhận IN từ Block 2 = {d2, d3}, OUT = {d4} ∪ ({d2, d3} - {d2}) = {d3, d4}

**Iteration 2:**
- Block 1 không đổi
- Block 2 giờ nhận từ cả Block 1 và Block 3 (vòng lặp), IN = {d1, d2} ∪ {d3, d4} = {d1, d2, d3, d4}
- OUT của Block 2 = {d3} ∪ ({d1, d2, d3, d4} - {d1}) = {d2, d3, d4}
- Block 3 cập nhật IN = {d2, d3, d4}, OUT = {d3, d4}

Không có thay đổi tiếp theo, đạt fixpoint.

---

## MODULE 2: CONTROL FLOW GRAPH (CFG)

### Luồng xử lý

**Bước 1: Định nghĩa pseudo-code**
```
Mã giả được định nghĩa dưới dạng slice of strings:
  1: x = 0
  2: y = 10
  3: if x >= y goto 8
  4: sum = 0
  5: sum = sum + x
  6: x = x + 1
  7: goto 3
  8: print sum
  9: return
```

**Bước 2: Xác định Leaders**
```
Áp dụng 3 quy tắc:
  1. Câu lệnh đầu tiên (line 1) là leader
  2. Đích của jump (line 3, 8) là leaders
  3. Câu lệnh sau jump (line 4, 8) là leaders

Kết quả: Leaders = {1, 3, 4, 8}
```

**Bước 3: Phân chia Basic Blocks**
```
Với mỗi leader, tạo một block mới:
  Block 1: từ line 1 đến trước line 3
  Block 2: từ line 3 đến trước line 4
  Block 3: từ line 4 đến trước line 8
  Block 4: từ line 8 đến hết
```

**Bước 4: Thêm cạnh (Edges)**
```
Phân tích câu lệnh cuối của mỗi block:
  Block 1 (line 2): không phải jump -> fall-through -> Block 2
  Block 2 (line 3): conditional jump -> 2 successors:
    - Nhảy đến line 8 -> Block 4
    - Fall-through -> Block 3
  Block 3 (line 7): unconditional jump -> Block 2
  Block 4 (line 9): return -> EXIT
```

**Bước 5: Hiển thị kết quả**

### Output đầu ra

```
=== Running CFG Construction ===

Input Pseudo-Code:
  1: x = 0
  2: y = 10
  3: if x >= y goto 8
  4: sum = 0
  5: sum = sum + x
  6: x = x + 1
  7: goto 3
  8: print sum
  9: return


=== Basic Block Identification ===
Leader rules:
  1. First statement is a leader
  2. Target of a jump is a leader
  3. Statement following a jump is a leader


=== Control Flow Graph ===

Block 1 (starts at line 1):
  Statements:
    1: x = 0
    2: y = 10
  Successors: [2]

Block 2 (starts at line 3):
  Statements:
    3: if x >= y goto 8
  Successors: [4 3]

Block 3 (starts at line 4):
  Statements:
    4: sum = 0
    5: sum = sum + x
    6: x = x + 1
    7: goto 3
  Successors: [2]

Block 4 (starts at line 8):
  Statements:
    8: print sum
    9: return
  Successors: [EXIT]


=== Graph Edges ===
Block 1 -> Block 2
Block 2 -> Block 4
Block 2 -> Block 3
Block 3 -> Block 2
Block 4 -> EXIT
```

### Giải thích cấu trúc CFG

Đồ thị này biểu diễn một vòng lặp while:

```
       +-------+
       | Block 1|
       | x = 0  |
       | y = 10 |
       +-------+
           |
           v
    +------+------+
    |   Block 2   |
    | if x >= y   |
    | goto 8      |
    +------+------+
        |      |
    True|      |False
        |      v
        |  +--------+
        |  | Block 3|
        |  | sum=0  |
        |  | sum+=x |
        |  | x++    |
        |  | goto 3 |
        |  +--------+
        |      |
        |      |
        +------+
           |
           v
       +-------+
       | Block 4|
       | print  |
       | return |
       +-------+
```

---

## MODULE 3: MEMORY MANAGEMENT

### Luồng xử lý

**Bước 1: Khởi tạo và đo baseline**
```
1. Gọi runtime.GC() hai lần để cleanup
2. Đọc runtime.MemStats vào m1
3. In thống kê bộ nhớ ban đầu
```

**Bước 2: Cấp phát bộ nhớ**
```
1. Tạo slice chứa 1000 con trỏ
2. Với mỗi phần tử:
   - Cấp phát array[1024]int (8KB mỗi array)
   - Khởi tạo dữ liệu trong array
3. Tổng cộng: 1000 * 8KB = ~8MB
```

**Bước 3: Đo bộ nhớ sau cấp phát**
```
1. Đọc runtime.MemStats vào m2
2. In thống kê
3. Tính và hiển thị mức tăng bộ nhớ
```

**Bước 4: Giải phóng references**
```
1. Set objects = nil
2. Tất cả 1000 arrays trở thành unreachable
3. In thông báo
```

**Bước 5: Kích hoạt GC và đo lại**
```
1. Gọi runtime.GC()
2. Đọc runtime.MemStats vào m3
3. In thống kê và bộ nhớ đã thu hồi
```

### Output đầu ra

```
=== Running Memory Allocation Demo ===

This demo shows Go's automatic memory management.
Unlike C/C++ where you manually malloc/free, Go uses GC.

=== Initial Memory Stats ===
  Alloc:            0.23 MB (bytes allocated and still in use)
  TotalAlloc:       0.23 MB (total bytes allocated over time)
  Sys:              8.45 MB (bytes obtained from OS)
  HeapAlloc:        0.23 MB (bytes allocated on heap)
  HeapSys:          4.19 MB (bytes obtained from OS for heap)
  NumGC:            2 (number of completed GC cycles)


=== Allocating Memory ===
Allocating 1000 objects of 1024 integers each...
Allocated 1000 objects (approximately 7.81 MB)

=== Memory Stats After Allocation ===
  Alloc:            8.05 MB (bytes allocated and still in use)
  TotalAlloc:       8.05 MB (total bytes allocated over time)
  Sys:             16.70 MB (bytes obtained from OS)
  HeapAlloc:        8.05 MB (bytes allocated on heap)
  HeapSys:         12.44 MB (bytes obtained from OS for heap)
  NumGC:            2 (number of completed GC cycles)

Memory increase:
  Alloc:      +7.82 MB
  TotalAlloc: +7.82 MB
  HeapAlloc:  +7.82 MB


=== Releasing References ===
Setting slice to nil to make objects unreachable...

Manually triggering garbage collection...
(In practice, Go's GC runs automatically based on heap size)

=== Memory Stats After GC ===
  Alloc:            0.26 MB (bytes allocated and still in use)
  TotalAlloc:       8.08 MB (total bytes allocated over time)
  Sys:             16.70 MB (bytes obtained from OS)
  HeapAlloc:        0.26 MB (bytes allocated on heap)
  HeapSys:         12.44 MB (bytes obtained from OS for heap)
  NumGC:            4 (number of completed GC cycles)

Memory after GC (vs after allocation):
  Alloc:     0.26 MB (was 8.05 MB, freed 7.79 MB)
  HeapAlloc: 0.26 MB (was 8.05 MB, freed 7.79 MB)


=== Summary ===
Go's garbage collector automatically reclaimed the memory
once we removed all references to the allocated objects.
No manual free() calls needed like in C!
```

### Giải thích metrics

**Alloc**: Bộ nhớ đang được sử dụng (tăng từ 0.23MB lên 8.05MB, sau GC giảm về 0.26MB)

**TotalAlloc**: Tổng bộ nhớ đã cấp phát (chỉ tăng, không giảm)

**HeapAlloc**: Tương tự Alloc nhưng chỉ tính heap

**NumGC**: Số lần GC chạy (tăng từ 2 lên 4)

**Sys, HeapSys**: Bộ nhớ OS cấp cho Go runtime (không trả lại ngay cho OS)

---

## MODULE 4: POINTER ANALYSIS

### Luồng xử lý

**Bước 1: Giới thiệu và so sánh với C**
```
1. In thông báo demo
2. Hiển thị code C vs Go để so sánh
3. Giải thích sự khác biệt
```

**Bước 2: Tạo linked list**
```
1. Gọi createLinkedList(10)
2. Hàm này:
   - Tạo node đầu tiên (head)
   - Lặp tạo 9 nodes còn lại
   - Liên kết mỗi node với node tiếp theo
3. Trả về con trỏ head
```

**Bước 3: Hiển thị linked list**
```
1. In 5 nodes đầu tiên
2. Với mỗi node hiển thị:
   - Số thứ tự
   - Giá trị (Value)
   - Địa chỉ bộ nhớ (Addr)
   - Địa chỉ node tiếp theo (Next)
3. In thông báo còn n nodes nữa
```

**Bước 4: Đo bộ nhớ trước giải phóng**
```
1. Gọi runtime.GC() hai lần
2. Đọc MemStats vào m1
3. In thống kê bộ nhớ
4. Đếm số nodes để verify
```

**Bước 5: Giải phóng linked list**
```
1. In thông báo sẽ set head = nil
2. Giải thích ý nghĩa (tất cả nodes unreachable)
3. So sánh với C (sẽ là memory leak)
4. Set head = nil
```

**Bước 6: Kích hoạt GC**
```
1. In thông báo triggering GC
2. Gọi runtime.GC() hai lần
3. Đọc MemStats vào m2
4. In thống kê và bộ nhớ thu hồi
```

**Bước 7: Tóm tắt**
```
In summary với các điểm chính:
- Đã tạo linked list với dynamic allocation
- Đã drop references (head = nil)
- GC tự động freed tất cả nodes
- Không cần free() như trong C
```

### Output đầu ra

```
=== Running Pointer and GC Demo ===

This demo shows pointer usage and garbage collection in Go.

=== Comparison with C ===
In C:
  struct Node* node = (struct Node*)malloc(sizeof(struct Node));
  // ... use node ...
  free(node);  // Manual memory management!

In Go:
  node := &Node{Value: 42}
  // ... use node ...
  // No free() needed - GC handles it automatically!

=== Creating Linked List (10 nodes) ===

Linked list created. First 5 nodes:
  Node 1: Value = 1, Addr = 0xc00001e0c0, Next = 0xc00001e0d8
  Node 2: Value = 2, Addr = 0xc00001e0d8, Next = 0xc00001e0f0
  Node 3: Value = 3, Addr = 0xc00001e0f0, Next = 0xc00001e108
  Node 4: Value = 4, Addr = 0xc00001e108, Next = 0xc00001e120
  Node 5: Value = 5, Addr = 0xc00001e120, Next = 0xc00001e138
  ... (5 more nodes)


=== Memory Stats (with linked list) ===
  Alloc:     0.18 KB
  HeapAlloc: 0.18 KB
  NumGC:     2

Total nodes in list: 10


=== Dropping Reference to List ===
Setting head = nil makes all nodes unreachable...
(In C, this would be a memory leak without manual free!)

Triggering garbage collection...


=== Memory Stats (after GC) ===
  Alloc:     0.16 KB
  HeapAlloc: 0.16 KB
  NumGC:     4

Memory reclaimed: 0.02 KB


=== Summary ===
Created linked list with dynamic allocation (like malloc in C)
Dropped all references to the list (head = nil)
GC automatically freed all unreachable nodes
No manual free() needed - Go manages memory automatically!

Key advantage: No memory leaks, no use-after-free bugs!
```

### Giải thích cấu trúc bộ nhớ

**Trước khi head = nil:**
```
[Stack]
  head (0xc00001e0c0)
    |
    v
[Heap]
  Node1 (0xc00001e0c0) -> Next (0xc00001e0d8)
    |
    v
  Node2 (0xc00001e0d8) -> Next (0xc00001e0f0)
    |
    v
  Node3 (0xc00001e0f0) -> Next (0xc00001e108)
    |
  ...
    |
    v
  Node10 -> Next (nil)
```

Tất cả nodes đều reachable từ head.

**Sau khi head = nil:**
```
[Stack]
  head = nil (không trỏ đến gì)

[Heap - Unreachable]
  Node1 -> Node2 -> Node3 -> ... -> Node10
  (Tất cả đều không thể truy cập)
```

GC phát hiện tất cả nodes không reachable và thu hồi bộ nhớ.

---

## SO SÁNH VỚI C

### Memory Management trong C vs Go

**Trong C:**
```c
// Phải malloc từng node
struct Node* node = malloc(sizeof(struct Node));
if (node == NULL) {
    // Handle error
}

// Phải free từng node
struct Node* current = head;
while (current != NULL) {
    struct Node* temp = current;
    current = current->next;
    free(temp);
}
```

**Trong Go:**
```go
// Cấp phát đơn giản
node := &Node{Value: 42}

// Giải phóng tự động
head = nil  // GC tự động free tất cả!
```

### Lợi ích của GC

**An toàn:**
- Không có memory leak
- Không có use-after-free
- Không có double-free
- Không có dangling pointers

**Năng suất:**
- Code ngắn gọn hơn
- Ít bugs hơn
- Developer focus vào logic

**Trade-off:**
- GC pause (có thể ảnh hưởng real-time systems)
- Memory overhead
- Khó dự đoán timing

---

## CÁCH SỬ DỤNG PROJECT

### Chạy tất cả demos

```bash
# Demo 1: Dataflow Analysis
go run main.go dataflow

# Demo 2: CFG Construction
go run main.go cfg

# Demo 3: Memory Management
go run main.go mem

# Demo 4: Pointer Analysis
go run main.go ptr
```

### Xem help

```bash
# Help chung
go run main.go

# Help cho subcommand cụ thể
go run main.go dataflow -help
go run main.go cfg -help
go run main.go mem -help
go run main.go ptr -help
```

---

## KẾT LUẬN

Project này minh họa 4 khía cạnh quan trọng trong program analysis:

1. **Dataflow Analysis**: Phân tích cách data flow qua chương trình
2. **CFG**: Biểu diễn luồng điều khiển
3. **Memory Management**: Quản lý bộ nhớ tự động vs thủ công
4. **Pointer Analysis**: Cấu trúc dữ liệu động và GC

Mỗi module đều có output rõ ràng, giúp hiểu sâu về cách hoạt động của từng kỹ thuật.

