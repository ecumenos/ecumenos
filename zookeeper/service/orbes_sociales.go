package service

import (
	"context"

	models "github.com/ecumenos/ecumenos/models/zookeeper"
)

func (s *Service) GetOrbisSociusCountries() []string {
	return s.settings.GetCountries(true)
}

func (s *Service) GetOrbisSociusRegions(countryCode string) []string {
	return s.settings.GetRegionsByCountryCode(countryCode)
}

func (s *Service) MakeCreateOrbisSociusLaunchRequest(ctx context.Context, ownerID int64, region, name, desc, url string) (*models.OrbisSociusLaunchRequest, error) {
	if err := s.settings.ValidateRegionCode(region); err != nil {
		return nil, err
	}

	return s.repo.InsertOrbisSociusLaunchRequest(ctx, ownerID, region, name, desc, url, models.PendingOrbisSociusLaunchRequest)
}

func (s *Service) GetOrbisSociusLanguages() []string {
	return s.settings.GetLanguages(true)
}
