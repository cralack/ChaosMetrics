package scheduler

type Scheduler interface {
	Schedule()
	Push(...*Task)
	Pull() *Task
}

type Task struct {
	Key   string
	Loc   string
	Retry uint
	Data  interface{}
}

type RiotDTOSchedule struct {
	RequestCh chan *Task
	WorkerCh  chan *Task
	ReqQueue  []*Task
}

var _ Scheduler = &RiotDTOSchedule{}

func NewSchdule() *RiotDTOSchedule {
	return &RiotDTOSchedule{
		RequestCh: make(chan *Task),
		WorkerCh:  make(chan *Task),
	}
}

func (s *RiotDTOSchedule) Push(reqs ...*Task) {
	for _, req := range reqs {
		s.RequestCh <- req
	}
}
func (s *RiotDTOSchedule) Pull() *Task {
	r := <-s.WorkerCh
	return r
}

func (s *RiotDTOSchedule) Schedule() {
	var (
		req *Task
		ch  chan *Task
	)
	
	for {
		// reque 不空则pop出 req
		if req == nil && len(s.ReqQueue) > 0 {
			// pop req
			req = s.ReqQueue[0]
			s.ReqQueue = s.ReqQueue[1:]
			ch = s.WorkerCh
		}
		// check req
		select {
		case r := <-s.RequestCh:
			s.ReqQueue = append(s.ReqQueue, r)
		case ch <- req:
			req = nil
			ch = nil
		}
	}
}
