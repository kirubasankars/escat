package main

type indices struct {
	Index string `json:"index"`
}

type templates struct {
	Name string `json:"name"`
}

type snapshots struct {
	ID string `json:"id"`
}

type catRequest struct {
	format     string
	verbose    bool
	arg1       string
	fields     string
	sortFields string
	def        bool
	debug      bool
}
