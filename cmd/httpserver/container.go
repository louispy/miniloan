package main

type Container struct {
}

type Opts struct{}

func NewContainer(o Opts) *Container {
	return &Container{}
}
