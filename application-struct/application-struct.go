package applicationstruct

type Applications struct {
	Id          int32
	AppName     string
	Description string
	ImageName   string
}

type AppVersionsResponse struct {
	Status      string
	AppVersions []string
}

type Environment struct {
	Id          int32
	Name        string
	Description string
}
