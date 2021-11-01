package handlers

import (
	"autumn-2021-intern-assignment/pkg/transaction"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type ItemsRepositoryInterface interface {
	GetUsersBalance(userID int, currency string) (*transaction.User, error)
	AddMoney(userID int, money float64) error
	WithdrawMoney(userID int, money float64) error
	TransferMoney(fromUserID int, toUserID int, money float64) error
	GetTransaction(userID int, orderBy string) ([]*transaction.Transaction, error)
}

type ItemsHandler struct {
	ItemRepo ItemsRepositoryInterface
	Logger   *zap.SugaredLogger
}

//находять в папке pkg/handlers
//mockgen -source=items.go -destination=items_mock.go -package=handlers ItemRepositoryInterface

func sendData(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger, data interface{}) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(dataJSON)
	logger.Infow("New request",
		"method", r.Method,
		"remote_addr", r.RemoteAddr,
		"url", r.URL.Path,
		"time", time.Now().Format(time.RFC3339),
	)
}

func receiveData(r *http.Request) (*transaction.User, int, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer r.Body.Close()

	userCurr := &transaction.User{}
	err = json.Unmarshal(b, &userCurr)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	//userCurr.Currency = r.FormValue("currency")

	return userCurr, http.StatusOK, nil
}

func sendSuccessStatus(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	status := make(map[string]string, 1)
	status["status"] = "success"
	sendData(w, r, logger, status)
}

func sendError(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger, errCurr error, status int) {
	data := make(map[string]string, 1)
	data["error"] = errCurr.Error()

	dataJSON, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Now(),
			"error", err.Error(),
		)
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	http.Error(w, string(dataJSON), status)
	logger.Errorf("New request",
		"method", r.Method,
		"remote_addr", r.RemoteAddr,
		"url", r.URL.Path,
		"time", time.Now(),
		"error", errCurr.Error(),
	)
}

func (h ItemsHandler) GetBalanceFromUser(w http.ResponseWriter, r *http.Request) {
	userCurr, status, err := receiveData(r)
	if err != nil {
		sendError(w, r, h.Logger, err, status)
		return
	}

	userCurr.Currency = r.FormValue("currency")

	tx, err := h.ItemRepo.GetUsersBalance(userCurr.UserID, userCurr.Currency)
	if err != nil {
		sendError(w, r, h.Logger, err, http.StatusInternalServerError)
		return
	}

	sendData(w, r, h.Logger, tx)
}

func (h ItemsHandler) IncreaseBalance(w http.ResponseWriter, r *http.Request) {
	userCurr, status, err := receiveData(r)
	if err != nil {
		sendError(w, r, h.Logger, err, status)
		return
	}

	err = h.ItemRepo.AddMoney(userCurr.UserID, userCurr.Balance)
	if err != nil {
		sendError(w, r, h.Logger, err, http.StatusInternalServerError)
		return
	}

	sendSuccessStatus(w, r, h.Logger)
}

func (h ItemsHandler) DecreaseBalance(w http.ResponseWriter, r *http.Request) {
	userCurr, status, err := receiveData(r)
	if err != nil {
		sendError(w, r, h.Logger, err, status)
		return
	}

	err = h.ItemRepo.WithdrawMoney(userCurr.UserID, userCurr.Balance)
	if err != nil {
		sendError(w, r, h.Logger, err, http.StatusInternalServerError)
		return
	}

	sendSuccessStatus(w, r, h.Logger)
}

func (h ItemsHandler) TransferBalance(w http.ResponseWriter, r *http.Request) {
	userCurr, status, err := receiveData(r)
	if err != nil {
		sendError(w, r, h.Logger, err, status)
		return
	}

	err = h.ItemRepo.TransferMoney(userCurr.UserID, userCurr.ToUserID, userCurr.Balance)
	if err != nil {
		sendError(w, r, h.Logger, err, http.StatusInternalServerError)
		return
	}

	sendSuccessStatus(w, r, h.Logger)
}

func (h ItemsHandler) ListTransaction(w http.ResponseWriter, r *http.Request) {
	userCurr, status, err := receiveData(r)
	if err != nil {
		sendError(w, r, h.Logger, err, status)
		return
	}

	info, err := h.ItemRepo.GetTransaction(userCurr.UserID, userCurr.Field)
	if err != nil {
		sendError(w, r, h.Logger, err, http.StatusInternalServerError)
		return
	}

	sendData(w, r, h.Logger, info)
}
