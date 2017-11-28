package vec

import (
	"fmt"
	"sync"
)

// GoVecMul Routine workers
// Create routines per job
func GoVecMul(NWorkers int, vec1, vec2, out []float32, fn MulFunc) {
	wg := sync.WaitGroup{}
	wg.Add(NWorkers)
	lasti := NWorkers - 1
	for i := 0; i < NWorkers; i++ { // Divide workload between cores?
		sz := len(vec1) / NWorkers
		go func(i int) {
			s := i * sz
			e := s + sz
			if i == lasti {
				e = len(vec1)
			}
			fn(
				vec1[s:e],
				vec2[s:e],
				out[s:e],
			)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

// WorkerJob data for worker to process
type WorkerJob struct {
	Vec1, Vec2, Out []float32
	Fn              MulFunc
	*sync.WaitGroup
}

// Worker worker channels
type Worker struct {
	id         int
	workerPool *WorkerPool
}

// Start starts a worker(go routine)
func (w *Worker) Start() {
	go func() {
		for {
			work := <-w.workerPool.In
			work.Fn(work.Vec1, work.Vec2, work.Out)
			work.Done()
		}
	}()
}

//func (w *Worker) work(j WorkerJob) {
//	w.workerPool.In <- j
//}

func (w *Worker) String() string {
	return fmt.Sprintf("[worker %d]", w.id)
}

// WorkerPool struct containing all the running workers
type WorkerPool struct {
	In      chan WorkerJob
	workers []*Worker
	wgPool  sync.Pool
}

// NewWorkerPool creates several go routines to process vectors
func NewWorkerPool(NWorkers int) *WorkerPool {
	wp := &WorkerPool{
		In:      make(chan WorkerJob, NWorkers),
		workers: []*Worker{}, // we don't even need this
		wgPool: sync.Pool{
			New: func() interface{} { return &sync.WaitGroup{} },
		},
	}
	// Prealloc
	wg := []interface{}{}
	for i := 0; i < NWorkers; i++ {
		wg = append(wg, wp.wgPool.Get())
	}
	for _, v := range wg { // Put back
		wp.wgPool.Put(v)
	}

	wp.Launch(NWorkers)
	return wp
}

// Launch launch workers (go routines)
func (wp *WorkerPool) Launch(NWorkers int) {
	//gin := make(chan _bgData, 8) // Single in
	for i := 0; i < NWorkers; i++ {
		newWorker := &Worker{
			id:         i + 1,
			workerPool: wp,
		}
		wp.workers = append(wp.workers, newWorker)
		newWorker.Start()
	}
}

// VecMul multiply a vector with elements spreaded in multiple routines
func (wp *WorkerPool) VecMul(vec1, vec2, out []float32, fn MulFunc) {
	var NWorkers = len(wp.workers)
	var sz = len(vec1) / NWorkers

	// waitgroup for this session
	wg := wp.wgPool.Get().(*sync.WaitGroup) // this is the alloc
	wg.Add(NWorkers)
	lasti := len(wp.workers) - 1
	for i := range wp.workers { // Divide workload between cores?
		s := i * sz
		e := s + sz
		if i == lasti {
			e = len(vec1)
		}
		wp.In <- WorkerJob{vec1[s:e], vec2[s:e], out[s:e], fn, wg} // Copy all
	}
	wg.Wait()

}
