package game

import (
	"fmt"
	"fzwds/src/enums"
	"log/slog"
	"reflect"
)

type Settings struct {
	Alphabet				enums.Alphabet		`json:"alphabet"`
	PromptMode				enums.PromptMode	`json:"promptMode"`
	WinCondition			enums.WinCondition	`json:"winCondition"`
	HealthInitial			int					`json:"healthInital"`
	HealthMax				int					`json:"healthMax"`
	HighlightInput			bool				`json:"highlightInput"`
	PromptLenMin			int					`json:"promptLenMin"`
	PromptLenMax			int					`json:"promptLenMax"`
	TurnDurationMin			int					`json:"turnDurationMin"`
	PromptStrikesMax		int					`json:"promptStrikesMax"`
	// TODO: add cfg for hints after each strike?
	// hints_enabled			bool
	// hint_chars_per_turn		int
}

func InitializeSettings() Settings {
	return Settings{
		Alphabet: 			enums.FullAlphabet,
		PromptMode: 		enums.Fuzzy,
		WinCondition: 		enums.Endless,
		HealthInitial: 		2,
		HealthMax: 			3,
		HighlightInput: 	false,
		PromptLenMin: 		2,
		PromptLenMax: 		3,
		TurnDurationMin: 	10,
		PromptStrikesMax:	2,
	}
}

func EasySettings() Settings {
	cfg := InitializeSettings()
	cfg.Alphabet = enums.EasyAlphabet
	cfg.TurnDurationMin = 20
	return cfg
}

func MediumSettings() Settings {
	cfg := InitializeSettings()
	cfg.Alphabet = enums.MediumAlphabet
	cfg.TurnDurationMin = 10
	return cfg
}

func DifficultSettings() Settings {
	cfg := InitializeSettings()
	cfg.Alphabet = enums.FullAlphabet
	cfg.TurnDurationMin = 5
	cfg.PromptLenMax = 4
	cfg.PromptStrikesMax = 2
	return cfg
}

func ExpertSettings() Settings {
	cfg := InitializeSettings()
	cfg.Alphabet = enums.FullAlphabet
	cfg.HealthInitial = 1
	cfg.HealthMax = 1
	cfg.TurnDurationMin = 5
	cfg.PromptLenMax = 5
	cfg.PromptStrikesMax = 1
	return cfg
}

func (s *Settings) GetSetting(propName string) any {
	switch propName {
	case "Alphabet":
		return s.Alphabet.String()
	case "PromptMode":
		return s.PromptMode.String()
	case "WinCondition":
		return s.WinCondition.String()
	case "HealthInitial":
		return s.HealthInitial
	case "HealthMax":
		return s.HealthMax
	case "HighlightInput":
		return s.HighlightInput
	case "PromptLenMin":
		return s.PromptLenMin
	case "PromptLenMax":
		return s.PromptLenMax
	case "TurnDurationMin":
		return s.TurnDurationMin
	case "PromptStrikesMax":
		return s.PromptStrikesMax
	}

	return nil
}

func (s *Settings) SetSetting(propName string, value any, schema SettingsSchema) error {
	slog.Debug("SetSetting",
		"propName", propName,
		"value", value,
		"valType", reflect.TypeOf(value).String())

	schema_item := schema.GetSchemaItem(propName)
	if schema_item == nil {
		slog.Error("Unknown setting", "propName", propName, "value", value)
		return fmt.Errorf("Unknown setting: %s", propName)
	}

	if !ValidateSettingValue(*schema_item, value) {
		slog.Error("Invalid setting value provided. Setting default",
			"propName", propName,
			"value", value,
			"default", schema_item.Default)
		// TODO: Since Default is of type any, default value is float64 for numeric settings
		// Consider changing int settings to float64 to avoid int parsing issues?
		s.SetSetting(propName, schema_item.Default, schema)
		return fmt.Errorf("Invalid value for %s: %v", propName, value)
	}

	switch propName {
	case "Alphabet":
		if vStr, ok := value.(string); ok {
			s.Alphabet = enums.ParseAlphabet(vStr)
		}
	case "PromptMode":
		if vStr, ok := value.(string); ok {
			s.PromptMode = enums.ParsePromptMode(vStr)
		}
	case "WinCondition":
		if vStr, ok := value.(string); ok {
			s.WinCondition = enums.ParseWinCond(vStr)
		}
	case "HealthInitial":
		if vInt, ok := value.(int); ok {
			s.HealthInitial = vInt
		}
	case "HealthMax":
		if vInt, ok := value.(int); ok {
			s.HealthMax = vInt
		}
	case "HighlightInput":
		if vbool, ok := value.(bool); ok {
			s.HighlightInput = vbool
		}
	case "PromptLenMin":
		if vInt, ok := value.(int); ok {
			s.PromptLenMin = vInt
		}
	case "PromptLenMax":
		if vInt, ok := value.(int); ok {
			s.PromptLenMax = vInt
		}
	case "TurnDurationMin":
		if vFloat, ok := value.(float64); ok {
			s.TurnDurationMin = int(vFloat)
		} else if vInt, ok := value.(int); ok {
			s.TurnDurationMin = vInt
		}
	case "PromptStrikesMax":
		if vFloat, ok := value.(float64); ok {
			slog.Debug("Setting PromptStrikesMax as FLOAT", "val", vFloat)
			s.PromptStrikesMax = int(vFloat)
		} else if vInt, ok := value.(int); ok {
			slog.Debug("Setting PromptStrikesMax as INT", "val", vInt)
			s.PromptStrikesMax = vInt
		}
	}

	if schema_item.BindTo != "" && schema_item.BindRule != "" {
		bind_config := schema.GetSchemaItem(schema_item.BindTo)
		if bind_config != nil && schema_item.Type == "int" && bind_config.Type == "int" {
			cur_val := s.GetSetting(propName).(int)
			bind_val := s.GetSetting(schema_item.BindTo).(int)
			switch schema_item.BindRule {
			case "<=":
				if cur_val > bind_val {
					s.SetSetting(schema_item.BindTo, cur_val, schema)
				}
			case ">=":
				if cur_val < bind_val {
					s.SetSetting(schema_item.BindTo, cur_val, schema)
				}
			}
		}
	}

	return nil
}

func (s *Settings) SetAlphabet(alphabet string, schema SettingsSchema) *Settings {
	s.SetSetting("Alphabet", alphabet, schema)
	return s
}

func (s *Settings) SetHealthInitial(health int, schema SettingsSchema) *Settings {
	s.SetSetting("HealthInitial", health, schema)
	return s
}

func (s *Settings) SetHealthMax(health int, schema SettingsSchema) *Settings {
	s.SetSetting("HealthMax", health, schema)
	return s
}

func (s *Settings) SetPromptLenMin(len int, schema SettingsSchema) *Settings {
	s.SetSetting("PromptLenMin", len, schema)
	return s
}

func (s *Settings) SetPromptLenMax(len int, schema SettingsSchema) *Settings {
	s.SetSetting("PromptLenMax", len, schema)
	return s
}

func (s *Settings) SetPromptMode(mode string, schema SettingsSchema) *Settings {
	s.SetSetting("PromptMode", mode, schema)
	return s
}

func (s *Settings) SetWinCondition(cond string, schema SettingsSchema) *Settings {
	s.SetSetting("WinCondition", cond, schema)
	return s
}

func (s *Settings) SetPromptStrikesMax(strikes int, schema SettingsSchema) *Settings {
	s.SetSetting("PromptStrikesMax", strikes, schema)
	return s
}

func (s *Settings) SetTurnDurationMin(duration int, schema SettingsSchema) *Settings {
	s.SetSetting("TurnDurationMin", duration, schema)
	return s
}

func (s *Settings) ValidateSettings(schema SettingsSchema) *Settings {
	return s.
		SetAlphabet(s.Alphabet.String(), schema).
		SetHealthInitial(s.HealthInitial, schema).
		SetHealthMax(s.HealthMax, schema).
		SetPromptLenMin(s.PromptLenMin, schema).
		SetPromptLenMax(s.PromptLenMax, schema).
		SetPromptMode(s.PromptMode.String(), schema).
		SetPromptStrikesMax(s.PromptStrikesMax, schema).
		SetTurnDurationMin(s.TurnDurationMin, schema).
		SetWinCondition(s.WinCondition.String(), schema)
}

func ValidateSettingValue(schema_item SettingsSchemaItem, value any) bool {
	switch schema_item.Type {
	case "int":
		var vInt int
		switch val := value.(type) {
		case int:
			vInt = val
		case float64:
			vInt = int(val)
		default:
			return false
		}

		if len(schema_item.ValidValues) > 0 {
			for _, vv := range schema_item.ValidValues {
				switch vvInt := vv.Value.(type) {
				case float64:
					if vInt == int(vvInt) {
						return true
					}
				case int:
					if vInt == vvInt {
						return true
					}
				}
			}
			return false
		}

		if schema_item.Min != nil && vInt < *schema_item.Min {
			return false
		}
		if schema_item.Max != nil && vInt > *schema_item.Max {
			return false
		}

		return true
	case "enum", "string":
		vStr, ok := value.(string)
		if !ok {
			return false
		}

		if len(schema_item.ValidValues) > 0 {
			for _, vv := range schema_item.ValidValues {
				if vvStr, ok := vv.Value.(string); ok {
					if vStr == vvStr {
						return true
					}
				}
			}
			return false
		}
		return true // allow any if not specified
	case "bool":
		_, ok := value.(bool)
		return ok
	}

	return false
}
