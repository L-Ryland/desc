package mongodb

import (
	"context"
	"server/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
)

var userdb *mongo.Collection

type DBUser struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string
	Password string `json:"-"`
	Heros    []int
	Role     int
}
type UserPayload struct {
	Name     string
	Password string
	Heros    []int
	Role     int
}

func init() {
	registerDBData(DBUser{})
}

func (DBUser) initTable() {
	userdb = db.Collection("user")
}

// 插入 user 表数据
func AddUser(user UserPayload) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	res, err := userdb.InsertOne(ctx, user)
	if err != nil {
		return "", util.Errorf("add user %s failed to exec.", user.Name).WithCause(err)
	}
	id := res.InsertedID.(primitive.ObjectID)

	return id.String(), nil
}

// 删除 user 表数据
func DeleteUserById(id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		util.Errorf(err.Error())
	}
	filter := bson.M{"_id": objId}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	_, err = userdb.DeleteOne(ctx, filter)
	if err != nil {
		return util.Errorf("delete user with ID %s failed", id).WithCause(err)
	}
	return nil
}
func DeleteUser(user DBUser) error {
	filter := bson.M{"name": user.Name}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	result, err := userdb.DeleteOne(ctx, filter)
	if err != nil {
		return util.Errorf("delete user with ID %s failed", user.Name).WithCause(err)
	}
	// ignore not found error
	if result.DeletedCount == 0 {
		util.Errorf("delete user %s failed", user.Name).WithCause(err).WithCode(codes.NotFound).Log()
	}
	return nil
}

// 更新 user 表数据
func UpdateUser(user DBUser) error {
	filter := bson.M{"name": user.Name}
	update := bson.M{"$set": user}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	result, err := userdb.UpdateOne(ctx, filter, update)
	if err != nil {
		return util.Errorf("update user %s failed", user.Name).WithCause(err)
	}
	if result.ModifiedCount == 0 {
		return util.Errorf("update user %s failed", user.Name).WithCause(err).WithCode(codes.NotFound)
	}
	return nil
}

func GetUserByName(uname string) (DBUser, error) {
	filter := bson.M{"name": uname}
	result := DBUser{}
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	err := userdb.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, util.Errorf("get %s user failed", uname).WithCause(err)
	}
	return result, nil
}

// 查询 user 表数据
func GetAllUsers() ([]DBUser, error) {
	filter := bson.M{} // 空的过滤条件，匹配所有文档
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutTime)
	defer cancel()
	cursor, err := userdb.Find(ctx, filter)
	if err != nil {
		return nil, util.Errorf("get all user failed").WithCause(err)
	}
	defer cursor.Close(context.Background())

	var users []DBUser
	if err := cursor.All(ctx, &users); err != nil {
		return nil, util.Errorf("get all user failed").WithCause(err)
	}
	return users, nil
}
