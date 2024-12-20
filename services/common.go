package services

import (
	"reson8-learning-api/models/dtos/operation"
	"reson8-learning-api/util"

	"gorm.io/gorm"
)

func saveAssociationUpdatesToDb[TOps any, TResult any](db *gorm.DB, course TResult, onesToAdd []TOps, onesToDelete []TOps, AssociationType string) (TResult, error) {
	if len(onesToAdd) > 0 {
		err := db.Model(&course).Association(AssociationType).Append(onesToAdd)
		if err != nil {
			return course, err
		}
	}
	if len(onesToDelete) > 0 {
		err := db.Model(&course).Association(AssociationType).Delete(onesToDelete)
		if err != nil {
			return course, err
		}
	}

	return course, nil
}

func filterOperations[T any](sliceToFilter []operation.Operable[T], logger *util.Logger) ([]T, []T) {
	onesToAdd := make([]T, 0, 16)
	onesToDelete := make([]T, 0, 16)
	for _, value := range sliceToFilter {
		switch value.Op {
		case operation.Insert:
			onesToAdd = append(onesToAdd, value.Val)

		case operation.Delete:
			onesToDelete = append(onesToDelete, value.Val)
		default:
			logger.Error().Msg("an update operable with an unknown operation type")
		}
	}

	return onesToAdd, onesToDelete
}
