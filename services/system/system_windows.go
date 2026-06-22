//go:build windows

package system

import "golang.org/x/sys/windows/registry"

func GetSystemLocale() string {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\International`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()
	v, _, err := k.GetStringValue("LocaleName")
	if err != nil {
		return ""
	}
	return v
}
