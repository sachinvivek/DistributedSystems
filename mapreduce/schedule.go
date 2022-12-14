package mapreduce

import (
	"fmt"
	"sync"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	fmt.Printf("SCHEDULING %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	

	var w sync.WaitGroup
	w.Add(ntasks)

	for i:=0; i < ntasks; i++{
	
		go func(i int, phase jobPhase){
			
			defer w.Done()
			for{
				work:= <- mr.registerChannel
				var args DoTaskArgs
				
				switch phase{
				
				case mapPhase:
				
					args = DoTaskArgs{mr.jobName, mr.files[i], phase,i,nios, }
	
				case reducePhase:

					args = DoTaskArgs{mr.jobName, "",phase, i, nios, }
				}
				success := call(work, "Worker.DoTask",&args,new(struct{}))
				
				go func() {
					mr.registerChannel <- work
				}()
				if success == true{
					break
				}
			}
				
		}(i, phase) 
	}
	w.Wait()
	
	fmt.Printf("%v PHASE DONE!!!!\n", phase)
}
