package mock_repository

import (
	"testing"

	"github.com/golang/mock/gomock"
	repository "github.com/traPtitech/traPortfolio/usecases/repository"
)

func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expected := []*repository.KnoQEvent{}
	knoQRepository := NewMockKnoqRepository(mockCtrl)
	knoQRepository.EXPECT().GetAll().Return(expected, nil)

	e, err := knoQRepository.GetAll()
	t.Log("result:", e)
	t.Log("err:", err)
}
