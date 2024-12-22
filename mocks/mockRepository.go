package mocks

import (
	"errors"
	"fmt"
	"reflect"
	"reson8-learning-api/util"
	"slices"
	"strconv"
	"strings"

	"gorm.io/gorm/clause"
)

type statement struct {
	Value      any
	Comparison string
}

type MockRepository[T any] struct {
	statements map[string]statement
	data       map[string]T
	idCounter  uint64
	order      string
}

func NewMockRepo[T any]() *MockRepository[T] {
	return &MockRepository[T]{statements: make(map[string]statement), data: make(map[string]T), idCounter: 0}
}

func (r *MockRepository[T]) Query() util.Repository[T] {
	return &MockRepository[T]{statements: make(map[string]statement), data: r.data, idCounter: 0}
}

func (r *MockRepository[T]) First(id string, shouldPreload bool) (T, error) {
	res, ok := r.data[id]
	if !ok {
		var temp T
		return temp, errors.New("notFoundError")
	}
	return res, nil
}

func (r *MockRepository[T]) Find(shouldPreload bool) ([]T, error) {
	res := make([]T, 0)
	for _, v := range r.data {
		valueOf := reflect.ValueOf(v)
		for sK, sV := range r.statements {
			f := valueOf.FieldByName(sK).Elem().Interface()
			if f != sV {
				break
			}
		}
		res = append(res, v)
	}
	slices.SortFunc(res, r.orderByReflection)
	return res, nil
}

func (r *MockRepository[T]) Count() (int64, error) {
	res, err := r.Find(false)
	if err != nil {
		return 0, err
	}
	return int64(len(res)), nil
}

func (r *MockRepository[T]) Create(val *T) (T, error) {
	currId := r.idCounter + 1
	currIdStr := strconv.FormatUint(currId, 10)
	f := reflect.ValueOf(val).Elem().FieldByName("ID")
	if !f.CanSet() {
		fmt.Printf("Field ID Kind %v\n", f.Kind().String())
		fmt.Printf("Field Can addr: %v\n", f.CanAddr())
		var temp T
		return temp, errors.New("idCantSetError")
	}
	f.SetString(currIdStr)
	r.data[currIdStr] = *val
	r.idCounter = currId
	return *val, nil
}

func (r *MockRepository[T]) Update(val *T) (T, error) {
	id := reflect.ValueOf(val).Elem().FieldByName("ID").String()
	r.data[id] = *val
	return *val, nil
}

func (r *MockRepository[T]) Delete(id string) error {
	_, ok := r.data[id]
	if !ok {
		return errors.New("notFoundError")
	}

	delete(r.data, id)

	return nil
}

func (r *MockRepository[T]) Where(query string, args ...any) util.Repository[T] {
	queryParts := strings.Split(query, " ")
	if len(queryParts) < 2 {
		return r
	}
	st := statement{Value: args, Comparison: queryParts[1]}
	r.statements[queryParts[0]] = st
	return r
}

// If needed implement
func (r *MockRepository[T]) Offset(offset int) util.Repository[T] {
	return r
}

func (r *MockRepository[T]) Limit(limit int) util.Repository[T] {
	return r
}

func (r *MockRepository[T]) Order(args ...any) util.Repository[T] {
	if len(args) > 1 {
		r.order = args[0].(string)
	}
	return r
}
func (r *MockRepository[T]) Clauses(conds ...clause.Expression) util.Repository[T] {
	return r
}

func (r *MockRepository[T]) orderByReflection(a T, b T) int {
	aV := reflect.ValueOf(a).FieldByName(r.order)
	bV := reflect.ValueOf(b).FieldByName(r.order)
	if !aV.IsValid() || !bV.IsValid() {
		return 1
	}

	return strings.Compare(aV.Elem().String(), bV.Elem().String())
}
