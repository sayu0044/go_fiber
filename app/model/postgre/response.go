package postgre

// MetaInfo -> informasi pagination & filter
type MetaInfo struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Total  int    `json:"total"`
	Pages  int    `json:"pages"`
	SortBy string `json:"sortBy"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

// AlumniData -> data wrapper untuk GetAll Alumni
type AlumniData struct {
	Items []Alumni `json:"items"`
	Meta  MetaInfo `json:"meta"`
}

// PekerjaanData -> data wrapper untuk GetAll Pekerjaan
type PekerjaanData struct {
	Items []PekerjaanAlumni `json:"items"`
	Meta  MetaInfo          `json:"meta"`
}

// GetAllAlumniResponse -> response konsisten dengan Success, Message, Data
type GetAllAlumniResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    AlumniData `json:"data"`
}

// GetAllPekerjaanResponse -> response konsisten dengan Success, Message, Data
type GetAllPekerjaanResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    PekerjaanData `json:"data"`
}

// Legacy responses (deprecated, untuk backward compatibility)
type AlumniResponse struct {
	Data []Alumni `json:"data"`
	Meta MetaInfo `json:"meta"`
}

type PekerjaanResponse struct {
	Data []PekerjaanAlumni `json:"data"`
	Meta MetaInfo          `json:"meta"`
}
