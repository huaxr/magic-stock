// @Contact:    huaxinrui
// @Time:       2019/9/26 下午9:25

package s_auth

import (
	"fmt"
	"strings"

	"code.byted.org/byte_security/platform_api/utils"

	"code.byted.org/byte_security/dal/auth"
	"code.byted.org/byte_security/platform_api/dao/d_auth"
	"code.byted.org/byte_security/platform_api/models"
)

type UserServiceIF interface {
	Create(event *auth.User) error
	Delete(id int) error
	Update(id int, m map[string]interface{}) error
	Query(where string, args []interface{}) (*auth.User, error)
	QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.User, error)
	Count(where string, args []interface{}) (int, error)
	GetUserEmail(user *auth.User) string
	GetUserEmailFromDB(name string) (string, error)
	GetUserEmailFromPeople(name string) (string, error)
	GetUserInfoInnerByUserID(userID uint) (*models.UserInfoInner, error)
	GetUserInfoInnerByUserName(userName string) (*models.UserInfoInner, error)
	GetUserInfoByNameLike(partName string) ([]models.PeopleEmp, error)
	GetUserInfoByNameEqual(name string) (models.PeopleEmp, error)
}

var UserServiceGlobal UserServiceIF

func init() {
	tmp := new(UserService)
	tmp.dao = d_auth.UserDao
	UserServiceGlobal = tmp
}

type UserService struct {
	dao d_auth.UserDaoIF
}

func (m *UserService) Create(app *auth.User) error {
	return m.dao.Create(app)
}

func (m *UserService) Delete(id int) error {
	return m.dao.Delete(id)
}

func (m *UserService) Update(id int, ma map[string]interface{}) error {
	return m.dao.Update(id, ma)
}

func (m *UserService) Query(where string, args []interface{}) (*auth.User, error) {
	return m.dao.Query(where, args)
}

func (m *UserService) QueryAll(where string, args []interface{}, offset, limit int) (*[]auth.User, error) {
	return m.dao.QueryAll(where, args, offset, limit)
}

func (m *UserService) Count(where string, args []interface{}) (int, error) {
	return m.dao.Count(where, args)
}

func (m *UserService) GetUserEmail(us *auth.User) string {
	if len(us.Email) > 0 {
		return us.Email
	}
	email, err := m.GetUserEmailFromPeople(us.UserName)
	if err != nil {
		return fmt.Sprintf("%s@bytedance.com", us.UserName)
	} else {
		return email
	}
}

func (m *UserService) GetUserEmailFromDB(name string) (string, error) {
	user, err := m.dao.Query("user_name = ?", []interface{}{name})
	if err != nil {
		return "", err
	} else if len(user.Email) > 0 && strings.Contains(user.Email, "@") {
		return user.Email, nil
	} else {
		return "", fmt.Errorf("email is empty")
	}
}

func (m *UserService) GetUserEmailFromPeople(name string) (string, error) {
	emp, err := m.GetUserInfoByNameEqual(name)
	if err != nil {
		return "", nil
	}
	return emp.Email, err
}

func (m *UserService) GetUserInfoInnerByUserID(userID uint) (*models.UserInfoInner, error) {
	return m.GetUserInfoInner("id = ?", []interface{}{userID})
}

func (m *UserService) GetUserInfoInnerByUserName(userName string) (*models.UserInfoInner, error) {
	return m.GetUserInfoInner("user_name = ?", []interface{}{userName})
}

func (m *UserService) GetUserInfoInner(where string, args []interface{}) (*models.UserInfoInner, error) {
	user, err := m.Query(where, args)
	if err != nil {
		return nil, err
	}
	userInfo := models.UserInfoInner{
		UserID:    int(user.ID),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		UserName:  user.UserName,
		RealName:  user.RealName,
		UserNum:   user.UserNum,
		GroupID:   user.GroupId,
		GroupName: "",
		AvatarURL: user.AvatarUrl,
		Leader:    user.Leader,
		Email:     user.Email,
	}
	group, err := GroupServiceGlobal.Query("group_id = ?", []interface{}{user.GroupId})
	if err != nil {
		return &userInfo, err
	}
	userInfo.GroupName = group.GroupName
	return &userInfo, err
}

func (m *UserService) GetUserInfoByNameLike(partName string) ([]models.PeopleEmp, error) {
	url := "https://open.byted.org/people/employee/?name_like=" + partName
	token := "Basic MTExOjU5QTFEMTlENUJGRjQ2QzdCRkVENDQ2QzBCNUU1RjY2"
	var response models.PeopleRspInfo
	err := utils.HttpGetWithToken(url, token, &response)
	return response.Employees, err
}

func (m *UserService) GetUserInfoByNameEqual(name string) (emp models.PeopleEmp, err error) {
	emps, e := m.GetUserInfoByNameLike(fmt.Sprintf("%s@", name))
	if e != nil {
		err = e
	} else if len(emps) == 1 {
		emp = emps[0]
	} else if len(emps) > 1 {
		for _, e := range emps {
			if e.Username == name {
				emp = e
			}
		}
	} else if len(emps) == 0 {
		err = fmt.Errorf("no employee of name:%s", name)
	}
	return
}
