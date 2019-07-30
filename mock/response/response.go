package response

// Response swagg-doc:model:interface
// data swagg-doc:attribute:payload
type Response struct {
	Data     interface{} `json:"data"`
	Metadata Metadata    `json:"metadata"`
	Error    string      `json:"error"`
}

// Metadata swagg-doc:model
type Metadata struct {
	Filter string `json:"filter"`
	Order  string `json:"order"`
	Page   int    `json:"page"`
}
