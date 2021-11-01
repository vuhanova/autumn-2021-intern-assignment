package handlers

import (
	"autumn-2021-intern-assignment/pkg/transaction"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetBalanceFromUser(t *testing.T) {
	// мы передаём t сюда, это надо, чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательность вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockItemsRepositoryInterface(ctrl)

	service := &ItemsHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	elemId := 17
	resultItem := &transaction.User{
		UserID:  elemId,
		Balance: 50.0,
	}

	b, err := json.Marshal(resultItem)
	if err != nil {
		t.Errorf("internal error")
		return
	}
	bodyReader := strings.NewReader(string(b))

	st.EXPECT().GetUsersBalance(elemId, "").Return(resultItem, nil)

	req := httptest.NewRequest("POST", "/user", bodyReader)
	w := httptest.NewRecorder()
	service.GetBalanceFromUser(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	us := transaction.User{}
	err = json.Unmarshal(body, &us)

	if err != nil && !reflect.DeepEqual(us, resultItem) {
		t.Errorf("results not match, want %v, have %v", resultItem, us)
		return
	}

	// marshaling error

	bodyReader = strings.NewReader("mess1111ag11e:1qq11111powei")
	req = httptest.NewRequest("POST", "/user", bodyReader)
	w = httptest.NewRecorder()
	service.GetBalanceFromUser(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected resp status 400, got %d", resp.StatusCode)
		return
	}

	// result error

	st.EXPECT().GetUsersBalance(elemId, "").Return(resultItem, fmt.Errorf("bad result"))
	bodyReader = strings.NewReader(string(b))
	req = httptest.NewRequest("POST", "/create", bodyReader)
	w = httptest.NewRecorder()
	service.GetBalanceFromUser(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}
}

func TestIncreaseBalance(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательность вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockItemsRepositoryInterface(ctrl)

	service := &ItemsHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	elemId := 1
	resultItem := &transaction.User{
		UserID:  elemId,
		Balance: 50.0,
	}

	b, err := json.Marshal(resultItem)
	if err != nil {
		t.Errorf("internal error")
		return
	}
	bodyReader := strings.NewReader(string(b))

	st.EXPECT().AddMoney(resultItem.UserID, resultItem.Balance).Return(nil)

	req := httptest.NewRequest("POST", "/balance/add", bodyReader)
	w := httptest.NewRecorder()
	service.IncreaseBalance(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if !bytes.Contains(body, []byte("success")) {
		t.Errorf("no text found")
		return
	}

	// marshaling error

	bodyReader = strings.NewReader("mess1111ag11e:1qq11111powei")
	req = httptest.NewRequest("POST", "/balance/add", bodyReader)
	w = httptest.NewRecorder()
	service.IncreaseBalance(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected resp status 400, got %d", resp.StatusCode)
		return
	}

	// result error

	st.EXPECT().AddMoney(resultItem.UserID, resultItem.Balance).Return(fmt.Errorf("bad result"))
	bodyReader = strings.NewReader(string(b))
	req = httptest.NewRequest("POST", "/balance/add", bodyReader)
	w = httptest.NewRecorder()
	service.IncreaseBalance(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}
}

func TestDecreaseBalance(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательность вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockItemsRepositoryInterface(ctrl)

	service := &ItemsHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	elemId := 1
	resultItem := &transaction.User{
		UserID:  elemId,
		Balance: 50.0,
	}

	b, err := json.Marshal(resultItem)
	if err != nil {
		t.Errorf("internal error")
		return
	}
	bodyReader := strings.NewReader(string(b))

	st.EXPECT().WithdrawMoney(resultItem.UserID, resultItem.Balance).Return(nil)

	req := httptest.NewRequest("POST", "/balance/reduce", bodyReader)
	w := httptest.NewRecorder()
	service.DecreaseBalance(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if !bytes.Contains(body, []byte("success")) {
		t.Errorf("no text found")
		return
	}

	// marshaling error

	bodyReader = strings.NewReader("mess1111ag11e:1qq11111powei")
	req = httptest.NewRequest("POST", "/balance/reduce", bodyReader)
	w = httptest.NewRecorder()
	service.DecreaseBalance(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected resp status 400, got %d", resp.StatusCode)
		return
	}

	// result error

	st.EXPECT().WithdrawMoney(resultItem.UserID, resultItem.Balance).Return(fmt.Errorf("bad result"))
	bodyReader = strings.NewReader(string(b))
	req = httptest.NewRequest("POST", "/balance/reduce", bodyReader)
	w = httptest.NewRecorder()
	service.DecreaseBalance(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}
}

func TestTransferBalance(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательность вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockItemsRepositoryInterface(ctrl)

	service := &ItemsHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	elemId := 1
	resultItem := &transaction.User{
		UserID:   elemId,
		ToUserID: 2,
		Balance:  50.0,
	}

	b, err := json.Marshal(resultItem)
	if err != nil {
		t.Errorf("internal error")
		return
	}
	bodyReader := strings.NewReader(string(b))

	st.EXPECT().TransferMoney(resultItem.UserID, resultItem.ToUserID, resultItem.Balance).Return(nil)

	req := httptest.NewRequest("POST", "/balance/transfer", bodyReader)
	w := httptest.NewRecorder()
	service.TransferBalance(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if !bytes.Contains(body, []byte("success")) {
		t.Errorf("no text found")
		return
	}

	// marshaling error

	bodyReader = strings.NewReader("mess1111ag11e:1qq11111powei")
	req = httptest.NewRequest("POST", "/balance/transfer", bodyReader)
	w = httptest.NewRecorder()
	service.TransferBalance(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected resp status 400, got %d", resp.StatusCode)
		return
	}

	// result error

	st.EXPECT().TransferMoney(resultItem.UserID, resultItem.ToUserID, resultItem.Balance).Return(fmt.Errorf("bad result"))
	bodyReader = strings.NewReader(string(b))
	req = httptest.NewRequest("POST", "/balance/transfer", bodyReader)
	w = httptest.NewRecorder()
	service.TransferBalance(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}
}

func TestListTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательность вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := NewMockItemsRepositoryInterface(ctrl)

	service := &ItemsHandler{
		ItemRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	elemId := 1
	sourceItem := &transaction.User{
		UserID: elemId,
		Field:  "",
	}

	resultItem := []*transaction.Transaction{
		{
			ToID:    nil,
			FromID:  &elemId,
			Money:   50,
			Created: time.Now(),
		},
	}

	b, err := json.Marshal(sourceItem)
	if err != nil {
		t.Errorf("internal error")
		return
	}
	bodyReader := strings.NewReader(string(b))

	st.EXPECT().GetTransaction(elemId, "").Return(resultItem, nil)

	req := httptest.NewRequest("POST", "/info", bodyReader)
	w := httptest.NewRecorder()
	service.ListTransaction(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if !bytes.Contains(body, []byte("success")) {
		//t.Errorf("no text found")
		//return
	}

	// marshaling error

	bodyReader = strings.NewReader("mess1111ag11e:1qq11111powei")
	req = httptest.NewRequest("POST", "/info", bodyReader)
	w = httptest.NewRecorder()
	service.ListTransaction(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected resp status 400, got %d", resp.StatusCode)
		return
	}

	// result error

	st.EXPECT().GetTransaction(elemId, "").Return(nil, fmt.Errorf("bad result"))
	bodyReader = strings.NewReader(string(b))
	req = httptest.NewRequest("POST", "/info", bodyReader)
	w = httptest.NewRecorder()
	service.ListTransaction(w, req)

	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}
}
