package question

import (
	"github.com/jinzhu/gorm"
	"github.com/rithikjain/LiveQnA/pkg"
	"github.com/rithikjain/LiveQnA/pkg/user"
)

type Repository interface {
	CreateQuestion(question *Question) (*Question, error)

	GetAllQuestions() (*[]Question, error)

	IncreaseUpVote(questionID float64, user *user.User) (*Question, error)

	DecreaseUpVote(questionID float64, user *user.User) (*Question, error)

	DeleteQuestion(questionID float64) error

	GetUser(userID float64) (*user.User, error)
}

type repo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}

func (r *repo) CreateQuestion(question *Question) (*Question, error) {
	result := r.DB.Create(question)
	if result.Error != nil {
		return nil, pkg.ErrDatabase
	}
	return question, nil
}

func (r *repo) GetAllQuestions() (*[]Question, error) {
	var questions []Question
	err := r.DB.Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return &questions, nil
}

func (r *repo) IncreaseUpVote(questionID float64, user *user.User) (*Question, error) {
	question := &Question{}
	result := r.DB.Where("ID = ?", questionID).First(question)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, pkg.ErrNotFound
	}
	if question.CreatedByEmail == user.Email {
		return nil, pkg.ErrNotAllowed
	} else {
		question.UpVotes += 1
		result = r.DB.Save(question)
		if result.Error != nil {
			return nil, pkg.ErrDatabase
		}
		return question, nil
	}
}

func (r *repo) DecreaseUpVote(questionID float64, user *user.User) (*Question, error) {
	question := &Question{}
	result := r.DB.Where("id = ?", questionID).First(question)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, pkg.ErrNotFound
	}
	if question.CreatedByEmail == user.Email {
		return nil, pkg.ErrNotAllowed
	} else {
		if question.UpVotes != 0 {
			question.UpVotes -= 1
			result = r.DB.Save(question)
			if result.Error != nil {
				return nil, pkg.ErrDatabase
			}
		}
		return question, nil
	}
}

func (r *repo) DeleteQuestion(questionID float64) error {
	question := &Question{}
	result := r.DB.Where("id = ?", questionID).First(question)
	if result.Error == gorm.ErrRecordNotFound {
		return pkg.ErrNotFound
	}
	result = r.DB.Delete(question)
	if result.Error != nil {
		return pkg.ErrDatabase
	}
	return nil
}

func (r *repo) GetUser(userID float64) (*user.User, error) {
	user := &user.User{}
	result := r.DB.Where("id = ?", userID).First(user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, pkg.ErrNotFound
	}
	return user, nil
}
