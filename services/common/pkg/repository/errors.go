package repository

import "errors"

var (
	ErrCreateEntity    = errors.New("unable to create entity")
	ErrGetEntityByID   = errors.New("unable to get entity by id")
	ErrUpdateEntity    = errors.New("unable to update entity")
	ErrDeleteEntity    = errors.New("unable to delete entity")
	ErrFetchEntities   = errors.New("unable to fetch entities")
	ErrDeleteWithCond  = errors.New("unable to delete entities with condition")
	ErrFindEntities    = errors.New("unable to find entities")
	ErrFindEntity      = errors.New("unable to find entity")
	ErrCountEntities   = errors.New("unable to count entities")
	ErrGetPageEntities = errors.New("unable to get page of entities")
	ErrBulkInsert      = errors.New("unable to bulk insert entities")
	ErrBulkUpdate      = errors.New("unable to bulk update entities")
	ErrTransaction     = errors.New("unable to begin transaction")
)

