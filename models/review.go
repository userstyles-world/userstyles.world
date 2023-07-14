package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"userstyles.world/modules/database"
)

var (
	errorReviewOutOfRange  = errors.New("rating is out of range")
	errorReviewLongComment = errors.New("comment is too long")
	errorReviewEmptyFields = errors.New("fields cannot be empty")
)

type Review struct {
	gorm.Model
	Rating  int
	Comment string

	User   User
	UserID uint

	Style   Style
	StyleID uint
}

func (Review) FindAllForStyle(id any) (q []Review, err error) {
	err = db().
		// Preload(clause.Associations).
		Preload("User").
		Model(modelReview).
		Order("id DESC").
		Find(&q, "style_id = ? ", id).
		Error

	if err != nil {
		return nil, err
	}

	return q, nil
}

func (r *Review) CreateForStyle() error {
	return db().Create(r).Error
}

func (r *Review) FindLastFromUser(styleID, userID any) error {
	return db().
		Preload("User").
		Model(modelReview).
		Order("id DESC").
		Last(&r, "style_id = ? and user_id = ?", styleID, userID).
		Error
}

// Validate verifies user-generated content.
func (r *Review) Validate() error {
	switch {
	case r.Rating < 0 || r.Rating > 5:
		return errorReviewOutOfRange
	case len(r.Comment) > 500:
		return errorReviewLongComment
	case r.Rating == 0 && len(r.Comment) == 0:
		return errorReviewEmptyFields
	default:
		return nil
	}
}

// Permalink returns a link to the review page.
func (r *Review) Permalink() string {
	return fmt.Sprintf("/styles/%d/reviews/%d", r.StyleID, r.ID)
}

// UpdateFromUser updates a review from its author.
func (r *Review) UpdateFromUser() error {
	return database.Conn.
		Select("updated_at", "rating", "comment").
		Where("id = ? AND user_id = ?", r.ID, r.UserID).
		Updates(r).
		Error
}

// GetReview returns a specific review, or an error if the review doesn't exist.
func GetReview(id int) (*Review, error) {
	var r Review
	err := database.Conn.Preload("User").First(&r, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// DeleteReviewFromUser deletes a review from user, otherwise returns an error.
func DeleteReviewFromUser(id, uid int) error {
	return database.Conn.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&Review{}).Error
		if err != nil {
			return err
		}

		err = tx.Where("id = ?", id).Delete(&Notification{}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// MatchReviewUser returns whether or not current user matches review's user.
func MatchReviewUser(id, uid int) bool {
	var i int64
	err := database.Conn.Debug().
		Model(modelReview).
		Where("id = ? AND user_id = ?", id, uid).
		Count(&i).
		Error

	return err == nil && i > 0
}

// NewReview is a helper for creating a new review.
func NewReview(uid, sid uint, rating, comment string) *Review {
	i, err := strconv.Atoi(rating)
	if err != nil {
		i = -1
	}

	return &Review{
		UserID:  uid,
		StyleID: sid,
		Comment: strings.TrimSpace(comment),
		Rating:  i,
	}
}

// NewReviewUpdate is a helper for updating a review.
func NewReviewUpdate(uid, sid, rid uint, rating, comment string) *Review {
	i, err := strconv.Atoi(rating)
	if err != nil {
		i = -1
	}

	return &Review{
		UserID:  uid,
		StyleID: sid,
		Model:   gorm.Model{ID: rid},
		Comment: strings.TrimSpace(comment),
		Rating:  i,
	}
}
