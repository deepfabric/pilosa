package test

import (
	"strings"

	"github.com/deepfabric/pilosa"
	"github.com/deepfabric/pilosa/pql"
)

// Executor represents a test wrapper for pilosa.Executor.
type Executor struct {
	*pilosa.Executor
}

// NewExecutor returns a new instance of Executor.
// The executor always matches the hostname of the first cluster node.
func NewExecutor(holder *pilosa.Holder, cluster *pilosa.Cluster) *Executor {
	executor := pilosa.NewExecutor(nil)
	e := &Executor{Executor: executor}
	e.Holder = holder
	e.Cluster = cluster
	e.Scheme = cluster.Nodes[0].Scheme
	e.Host = cluster.Nodes[0].Host
	return e
}

// MustParse parses s into a PQL query. Panic on error.
func MustParse(s string) *pql.Query {
	q, err := pql.NewParser(strings.NewReader(s)).Parse()
	if err != nil {
		panic(err)
	}
	return q
}
