package error

import (
	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"net/http"
)

func ErrXtremeActivityNotFound() {
	xtremeres.Error(http.StatusNotFound, "Object not found", "", false, nil)
}

func ErrXtremeActivitySave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save activity", internalMsg, false, nil)
}

func ErrXtremeActivityActionType() {
	xtremeres.Error(http.StatusInternalServerError, "Object action doesn't exists!!", "", false, nil)
}
