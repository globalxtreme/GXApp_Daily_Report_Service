package core

import "gorm.io/gorm"

type TransactionInterface interface {
	SetTransaction(tx *gorm.DB)
}

// TODO: Re-enable this code after installing github.com/globalxtreme/go-identifier module (If you use GX Identifier for authorization)
//type EmployeeIdentifierInterface interface {
//	SetEmployeeIdentifier(employee data.EmployeeIdentifierData)
//}
