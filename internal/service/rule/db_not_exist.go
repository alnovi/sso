package rule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-playground/validator/v10"

	"github.com/alnovi/sso/pkg/database/postgres"
)

// DatabaseNotExist - use validate:"db_not_exist=users;id"
type DatabaseNotExist struct {
	db *postgres.Client
	qb sq.StatementBuilderType
}

func NewDatabaseNotExist(db *postgres.Client) *DatabaseNotExist {
	return &DatabaseNotExist{db: db, qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *DatabaseNotExist) Tag() string {
	return "db_not_exist"
}

func (r *DatabaseNotExist) ErrMsg() string {
	return "Это значение уже занято"
}

func (r *DatabaseNotExist) CallIfNull() bool {
	return true
}

func (r *DatabaseNotExist) Validate(fl validator.FieldLevel) bool {
	params := fl.Param()
	if len(params) == 0 {
		panic(fmt.Errorf("empty parametr for validation tag [%s]", fl.GetTag()))
	}

	splitParams := strings.Split(params, ";")
	if len(splitParams) != 2 { //nolint:mnd
		panic(fmt.Errorf("invalid parametr format for validation tag [%s]", fl.GetTag()))
	}

	dbName := splitParams[0]
	columnName := splitParams[1]

	builder := r.qb.Select("1").
		From(dbName).
		Where(sq.Eq{columnName: fl.Field().Interface()})

	query, args, err := builder.ToSql()
	if err != nil {
		panic(err)
	}

	exist := ""

	err = r.db.QueryRow(context.Background(), query, args...).Scan(&exist)
	if err != nil {
		return errors.Is(err, sql.ErrNoRows)
	}

	return exist != "1"
}
