package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Snapshot struct {
	Version     int64             `json:"version"`
	Assignments map[string]string `json:"assignments"`
}

func main() {
	filePtr := flag.String("file", "state.json", "Path to the state JSON file")
	nodePtr := flag.String("node", "", "Filter tasks by node ID")
	flag.Parse()

	data, err := os.ReadFile(*filePtr)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", *filePtr, err)
		os.Exit(1)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cluster State Version: %d\n", snap.Version)
	fmt.Println("Assignments:")

	count := 0
	for task, node := range snap.Assignments {
		if *nodePtr == "" || *nodePtr == node {
			fmt.Printf("- Task: %s -> Node: %s\n", task, node)
			count++
		}
	}

	if count == 0 {
		if *nodePtr != "" {
			fmt.Printf("No tasks found for node: %s\n", *nodePtr)
		} else {
			fmt.Println("No tasks assigned in this state.")
		}
	} else {
		fmt.Printf("Total tasks matched: %d\n", count)
	}
}