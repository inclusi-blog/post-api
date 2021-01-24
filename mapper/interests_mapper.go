package mapper

//go:generate mockgen -source=interests_mapper.go -destination=./../mocks/mock_interests_mapper.go -package=mocks
import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/models/db"
	"post-api/models/response"
)

type InterestsMapper interface {
	MapUserFollowedInterest(ctx context.Context, dbCategoriesAndInterests []db.CategoryAndInterest) []response.CategoryAndInterest
}

type interestsMapper struct {
}

func (mapper interestsMapper) MapUserFollowedInterest(ctx context.Context, dbCategoriesAndInterests []db.CategoryAndInterest) []response.CategoryAndInterest {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsMapper").WithField("method", "MapUserFollowedInterest")
	logger.Info("Mapping db categories and interests with user followed interests")

	var exploreResponse []response.CategoryAndInterest
	for _, dbCategory := range dbCategoriesAndInterests {
		var interests []response.InterestWithIcon
		for _, interest := range dbCategory.Interests {
			interests = append(interests, response.InterestWithIcon{
				Name:             interest.Name,
				Image:            interest.Image,
				IsFollowedByUser: false,
			})
		}
		exploreResponse = append(exploreResponse, response.CategoryAndInterest{
			Category:  dbCategory.Category,
			Interests: interests,
		})
	}

	logger.Info("Returning explore response with user followed interest status")
	return exploreResponse
}

func NewInterestsMapper() InterestsMapper {
	return interestsMapper{}
}
