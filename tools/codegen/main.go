package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	idlPath := "../../idl/starter_program.json"
	outputPath := "../../pkg/generated/starterprogram"

	fmt.Println("Generating code from IDL...")
	fmt.Printf("IDL: %s\n", idlPath)
	fmt.Printf("Output: %s\n", outputPath)

	if err := os.MkdirAll(outputPath, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	cmd := exec.Command("carbon", "codegen", "--idl", idlPath, "--output", outputPath, "--package", "starterprogram")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("codegen failed: %v", err)
	}

	fmt.Println("Code generation completed successfully!")
}
