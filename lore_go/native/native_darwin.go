// Copyright Epic Games, Inc. All Rights Reserved.

//go:build darwin

package native

import "github.com/ebitengine/purego"

// loadLibrary loads a dynamic library using Unix dlopen
func loadLibrary(path string) (uintptr, error) {
	return purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

// loadSymbol loads a symbol from a library handle using Unix dlsym
func loadSymbol(handle uintptr, name string) (uintptr, error) {
	return purego.Dlsym(handle, name)
}

func nativeLibraryFileName() string { return "liblore.dylib" }

// loreFuncWithCallback is the function signature for Lore functions on Unix
// On Unix, the callback struct is passed by value (flattened into separate arguments)
type loreFuncWithCallback func(globalsPtr, argsPtr, callbackUserContext, callbackFuncPtr uintptr) int32

// callLoreFuncWithCallback invokes a Lore function with the callback config
// On Unix, the callback struct fields are passed as separate arguments
// The callbackConfigPtr parameter is unused on Unix
func callLoreFuncWithCallback(fn loreFuncWithCallback, globalsPtr, argsPtr, callbackConfigPtr, callbackUserContext, callbackFuncPtrValue uintptr) int32 {
	return fn(globalsPtr, argsPtr, callbackUserContext, callbackFuncPtrValue)
}
