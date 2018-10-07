package worker

type Source interface {
	Produce()
}

type SourceWorkerPool struct {
	WorkerPool
	handler       Source
}

func NewSourceWorkerPool(name string, maxworkers int, handler Source) *SourceWorkerPool {
	wp := &SourceWorkerPool{}
	wp.maxworkers = maxworkers
	wp.handler = handler
	wp.name = name
	return wp
}

func (wp *SourceWorkerPool) Start() {
	for i := 0; i < wp.maxworkers; i++ {
		wp.wg.Add(1)
		go func() {
			wp.handler.Produce()
			wp.wg.Done()
		}()
	}
}

func (wp *SourceWorkerPool) Stop() {
	wp.wg.Wait()
}
