package vec

import (
	"fmt"
	"sync"
)

// WorkerJob data for worker to process
type WorkerJob struct {
	Vec1, Vec2, Out []float32
	Fn              MulFunc
	wg              *sync.WaitGroup
}

// Worker worker channels
type Worker struct {
	id         int
	dbgState   string
	workerPool *WorkerPool
	//Done chan int
}

// Start starts a worker(go routine)
func (w *Worker) Start() {
	go func() {
		for {
			work := <-w.workerPool.In
			work.Fn(work.Vec1, work.Vec2, work.Out)
			work.wg.Done()
			// work.Done() //per job waitgroup (Safe)
			//w.pool.Done() // Per pool waitgroup

		}
	}()
}

func (w *Worker) work(j WorkerJob) {
	w.workerPool.In <- j
}

func (w *Worker) String() string {
	return fmt.Sprintf("[worker %d (%s)]", w.id, w.dbgState)
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
	wp.wgPool.Put(wp.wgPool.Get())
	wp.wgPool.Put(wp.wgPool.Get())
	wp.wgPool.Put(wp.wgPool.Get())
	wp.wgPool.Put(wp.wgPool.Get())

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

	// Safer
	//wp.Lock()
	//defer wp.Unlock()
	// Lock all process
	// Contextual would be good here
	// send to bgRoutines
	wg := wp.wgPool.Get().(*sync.WaitGroup) // this is the alloc
	wg.Add(NWorkers)
	for i := range wp.workers { // Divide workload between cores?
		s := i * sz
		e := s + sz
		wp.In <- WorkerJob{vec1[s:e], vec2[s:e], out[s:e], fn, wg} // Copy all
	}
	wg.Wait()

}
