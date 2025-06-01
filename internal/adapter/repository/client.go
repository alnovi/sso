package repository

import (
	"context"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
)

const ClientTable = "clients"

var clientFields = []string{"id", "name", "icon", "secret", "callback", "is_system", "created_at", "updated_at", "deleted_at"}

func (r *Repository) Clients(ctx context.Context, opts ...OptSelect) ([]*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.Clients")
	defer span.End()

	clients := make([]*entity.Client, 0)

	builder := r.qb.Select(clientFields...).From(ClientTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQuery(ctx, &clients, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return clients, nil
}

func (r *Repository) ClientsCount(ctx context.Context, opts ...OptSelect) (int, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.ClientsCount")
	defer span.End()

	count := 0

	builder := r.qb.Select("COUNT (*)").From(ClientTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return count, err
	}

	err = r.checkErr(r.db.QueryRow(ctx, query, args...).Scan(&count))
	if err != nil {
		helper.SpanError(span, err)
	}

	return count, err
}

func (r *Repository) ClientById(ctx context.Context, id string, opts ...OptSelect) (*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.ClientById", helper.SpanAttr(
		attribute.String("client.id", id),
	))
	defer span.End()

	client := new(entity.Client)

	if id == "" {
		helper.SpanError(span, ErrNoResult)
		return nil, ErrNoResult
	}

	builder := r.qb.Select(clientFields...).
		From(ClientTable).
		Where(sq.Eq{"id": id})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQueryRow(ctx, client, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return client, nil
}

func (r *Repository) ClientByIds(ctx context.Context, ids []string, opts ...OptSelect) ([]*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.ClientByIds", helper.SpanAttr(
		attribute.String("ids", strings.Join(ids, ", ")),
	))
	defer span.End()

	clients := make([]*entity.Client, 0)

	builder := r.qb.Select(clientFields...).
		From(ClientTable).
		Where(sq.Eq{"id": ids})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQuery(ctx, &clients, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return clients, nil
}

func (r *Repository) ClientCreate(ctx context.Context, client *entity.Client) error {
	ctx, span := helper.SpanStart(ctx, "Repository.ClientCreate")
	defer span.End()

	now := time.Now()

	if client.CreatedAt.IsZero() {
		client.CreatedAt = now
	}

	if client.UpdatedAt.IsZero() {
		client.UpdatedAt = now
	}

	client.Id = strings.ToLower(client.Id)
	client.DeletedAt = nil

	span.SetAttributes(attribute.String("client.id", client.Id))

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
		helper.SpanError(span, err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err = r.checkErr(err); err != nil {
		helper.SpanError(span, err)
		return err
	}

	return nil
}

func (r *Repository) ClientUpdate(ctx context.Context, client *entity.Client) error {
	ctx, span := helper.SpanStart(ctx, "Repository.ClientUpdate", helper.SpanAttr(
		attribute.String("client.id", client.Id),
	))
	defer span.End()

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
		helper.SpanError(span, err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err = r.checkErr(err); err != nil {
		helper.SpanError(span, err)
		return r.checkErr(err)
	}

	return nil
}

func (r *Repository) ClientDelete(ctx context.Context, client *entity.Client) error {
	ctx, span := helper.SpanStart(ctx, "Repository.ClientDelete", helper.SpanAttr(
		attribute.String("id", client.Id),
	))
	defer span.End()

	now := time.Now()
	client.DeletedAt = &now

	builder := r.qb.Update(ClientTable).
		Set("deleted_at", client.DeletedAt).
		Where(sq.Eq{"id": client.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err = r.checkErr(err); err != nil {
		helper.SpanError(span, err)
		return err
	}

	return nil
}

func (r *Repository) ClientDeleteForce(ctx context.Context, client *entity.Client) error {
	ctx, span := helper.SpanStart(ctx, "Repository.ClientDeleteForce", helper.SpanAttr(
		attribute.String("client.id", client.Id),
	))
	defer span.End()

	builder := r.qb.Delete(ClientTable).Where(sq.Eq{"id": client.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err = r.checkErr(err); err != nil {
		helper.SpanError(span, err)
		return err
	}

	return nil
}
