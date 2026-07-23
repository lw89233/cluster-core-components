package main

import (
	"fmt"

	"github.com/lw89233/cluster-core-components/pkg/reconciler"
)

func main() {
	desiredState := reconciler.State{
		Tasks: map[string]bool{
			"k8s-pod-1": true,
			"k8s-pod-2": true,
			"k8s-pod-3": true,
		},
	}

	actualState := reconciler.State{
		Tasks: map[string]bool{
			"k8s-pod-1": true,
			"k8s-pod-4": true,
		},
	}

	r := reconciler.New()
	actions := r.ComputeActions(desiredState, actualState)

	for _, action := range actions {
		fmt.Printf("%s\t%s\n", action.Type, action.TaskID)
	}
}