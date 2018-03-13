package trends

// Record from git search
type Record struct {
	ID          int64
	Name        string
	URL         string
	Description string
	CloneURL    string
	Stars       int
}

// Results from git search
type Results struct {
	Outputs []Record
}
