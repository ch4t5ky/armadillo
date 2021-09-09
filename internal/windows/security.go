package windows

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	advAPI32DLL          = windows.NewLazyDLL("Advapi32.dll")
	setNamedSecurityInfo = advAPI32DLL.NewProc("SetNamedSecurityInfo")
)

func SetNamedSecurityInfo(filepath string) {
	/*
		DWORD SetNamedSecurityInfo(
		  LPSTR                pObjectName,
		  SE_OBJECT_TYPE       ObjectType,
		  SECURITY_INFORMATION SecurityInfo,
		  PSID                 psidOwner,
		  PSID                 psidGroup,
		  PACL                 pDacl,
		  PACL                 pSacl
		);
	*/

	information := windows.SECURITY_INFORMATION(windows.OWNER_SECURITY_INFORMATION)
	retCode, _, err := setNamedSecurityInfo.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(filepath))),
		uintptr(windows.SE_FILE_OBJECT),
		uintptr(information),
	)
	if retCode != 0 {
		fmt.Println("Function succeed: ", err.Error())
	}
}
