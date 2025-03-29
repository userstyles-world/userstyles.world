package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"userstyles.world/modules/database"
	"userstyles.world/modules/util"
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

// GetReviewsForUser tries to return reviews sorted from most to least recent.
func GetReviewsForUser(db *gorm.DB, id int) (r []Review, err error) {
	err = db.Preload("Style").Order("id DESC").Find(&r, "user_id = ?", id).Error
	return r, err
}

// DeleteReviewsForStyle tries to delete all reviews for a specific style.
func DeleteReviewsForStyle(db *gorm.DB, id int) (err error) {
	return db.Delete(&Review{}, "style_id = ?", id).Error
}

func FindAllForStyle(id any) (q []Review, err error) {
	err = db().
		Preload(clause.Associations).
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
	slug := util.Slug(r.Style.Name)
	return fmt.Sprintf("/styles/%d-%s/reviews/%d", r.StyleID, slug, r.ID)
}

// UpdateFromUser updates a review from its author.
func (r *Review) UpdateFromUser() error {
	return database.Conn.
		Select("updated_at", "rating", "comment").
		Where("id = ? AND user_id = ?", r.ID, r.UserID).
		Updates(r).
		Error
}

// GetReviewForStyle tries to return a specific review for a specific style.
func GetReviewForStyle(db *gorm.DB, rid, sid int) (r *Review, err error) {
	err = db.
		Preload(clause.Associations).
		Where("id = ? AND style_id = ?", rid, sid).
		First(&r).Error
	if err != nil {
		return nil, err
	}

	return r, nil
}

// DeleteReviewFromUser deletes a review from user, otherwise returns an error.
func DeleteReviewFromUser(id, uid int) error {
	return database.Conn.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id = ? AND user_id = ?", id, uid).Delete(&Review{}).Error
		if err != nil {
			return err
		}

		err = tx.Where("review_id = ?", id).Delete(&Notification{}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// MatchReviewForStyle returns whether or not a review belongs to a style.
func MatchReviewForStyle(db *gorm.DB, rid, sid int) bool {
	var i int64
	err := db.
		Model(modelReview).
		Where("id = ? AND style_id = ?", rid, sid).
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

// Prepare is a helper for updating a review.
func (r *Review) Prepare(rating, comment string) {
	i, err := strconv.Atoi(rating)
	if err != nil {
		i = -1
	}

	r.Rating = i
	r.Comment = strings.TrimSpace(comment)
}
