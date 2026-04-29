package error

import (
	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"net/http"
)

// TODO: Hanya contoh. nanti langsung hapus saja

func ErrXtremeTestingGet(internalMsg string) {
	xtremeres.Error(http.StatusNotFound, "Testing not found", internalMsg, false, nil)
}

func ErrXtremeTestingSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save testing", internalMsg, false, nil)
}

func ErrXtremeTestingDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete testing", internalMsg, false, nil)
}

func ErrXtremeTestingSubSave(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to save testing sub", internalMsg, false, nil)
}

func ErrXtremeTestingSubDelete(internalMsg string) {
	xtremeres.Error(http.StatusInternalServerError, "Unable to delete testing sub", internalMsg, false, nil)
}
