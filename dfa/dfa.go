package dfa

const ErrorToken = -1

type transition struct {
	checkHandler func(rune)bool
	dstState int
}

type DFA struct {
	value string

	// default error state : -1
	initState int
	currentState int
	totalState []int
	finalState map[int]int
	transition map[int][]transition

	// create result
	resultHandler func(state int, value string)interface{}
	result []interface{}
}

func NewDFA(initState int, resultHandler func(int, string)interface{}) *DFA {
	return &DFA{
		value: "",
		initState:    initState,
		currentState: initState,
		totalState:   []int{ initState },
		finalState:   make(map[int]int),
		transition:   make(map[int][]transition),

		resultHandler: resultHandler,
	}
}

func (d *DFA) AddState(state int, stateType int) {
	d.totalState = append(d.totalState, state)

	if stateType != -1 {
		d.finalState[state] = stateType
	}
}

func (d *DFA) AddTransition(srcState int, dstState int ,checkHandler func(rune)bool) {
	trans, ok := d.transition[srcState]

	if !ok {
		var tempTrans []transition
		d.transition[srcState] = append(tempTrans, transition{
			checkHandler: checkHandler,
			dstState: dstState,
		})
		return
	}

	d.transition[srcState] = append(trans, transition{
		checkHandler: checkHandler,
		dstState:     dstState,
	})
}

func (d *DFA) Input(r rune) bool {
	trans := d.transition[d.currentState]

	// state move
	for _, t := range trans {
		if t.checkHandler(r) {
			d.currentState = t.dstState
			d.value += string(r)
			return false
		}
	}

	// final state
	for k, v := range d.finalState {
		if d.currentState == k {
			d.result = append(d.result, d.resultHandler(v, d.value))
			d.Reset()
			return true
		}
	}

	d.errorToken()

	return false
}

func (d *DFA) Verify() {
	for k, v := range d.finalState {
		if d.currentState == k {
			d.result = append(d.result, d.resultHandler(v, d.value))
			d.Reset()
			return
		}
	}

	d.errorToken()
}

func (d *DFA) Reset() {
	d.value = ""
	d.currentState = d.initState
}

func (d *DFA) GetResult() []interface{} {
	return d.result
}

func (d *DFA) errorToken() {
	d.result = append(d.result, d.resultHandler(ErrorToken, d.value))
	d.Reset()
}
