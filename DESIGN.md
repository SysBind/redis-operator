# Operator Design

## State Pattern

The operator leverages the [state design pattern](https://golangbyexample.com/state-design-pattern-go/) to
devide the logic into states and transitions, sample flow:

<new.boot> -> booting
<booting> -> init
<init.create> -> creating
<creating> -> stable

Some state transitions triggered by calling function (boot, create)
And others (transition from booting, creating, etc.) are triggered by external cluster events.    
 
All the code is under the cluster/ subdir, each state is a go file: newState.go, initState.go etc.