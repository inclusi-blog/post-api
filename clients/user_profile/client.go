package user_profile

//go:generate mockgen -source=client.go -destination=./../../mocks/mock_user_profile_client.go -package=mocks
import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/http/request"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/configuration"
	"post-api/utils"
)

type Client interface {
	FetchUserFollowingInterests(ctx context.Context) ([]string, *golaerror.Error)
}

type client struct {
	requestBuilder request.HttpRequestBuilder
	config         *configuration.ConfigData
}

func (client client) FetchUserFollowingInterests(ctx context.Context) ([]string, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "Client").WithField("method", "FetchUserFollowingInterests")
	logger.Info("Calling user profile client to get user following interests")

	config := client.config
	fetchUserFollowingInterestsURL := config.UserProfileBaseUrl + config.FetchUserFollowedInterests

	var userFollowingInterests []string

	err := client.requestBuilder.NewRequestWithContext(ctx).ResponseAs(&userFollowingInterests).Get(fetchUserFollowingInterestsURL)

	if err != nil {
		logger.Errorf("Error occurred while fetching user following interest. Error: %v", err)
		golaError := utils.GetGolaError(err)
		if golaError.ErrorCode == "ERR_USER_PROFILE_NO_INTERESTS_FOLLOWED" {
			logger.Error("user is not following any interests")
			return nil, nil
		}
		return nil, golaError
	}

	return userFollowingInterests, nil
}

func NewClient(httpRequestBuilder request.HttpRequestBuilder, data *configuration.ConfigData) Client {
	return client{requestBuilder: httpRequestBuilder, config: data}
}
