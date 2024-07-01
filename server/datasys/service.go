package datasys

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/mongodb"
	"server/usersys"
	"server/util"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	// "github.com/sirupsen/logrus"
)

func HandleAddWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// _, role, err := util.GetUser(w, r)
	// if err != nil {
	// 	logrus.Error(util.Errorf("get user failed").WithCause(err))
	// 	return
	// }
	// if role < util.RoleManager {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	webJson := r.Body
	decoder := json.NewDecoder(webJson)
	var webData mongodb.WebData
	if err := decoder.Decode(&webData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	if _, err := mongodb.AddWebData(webData); err != nil {
		if util.HaveErrorCode(err, codes.AlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}

func HandleSearchWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tagsString := mux.Vars(r)["tags"]
	tags := strings.Split(tagsString, ",")

	webDatas, err := mongodb.GetWebDataByTags(tags)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, util.EncodeJson(webDatas))
}

func HandleDeleteWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// _, role, err := util.GetUser(w, r)
	// if err != nil {
	// 	logrus.Error(util.Errorf("get user failed").WithCause(err))
	// 	return
	// }
	// if role < util.RoleManager {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	err = mongodb.DeleteWebData(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}

func HandlePatchWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// _, role, err := util.GetUser(w, r)
	// if err != nil {
	// 	logrus.Error(util.Errorf("get user failed").WithCause(err))
	// 	return
	// }
	// if role < util.RoleManager {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	idString := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	webJson := r.Body
	decoder := json.NewDecoder(webJson)
	var webData mongodb.WebData
	if err := decoder.Decode(&webData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	webData.ID = id

	err = mongodb.UpdateWebData(webData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}

func HandleAddTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_, _, err := util.GetUser(w, r)
	if err != nil {
		logrus.Error(util.Errorf("get user failed").WithCause(err))
		return
	}
	// if role < util.RoleManager {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	tagJson := r.Body
	decoder := json.NewDecoder(tagJson)
	var tagData mongodb.Tag
	if err := decoder.Decode(&tagData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	if err := mongodb.AddTag(tagData); err != nil {
		if util.HaveErrorCode(err, codes.AlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}

func HandleGetAllTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	filter := bson.M{"category": bson.M{"$exists": false}}
	tags, err := mongodb.GetAllTags(filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, util.EncodeJson(tags))
}

func HandleDeleteTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	_, role, err := util.GetUser(w, r)
	if err != nil {
		logrus.Error(util.Errorf("get user failed").WithCause(err))
		return
	}
	if role < util.RoleManager {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	name := mux.Vars(r)["name"]

	err = mongodb.DeleteTag(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}

func HandleReorderTag(w http.ResponseWriter, r *http.Request) {
	user := usersys.GetUserFromCookie(w, r)
	if user == nil {
		return
	}
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := mux.Vars(r)["name"]
	tagJson := r.Body
	decoder := json.NewDecoder(tagJson)
	var tagData mongodb.Tag
	if err := decoder.Decode(&tagData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	tagData.Name = name
	err := mongodb.UpdateTag(name, tagData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusNotImplemented)
}

func HandleGetAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	categories, err := mongodb.GetAllCategories()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, util.EncodeJson(categories))
}

func HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := mux.Vars(r)["id"]

	categoryJson := r.Body
	decoder := json.NewDecoder(categoryJson)
	var categoryData mongodb.Category
	if err := decoder.Decode(&categoryData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	err := mongodb.UpdateCategory(id, categoryData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "success")
}
