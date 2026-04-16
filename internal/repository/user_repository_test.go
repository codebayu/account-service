package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/codebayu/account-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	db   *gorm.DB
	repo UserRepository
}

func (s *RepositoryTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	assert.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	s.db, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(s.T(), err)

	s.repo = NewUserRepository(s.db)
}

func (s *RepositoryTestSuite) TearDownTest() {
	db, _ := s.db.DB()
	db.Close()
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestCreate() {
	user := &models.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "id"}).AddRow("uuid-123", 1))
	s.mock.ExpectCommit()

	err := s.repo.Create(user)
	assert.NoError(s.T(), err)
}

func (s *RepositoryTestSuite) TestFindByEmail_Success() {
	email := "test@example.com"
	rows := sqlmock.NewRows([]string{"id", "email", "name"}).
		AddRow(1, email, "Test User")

	s.mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email, 1). // GORM adds LIMIT 1
		WillReturnRows(rows)

	user, err := s.repo.FindByEmail(email)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), email, user.Email)
}

func (s *RepositoryTestSuite) TestFindByEmail_NotFound() {
	email := "notfound@example.com"
	s.mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email, 1). // GORM adds LIMIT 1
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := s.repo.FindByEmail(email)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
}

func (s *RepositoryTestSuite) TestFindByUUID_Success() {
	uuid := "uuid-123"
	rows := sqlmock.NewRows([]string{"id", "uuid", "name"}).
		AddRow(1, uuid, "Test User")

	s.mock.ExpectQuery(`SELECT \* FROM "users" WHERE uuid = \$1`).
		WithArgs(uuid, 1).
		WillReturnRows(rows)

	user, err := s.repo.FindByUUID(uuid)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), uuid, user.UUID)
}

func (s *RepositoryTestSuite) TestFindByUUID_NotFound() {
	uuid := "notfound-uuid"
	s.mock.ExpectQuery(`SELECT \* FROM "users" WHERE uuid = \$1`).
		WithArgs(uuid, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := s.repo.FindByUUID(uuid)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), user)
}
