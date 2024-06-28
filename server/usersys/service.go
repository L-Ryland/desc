package usersys

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/mongodb"
	"server/util"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

type authRequest struct {
	Username string
	Password string
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request authRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "无法解析请求体", http.StatusBadRequest)
		return
	}
	logrus.Infof("register:%v", request)

	if err := Register(request.Username, request.Password); err != nil {
		if util.HaveErrorCode(err, codes.InvalidArgument) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		if util.HaveErrorCode(err, codes.AlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	if err := util.AddSession(w, r, request.Username, util.RolePlayer); err != nil {
		util.Errorf("save session error:%s", request.Username).WithCause(err).Log()
		fmt.Fprint(w, err.Error())
		return
	}

	user, err := getUser(request.Username)
	if err != nil {
		fmt.Fprint(w, util.Errorf("register got some internal error.").WithCause(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, util.EncodeJson(user))
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var request authRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "无法解析请求体", http.StatusBadRequest)
		return
	}
	logrus.Infof("login:%v", request)

	user, err := Login(request.Username, request.Password)
	if err != nil {
		if util.HaveErrorCode(err, codes.InvalidArgument) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		if util.HaveErrorCode(err, codes.NotFound) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	if err := util.AddSession(w, r, user.Name, user.Role); err != nil {
		util.Errorf("save session error:%s.", request.Username).WithCause(err).Log()
		fmt.Fprint(w, err.Error())
		return
	}

	// User login successful
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, util.EncodeJson(user))
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := util.RemoveSession(w, r); err != nil {
		util.Errorf("remove session error").WithCause(err).Log()
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	users, err := mongodb.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, util.EncodeJson(users))
}

func HandleAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userJson := r.Body
	decoder := json.NewDecoder(userJson)
	var userData mongodb.UserPayload
	if err := decoder.Decode(&userData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	if _, err := mongodb.AddUser(userData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	resp := map[string]bool{"success": true}
	fmt.Fprint(w, util.EncodeJson(resp))
}

func HandleRemoveUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := mux.Vars(r)["id"]
	err := mongodb.DeleteUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(w, err.Error())
	}
	w.WriteHeader(http.StatusOK)
	resp := map[string]bool{"success": true}
	fmt.Fprint(w, util.EncodeJson(resp))
}

func GetUserFromCookie(w http.ResponseWriter, r *http.Request) *user {
	userName, _, err := util.GetUser(w, r)
	if err != nil {
		logrus.Error(util.Errorf("get hero failed").WithCause(err))
		return nil
	}
	user, err := getUser(userName)
	if err != nil {
		logrus.Error(err)
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, util.Errorf("User %s not found", userName).Error())
		return nil
	}
	return user
}
