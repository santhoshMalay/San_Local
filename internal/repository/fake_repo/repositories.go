package fake_repo

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/repository"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"time"
)

func New() *repository.Repositories {
	result := &repository.Repositories{
		Courses: NewCourses(),
		Users:   newUsers(),
	}
	
	err := result.Users.Insert(context.Background(), &SampleUser)
	if err != nil {
		panic(err)
	}
	
	return result
}

var SampleUser = core.User{
	Id:               "1582550893222432768",
	Email:            "doe.j@example.com",
	FirstName:        "John",
	LastName:         "Doe",
	DisplayName:      "JonnyD",
	RegistrationDate: time.Date(2017, time.July, 21, 17, 32, 28, 0, time.UTC),
	Roles:            []security.Role{security.Student},
}
