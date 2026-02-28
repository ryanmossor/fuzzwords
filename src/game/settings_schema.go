package game

type ValidValue struct {
	Value			any		`json:"value"`
	Description		string 	`json:"description"`
}

type SettingsSchemaItem struct {
	PropName		string 			`json:"propName"`
	DisplayName		string 			`json:"displayName"`
	Type			string			`json:"type"`
	Disabled		bool 			`json:"disabled"`
	Description		string 			`json:"description,omitempty"`
	Default			any 			`json:"default"`
	Min				*int 			`json:"min,omitempty"`
	Max				*int 			`json:"max,omitempty"`
	ValidValues		[]ValidValue 	`json:"validValues,omitempty"`
	BindTo			string			`json:"bindTo,omitempty"`
	BindRule		string			`json:"bindRule,omitempty"`
}

type SettingsSchema []SettingsSchemaItem

func (sc SettingsSchema) GetSchemaItem(propName string) *SettingsSchemaItem {
	for i, schema_item := range sc {
		if schema_item.PropName == propName {
			return &sc[i]
		}
	}

	return nil
}
