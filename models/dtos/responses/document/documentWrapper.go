package document

import "reson8-learning-api/models/entities"

type DocumentWrapped struct {
	entities.Document
	Path string
}
