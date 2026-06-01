// Copyright Epic Games, Inc. All Rights Reserved.

//go:build windows

package native

import "syscall"

// loadLibrary loads a dynamic library using Windows LoadLibrary
func loadLibrary(path string) (uintptr, error) {
	handle, err := syscall.LoadLibrary(path)
	return uintptr(handle), err
}

// loadSymbol loads a symbol from a library handle using Windows GetProcAddress
func loadSymbol(handle uintptr, name string) (uintptr, error) {
	return syscall.GetProcAddress(syscall.Handle(handle), name)
}

func nativeLibraryFileName() string { return "lore.dll" }

// loreFuncWithCallback is the function signature for Lore functions on Windows
// On Windows x64, structs > 8 bytes passed by value are actually passed by reference
type loreFuncWithCallback func(globalsPtr, argsPtr, callbackConfigPtr uintptr) int32

// callLoreFuncWithCallback invokes a Lore function with the callback config
// On Windows, the callback struct is passed as a pointer
// The callbackUserContext and callbackFuncPtrValue parameters are unused on Windows
func callLoreFuncWithCallback(fn loreFuncWithCallback, globalsPtr, argsPtr, callbackConfigPtr, callbackUserContext, callbackFuncPtrValue uintptr) int32 {
	return fn(globalsPtr, argsPtr, callbackConfigPtr)
}
