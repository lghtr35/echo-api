package util

import "gorm.io/gorm/clause"

type Repository[T any] interface {
	Query() Repository[T]

	First(id string, shouldPreload bool) (T, error)
	Find(shouldPreload bool) ([]T, error)
	Count() (int64, error)

	Create(val *T) (T, error)
	Update(val *T) (T, error)
	Delete(id string) error

	Where(query string, args ...any) Repository[T]
	Offset(offset int) Repository[T]
	Limit(limit int) Repository[T]
	Order(args ...any) Repository[T]
	Clauses(conds ...clause.Expression) Repository[T]
}
