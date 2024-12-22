package util

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository[T any] struct {
	db       *gorm.DB
	preloads []string
}

func NewGormRepository[T any](db *gorm.DB, preloads []string) *GormRepository[T] {
	return &GormRepository[T]{db: db, preloads: preloads}
}

func (r *GormRepository[T]) Query() Repository[T] {
	var temp T
	return &GormRepository[T]{db: r.db.Model(&temp), preloads: r.preloads}
}

func (r *GormRepository[T]) First(id string, shouldPreload bool) (T, error) {
	var val T
	if shouldPreload {
		for _, v := range r.preloads {
			r.db = r.db.Preload(v)
		}
	}
	res := r.db.First(&val, id)
	if res.Error != nil {
		return val, res.Error
	}
	return val, nil
}

func (r *GormRepository[T]) Find(shouldPreload bool) ([]T, error) {
	var vals []T
	if shouldPreload {
		for _, v := range r.preloads {
			r.db = r.db.Preload(v)
		}
	}
	res := r.db.Find(&vals)
	if res.Error != nil {
		return vals, res.Error
	}
	return vals, nil
}

func (r *GormRepository[T]) Where(query string, args ...any) Repository[T] {
	r.db = r.db.Where(query, args...)

	return r
}

func (r *GormRepository[T]) Create(val *T) (T, error) {
	res := r.db.Create(val)
	if res.Error != nil {
		var temp T
		return temp, res.Error
	}
	res = r.db.Save(&val)
	if res.Error != nil {
		var temp T
		return temp, res.Error
	}

	return *val, nil
}

func (r *GormRepository[T]) Delete(id string) error {
	var temp T
	res := r.db.Select(clause.Associations).Delete(&temp, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *GormRepository[T]) Update(val *T) (T, error) {
	res := r.db.Save(val)
	if res.Error != nil {
		var temp T
		return temp, res.Error
	}
	return *val, nil
}

func (r *GormRepository[T]) Count() (int64, error) {
	var count int64
	res := r.db.Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}

	return count, nil
}

func (r *GormRepository[T]) Offset(offset int) Repository[T] {
	r.db = r.db.Offset(offset)
	return r
}

func (r *GormRepository[T]) Limit(limit int) Repository[T] {
	r.db = r.db.Limit(limit)
	return r
}
func (r *GormRepository[T]) Order(args ...any) Repository[T] {
	r.db = r.db.Order(args)
	return r
}
func (r *GormRepository[T]) Clauses(conds ...clause.Expression) Repository[T] {
	r.db = r.db.Clauses(conds...)
	return r
}
