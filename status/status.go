package status

var ServerStatus = SERVER_DONE

const (
	SERVER_DONE      = 0
	SERVER_RUN       = 1
	SERVER_WAIT_DONE = 2
)

var statusMap = map[int]struct{}{
	SERVER_DONE:      struct{}{},
	SERVER_RUN:       struct{}{},
	SERVER_WAIT_DONE: struct{}{},
}

func Update(s int) {
	if _, ok := statusMap[s]; !ok {
		s = SERVER_DONE
	}

	ServerStatus = s
}

func Get() int {
	return ServerStatus
}
