package mockrepo

import "github.com/zhuravlev-pe/course-watch/internal/repository"

func New() *repository.Repositories {
	return &repository.Repositories{
		Courses: newCourses(),
	}
}
