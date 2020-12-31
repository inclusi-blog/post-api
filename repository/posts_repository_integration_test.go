package repository

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"post-api/dbhelper"
	"post-api/models"
	"post-api/models/db"
	"post-api/repository/helper"
	"testing"
)

type PostsRepositoryIntegrationTest struct {
	suite.Suite
	db              neo4j.Session
	driver          neo4j.Driver
	adminDb         neo4j.Session
	adminDriver     neo4j.Driver
	goContext       context.Context
	postsRepository PostsRepository
	dbHelper        helper.DbHelper
}

func (suite *PostsRepositoryIntegrationTest) SetupTest() {
	err := godotenv.Load("../docker-compose-test.env")
	suite.Nil(err)
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	logger := logging.GetLogger(context.Background())
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }
	suite.driver, err = neo4j.NewDriver(dbhelper.BuildConnectionString(), neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), ""), configForNeo4j40)
	suite.Nil(err)
	suite.adminDriver, err = neo4j.NewDriver(dbhelper.BuildConnectionString(), neo4j.BasicAuth(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"), ""), configForNeo4j40)
	suite.Nil(err)
	suite.NotNil(suite.adminDriver)
	suite.NotNil(suite.driver)

	logger.Info("logging")
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: os.Getenv("DB_SERVICE_NAME")}
	suite.db, err = suite.driver.NewSession(sessionConfig)
	suite.Nil(err)
	suite.adminDb, err = suite.adminDriver.NewSession(sessionConfig)
	suite.Nil(err)

	suite.postsRepository = NewPostsRepository(suite.db)
	suite.dbHelper = helper.NewDbHelper(suite.adminDb)
	suite.insertInterestEntries()
	suite.createSampleUser("some-user")
}

func (suite *PostsRepositoryIntegrationTest) TearDownTest() {
	suite.ClearPostsData()
	err := suite.driver.Close()
	suite.Nil(err)
	err = suite.adminDriver.Close()
	suite.Nil(err)
	err = suite.db.Close()
	suite.Nil(err)
	err = suite.adminDb.Close()
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) ClearPostsData() {
	e := suite.dbHelper.ClearAll()
	if e != nil {
		assert.Error(suite.T(), e)
	}
}

func TestPostsRepositoryIntegrationTest(t *testing.T) {
	suite.Run(t, new(PostsRepositoryIntegrationTest))
}

func (suite *PostsRepositoryIntegrationTest) TestCreatePost_WhenSuccessfullyStoredInDB() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)

	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCreatePost_WhenSamePostIsSavedDBShouldReturnError() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)
	err = suite.postsRepository.CreatePost(suite.goContext, post)
	suite.NotNil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestLikePost_WhenUserLikes() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestLikePost_WhenThereIsNoPost() {
	err := suite.postsRepository.LikePost("1q2w3e4r5t6y", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q2w3e4r5t6y")
	suite.NotNil(err)
	suite.Equal("no results found", err.Error())
	suite.False(isLiked)
}

func (suite *PostsRepositoryIntegrationTest) TestIsPostLikedByPerson_WhenThereArePostAndLikes() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestIsPostLikedByPerson_WhenThereIsNoPostAndLikes() {
	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q2w3e4r5t6y")
	suite.NotNil(err)
	suite.Equal("no results found", err.Error())
	suite.False(isLiked)
}

func (suite *PostsRepositoryIntegrationTest) TestUnlikePost_WhenUserLikes() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.UnlikePost(suite.goContext, "some-user", "1q323e4r4r43")
	suite.Nil(err)

	isLiked, err = suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.False(isLiked)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestUnlikePost_WhenThereIsNoPost() {
	err := suite.postsRepository.UnlikePost(suite.goContext, "some-user", "1q2w3e4r5t6y")
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q2w3e4r5t6y")
	suite.NotNil(err)
	suite.Equal("no results found", err.Error())
	suite.False(isLiked)
}

func (suite *PostsRepositoryIntegrationTest) TestCommentPost_WhenUserCommentOnAPost() {
	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.CommentPost(suite.goContext, "third-user", "this is awesome", "1q323e4r4r43")
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCommentPost_WhenThereIsNoPost() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.CommentPost(suite.goContext, "third-user", "this is awesome", "1q323e4r4r43")
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCommentPost_WhenThereIsUserButNoPost() {
	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")

	err := suite.postsRepository.CommentPost(suite.goContext, "third-user", "this is awesome", "1q323e4r4r43")
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCommentPost_WhenThereIsNoUserAndNoPost() {
	err := suite.postsRepository.CommentPost(suite.goContext, "third-user", "this is awesome", "1q323e4r4r43")
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestGetLikesCountByPostID_WhenThereIsAPost() {
	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "second-user", suite.goContext)
	suite.Nil(err)

	isLiked, err = suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	likesCount, err := suite.postsRepository.GetLikesCountByPostID(suite.goContext, "1q323e4r4r43")
	suite.Nil(err)
	suite.Equal(int64(2), likesCount)
}

func (suite *PostsRepositoryIntegrationTest) TestGetLikesCountByPostID_WhenThereIsNoPost() {
	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	err := suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "second-user", suite.goContext)
	suite.Nil(err)

	likesCount, err := suite.postsRepository.GetLikesCountByPostID(suite.goContext, "1q323e4r4r43")
	suite.Nil(err)
	suite.Equal(int64(0), likesCount)
}

func (suite *PostsRepositoryIntegrationTest) TestGetLikesCountByPostID_WhenThereIsNoUserWhenLikingPost() {
	suite.createSampleUser("third-user")
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "second-user", suite.goContext)
	suite.Nil(err)

	isLiked, err = suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	likesCount, err := suite.postsRepository.GetLikesCountByPostID(suite.goContext, "1q323e4r4r43")
	suite.Nil(err)
	suite.Equal(int64(1), likesCount)
}

func (suite *PostsRepositoryIntegrationTest) insertInterestEntries() {
	interests := []string{
		"CREATE (interest:Category:Interest{name: 'Art'})",
		"CREATE (interest:Category:Interest{name: 'Entertainment'})",
		"CREATE (interest:Category:Interest{name: 'Culture'})",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Poem'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Short stories'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Books'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Literature'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Grammar'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Comics'})-[:BELONGS_TO]->(art)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Movies'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Series'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Anime'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Cartoon'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Animation'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Cooking'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Food'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Agriculture'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Festival'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Language'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Philosophy'})-[:BELONGS_TO]->(culture)",
	}
	for _, query := range interests {
		_, err := suite.adminDb.Run(query, nil)
		suite.Nil(err)
	}
}

func (suite *PostsRepositoryIntegrationTest) createSampleUser(userId string) {
	_, err := suite.adminDb.Run("CREATE (person:Person{ userId: $userId})", map[string]interface{}{
		"userId": userId,
	})
	suite.Nil(err)
}
