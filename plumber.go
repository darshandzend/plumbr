/*

What you do:
1. Define Individual functions (with parameters)
2, BIG ASSUMPTION : functions are of type f(string) string
3. in, out := plumber.Pipeline(funcA, funcB, funcC)
4. Start putting things in "in"
5. Start extracting things from "out"

*/

package plumber

type unitT byte //TODO: Lots and lots of reflection ahead!
type streamT []unitT
type inChanT chan streamT //TODO: support for buffers TODO: Enforce directions
type outChanT chan streamT
type funcT func(streamT) streamT

type pipeline struct {
	in    inChanT
	out   outChanT
	funcs []funcT
}

type Pipeliner interface {
	AddNode(funcT, uint)
	DeleteNode(funcT, uint)

	//TODO: Add methods
	//Ge/SetNthNode()
	//GetNthChan()
	//etc
}

func New(funcs ...funcT) Pipeliner {
	in := make(inChanT)
	out := make(outChanT)
	p := pipeline{in, out, funcs}
	go p.run()
	return p
}

//TODO: Implement methods

func (p pipeline) run() {
	inChanI := p.in
	outChanI := make(outChanT)
	for _, fI := range p.funcs {
		go node(outChanI, inChanI, fI)
		inChanI = (outChanI).(inChanT)
		outChanI = make(outChanT)
	}
	p.out = outChanI
}

func node(outChanI outChanT, inChanI inChanT, fI funcT) {
	for val := range inChanI {
		outChanI <- fI(val)
	}
	close(inChanI)
}

/*




















 */
