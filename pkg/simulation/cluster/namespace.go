package cluster

import (
	"fmt"

	"github.com/howardjohn/pilot-load/pkg/simulation/app"
	"github.com/howardjohn/pilot-load/pkg/simulation/config"
	"github.com/howardjohn/pilot-load/pkg/simulation/model"
)

type NamespaceSpec struct {
	Name     string
	Services []model.ServiceArgs
}

type Namespace struct {
	Spec      *NamespaceSpec
	ns        *KubernetesNamespace
	sa        map[string]*app.ServiceAccount
	sidecar   *config.Sidecar
	workloads []*app.Workload
}

var _ model.Simulation = &Namespace{}

func NewNamespace(s NamespaceSpec) *Namespace {
	ns := &Namespace{Spec: &s}

	ns.ns = NewKubernetesNamespace(KubernetesNamespaceSpec{
		Name: s.Name,
	})
	ns.sa = map[string]*app.ServiceAccount{
		"default": app.NewServiceAccount(app.ServiceAccountSpec{
			Namespace: ns.Spec.Name,
			Name:      "default",
		}),
	}
	ns.sidecar = config.NewSidecar(config.SidecarSpec{Namespace: s.Name})
	for i, w := range s.Services {
		ns.workloads = append(ns.workloads, app.NewWorkload(app.WorkloadSpec{
			App:            fmt.Sprintf("app-%d", i),
			Node:           "node",
			Namespace:      ns.Spec.Name,
			ServiceAccount: "default",
			Instances:      w.Instances,
		}))
	}
	return ns
}

func (n *Namespace) getSims() []model.Simulation {
	sims := []model.Simulation{n.ns, n.sidecar}
	for _, sa := range n.sa {
		sims = append(sims, sa)
	}
	for _, w := range n.workloads {
		sims = append(sims, w)
	}
	return sims
}

func (n *Namespace) Run(ctx model.Context) error {
	return model.AggregateSimulation{n.getSims()}.Run(ctx)
}

func (n *Namespace) Cleanup(ctx model.Context) error {
	return model.AggregateSimulation{n.getSims()}.Cleanup(ctx)
}
