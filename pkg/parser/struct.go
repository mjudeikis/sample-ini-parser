package parser

// IniSection defines ini file [section]
type IniSection struct {
	Name          string        `json:"name"`
	KeyValuePairs []IniKeyValue `json:"keyValuePairs"`
}

// IniKeyValue defines ini file key=value
type IniKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// IniFile defines end-to-end ini file
type IniFile struct {
	FileName string       `json:"fileName"`
	Sections []IniSection `json:"sections"`
}
