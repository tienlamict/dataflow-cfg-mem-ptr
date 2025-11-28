# BÀI LUẬN THUYẾT TRÌNH
## Phân Tích Chương Trình: Dataflow Analysis, CFG, Memory Management và Pointer Analysis

### Giới thiệu tổng quan

Project **"dataflow-cfg-mem-ptr"** là một bộ công cụ giáo dục được viết bằng Go, minh họa các khái niệm cơ bản trong phân tích chương trình (Program Analysis) và quản lý bộ nhớ. Project bao gồm 4 bài tập chính, mỗi bài tập tập trung vào một khía cạnh quan trọng của phân tích và thực thi chương trình.

Các module trong project:
1. **Dataflow Analysis** - Phân tích luồng dữ liệu với reaching definitions
2. **Control Flow Graph (CFG)** - Xây dựng đồ thị luồng điều khiển
3. **Memory Management** - Quản lý bộ nhớ và garbage collection
4. **Pointer Analysis** - Phân tích con trỏ và cấu trúc dữ liệu động

---

## BÀI TẬP 1: DATAFLOW ANALYSIS - Phân Tích Reaching Definitions

### 1.1. Mục đích và ý nghĩa

Dataflow analysis là một kỹ thuật cơ bản trong biên dịch và phân tích chương trình, được sử dụng để thu thập thông tin về cách dữ liệu di chuyển qua các điểm khác nhau trong chương trình. Bài tập này minh họa **Reaching Definitions Analysis** - một dạng phân tích luồng dữ liệu tiến (forward dataflow analysis).

**Reaching Definitions** trả lời câu hỏi: "Tại một điểm trong chương trình, những định nghĩa (gán giá trị) nào của biến có thể 'đạt được' điểm đó mà không bị gán đè lại?"

### 1.2. Các khái niệm chính

#### 1.2.1. Basic Block và CFG
- **Basic Block**: Một chuỗi các câu lệnh liên tiếp mà chỉ có thể vào qua câu lệnh đầu tiên và ra qua câu lệnh cuối cùng
- **CFG (Control Flow Graph)**: Đồ thị biểu diễn luồng điều khiển giữa các basic blocks

#### 1.2.2. Các tập hợp GEN, KILL, IN, OUT

**GEN (Generate)**: Tập các định nghĩa được tạo ra trong block
- Ví dụ: Block có `d1: x = 5` thì GEN chứa `d1`

**KILL (Kill)**: Tập các định nghĩa bị loại bỏ bởi block
- Ví dụ: Nếu block có `d3: x = x + y` thì KILL chứa tất cả định nghĩa cũ của biến `x` (như `d1`)

**IN**: Tập các định nghĩa đạt được điểm vào của block
- Công thức: `IN[n] = ∪ OUT[pred]` (hợp của OUT từ tất cả predecessor)

**OUT**: Tập các định nghĩa đạt được điểm ra của block
- Công thức: `OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])`

### 1.3. Thuật toán Fixpoint Iteration

Thuật toán lặp đến điểm bất động (fixpoint) để tính toán IN/OUT:

```
1. Khởi tạo: Đặt tất cả IN, OUT = rỗng
2. Lặp cho đến khi không có thay đổi:
   - Với mỗi block n:
     a. Tính IN[n] = ∪ OUT[pred] cho tất cả predecessor
     b. Tính OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])
3. Khi không còn thay đổi → đạt fixpoint
```

### 1.4. Ví dụ cụ thể trong code

Chương trình demo một CFG với 3 blocks và có vòng lặp:

```
Block 1: d1: x = 5, d2: y = 3
Block 2: d3: x = x + y
Block 3: d4: y = x - 1
```

**Cấu trúc CFG**: Block 1 → Block 2 → Block 3 → Block 2 (vòng lặp)

**Tính toán GEN/KILL**:
- Block 1: GEN = {d1, d2}, KILL = {}
- Block 2: GEN = {d3}, KILL = {d1} (vì d3 định nghĩa lại x)
- Block 3: GEN = {d4}, KILL = {d2} (vì d4 định nghĩa lại y)

**Quá trình lặp**:

*Iteration 1*:
- Block 1: IN={}, OUT={d1, d2}
- Block 2: IN={d1, d2}, OUT={d2, d3}
- Block 3: IN={d2, d3}, OUT={d3, d4}

*Iteration 2*:
- Block 1: IN={}, OUT={d1, d2}
- Block 2: IN={d1, d2, d3, d4}, OUT={d2, d3, d4}
- Block 3: IN={d2, d3, d4}, OUT={d3, d4}

Sau 2 iterations, thuật toán đạt fixpoint vì không còn thay đổi.

### 1.5. Ứng dụng thực tế

- **Tối ưu hóa compiler**: Loại bỏ mã chết (dead code elimination)
- **Phát hiện biến chưa khởi tạo**: Kiểm tra xem biến có định nghĩa trước khi sử dụng không
- **Constant propagation**: Lan truyền giá trị hằng số
- **Register allocation**: Phân bổ thanh ghi hiệu quả

---

## BÀI TẬP 2: CONTROL FLOW GRAPH (CFG) - Xây Dựng Đồ Thị Luồng Điều Khiển

### 2.1. Mục đích và ý nghĩa

Control Flow Graph là một công cụ cơ bản trong phân tích chương trình, biểu diễn tất cả các đường đi có thể có mà luồng điều khiển có thể đi qua trong một chương trình. CFG là nền tảng cho nhiều kỹ thuật phân tích và tối ưu hóa khác.

### 2.2. Thuật toán xây dựng CFG

#### 2.2.1. Bước 1: Xác định Leaders

**Leaders** là các câu lệnh đánh dấu điểm bắt đầu của basic blocks, được xác định theo 3 quy tắc:

1. **Quy tắc 1**: Câu lệnh đầu tiên của chương trình là một leader
2. **Quy tắc 2**: Đích của lệnh nhảy (jump target) là một leader
3. **Quy tắc 3**: Câu lệnh ngay sau lệnh nhảy là một leader

#### 2.2.2. Bước 2: Phân chia Basic Blocks

Sau khi xác định leaders:
- Mỗi leader bắt đầu một basic block mới
- Block kéo dài từ leader cho đến (nhưng không bao gồm) leader tiếp theo
- Block cuối cùng kết thúc tại cuối chương trình

#### 2.2.3. Bước 3: Thêm các cạnh (Edges)

Cạnh trong CFG biểu diễn luồng điều khiển có thể:

**Unconditional Jump** (`goto`):
- Chỉ có 1 successor: block chứa đích nhảy

**Conditional Jump** (`if...goto`):
- Có 2 successors:
  1. Block chứa đích nhảy (nếu điều kiện đúng)
  2. Block tiếp theo (nếu điều kiện sai - fall-through)

**Không phải Jump**:
- Successor: Block tiếp theo (fall-through)

**Return statement**:
- Không có successor (thoát chương trình)

### 2.3. Ví dụ cụ thể

Pseudo-code mô phỏng vòng lặp:

```
1: x = 0              // Khởi tạo
2: y = 10             // Khởi tạo
3: if x >= y goto 8   // Kiểm tra điều kiện
4: sum = 0            // Vào thân vòng lặp
5: sum = sum + x      // Tính toán
6: x = x + 1          // Tăng counter
7: goto 3             // Quay lại kiểm tra
8: print sum          // Sau vòng lặp
9: return             // Kết thúc
```

**Xác định Leaders**:
- Line 1: Leader (quy tắc 1 - câu lệnh đầu tiên)
- Line 3: Leader (quy tắc 3 - sau statement 2 có thể rơi xuống)
- Line 4: Leader (quy tắc 2 - sau conditional jump; quy tắc 3 - đích của jump từ 7)
- Line 8: Leader (quy tắc 2 - đích của conditional jump từ 3)

**Kết quả CFG**:

**Block 1** (lines 1-2):
- Statements: `x = 0`, `y = 10`
- Successors: [Block 2]

**Block 2** (line 3):
- Statement: `if x >= y goto 8`
- Successors: [Block 4, Block 3] (nhảy đến 8 hoặc rơi xuống 4)

**Block 3** (lines 4-6):
- Statements: `sum = 0`, `sum = sum + x`, `x = x + 1`
- Followed by: `goto 3`
- Successors: [Block 2] (quay lại kiểm tra điều kiện)

**Block 4** (lines 8-9):
- Statements: `print sum`, `return`
- Successors: [EXIT]

### 2.4. Cấu trúc dữ liệu trong code

```go
type Block struct {
    ID    int        // ID của block
    Label int        // Số dòng của câu lệnh đầu tiên
    Stmts []string   // Danh sách câu lệnh
    Succs []int      // Danh sách ID của blocks kế tiếp
}

type CFG struct {
    Blocks []*Block  // Danh sách tất cả các blocks
}
```

### 2.5. Ứng dụng

- **Dataflow Analysis**: CFG là input cho các phân tích dataflow
- **Optimization**: Dead code elimination, loop optimization
- **Testing**: Test path coverage, branch coverage
- **Static Analysis**: Phát hiện lỗi, security vulnerabilities
- **Program Understanding**: Hiểu cấu trúc và logic chương trình

---

## BÀI TẬP 3: MEMORY MANAGEMENT - Quản Lý Bộ Nhớ và Garbage Collection

### 3.1. Mục đích và ý nghĩa

Quản lý bộ nhớ là một khía cạnh quan trọng trong lập trình hệ thống. Bài tập này so sánh cách quản lý bộ nhớ thủ công (như trong C/C++) với quản lý bộ nhớ tự động thông qua Garbage Collection (như trong Go, Java).

### 3.2. Stack vs Heap

#### Stack Allocation
- **Cấp phát nhanh**: Chỉ cần di chuyển stack pointer
- **Tự động giải phóng**: Khi function return
- **Kích thước cố định**: Phải biết trước kích thước
- **Lifetime ngắn**: Chỉ tồn tại trong scope

#### Heap Allocation
- **Cấp phát chậm hơn**: Cần tìm không gian phù hợp
- **Quản lý phức tạp**: Cần theo dõi khi nào giải phóng
- **Kích thước động**: Có thể cấp phát bất kỳ kích thước nào
- **Lifetime dài**: Tồn tại cho đến khi được giải phóng

### 3.3. Manual Memory Management (C/C++)

#### Ưu điểm:
- **Kiểm soát tuyệt đối**: Lập trình viên quyết định khi nào cấp phát/giải phóng
- **Hiệu năng tối ưu**: Không có overhead của GC
- **Dự đoán được**: Không có pause do GC

#### Nhược điểm:
- **Memory Leaks**: Quên gọi `free()` → rò rỉ bộ nhớ
- **Dangling Pointers**: Truy cập vào bộ nhớ đã `free()` → undefined behavior
- **Double Free**: Gọi `free()` 2 lần → crash
- **Use-After-Free**: Sử dụng sau khi đã giải phóng → security bugs

Ví dụ trong C:
```c
// Cấp phát
int* arr = (int*)malloc(1000 * sizeof(int));
if (arr == NULL) { /* handle error */ }

// Sử dụng
for (int i = 0; i < 1000; i++) {
    arr[i] = i;
}

// PHẢI nhớ giải phóng
free(arr);  // Nếu quên → memory leak!
```

### 3.4. Automatic Memory Management (Go)

Go sử dụng **Garbage Collection** để tự động thu hồi bộ nhớ.

#### Nguyên lý hoạt động:

**Reachability-based GC**:
1. **Mark Phase**: Đánh dấu tất cả objects có thể truy cập từ root
2. **Sweep Phase**: Thu hồi tất cả objects không được đánh dấu

**Roots** bao gồm:
- Biến global
- Biến local trên stack
- Registers

### 3.5. Demo trong code

Chương trình demo các bước:

#### Bước 1: Baseline
```go
runtime.GC()  // Chạy GC để có baseline sạch
var m1 runtime.MemStats
runtime.ReadMemStats(&m1)
```

#### Bước 2: Cấp phát 1000 objects lớn
```go
objects := make([]*[1024]int, 1000)
for i := 0; i < 1000; i++ {
    objects[i] = new([1024]int)  // Mỗi object ~8KB
    // Khởi tạo data
}
// Tổng cộng ~8MB
```

#### Bước 3: Đo bộ nhớ sau cấp phát
```go
var m2 runtime.MemStats
runtime.ReadMemStats(&m2)
// m2.HeapAlloc - m1.HeapAlloc ≈ 8MB
```

#### Bước 4: Giải phóng tham chiếu
```go
objects = nil  // Làm tất cả objects không còn reachable
// Trong C: Đây sẽ gây memory leak nếu không free()!
```

#### Bước 5: Kích hoạt GC và đo lại
```go
runtime.GC()
var m3 runtime.MemStats
runtime.ReadMemStats(&m3)
// m3.HeapAlloc << m2.HeapAlloc (bộ nhớ đã được thu hồi!)
```

### 3.6. Các metrics quan trọng

**runtime.MemStats** cung cấp nhiều thông tin:

- **Alloc**: Bytes hiện đang được sử dụng
- **TotalAlloc**: Tổng bytes đã cấp phát từ đầu chương trình
- **Sys**: Bytes lấy từ OS
- **HeapAlloc**: Bytes được cấp phát trên heap
- **HeapSys**: Bytes heap lấy từ OS
- **NumGC**: Số lần GC đã chạy

### 3.7. So sánh C vs Go

| Khía cạnh | C | Go |
|-----------|---|-----|
| Cấp phát | `malloc()` | `new()`, `make()` |
| Giải phóng | `free()` (thủ công) | GC (tự động) |
| Memory leak | Rất dễ xảy ra | Hiếm khi xảy ra |
| Use-after-free | Có thể xảy ra | Không xảy ra |
| Performance | Dự đoán được | Có GC pause |
| Complexity | Cao (phải quản lý) | Thấp (GC tự động) |

### 3.8. Trade-offs của GC

#### Ưu điểm:
- ✓ An toàn: Không có memory leak, use-after-free
- ✓ Đơn giản: Không cần quản lý thủ công
- ✓ Năng suất: Developer tập trung vào logic

#### Nhược điểm:
- ✗ GC pause: Chương trình tạm dừng khi GC chạy
- ✗ Memory overhead: GC cần metadata
- ✗ Không dự đoán: Khó biết khi nào GC chạy

---

## BÀI TẬP 4: POINTER ANALYSIS - Phân Tích Con Trỏ và Cấu Trúc Động

### 4.1. Mục đích và ý nghĩa

Pointer analysis nghiên cứu cách con trỏ được sử dụng để xây dựng các cấu trúc dữ liệu động (dynamic data structures) như linked list, tree, graph. Bài tập này minh họa cách Go quản lý bộ nhớ cho các cấu trúc này tự động, so với C yêu cầu quản lý thủ công.

### 4.2. Con trỏ và cấu trúc dữ liệu động

#### Linked List
Cấu trúc cơ bản:
```go
type Node struct {
    Value int    // Dữ liệu
    Next  *Node  // Con trỏ đến node tiếp theo
}
```

Mỗi node:
- Chứa data (Value)
- Chứa pointer đến node kế tiếp (Next)
- Node cuối có Next = nil

### 4.3. Cấp phát động trong C vs Go

#### Trong C:
```c
// Cấp phát mỗi node
struct Node* node1 = (struct Node*)malloc(sizeof(struct Node));
node1->value = 1;
node1->next = NULL;

struct Node* node2 = (struct Node*)malloc(sizeof(struct Node));
node2->value = 2;
node2->next = NULL;

node1->next = node2;  // Liên kết

// PHẢI nhớ giải phóng từng node
free(node2);
free(node1);
// Nếu quên → memory leak!
```

**Vấn đề**: Phải nhớ giải phóng tất cả nodes, đúng thứ tự!

#### Trong Go:
```go
// Cấp phát node rất đơn giản
node1 := &Node{Value: 1}
node2 := &Node{Value: 2}
node1.Next = node2

// Không cần free()!
// Khi không còn reference → GC tự động thu hồi
```

**Lợi ích**: Không cần lo lắng về giải phóng bộ nhớ!

### 4.4. Demo: Tạo và giải phóng Linked List

#### Bước 1: Tạo linked list với 10 nodes

```go
func createLinkedList(n int) *Node {
    if n <= 0 {
        return nil
    }
    
    head := &Node{Value: 1}  // Node đầu tiên
    current := head
    
    for i := 2; i <= n; i++ {
        newNode := &Node{Value: i}
        current.Next = newNode  // Liên kết
        current = newNode
    }
    
    return head
}
```

Kết quả: Một chuỗi liên kết
```
head -> Node(1) -> Node(2) -> Node(3) -> ... -> Node(10) -> nil
```

#### Bước 2: Hiển thị nodes và địa chỉ bộ nhớ

```
Node 1: Value = 1, Addr = 0xc00001e0c0, Next = 0xc00001e0d8
Node 2: Value = 2, Addr = 0xc00001e0d8, Next = 0xc00001e0f0
Node 3: Value = 3, Addr = 0xc00001e0f0, Next = 0xc00001e108
...
```

Mỗi node nằm ở địa chỉ khác nhau trên heap. Con trỏ Next trỏ đến địa chỉ của node kế tiếp.

#### Bước 3: Đo bộ nhớ trước giải phóng

```go
runtime.GC()
var m1 runtime.MemStats
runtime.ReadMemStats(&m1)
// m1.HeapAlloc chứa tất cả 10 nodes
```

#### Bước 4: Giải phóng tất cả nodes

**Trong C**: Phải viết vòng lặp để free từng node:
```c
struct Node* current = head;
while (current != NULL) {
    struct Node* temp = current;
    current = current->next;
    free(temp);  // Free từng node
}
head = NULL;
```

**Trong Go**: Chỉ cần 1 dòng!
```go
head = nil  // Tất cả nodes trở thành unreachable!
```

Khi `head = nil`:
- Không còn cách nào truy cập vào node đầu tiên
- Do đó không thể truy cập vào bất kỳ node nào
- Tất cả nodes trở thành **unreachable**
- GC sẽ tự động thu hồi TẤT CẢ các nodes!

#### Bước 5: Kích hoạt GC và kiểm tra

```go
runtime.GC()
var m2 runtime.MemStats
runtime.ReadMemStats(&m2)
// m2.HeapAlloc << m1.HeapAlloc
// Bộ nhớ của tất cả 10 nodes đã được thu hồi!
```

### 4.5. Reachability và Garbage Collection

#### Khái niệm Reachability

**Reachable**: Object có thể truy cập từ roots
- Global variables
- Stack variables (local variables)
- Registers

**Unreachable**: Object không thể truy cập từ bất kỳ root nào

#### GC Algorithm (Simplified)

1. **Mark Phase**:
   - Bắt đầu từ tất cả roots
   - Đi theo mọi pointer và đánh dấu (mark) mọi object gặp được
   - Sử dụng depth-first hoặc breadth-first traversal

2. **Sweep Phase**:
   - Quét qua toàn bộ heap
   - Thu hồi tất cả objects không được mark
   - Reset marks cho lần GC tiếp theo

#### Ví dụ với Linked List

**Trước khi `head = nil`**:
```
[Stack]
  head -> Node1 -> Node2 -> Node3 -> ... -> Node10
```
- `head` là root (stack variable)
- Node1 reachable từ `head`
- Node2 reachable từ Node1
- ... tất cả đều reachable

**Sau khi `head = nil`**:
```
[Stack]
  head = nil

[Heap - Unreachable]
  Node1 -> Node2 -> Node3 -> ... -> Node10
```
- `head` không trỏ đến gì
- Node1 không thể truy cập từ bất kỳ root nào
- Tất cả nodes khác cũng không thể truy cập
- GC mark phase: không node nào được mark
- GC sweep phase: thu hồi tất cả!

### 4.6. Ưu điểm của automatic memory management với pointers

#### Trong C - Dễ mắc lỗi:

**Lỗi 1: Memory Leak**
```c
Node* head = createList();
head = NULL;  // Forgot to free! → Memory leak
```

**Lỗi 2: Use-After-Free**
```c
Node* node = createNode();
free(node);
printf("%d\n", node->value);  // Undefined behavior!
```

**Lỗi 3: Double Free**
```c
free(node);
free(node);  // Crash!
```

**Lỗi 4: Dangling Pointer**
```c
Node* temp = head->next;
free(head);
temp->value = 5;  // temp points to freed memory!
```

#### Trong Go - An toàn:

```go
head := createList()
head = nil  // All nodes automatically freed - no leak!

node := &Node{Value: 42}
// Even if we lose reference, GC handles it
// No use-after-free possible
// No double-free possible
// No dangling pointers
```

### 4.7. Performance considerations

#### GC Overhead
- **Extra memory**: GC needs metadata
- **Pause time**: GC runs periodically and pauses program
- **CPU usage**: GC uses CPU cycles

#### When manual management might be better:
- Real-time systems (need predictable latency)
- Systems with very tight memory constraints
- High-performance gaming, HPC

#### When GC is better:
- Most application development
- Servers, web services
- Rapid prototyping
- Maintainability is important

### 4.8. Advanced pointer analysis

Bài tập này là nền tảng cho các kỹ thuật phân tích phức tạp hơn:

- **Alias Analysis**: Xác định khi nào 2 pointers trỏ đến cùng object
- **Shape Analysis**: Xác định cấu trúc (shape) của data structures
- **Escape Analysis**: Xác định object có "escape" khỏi scope không
- **Points-To Analysis**: Xác định pointer có thể trỏ đến objects nào

---

## KẾT LUẬN

### Tổng kết 4 bài tập

Project này cung cấp một cái nhìn toàn diện về program analysis và memory management:

#### 1. Dataflow Analysis
- **Học được**: Cách thông tin lan truyền qua chương trình
- **Ứng dụng**: Compiler optimization, static analysis
- **Kỹ thuật**: Fixpoint iteration, transfer functions

#### 2. Control Flow Graph
- **Học được**: Biểu diễn luồng điều khiển của chương trình
- **Ứng dụng**: Nền tảng cho dataflow analysis, testing, optimization
- **Kỹ thuật**: Leader identification, basic block splitting

#### 3. Memory Management
- **Học được**: So sánh manual vs automatic memory management
- **Trade-off**: Control vs safety, performance vs productivity
- **Kỹ thuật**: Heap allocation, garbage collection

#### 4. Pointer Analysis
- **Học được**: Cách con trỏ xây dựng cấu trúc động
- **Ưu điểm GC**: Không có memory leak, use-after-free bugs
- **Kỹ thuật**: Reachability analysis, automatic memory reclamation

### Mối liên hệ giữa các bài tập

```
CFG (Bài 2)
    ↓
Dataflow Analysis (Bài 1) ← Hoạt động trên CFG
    ↓
Pointer Analysis (Bài 4) ← Một dạng dataflow analysis
    ↓
Memory Management (Bài 3) ← Quản lý objects được pointer trỏ đến
```

### Ý nghĩa thực tiễn

**Đối với Compiler Engineers**:
- Hiểu cách optimize code
- Implement static analysis tools
- Design garbage collectors

**Đối với Software Developers**:
- Viết code an toàn hơn
- Hiểu performance implications
- Debug memory issues

**Đối với Security Researchers**:
- Phát hiện vulnerabilities (buffer overflow, use-after-free)
- Analyze malware
- Build security tools

### Hướng phát triển

Từ project cơ bản này có thể mở rộng:

1. **More dataflow analyses**: Live variables, available expressions
2. **Interprocedural analysis**: Phân tích qua nhiều functions
3. **Pointer analysis**: Anderson's, Steensgaard's algorithms
4. **Advanced optimizations**: Loop optimization, inlining
5. **Concurrency analysis**: Data races, deadlocks

### Lời kết

Project này là một công cụ giáo dục xuất sắc, kết hợp lý thuyết và thực hành một cách hiệu quả. Code được viết rõ ràng với nhiều comments bằng tiếng Việt và tiếng Anh, giúp người học dễ dàng hiểu các khái niệm phức tạp. Việc sử dụng Go thay vì C cũng giúp tập trung vào ý tưởng chính mà không bị phân tâm bởi các vấn đề quản lý bộ nhớ thủ công.

---

## PHỤ LỤC: Hướng dẫn chạy project

### Cài đặt
```bash
cd dataflow-cfg-mem-ptr
go mod tidy
```

### Chạy các demo

**Dataflow Analysis**:
```bash
go run ./cmd/tool dataflow
```

**CFG Construction**:
```bash
go run ./cmd/tool cfg
```

**Memory Management**:
```bash
go run ./cmd/tool mem
```

**Pointer & GC**:
```bash
go run ./cmd/tool ptr
```

### Cấu trúc project
```
dataflow-cfg-mem-ptr/
├── cmd/tool/main.go           # Entry point
├── internal/
│   ├── dataflow/reaching.go   # Bài tập 1
│   ├── cfg/cfg.go            # Bài tập 2
│   ├── mem/mem.go            # Bài tập 3
│   └── ptr/ptr.go            # Bài tập 4
├── go.mod
└── README.md
```

---

**Người thực hiện**: AI Assistant
**Ngày**: 28/11/2025
**Tài liệu tham khảo**: Source code của project dataflow-cfg-mem-ptr


