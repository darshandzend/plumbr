package plumbr

import (
	"errors"
	"log"
)

type (
	unitT   byte
	streamT []unitT
	inT     <-chan streamT //TODO: support for buffers
	outT    chan<- streamT
	funcT   func(streamT) streamT
)

type node struct {
	in  <-chan interface{}
	f   func(interface{}) interface{}
	out chan<- interface{}
}

func newNode(in <-chan interface{},
	fn func(interface{}) interface{}, //TODO: should be interface{}
	out chan<- interface{},
) node {
	log.Println("Creating node", in, fn, out)
	return node{in, fn, out}
}

func (n node) runNode() {
	log.Printf("node %v is waitng on chan %v\n", n, n.in)
	for input := range n.in {
		log.Println(n, "recieved", input)
		n.out <- n.f(input)
	}
	close(n.out)
}

type Pipeliner interface {
	Run()
}

type pipeline struct {
	in    <-chan interface{}
	nodes []node
	out   chan<- interface{}
}

func Pipeline(
	in <-chan interface{},
	out chan<- interface{},
	fns ...func(interface{}) interface{},
) (Pipeliner, error) {

	if len(fns) == 0 {
		return nil, errors.New(ERR_ZERO_NODES)
	}

	var nodes []node

	ini := in
	outi := make(chan interface{})
	for _, fn := range fns[:len(fns)-1] {
		nodes = append(nodes, newNode(ini, fn, outi))
		ini = outi
		outi = make(chan interface{})
	}
	//wire the last one
	nodes = append(nodes, newNode(ini, fns[len(fns)-1], out))

	log.Println("Created pipeline with", in, nodes, out)
	return pipeline{in, nodes, out}, nil
}

func (p pipeline) Run() {
	for _, n := range p.nodes {
		log.Println("Starting to run ", n)
		go n.runNode()
		//time.Sleep(time.Second)
	}
}

/*




















 */
