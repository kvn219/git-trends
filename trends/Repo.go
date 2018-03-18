package trends

// Record from git search
type Record struct {
	ID          *int64  `json:"id,omitempty"`
	Name        *string `json:"name"`
	URL         *string `json:"url"`
	Description *string `json:"description"`
	CloneURL    *string `json:"clone_url"`
	Stars       *int    `json:"stars"`
	ForksCount  *int    `json:"forks_count"`
}

// Results from git search
type Results struct {
	Outputs []Record
}

// UIRecord is a type for the UI.
type UIRecord struct {
	Name       string
	Stars      int
	ForksCount int
	URL        string
}
