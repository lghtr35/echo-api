package document

import "echo-api/models/entities"

type DocumentWrapped struct {
	entities.Document
	Path string
}
