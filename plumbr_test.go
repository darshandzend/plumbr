package plumbr

import (
	"log"
	"testing"
)

type NullWriter int

const supressLogs = true

func (NullWriter) Write([]byte) (int, error) { return 0, nil }
func init() {
	//supress logs
	if supressLogs == true {
		log.SetOutput(new(NullWriter))
	}
}

func TestPipelineManyNodes(t *testing.T) {

	in := make(chan interface{})
	out := make(chan interface{})
	func1 := func(i1 interface{}) interface{} {
		in1 := i1.(string)
		out1 := in1 + " pants"
		return out1
	}
	func2 := func(i2 interface{}) interface{} {
		in2 := i2.(string)
		out2 := in2 + " tiny"
		return out2
	}
	func3 := func(i3 interface{}) interface{} {
		in3 := i3.(string)
		out3 := in3 + " ants"
		return out3
	}

	log.Println("Create Pipeline with", in, out, func1, func2, func3)
	p, err := Pipeline(in, out, func1, func2, func3)
	if err != nil {
		t.Errorf("Pipeline() returned error: %q\n", err)
	}
	p.Run()

	invals := []string{"Ice", "Fire", "Water", "Thunder"}
	go func() {
		for _, ival := range invals {
			log.Println("I sent", ival)
			in <- ival
		}
		close(in)
	}()

	i := 0
	//not ranging over invals so that close() is tested too
	for oval := range out {
		got := oval.(string)
		log.Println(got)

		expected := (invals[i] + " pants tiny ants")
		if got != expected {
			t.Errorf("Output mismatch: Got %q, expected %q\n",
				got, expected)
		}
		i++
	}

}

func TestPipelineSingleNode(t *testing.T) {

	in := make(chan interface{})
	out := make(chan interface{})
	func1 := func(i1 interface{}) interface{} {
		in1 := i1.(string)
		out1 := in1 + " pants"
		return out1
	}

	log.Println("Create Pipeline with", in, out, func1)
	p, err := Pipeline(in, out, func1)
	if err != nil {
		t.Errorf("Pipeline() returned error: %q\n", err)
	}
	p.Run()

	invals := []string{"Ice", "Fire", "Water", "Thunder"}
	go func() {
		for _, ival := range invals {
			in <- ival
		}
		close(in)
	}()

	i := 0
	//not ranging over invals so that close() is tested too
	for oval := range out {
		got := oval.(string)
		log.Println(got)

		expected := (invals[i] + " pants")
		if got != expected {
			t.Errorf("Output mismatch: Got %q, expected %q\n",
				got, expected)
		}
		i++
	}

}

func TestPipelineZeroNodes(t *testing.T) {

	in := make(chan interface{})
	out := make(chan interface{})

	log.Println("Create Pipeline with", in, out)
	_, err := Pipeline(in, out)
	expected := ERR_ZERO_NODES + " asf"
	if err == nil {
		t.Errorf("Pipeline doesn't throw error on zero nodes")
	} else if err.Error() != expected {
		t.Errorf("Error mismatch: Got %q, expected %q\n", err, expected)
	}
}
