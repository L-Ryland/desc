package usersys

import (
	"flag"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"

	"server/mongodb"
	"server/util"
)

// user is cold data
type user struct {
	Name     string         `json:"name"`
	password string         `json:"-"`
	Heros    []int          `json:"heros"`
	Role     util.RoleLevel `json:"role"`
}

var testFlag = flag.Bool("user.test", false, "test without database")
var managerFlag = flag.Bool("user.auto-manager", false, "new user as manager")

var userList sync.Map // map<name string,*user>

const teamNameSubString = "的团队"

func Init() {
	u, err := getUser("user1")
	if err == nil {
		if u != nil {
			return
		}
		logrus.Info("load empty testuser")
	} else {
		logrus.Warnf("load user1 error with: %v. Will try to init it.", err)
	}

	err = Register("user1", "aassdd")
	if err != nil {
		panic(err)
	}
	u, err = getUser("user1")
	if err != nil {
		panic(err)
	}
}

func Register(username string, password string) error {
	if username == "" || password == "" {
		return util.Errorf("Invalid username or password").WithCode(codes.InvalidArgument)
	}

	// Store the password securely (e.g., using bcrypt)
	u := &user{
		username,
		password,
		[]int{},
		util.RolePlayer,
	}

	if *managerFlag {
		u.Role = util.RoleManager
	}

	if err := newUser(u); err != nil {
		return util.Errorf("failed to new user %s.", username).WithCause(err)
	}

	return nil
}

func Login(username string, password string) (*user, error) {
	if username == "" || password == "" {
		return nil, util.Errorf("Invalid username or password").WithCode(codes.InvalidArgument)
	}

	user, err := getUser(username)
	if err != nil {
		return nil, util.Errorf("Invalid username or password").WithCode(codes.NotFound).WithCause(err)
	}
	if user.password != password {
		return nil, util.Errorf("Invalid username or password").WithCode(codes.PermissionDenied)
	}

	return user, nil
}

func getUser(name string) (*user, error) {
	if *testFlag {
		u, ok := userList.Load(name)
		if !ok {
			return nil, util.Errorf("user %s notfound", name)
		}
		return u.(*user), nil
	}
	dbu, err := mongodb.GetUserByName(name)
	if err != nil {
		return nil, util.Errorf("user %s notfound", name).WithCause(err)
	}
	return &user{
		Name:     dbu.Name,
		password: dbu.Password,
		Heros:    dbu.Heros,
		Role:     util.RoleLevel(dbu.Role),
	}, nil
}

func newUser(u *user) error {
	if *testFlag {
		if _, ok := userList.LoadOrStore(u.Name, u); ok {
			return util.Errorf("Username already exists").WithCode(codes.AlreadyExists)
		}
		return nil
	}
	_, err := mongodb.AddUser(mongodb.UserPayload{
		Name:     u.Name,
		Password: u.password,
		Heros:    u.Heros,
		Role:     int(u.Role),
	})
	return err
}
