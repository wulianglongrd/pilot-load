package cluster

import (
	"github.com/howardjohn/pilot-load/pkg/simulation/model"
)

type ClusterSpec struct {
	Namespaces map[string]model.NamespaceArgs
}

type Cluster struct {
	Spec       *ClusterSpec
	namespaces map[string]*Namespace
}

var _ model.Simulation = &Cluster{}

func NewCluster(s ClusterSpec) *Cluster {
	cluster := &Cluster{Spec: &s, namespaces: map[string]*Namespace{}}

	for name, ns := range s.Namespaces {
		cluster.namespaces[name] = NewNamespace(NamespaceSpec{Name: name, Services: ns.Services})
	}
	return cluster
}

func (c *Cluster) getSims() []model.Simulation {
	sims := []model.Simulation{}
	for _, ns := range c.namespaces {
		sims = append(sims, ns)
	}
	return sims
}

func (n *Cluster) Run(ctx model.Context) error {
	return model.AggregateSimulation{n.getSims()}.Run(ctx)
}

func (n *Cluster) Cleanup(ctx model.Context) error {
	return model.AggregateSimulation{n.getSims()}.Cleanup(ctx)
}
