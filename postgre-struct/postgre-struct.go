package postgrestruct

type PostgresApp struct {
	Id          int32
	Status      string
	Description string
	ServerName  string
	AdminUser   string
	Password    string
	Version     string
	Environment string
	SizeDisk    string
	StorageType string
	SizeCPU     string
	SizeMemory  string
	Issues      string
}

type ErrorResponse struct {
	Message    string
	StatusCode string
}

type SuccessResponse struct {
	Message    string
	StatusCode string
	Status     string
}
