package scanner

type transition struct {
	checkHandler func(rune)bool
	dstState int
}

type DFA struct {
	value string
	initState int
	currentState int
	totalState []int
	finalState map[int]int
	transition map[int][]transition
	result []Token
}

func NewDFA(initState int) *DFA {
	return &DFA{
		value: "",
		initState:    initState,
		currentState: initState,
		totalState:   []int{ initState },
		finalState:   make(map[int]int),
		transition:   make(map[int][]transition),
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
			d.result = append(d.result, *NewToken(v, d.value))
			d.Reset()
			return true
		}
	}

	d.result = append(d.result, *NewToken(NONTOKEN, d.value))

	d.Reset()
	
	return false
}

func (d *DFA) Verify() {
	for k, v := range d.finalState {
		if d.currentState == k {
			d.result = append(d.result, *NewToken(v, d.value))
			d.Reset()
			return
		}
	}

	d.result = append(d.result, *NewToken(NONTOKEN, d.value))
	d.Reset()
}

func (d *DFA) Reset() {
	d.value = ""
	d.currentState = d.initState
}
