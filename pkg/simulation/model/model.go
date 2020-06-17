package model

import (
	"fmt"

	"github.com/howardjohn/pilot-load/pkg/kube"
	"github.com/howardjohn/pilot-load/pkg/simulation/util"
)

type Simulation interface {
	// Run starts the simulation. If the simulation is long lived, this should be done asynchronously,
	// watching ctx.Done() for termination.
	Run(ctx Context) error
	// Cleanup tears down the simulation.
	// TODO do not pass context. Simulations should store it and then cancel the context. This means we should always pass a new ctx.
	Cleanup(ctx Context) error
}

type Args struct {
	PilotAddress string
	NodeMetadata string
	KubeConfig   string
}

type Context struct {
	//context.Context
	Args   Args
	Client *kube.Client
}

type AggregateSimulation struct {
	Simulations []Simulation
}

var _ Simulation = AggregateSimulation{}

func (a AggregateSimulation) Run(ctx Context) error {
	for _, s := range a.Simulations {
		// TODO pass unique context so we can cancel independently
		if err := s.Run(ctx); err != nil {
			return fmt.Errorf("failed running simulation %T: %v", s, err)
		}
	}
	return nil
}

func (a AggregateSimulation) Cleanup(ctx Context) error {
	var err error
	for _, s := range a.Simulations {
		if err := s.Cleanup(ctx); err != nil {
			err = util.AddError(err, fmt.Errorf("failed cleaning simulation %T: %v", s, err))
		}
	}
	return err
}