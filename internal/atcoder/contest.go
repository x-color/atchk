package atcoder

type Contests map[string]*Contest

type Contest struct {
	tasks map[string]*Task
}

type Task struct {
	fingerprint string
	samples     []*Sample
}

type Sample struct {
	id     int
	input  string
	output string
}

func (smp *Sample) Exec() bool {
	return true
}
