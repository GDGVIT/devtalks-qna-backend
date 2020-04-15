package question

import (
	"github.com/rithikjain/LiveQnA/pkg"
	"github.com/rithikjain/LiveQnA/pkg/user"
)

type Service interface {
	CreateQuestion(question *Question) (*Question, error)

	ViewAllQuestions() (*[]Question, error)

	ViewAllQuestionsByUpVotes() (*[]Question, error)

	IncreaseUpVote(questionID float64, user *user.User) (*Question, error)

	DecreaseUpVote(questionID float64, user *user.User) (*Question, error)

	DeleteQuestion(questionID float64) error

	GetUser(userID float64) (*user.User, error)

	GetRepo() Repository
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateQuestion(question *Question) (*Question, error) {
	que, err := s.repo.CreateQuestion(question)
	if err != nil {
		return nil, err
	}
	return que, nil
}

func (s *service) ViewAllQuestions() (*[]Question, error) {
	return s.repo.GetAllQuestions()
}

func (s *service) ViewAllQuestionsByUpVotes() (*[]Question, error) {
	return s.repo.GetAllQuestionsByUpVotes()
}

func (s *service) IncreaseUpVote(questionID float64, user *user.User) (*Question, error) {
	alreadyUpVoted, err := s.repo.HasAlreadyUpVoted(questionID, user.Email)
	if err != nil {
		return nil, err
	}
	if !alreadyUpVoted {
		que, err := s.repo.IncreaseUpVote(questionID, user)
		if err != nil {
			return nil, err
		}
		return que, nil
	} else {
		return nil, pkg.ErrNotAllowed
	}
}

func (s *service) DecreaseUpVote(questionID float64, user *user.User) (*Question, error) {
	alreadyUpVoted, err := s.repo.HasAlreadyUpVoted(questionID, user.Email)
	if err != nil {
		return nil, err
	}
	if alreadyUpVoted {
		que, err := s.repo.DecreaseUpVote(questionID, user)
		if err != nil {
			return nil, err
		}
		return que, nil
	} else {
		return nil, pkg.ErrNotAllowed
	}
}

func (s *service) DeleteQuestion(questionID float64) error {
	err := s.repo.DeleteQuestion(questionID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUser(userID float64) (*user.User, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetRepo() Repository {
	return s.repo
}
