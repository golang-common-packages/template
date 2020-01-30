package database

// Database model for database config
type Database struct {
	MongoDB MongoDB `json:"mongodb"`
}

// MongoDB model for MongoDB config
type MongoDB struct {
	User     string   `json:"user"`
	Password string   `json:"password"`
	Hosts    []string `json:"hosts"`
	DB       string   `json:"db"`
	Options  []string `json:"options"`
}
