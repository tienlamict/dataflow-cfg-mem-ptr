package main

import (
	"dataflow-cfg-mem-ptr/internal/cfg"
	"dataflow-cfg-mem-ptr/internal/dataflow"
	"dataflow-cfg-mem-ptr/internal/mem"
	"dataflow-cfg-mem-ptr/internal/ptr"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]

	// Check for -help flag
	if len(os.Args) > 2 && (os.Args[2] == "-help" || os.Args[2] == "--help") {
		printSubcommandHelp(subcommand)
		return
	}

	switch subcommand {
	case "dataflow":
		fmt.Println("=== Running Dataflow Analysis (Reaching Definitions) ===")
		fmt.Println()
		dataflow.RunReachingDefinitions()
	case "cfg":
		fmt.Println("=== Running CFG Construction ===")
		fmt.Println()
		cfg.RunCFGDemo()
	case "mem":
		fmt.Println("=== Running Memory Allocation Demo ===")
		fmt.Println()
		mem.RunMemoryDemo()
	case "ptr":
		fmt.Println("=== Running Pointer and GC Demo ===")
		fmt.Println()
		ptr.RunPointerDemo()
	default:
		fmt.Printf("Unknown subcommand: %s\n\n", subcommand)
		fmt.Println()
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run ./cmd/tool <subcommand> [-help]")
	fmt.Println("\nAvailable subcommands:")
	fmt.Println("  dataflow   - Run reaching definitions analysis")
	fmt.Println("  cfg        - Build and display Control Flow Graph")
	fmt.Println("  mem        - Demonstrate memory allocation and GC")
	fmt.Println("  ptr        - Demonstrate pointers and garbage collection")
	fmt.Println("\nUse '<subcommand> -help' for more information about a subcommand.")
}

func printSubcommandHelp(subcommand string) {
	switch subcommand {
	case "dataflow":
		fmt.Println("dataflow - Reaching Definitions Analysis")
		fmt.Println("\nThis command demonstrates dataflow analysis using reaching definitions.")
		fmt.Println("It constructs a simple CFG with 3 basic blocks, computes GEN/KILL sets,")
		fmt.Println("and iterates to fixpoint to calculate IN/OUT sets for each block.")
		fmt.Println("\nFormulas used:")
		fmt.Println("  IN[n] = ∪ OUT[pred]")
		fmt.Println("  OUT[n] = GEN[n] ∪ (IN[n] - KILL[n])")
	case "cfg":
		fmt.Println("cfg - Control Flow Graph Construction")
		fmt.Println("\nThis command demonstrates CFG construction from pseudo-code.")
		fmt.Println("It takes pseudo-code with labels and conditional jumps, identifies")
		fmt.Println("basic block leaders, splits into blocks, and displays the graph structure.")
	case "mem":
		fmt.Println("mem - Memory Allocation and GC Demo")
		fmt.Println("\nThis command demonstrates Go's memory management and garbage collection.")
		fmt.Println("It allocates large objects, displays memory statistics before and after,")
		fmt.Println("releases references, and triggers GC to show memory reclamation.")
	case "ptr":
		fmt.Println("ptr - Pointer and Garbage Collection Demo")
		fmt.Println("\nThis command demonstrates pointer usage and automatic memory management.")
		fmt.Println("It creates a linked list, displays nodes, removes references, and triggers")
		fmt.Println("GC to show how Go automatically frees unreachable memory (vs manual free in C).")
	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
	}
}
