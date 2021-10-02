package windows

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

func GrantSid(accessPermissions uint32, sid *windows.SID) ExplicitAccess {
	return ExplicitAccess{
		AccessPermissions: accessPermissions,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       windows.SUB_CONTAINERS_AND_OBJECTS_INHERIT,
		Trustee: Trustee{
			TrusteeForm: windows.TRUSTEE_IS_SID,
			Name:        (*uint16)(unsafe.Pointer(sid)),
		},
	}
}
