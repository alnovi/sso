package repository

import (
	sq "github.com/Masterminds/squirrel"
)

type OptSelect func(builder sq.SelectBuilder) sq.SelectBuilder

func NotDeleted() OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(sq.Eq{"deleted_at": nil})
	}
}

func NotSystem() OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(sq.Eq{"is_system": false})
	}
}

func ForUpdate() OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Suffix("FOR UPDATE")
	}
}

func IsNotNull(field string) OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(sq.NotEq{field: nil})
	}
}

func Secret(val string) OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(sq.Eq{"secret": val})
	}
}

func Class(val string) OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(sq.Eq{"class": val})
	}
}

func IP(val string) OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(sq.Eq{"ip": val})
	}
}

func Agent(val string) OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(sq.Eq{"agent": val})
	}
}

func OrderAsc(field string) OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.OrderBy(field + " asc")
	}
}

func OrderDesc(field string) OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.OrderBy(field + " desc")
	}
}

func SelectWhere(raw sq.Sqlizer) OptSelect {
	return func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(raw)
	}
}
