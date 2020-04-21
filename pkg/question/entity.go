package question

import "github.com/jinzhu/gorm"

type Question struct {
	gorm.Model
	Question       string `json:"question"`
	CreatedByEmail string `json:"created_by_email"`
	UpVotes        uint   `json:"up_votes"`
}

type UpVoteDetail struct {
	gorm.Model
	Email      string
	QuestionID uint
}
