package storage

type LogEntry struct {
	Service string
	IP      string
	Path    *string
	Status  *string
	Method  *string
	Reason  *string
}

type Ban struct {
	IP     string
	Reason *string
}
