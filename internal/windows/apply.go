package windows

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

func Apply(name string, replace, inherit bool, entries ...ExplicitAccess) error {
	var oldAcl windows.Handle
	if !replace {
		var secDesc windows.Handle
		GetNamedSecurityInfo(
			name,
			windows.SE_FILE_OBJECT,
			windows.DACL_SECURITY_INFORMATION,
			nil,
			nil,
			&oldAcl,
			nil,
			&secDesc,
		)
		defer windows.LocalFree(secDesc)
	}
	var acl windows.Handle
	if err := SetEntriesInAcl(
		entries,
		oldAcl,
		&acl,
	); err != nil {
		return err
	}
	defer windows.LocalFree((windows.Handle)(unsafe.Pointer(acl)))
	var secInfo uint32
	if !inherit {
		secInfo = windows.PROTECTED_DACL_SECURITY_INFORMATION
	} else {
		secInfo = windows.UNPROTECTED_DACL_SECURITY_INFORMATION
	}
	return SetNamedSecurityInfo(name, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|secInfo, nil, nil, acl, 0)
}
