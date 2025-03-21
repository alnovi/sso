package repository

import (
	"context"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/alnovi/sso/internal/entity"
)

const ClientTable = "clients"

var clientFields = []string{"id", "name", "icon", "secret", "callback", "is_system", "created_at", "updated_at", "deleted_at"}

func (r *Repository) Clients(ctx context.Context, opts ...OptSelect) ([]*entity.Client, error) {
	clients := make([]*entity.Client, 0)

	builder := r.qb.Select(clientFields...).From(ClientTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQuery(ctx, &clients, query, args...)

	return clients, r.checkErr(err)
}

func (r *Repository) ClientsCount(ctx context.Context, opts ...OptSelect) (int, error) {
	count := 0

	builder := r.qb.Select("COUNT (*)").From(ClientTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return count, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&count)

	return count, r.checkErr(err)
}

func (r *Repository) ClientById(ctx context.Context, id string, opts ...OptSelect) (*entity.Client, error) {
	client := new(entity.Client)

	if id == "" {
		return nil, ErrNoResult
	}

	builder := r.qb.Select(clientFields...).
		From(ClientTable).
		Where(sq.Eq{"id": id})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, client, query, args...)

	return client, err
}

func (r *Repository) ClientByIds(ctx context.Context, ids []string, opts ...OptSelect) ([]*entity.Client, error) {
	clients := make([]*entity.Client, 0)

	builder := r.qb.Select(clientFields...).
		From(ClientTable).
		Where(sq.Eq{"id": ids})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQuery(ctx, &clients, query, args...)

	return clients, err
}

func (r *Repository) ClientCreate(ctx context.Context, client *entity.Client) error {
	now := time.Now()

	if client.CreatedAt.IsZero() {
		client.CreatedAt = now
	}

	if client.UpdatedAt.IsZero() {
		client.UpdatedAt = now
	}

	client.Id = strings.ToLower(client.Id)
	client.DeletedAt = nil

	builder := r.qb.Insert(ClientTable).
		Columns(clientFields...).
		Values(
			client.Id,
			client.Name,
			client.Icon,
			client.Secret,
			client.Callback,
			client.IsSystem,
			client.CreatedAt,
			client.UpdatedAt,
			client.DeletedAt,
		)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) ClientUpdate(ctx context.Context, client *entity.Client) error {
	client.UpdatedAt = time.Now()

	builder := r.qb.Update(ClientTable).
		Set("name", client.Name).
		Set("icon", client.Icon).
		Set("callback", client.Callback).
		Set("secret", client.Secret).
		Set("updated_at", client.UpdatedAt).
		Set("deleted_at", client.DeletedAt).
		Where(sq.Eq{"id": client.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) ClientDelete(ctx context.Context, client *entity.Client) error {
	now := time.Now()
	client.DeletedAt = &now

	builder := r.qb.Update(ClientTable).
		Set("deleted_at", client.DeletedAt).
		Where(sq.Eq{"id": client.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) ClientDeleteForce(ctx context.Context, client *entity.Client) error {
	builder := r.qb.Delete(ClientTable).Where(sq.Eq{"id": client.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}
