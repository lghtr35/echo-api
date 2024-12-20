package document

import (
	"mime/multipart"
)

type CreateDocumentRequestBase struct {
	UserID          string  `form:"userID" binding:"required"`
	Location        string  `form:"location" binding:"required"`
	IsReadableByAll bool    `form:"isReadableByAll" binding:"required"`
	ContextID       string  `form:"contextID" binding:"required"`
	EntityType      *string `form:"entityType"`
	EntityID        *string `form:"entityID"`
}

type CreateDocumentsMultipartRequest struct {
	Files []*multipart.FileHeader `form:"files[]" binding:"required"`
	CreateDocumentRequestBase
}

type CreateDocumentMultipartRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	CreateDocumentRequestBase
}

type CreateNoteDocumentsRequest struct {
	NoteID string `form:"entityID" binding:"required"`
	CreateDocumentsMultipartRequest
}
