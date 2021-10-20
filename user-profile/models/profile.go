package models

import "github.com/google/uuid"

type Profile struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Name       *string   `json:"name" db:"name"`
	Username   string    `json:"username" db:"username"`
	Email      string    `json:"email" db:"email"`
	About      *string   `json:"about" db:"about"`
	ProfilePic *string   `json:"profile_pic" db:"profile_pic"`
	SocialLinks
}

type SocialLinks struct {
	FacebookURL *string `json:"facebook_url" db:"facebook_url"`
	LinkedInURL *string `json:"linked_in_url" db:"linked_in_url"`
	TwitterURL  *string `json:"twitter_url" db:"twitter_url"`
}
