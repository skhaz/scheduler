package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/skhaz/scheduler/model"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setup() (conn *sql.DB, mock sqlmock.Sqlmock, repository TriggerRepository) {
	var err error

	conn, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialector := postgres.New(postgres.Config{Conn: conn, PreferSimpleProtocol: true})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	repository = TriggerRepository{}

	repository.Configure(db)

	return
}

func TestListWorkspaces(t *testing.T) {
	var err error
	conn, mock, repository := setup()
	defer conn.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "triggers"`)).WillReturnRows(sqlmock.NewRows([]string{}))

	var arr any
	arr, err = repository.List(time.Now(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, arr)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetWorkspace(t *testing.T) {
	var err error
	conn, mock, repository := setup()
	defer conn.Close()

	m := model.Trigger{
		ID:        uuid.New(),
		Name:      randstr.String(16),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at", "deleted_at"}).
		AddRow(m.ID, m.Name, m.CreatedAt, m.UpdatedAt, m.DeletedAt)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "triggers"`)).
		WillReturnRows(rows)

	var e any
	e, err = repository.Get(m.ID)
	assert.NoError(t, err)
	assert.NotNil(t, e)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateWorkspace(t *testing.T) {
	var err error
	conn, mock, repository := setup()
	defer conn.Close()

	m := model.Trigger{
		ID:        uuid.New(),
		Name:      randstr.String(16),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "triggers"`)).
		WithArgs(m.Name, m.CreatedAt, m.UpdatedAt, m.DeletedAt, m.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(m.ID))
	mock.ExpectCommit()

	var e any
	e, err = repository.Create(&m)
	assert.NoError(t, err)
	assert.NotNil(t, e)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateWorkspace(t *testing.T) {
	var err error
	conn, mock, repository := setup()
	defer conn.Close()

	m := model.Trigger{
		ID:        uuid.New(),
		Name:      randstr.String(16),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "triggers" SET`)).
		WithArgs(m.Name, AnyTime{}, AnyTime{}, m.ID, m.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var e bool
	e, err = repository.Update(m.ID, &m)
	assert.NoError(t, err)
	assert.True(t, e)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateWorkspaceError(t *testing.T) {
	var err error
	conn, mock, repository := setup()
	defer conn.Close()

	m := model.Trigger{
		ID:        uuid.New(),
		Name:      randstr.String(16),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "triggers" SET`)).
		WithArgs(m.Name, AnyTime{}, AnyTime{}, m.ID, m.ID).
		WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	var e bool
	e, err = repository.Update(m.ID, &m)
	assert.Error(t, err)
	assert.False(t, e)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteWorkspace(t *testing.T) {
	var err error
	conn, mock, repository := setup()
	defer conn.Close()

	uid := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "triggers" SET`)).
		WithArgs(AnyTime{}, uid).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var e bool
	e, err = repository.Delete(uid)
	assert.NoError(t, err)
	assert.True(t, e)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteWorkspaceError(t *testing.T) {
	var err error
	conn, mock, repository := setup()
	defer conn.Close()

	uid := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "triggers" SET`)).
		WithArgs(AnyTime{}, uid).
		WillReturnError(gorm.ErrMissingWhereClause)
	mock.ExpectRollback()

	var e bool
	e, err = repository.Delete(uid)
	assert.Error(t, err)
	assert.False(t, e)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
