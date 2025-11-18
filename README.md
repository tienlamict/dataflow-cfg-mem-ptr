# Dataflow, CFG, Memory & Pointer Analysis Demo

This repository demonstrates fundamental concepts in program analysis and memory management using Go, including:

- **Dataflow Analysis** - Reaching definitions with fixpoint iteration
- **Control Flow Graph (CFG)** - Basic block identification and graph construction  
- **Memory Management** - Heap allocation and garbage collection
- **Pointer Analysis** - Dynamic data structures and automatic memory reclamation

## Purpose

This project provides hands-on examples that complement theoretical slides on program analysis topics. Each module demonstrates a specific concept:

### 1. Dataflow Analysis (`internal/dataflow/`)
Implements **reaching definitions** analysis on a simple 3-block CFG. Demonstrates:
- GEN/KILL set computation
- Fixpoint iteration using transfer functions
- IN/OUT set calculation
- Forward dataflow analysis

**Key formulas:**
```
IN[n]  = ∪ OUT[pred] for all predecessors of n
OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])
```

### 2. Control Flow Graph (`internal/cfg/`)
Builds a CFG from pseudo-code with labels and jumps. Demonstrates:
- Leader identification (first statement, jump targets, statements after jumps)
- Basic block partitioning
- Edge construction for successors
- Graph visualization

### 3. Memory Management (`internal/mem/`)
Shows Go's automatic memory management. Demonstrates:
- Heap allocation of large objects
- Memory statistics before/after allocation
- Garbage collection triggering
- Memory reclamation

**Comparison:** Unlike C/C++ with `malloc`/`free`, Go automatically manages memory through garbage collection.

### 4. Pointer Analysis (`internal/ptr/`)
Creates dynamic data structures (linked list) using pointers. Demonstrates:
- Dynamic memory allocation
- Pointer-based structures
- Automatic cleanup when references are dropped
- GC vs manual memory management

**Key insight:** When `head = nil`, all nodes become unreachable and GC automatically frees them (no `free()` needed like in C).

## Requirements

- Go 1.24.2 or higher
- Standard library only (no external dependencies)

## Installation

```bash
# Clone or navigate to repository
cd dataflow-cfg-mem-ptr

# Initialize module (already done via go.mod)
go mod tidy
```

## Usage

The CLI tool provides four subcommands. Use `-help` with any command for details.

### Dataflow Analysis
```bash
go run ./cmd/tool dataflow
```

**Example Output:**
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
```

### Control Flow Graph
```bash
go run ./cmd/tool cfg
```

**Example Output:**
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
  Successors: [2]

Block 4 (starts at line 8):
  Statements:
    8: print sum
    9: return
  Successors: [EXIT]
```

### Memory Management
```bash
go run ./cmd/tool mem
```

**Example Output:**
```
=== Running Memory Allocation Demo ===

=== Initial Memory Stats ===
  Alloc:            0.23 MB (bytes allocated and still in use)
  TotalAlloc:       0.23 MB (total bytes allocated over time)
  HeapAlloc:        0.23 MB (bytes allocated on heap)

=== Allocating Memory ===
Allocating 1000 objects of 1024 integers each...
Allocated 1000 objects (approximately 7.81 MB)

=== Memory Stats After Allocation ===
  Alloc:            8.05 MB (bytes allocated and still in use)
  HeapAlloc:        8.05 MB (bytes allocated on heap)

Memory increase:
  Alloc:      +7.82 MB
  HeapAlloc:  +7.82 MB

=== Releasing References ===
Setting slice to nil to make objects unreachable...
Manually triggering garbage collection...

=== Memory Stats After GC ===
  Alloc:            0.26 MB
  HeapAlloc:        0.26 MB

Memory after GC (vs after allocation):
  Alloc:     0.26 MB (was 8.05 MB, freed 7.79 MB)
  HeapAlloc: 0.26 MB (was 8.05 MB, freed 7.79 MB)
```

### Pointer & GC Demo
```bash
go run ./cmd/tool ptr
```

**Example Output:**
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

=== Dropping Reference to List ===
Setting head = nil makes all nodes unreachable...
(In C, this would be a memory leak without manual free!)

Triggering garbage collection...

Memory reclaimed: 0.42 KB

=== Summary ===
✓ Created linked list with dynamic allocation
✓ Dropped all references to the list
✓ GC automatically freed all unreachable nodes
✓ No manual free() needed!
```

## Makefile (Optional)

For convenience, you can create a `Makefile`:

```makefile
.PHONY: run-dataflow run-cfg run-mem run-ptr help

help:
	@echo "Available targets:"
	@echo "  run-dataflow - Run dataflow analysis demo"
	@echo "  run-cfg      - Run CFG construction demo"
	@echo "  run-mem      - Run memory management demo"
	@echo "  run-ptr      - Run pointer and GC demo"

run-dataflow:
	go run ./cmd/tool dataflow

run-cfg:
	go run ./cmd/tool cfg

run-mem:
	go run ./cmd/tool mem

run-ptr:
	go run ./cmd/tool ptr
```

Then run with: `make run-dataflow`, `make run-cfg`, etc.

## Project Structure

```
dataflow-cfg-mem-ptr/
├── cmd/tool/              # CLI entry point
│   └── main.go            # Main program with subcommand routing
├── internal/dataflow/     # Dataflow analysis (reaching definitions)
│   └── reaching.go        # Fixpoint iteration implementation
├── internal/cfg/          # Control Flow Graph construction
│   └── cfg.go             # Basic block splitting and graph building
├── internal/mem/          # Memory management demo
│   └── mem.go             # Allocation and GC demonstration
├── internal/ptr/          # Pointer demo
│   └── ptr.go             # Linked list with automatic cleanup
├── go.mod                 # Module definition
└── README.md              # This file
```

## Key Concepts

### Dataflow Analysis
- **Forward analysis**: Information flows from entry to exit
- **Fixpoint**: Iterating until IN/OUT sets stop changing
- **Transfer function**: How information changes across a block

### CFG Construction  
- **Leaders**: Statements that start basic blocks
- **Basic blocks**: Maximal sequences of consecutive statements
- **Edges**: Control flow between blocks (fall-through, jumps)

### Memory Management
- **Stack vs Heap**: Local variables vs dynamic allocation
- **Garbage Collection**: Automatic memory reclamation
- **Reachability**: Objects with no references are collected

### Pointers
- **Reference semantics**: Pointers refer to memory locations
- **Dynamic structures**: Linked lists, trees, graphs
- **Automatic cleanup**: No manual `free()` needed

## License

MIT License - Feel free to use for educational purposes.

## Notes

- All examples use only Go standard library
- Code includes detailed comments linking to slide concepts
- Examples are simplified for educational clarity
- Memory statistics may vary based on Go runtime version and system

