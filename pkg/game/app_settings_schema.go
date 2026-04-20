package game

type ValidValue struct {
	Value			any		`json:"value"`
	DisplayText		string 	`json:"displayText"`
	Description		string 	`json:"description,omitempty"`
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

type SettingsSchema struct {
	Prefs	[]SettingsSchemaItem
	Game	[]SettingsSchemaItem
}

func (sc SettingsSchema) GetSchemaItem(propName string) *SettingsSchemaItem {
	for i, schema_item := range sc.Game {
		if schema_item.PropName == propName {
			return &sc.Game[i]
		}
	}
	for j, schema_item := range sc.Prefs {
		if schema_item.PropName == propName {
			return &sc.Prefs[j]
		}
	}

	return nil
}
