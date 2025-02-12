package mock

import "github.com/Noviiich/golang-url-shortener/internal/core/domain"

var MockLinkData []domain.Link = []domain.Link{
	{Id: "testid1", OriginalURL: "https://example.com/link1"},
	{Id: "testid2", OriginalURL: "https://example.com/link2"},
	{Id: "testid3", OriginalURL: "https://example.com/link3"},
}

var MockStatsData []domain.Stats = []domain.Stats{
	{Id: "testid1", Platform: domain.PlatformUnknown},
	{Id: "testid2", Platform: domain.PlatformInstagram},
	{Id: "testid3", Platform: domain.PlatformTwitter},
}
