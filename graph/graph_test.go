package graph_test

import (
	"context"
	"testing"

	g "github.com/tmc/langgraphgo/graph"
)

// TestMessageGraph_AddConditionalEdge tests adding a conditional edge to the graph.
func TestMessageGraph_Compile(t *testing.T) {
	graph := g.NewMessageGraph()
	graph.AddNode("node1", func(_ context.Context, state interface{}) (interface{}, error) {
		return state, nil
	})
	graph.AddNode(g.END, func(_ context.Context, state interface{}) (interface{}, error) {
		return state, nil
	})

	graph.AddEdge("node1", g.END)
	graph.SetEntryPoint("node1")

	runnable, err := graph.Compile()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if runnable == nil {
		t.Errorf("Expected non-nil runnable")
	}
}

// TestRunnable_Invoke tests invoking a simple graph.
func TestRunnable_Invoke(t *testing.T) {
	type State struct {
		visited bool
	}
	graph := g.NewMessageGraph()
	graph.AddNode("node1", func(_ context.Context, state interface{}) (interface{}, error) {
		agentState, _ := state.(State)
		agentState.visited = true
		return agentState, nil
	})
	graph.SetEntryPoint("node1")
	graph.AddNode(g.END, func(_ context.Context, state interface{}) (interface{}, error) {
		return state, nil
	})

	graph.AddEdge("node1", g.END)

	runnable, err := graph.Compile()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	state := State{visited: false}
	result, err := runnable.Invoke(context.Background(), state)
	stateResult, ok := result.(State)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !ok || !stateResult.visited {
		t.Errorf("Expected 'visited' to be true")
	}
}

// TestRunnable_InvokeWithConditionalEdge tests invoking a graph with a conditional edge.
func TestRunnable_InvokeWithConditionalEdge(t *testing.T) {
	type State struct {
		condition bool
		visited   string
	}
	graph := g.NewMessageGraph()
	state := State{}
	graph.AddNode("node1", func(_ context.Context, state interface{}) (interface{}, error) {
		agentState, _ := state.(State)
		agentState.condition = true
		return agentState, nil
	})
	graph.AddNode("node2", func(_ context.Context, state interface{}) (interface{}, error) {
		agentState, _ := state.(State)
		agentState.visited = "node2"
		return agentState, nil
	})
	graph.AddNode("node3", func(_ context.Context, state interface{}) (interface{}, error) {
		agentState, _ := state.(State)
		agentState.visited = "node3"
		return agentState, nil
	})
	graph.AddNode(g.END, func(_ context.Context, state interface{}) (interface{}, error) {
		return state, nil
	})
	graph.AddConditionalEdge("node1", "node2", "node3", func(_ context.Context, state interface{}) (bool, error) {
		agentState, ok := state.(State)
		condition := agentState.condition
		if !ok {
			return false, nil
		}
		return condition, nil
	})
	graph.AddEdge("node2", g.END)
	graph.AddEdge("node3", g.END)

	graph.SetEntryPoint("node1")
	runnable, err := graph.Compile()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	result, err := runnable.Invoke(context.Background(), state)
	agentResult, _ := result.(State)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if agentResult.visited != "node2" {
		t.Errorf("Expected 'visited' to be 'node2', got '%v'", agentResult.visited)
	}
}
