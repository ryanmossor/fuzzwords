package game

import (
	"encoding/json"
	"fmt"
	"fzwds/pkg/enums"
	"fzwds/pkg/utils"
	"log/slog"
	"os"
	"reflect"
)

type GeneralPreferences struct {
	AnimationsEnabled		bool		`json:"animationsEnabled"`
	BellEnabled				bool		`json:"bellEnabled"`
	HealthDisplay			string		`json:"healthDisplay"`
	// TODO: color theme, health bar color?
}

type GameSettings struct {
	Dictionary				enums.Dictionary	`json:"dictionary"`
	Alphabet				enums.Alphabet		`json:"alphabet"`
	PromptMode				enums.PromptMode	`json:"promptMode"`
	WinCondition			enums.WinCondition	`json:"winCondition"`
	HealthInitial			int					`json:"healthInital"`
	HealthMax				int					`json:"healthMax"`
	HighlightInput			bool				`json:"highlightInput"`
	PromptLenMin			int					`json:"promptLenMin"`
	PromptLenMax			int					`json:"promptLenMax"`
	TurnDurationMin			int					`json:"turnDurationMin"`
	PromptStrikes			int					`json:"promptStrikes"`
	PokemonGens				[]int				`json:"pokemonGens"`
}

type Settings struct {
	Prefs		GeneralPreferences	`json:"prefs"`
	Game 		GameSettings		`json:"game"`
}

func GetDefaultGeneralPreferences() GeneralPreferences {
	return GeneralPreferences {
		AnimationsEnabled: 	true,
		BellEnabled: 		false,
		HealthDisplay:		"● ;◯ ",
	}
}

func GetDefaultGameSettings() GameSettings {
	return GameSettings {
		Dictionary:			enums.English,
		Alphabet: 			enums.AlphabetEasy,
		PromptMode: 		enums.PromptModeFuzzy,
		WinCondition: 		enums.WinConditionEndless,
		HealthInitial: 		2,
		HealthMax: 			3,
		HighlightInput: 	false,
		PromptLenMin: 		2,
		PromptLenMax: 		3,
		TurnDurationMin: 	10,
		PromptStrikes:		2,
	}
}

func GetDefaultSettings() Settings {
	return Settings {
		Prefs: GetDefaultGeneralPreferences(),
		Game: GetDefaultGameSettings(),
	}
}

func (s *Settings) GetSetting(propName string) any {
	switch propName {
	case "Dictionary":
		return s.Game.Dictionary.String()
	case "Alphabet":
		return s.Game.Alphabet.String()
	case "PromptMode":
		return s.Game.PromptMode.String()
	case "WinCondition":
		return s.Game.WinCondition.String()
	case "HealthInitial":
		return s.Game.HealthInitial
	case "HealthMax":
		return s.Game.HealthMax
	case "HighlightInput":
		return s.Game.HighlightInput
	case "PromptLenMin":
		return s.Game.PromptLenMin
	case "PromptLenMax":
		return s.Game.PromptLenMax
	case "TurnDurationMin":
		return s.Game.TurnDurationMin
	case "PromptStrikes":
		return s.Game.PromptStrikes
	case "AnimationsEnabled":
		return s.Prefs.AnimationsEnabled
	case "BellEnabled":
		return s.Prefs.BellEnabled
	case "HealthDisplay":
		return s.Prefs.HealthDisplay
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
		s.SetSetting(propName, schema_item.Default, schema)
		return fmt.Errorf("Invalid value for %s: %v", propName, value)
	}

	switch propName {
	case "Dictionary":
		if vStr, ok := value.(string); ok {
			s.Game.Dictionary = enums.ParseDictionary(vStr)
		}
	case "Alphabet":
		if vStr, ok := value.(string); ok {
			s.Game.Alphabet = enums.ParseAlphabet(vStr)
		}
	case "PromptMode":
		if vStr, ok := value.(string); ok {
			s.Game.PromptMode = enums.ParsePromptMode(vStr)
		}
	case "WinCondition":
		if vStr, ok := value.(string); ok {
			s.Game.WinCondition = enums.ParseWinCond(vStr)
		}
	case "HealthInitial":
		if vInt, ok := utils.ParseInt(value); ok {
			s.Game.HealthInitial = vInt
		}
	case "HealthMax":
		if vInt, ok := utils.ParseInt(value); ok {
			s.Game.HealthMax = vInt
		}
	case "HighlightInput":
		if vbool, ok := value.(bool); ok {
			s.Game.HighlightInput = vbool
		}
	case "PromptLenMin":
		if vInt, ok := utils.ParseInt(value); ok {
			s.Game.PromptLenMin = vInt
		}
	case "PromptLenMax":
		if vInt, ok := utils.ParseInt(value); ok {
			s.Game.PromptLenMax = vInt
		}
	case "TurnDurationMin":
		if vInt, ok := utils.ParseInt(value); ok {
			s.Game.TurnDurationMin = vInt
		}
	case "PromptStrikes":
		if vInt, ok := utils.ParseInt(value); ok {
			s.Game.PromptStrikes = vInt
		}
	case "AnimationsEnabled":
		if vbool, ok := value.(bool); ok {
			s.Prefs.AnimationsEnabled = vbool
		}
	case "BellEnabled":
		if vbool, ok := value.(bool); ok {
			s.Prefs.BellEnabled = vbool
		}
	case "HealthDisplay":
		if vStr, ok := value.(string); ok {
			s.Prefs.HealthDisplay = vStr
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

func (s *Settings) SetDictionary(dictionary string, schema SettingsSchema) *Settings {
	s.SetSetting("Dictionary", dictionary, schema)
	return s
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

func (s *Settings) SetPromptStrikes(strikes int, schema SettingsSchema) *Settings {
	s.SetSetting("PromptStrikes", strikes, schema)
	return s
}

func (s *Settings) SetTurnDurationMin(duration int, schema SettingsSchema) *Settings {
	s.SetSetting("TurnDurationMin", duration, schema)
	return s
}

func (s *Settings) SetHealthDisplay(display string, schema SettingsSchema) *Settings {
	s.SetSetting("HealthDisplay", display, schema)
	return s
}

func (s *Settings) SetAnimationsEnabled(enabled bool, schema SettingsSchema) *Settings {
	s.SetSetting("AnimationsEnabled", enabled, schema)
	return s
}

func (s *Settings) SetBellEnabled(enabled bool, schema SettingsSchema) *Settings {
	s.SetSetting("BellEnabled", enabled, schema)
	return s
}

func (s *Settings) ValidateSettings(schema SettingsSchema) *Settings {
	return s.
		SetDictionary(s.Game.Dictionary.String(), schema).
		SetAlphabet(s.Game.Alphabet.String(), schema).
		SetHealthInitial(s.Game.HealthInitial, schema).
		SetHealthMax(s.Game.HealthMax, schema).
		SetPromptLenMin(s.Game.PromptLenMin, schema).
		SetPromptLenMax(s.Game.PromptLenMax, schema).
		SetPromptMode(s.Game.PromptMode.String(), schema).
		SetPromptStrikes(s.Game.PromptStrikes, schema).
		SetTurnDurationMin(s.Game.TurnDurationMin, schema).
		SetWinCondition(s.Game.WinCondition.String(), schema).
		SetAnimationsEnabled(s.Prefs.AnimationsEnabled, schema).
		SetBellEnabled(s.Prefs.BellEnabled, schema).
		SetHealthDisplay(s.Prefs.HealthDisplay, schema)
}

func ValidateSettingValue(schema_item SettingsSchemaItem, value any) bool {
	switch schema_item.Type {
	case "int":
		vInt, ok := utils.ParseInt(value)
		if !ok {
			return false
		}

		if len(schema_item.ValidValues) > 0 {
			for _, vv := range schema_item.ValidValues {
				vvInt, ok := utils.ParseInt(vv.Value)
				if ok && vInt == vvInt {
					return true
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

func (s Settings) WriteSettings(path string) {
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		slog.Error("Error marshaling settings", "error", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		slog.Error("Error writing settings.json", "error", err)
	}
}
