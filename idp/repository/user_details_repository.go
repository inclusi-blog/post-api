package repository

// mockgen -source=repository/user_details_repository.go -destination=mocks/mock_user_details_repository.go -package=mocks
import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/inclusi-blog/gola-utils/mask_util"
	"github.com/jmoiron/sqlx"
	"post-api/idp/models/db"
	"post-api/utils"
	"strings"
)

type UserDetailsRepository interface {
	SaveUserDetails(details db.SaveUserDetails, context context.Context) error
	IsEmailAvailable(email string, ctx context.Context) (bool, error)
	IsUserNameAvailable(username string, ctx context.Context) (bool, error)
	IsUserNameAndEmailAvailable(username, email string, ctx context.Context) (bool, error)
	GenerateUsername(ctx context.Context, email string) (string, error)
	GetUserProfile(email string, ctx context.Context) (db.UserProfile, error)
	GetPassword(email string, ctx context.Context) (string, error)
	UpdateName(ctx context.Context, name string, id uuid.UUID) error
	UpdateUsername(ctx context.Context, username string, id uuid.UUID) error
	UpdateAbout(ctx context.Context, about string, id uuid.UUID) error
	UpdateTwitterURL(ctx context.Context, twitterURL string, id uuid.UUID) error
	UpdateLinkedInURL(ctx context.Context, linkedinURL string, id uuid.UUID) error
	UpdateFacebookURL(ctx context.Context, facebookURL string, id uuid.UUID) error
	UpdateProfileImage(ctx context.Context, imageKey string, id uuid.UUID) error
	UpdatePassword(ctx context.Context, hashedPassword, email string) error
}

type userDetailsRepository struct {
	db *sqlx.DB
}

const (
	SaveUser                  = "insert into users(id, username, email, password, is_active, role_id)values($1,$2,$3,$4,$5, (select id from roles where name = 'User'))"
	UserExistence             = "select count(*) from users where email = $1"
	UsernameExistence         = "select count(*) from users where username = $1"
	UsernameAndEmailExistence = "select count(*) from users where username = $1 or email = $2"
	FetchUserDetails          = "select id, username, email, is_active from users where email = $1"
	FetchUserPassword         = "select password from users where email = $1"
	UpdateAbout               = "update users set about = $1 where id = $2"
	UpdateName                = "update users set name = $1 where id = $2"
	UpdateUsername            = "update users set username = $1 where id = $2"
	UpdateTwitter             = "insert into social_links(id, twitter, user_id)values (uuid_generate_v4(), $1, $2) on conflict (user_id) do update set twitter = $3"
	UpdateLinkedIn            = "insert into social_links(id, linkedin, user_id)values (uuid_generate_v4(), $1, $2) on conflict (user_id) do update set linkedin = $3"
	UpdateFacebook            = "insert into social_links(id, facebook, user_id)values (uuid_generate_v4(), $1, $2) on conflict (user_id) do update set facebook = $3"
	UpdateImage               = "update users set avatar = $1 where id = $2"
	UsernameCount             = "select count(*) as count from users where user = $1"
	UpdatePassword            = "update users set password = $1 where email = $2"
)

func (repository userDetailsRepository) SaveUserDetails(details db.SaveUserDetails, context context.Context) error {
	logger := logging.GetLogger(context)
	log := logger.WithField("key", "UserDetailsRepository").WithField("method", "SaveUserDetails")

	username := details.Username
	log.Infof("Saving user details for user %v", username)

	result, err := repository.db.ExecContext(context, SaveUser, details.ID, username, details.Email, details.Password, details.IsActive)

	if err != nil {
		log.Errorf("Unable to store user details for username %v . Error %v", username, err)
		return err
	}

	affectedRows, _ := result.RowsAffected()
	if affectedRows > 1 {
		log.Errorf("More than one row affected while saving user %v", username)
		return errors.New("more than one row affected")
	}

	log.Infof("Successfully saved user for username %v", username)
	return nil
}

func (repository userDetailsRepository) IsEmailAvailable(email string, ctx context.Context) (bool, error) {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("key", "UserDetailsRepository").WithField("method", "IsEmailAvailable")

	var userCount int

	log.Infof("Fetching user existence on user registration for email %v", email)
	err := repository.db.GetContext(ctx, &userCount, UserExistence, email)

	if err != nil {
		log.Errorf("Error occurred while fetching user existence %v", err)
		return false, err
	}

	if userCount == 1 {
		log.Infof("User already exists in gola for user email %v", email)
		return true, nil
	}

	log.Infof("User not exists in gola for user email %v", email)
	return false, nil
}

func (repository userDetailsRepository) IsUserNameAvailable(username string, ctx context.Context) (bool, error) {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("key", "UserDetailsRepository").WithField("method", "IsUserNameAvailable")

	var userCount int

	log.Infof("Fetching username availability on user registration for username %v", username)
	err := repository.db.GetContext(ctx, &userCount, UsernameExistence, username)

	if err != nil {
		log.Errorf("Error occurred while fetching username availability %v", err)
		return false, err
	}

	if userCount == 1 {
		log.Infof("Username already exists in gola for user username %v", username)
		return true, nil
	}

	log.Infof("Username not exists in gola for user username %v", username)
	return false, nil
}

func (repository userDetailsRepository) IsUserNameAndEmailAvailable(username, email string, ctx context.Context) (bool, error) {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("key", "UserDetailsRepository").WithField("method", "IsUserNameAndEmailAvailable")

	var userCount int

	log.Infof("Fetching username or email availability on user registration for username %v", username)
	err := repository.db.GetContext(ctx, &userCount, UsernameAndEmailExistence, username, email)

	if err != nil {
		log.Errorf("Error occurred while fetching username or email availability %v", err)
		return false, err
	}

	if userCount == 1 {
		log.Infof("Username or email already exists in gola for user username %v", username)
		return true, nil
	}

	log.Infof("Username or email not exists in gola for user username %v", username)
	return false, nil
}

func (repository userDetailsRepository) GetUserProfile(email string, ctx context.Context) (db.UserProfile, error) {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "GetUserProfile")

	maskEmail := mask_util.MaskEmail(ctx, email)
	logger.Infof("Fetching user profile details for the given email %v", maskEmail)

	var userProfileDetails db.UserProfile

	err := repository.db.GetContext(ctx, &userProfileDetails, FetchUserDetails, email)

	if err != nil {
		logger.Errorf("Error occurred while fetching user profile data for email %v .%v", maskEmail, err)
		return db.UserProfile{}, err
	}

	logger.Infof("User found for email %v", maskEmail)

	return userProfileDetails, nil
}

func (repository userDetailsRepository) GetPassword(email string, ctx context.Context) (string, error) {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "GetPassword")
	maskEmail := mask_util.MaskEmail(ctx, email)
	logger.Infof("Fetching user pass for authentication %v", maskEmail)

	var plainTextPassword string
	err := repository.db.Get(&plainTextPassword, FetchUserPassword, email)

	if err != nil {
		logger.Errorf("Unable to fetch user credentials for email %v .%v", maskEmail, err)
		return "", err
	}

	logger.Infof("Successfully fetched user credentials for emai %v", maskEmail)

	return plainTextPassword, nil
}

func (repository userDetailsRepository) UpdateName(ctx context.Context, name string, id uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "UpdateUserDetails")
	logger.Info("updating user details for user")

	_, err := repository.db.ExecContext(ctx, UpdateName, name, id)
	if err != nil {
		logger.Error("unable to update name for user %v", id)
		return err
	}

	return nil
}

func (repository userDetailsRepository) UpdateUsername(ctx context.Context, username string, id uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "UpdateUserDetails")
	logger.Info("updating user details for user")

	_, err := repository.db.ExecContext(ctx, UpdateUsername, username, id)
	if err != nil {
		logger.Error("unable to update username for user %v", id)
		return err
	}

	return nil
}

func (repository userDetailsRepository) UpdateAbout(ctx context.Context, about string, id uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "UpdateUserDetails")
	logger.Info("updating user details for user")

	_, err := repository.db.ExecContext(ctx, UpdateAbout, about, id)
	if err != nil {
		logger.Error("unable to update about for user %v", id)
		return err
	}

	return nil
}

func (repository userDetailsRepository) UpdateTwitterURL(ctx context.Context, twitterURL string, id uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "UpdateTwitterURL")
	logger.Info("updating user twitter url")

	_, err := repository.db.ExecContext(ctx, UpdateTwitter, twitterURL, id, twitterURL)
	if err != nil {
		logger.Error("unable to update about for user %v", id)
		return err
	}

	return nil
}

func (repository userDetailsRepository) UpdateFacebookURL(ctx context.Context, facebookURL string, id uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "UpdateTwitterURL")
	logger.Info("updating user twitter url")

	_, err := repository.db.ExecContext(ctx, UpdateFacebook, facebookURL, id, facebookURL)
	if err != nil {
		logger.Error("unable to update about for user %v", id)
		return err
	}

	return nil
}

func (repository userDetailsRepository) UpdateLinkedInURL(ctx context.Context, linkedinURL string, id uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "UpdateTwitterURL")
	logger.Info("updating user twitter url")

	_, err := repository.db.ExecContext(ctx, UpdateLinkedIn, linkedinURL, id, linkedinURL)
	if err != nil {
		logger.Error("unable to update about for user %v", id)
		return err
	}

	return nil
}

func (repository userDetailsRepository) UpdateProfileImage(ctx context.Context, imageKey string, id uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "UpdateTwitterURL")
	logger.Info("updating user image url")

	_, err := repository.db.ExecContext(ctx, UpdateImage, imageKey, id)
	if err != nil {
		logger.Error("unable to update avatar for user %v", id)
		return err
	}

	return nil
}

func (repository userDetailsRepository) GenerateUsername(ctx context.Context, email string) (string, error) {
	emailSlug := strings.Split(email, "@")[0]
	var counter int64
	err := repository.db.GetContext(ctx, &counter, UsernameCount, emailSlug)
	if err != nil {
		return "", err
	}
	num := utils.GenRandNum(100, 999)
	return fmt.Sprintf("%s%d%d", emailSlug, num, counter+1), nil
}

func (repository userDetailsRepository) UpdatePassword(ctx context.Context, hashedPassword string, email string) error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsRepository").WithField("method", "UpdatePassword")
	_, err := repository.db.ExecContext(ctx, UpdatePassword, hashedPassword, email)
	if err != nil {
		maskEmail := mask_util.MaskEmail(ctx, email)
		logger.Errorf("unable to update password for user %v .Error %v", maskEmail, err)
		return err
	}

	return nil
}

func NewUserDetailsRepository(db *sqlx.DB) UserDetailsRepository {
	return userDetailsRepository{
		db: db,
	}
}
