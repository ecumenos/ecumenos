package appsettings

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func (m *configurations) LoadLocales(ctx context.Context, path string) error {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlFile, m); err != nil {
		return err
	}

	return nil
}

func (m *configurations) GetLanguages(onlyEnabled bool) []string {
	languages := make([]string, 0, len(m.Languages))
	for _, lang := range m.Languages {
		if onlyEnabled && !lang.Enabled {
			continue
		}
		languages = append(languages, lang.LangaugeCode)
	}

	return languages
}

func (m *configurations) ValidateLanguageCode(code string) error {
	for _, c := range m.Languages {
		if code == c.LangaugeCode {
			if !c.Enabled {
				return fmt.Errorf(`not allowed to take "%s" language the language is disabled`, code)
			}

			return nil
		}
	}

	return fmt.Errorf("not found language (language code=%s)", code)
}
