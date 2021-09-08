package graph

type Permission struct {
	Name string `json:"Permission.name"`
}

type Role struct {
	Name        string        `json:"Role.name"`
	Permissions []*Permission `json:"Role.permissions"`
	ETag        string        `json:"Role.etag"`
}

type Resource struct {
	Name   string    `json:"Resource.name"`
	Policy *Policy   `json:"Resource.policy"`
	Parent *Resource `json:"Resource.parent"`
	ETag   string    `json:"Resource.etag"`
}

type Policy struct {
	Bindings []*Binding `json:"Policy.bindings"`
	Version  int32      `json:"Policy.version"`
	ETag     string     `json:"Policy.etag"`
}

type Binding struct {
	Role    *Role    `json:"Binding.role"`
	Members []Member `json:"Binding.members"`
}

type Member struct {
	Group   string `json:"Group.name"`
	Subject string `json:"Subject.name"`
}

type Group struct {
	Name    string   `json:"Group.name"`
	Members []Member `json:"Group.members"`
	ETag    string   `json:"Group.etag"`
}

type Subject struct {
	Name string `json:"Subject.name"`
}
