package scheduler

type Scheduler interface {
	Schedule()
	Push(...*Task)
	Pull() *Task
}

type Task struct {
	Key      string
	Loc      string
	ApiToken string
	Priority int
	Retry    uint
	Data     interface{}
}

type RiotDTOSchedule struct {
	requestCh   chan *Task
	workerCh    chan *Task
	reqQueue    []*Task
	priReqQueue []*Task
}

var _ Scheduler = &RiotDTOSchedule{}

func NewSchdule() *RiotDTOSchedule {
	return &RiotDTOSchedule{
		requestCh: make(chan *Task),
		workerCh:  make(chan *Task),
	}
}

func (s *RiotDTOSchedule) Push(reqs ...*Task) {
	for _, req := range reqs {
		s.requestCh <- req
	}
}
func (s *RiotDTOSchedule) Pull() *Task {
	r := <-s.workerCh
	return r
}

func (s *RiotDTOSchedule) Schedule() {
	var (
		req *Task
		ch  chan *Task
	)

	for {
		//priority queue first
		if req == nil && len(s.priReqQueue) > 0 {
			req = s.priReqQueue[0]
			s.priReqQueue = s.priReqQueue[1:]
			ch = s.workerCh
		}
		// reque 不空则pop出 req
		if req == nil && len(s.reqQueue) > 0 {
			// pop req
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerCh
		}

		// check req
		// schedule
		select {
		case r := <-s.requestCh:
			if r.Priority > 0 {
				s.priReqQueue = append(s.priReqQueue, r)
			} else {
				s.reqQueue = append(s.reqQueue, r)
			}
		case ch <- req:
			req = nil
			ch = nil
		}
	}
}
