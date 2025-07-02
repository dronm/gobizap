package fields

type Ref struct {
	Keys map[string]interface{} `json:"keys"`
	Descr string `json:"descr"`
	DataType string `json:"dataType"`
}

