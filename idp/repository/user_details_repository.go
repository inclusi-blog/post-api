package repository

// mockgen -source=repository/user_details_repository.go -destination=mocks/mock_user_details_repository.go -package=mocks
import (
	"context"
	"errors"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/mask_util"
	"github.com/jmoiron/sqlx"
	"post-api/idp/models/db"
)

type UserDetailsRepository interface {
	SaveUserDetails(details db.SaveUserDetails, context context.Context) error
	IsEmailAvailable(email string, ctx context.Context) (bool, error)
	IsUserNameAvailable(username string, ctx context.Context) (bool, error)
	IsUserNameAndEmailAvailable(username, email string, ctx context.Context) (bool, error)
	GetUserProfile(email string, ctx context.Context) (db.UserProfile, error)
	GetPassword(email string, ctx context.Context) (string, error)
}

type userDetailsRepository struct {
	db *sqlx.DB
}

const (
	SaveUser                  = "INSERT INTO user_details(UUID, USERNAME, EMAIL, PASSWD, IS_ACTIVE)VALUES($1,$2,$3,$4,$5)"
	UserExistence             = "SELECT COUNT(*) FROM user_details WHERE EMAIL = $1"
	UsernameExistence         = "SELECT COUNT(*) FROM user_details WHERE USERNAME = $1"
	UsernameAndEmailExistence = "SELECT COUNT(*) FROM user_details WHERE USERNAME = $1 OR EMAIL = $2"
	FetchUserDetails          = "SELECT UUID, USERNAME, EMAIL, IS_ACTIVE FROM user_details WHERE EMAIL = $1"
	FetchUserPassword         = "SELECT PASSWD FROM user_details WHERE EMAIL = $1"
)

func (repository userDetailsRepository) SaveUserDetails(details db.SaveUserDetails, context context.Context) error {
	logger := logging.GetLogger(context)
	log := logger.WithField("key", "UserDetailsRepository").WithField("method", "SaveUserDetails")

	username := details.Username
	log.Infof("Saving user details for user %v", username)

	result, err := repository.db.ExecContext(context, SaveUser, details.UUID, username, details.Email, details.Password, details.IsActive)

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

func NewUserDetailsRepository(db *sqlx.DB) UserDetailsRepository {
	return userDetailsRepository{
		db: db,
	}
}
