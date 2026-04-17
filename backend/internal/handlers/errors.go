package handlers

var (
	ErrInvalidRequestBody = "invalid request body"
	ErrEntityNotFound     = "%s not found"
	ErrFailedToFetch      = "failed to fetch %s"
	ErrFailedToCreate     = "failed to create %s"
	ErrFailedToUpdate     = "failed to update %s"
	ErrFailedToDelete     = "failed to delete %s"
	ErrIDRequired         = "%s id required"
	ErrAtLeastOneField    = "at least one field must be provided"
	ErrNameCannotBeEmpty  = "name cannot be empty"
)
