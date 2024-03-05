package appsettings

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func (m *configurations) LoadRegions(ctx context.Context, path string) error {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlFile, m); err != nil {
		return err
	}

	return nil
}

func (m *configurations) GetCountries(onlyEnabled bool) []string {
	countries := make([]string, 0, len(m.Countries))
	for _, c := range m.Countries {
		if onlyEnabled && !c.Enabled {
			continue
		}
		countries = append(countries, c.CountryCode)
	}

	return countries
}

func (m *configurations) GetRegionsByCountryCode(code string) []string {
	for _, c := range m.Countries {
		if c.CountryCode == code {
			return c.Regions
		}
	}

	return nil
}

func (m *configurations) ValidateCountryCode(code string) error {
	for _, c := range m.Countries {
		if code == c.CountryCode {
			if !c.Enabled {
				return fmt.Errorf(`not allowed to take "%s" country the country is disabled`, code)
			}

			return nil
		}
	}

	return fmt.Errorf("not found country (country code=%s)", code)
}

func (m *configurations) ValidateRegionCode(code string) error {
	for _, c := range m.Countries {
		for _, r := range c.Regions {
			if code == r {
				if !c.Enabled {
					return fmt.Errorf(`not allowed to take "%s" region because "%s" is disabled`, code, c.CountryCode)
				}

				return nil
			}
		}
	}

	return fmt.Errorf("not found region (region=%s)", code)
}
