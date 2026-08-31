// Copyright Epic Games, Inc. All Rights Reserved.

package types

import (
	"encoding/hex"
	"fmt"
	"runtime"
	"unsafe"
)

type LoreTraceLocationArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreTraceLocationArray = []LoreTraceLocation

func (arr LoreTraceLocationArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreTraceLocationArrayFFI) Get(index int) LoreTraceLocation {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreTraceLocation)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreTraceLocationArrayFFI) Clone() []LoreTraceLocation {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreTraceLocation)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreTraceLocation, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreTraceLocationArray(arr []LoreTraceLocation) (LoreTraceLocationArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreTraceLocationArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreTraceLocationFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreTraceLocation(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreTraceLocationArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreInstanceIdArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreInstanceIdArray = []LoreInstanceId

func (arr LoreInstanceIdArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreInstanceIdArrayFFI) Get(index int) LoreInstanceId {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreInstanceId)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreInstanceIdArrayFFI) Clone() []LoreInstanceId {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreInstanceId)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreInstanceId, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreInstanceIdArray(arr []LoreInstanceId) (LoreInstanceIdArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreInstanceIdArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type is FFI-compatible (primitive, enum, or alias); copy directly.
	ffiArray := make([]LoreInstanceId, len(arr))
	copy(ffiArray, arr)

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		pinner.Unpin()
	}

	return LoreInstanceIdArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreBranchPointArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreBranchPointArray = []LoreBranchPoint

func (arr LoreBranchPointArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreBranchPointArrayFFI) Get(index int) LoreBranchPoint {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreBranchPoint)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreBranchPointArrayFFI) Clone() []LoreBranchPoint {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreBranchPoint)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreBranchPoint, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreBranchPointArray(arr []LoreBranchPoint) (LoreBranchPointArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreBranchPointArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreBranchPointFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreBranchPoint(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreBranchPointArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreMetadataTypeArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreMetadataTypeArray = []LoreMetadataType

func (arr LoreMetadataTypeArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreMetadataTypeArrayFFI) Get(index int) LoreMetadataType {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreMetadataType)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreMetadataTypeArrayFFI) Clone() []LoreMetadataType {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreMetadataType)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreMetadataType, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreMetadataTypeArray(arr []LoreMetadataType) (LoreMetadataTypeArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreMetadataTypeArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type is FFI-compatible (primitive, enum, or alias); copy directly.
	ffiArray := make([]LoreMetadataType, len(arr))
	copy(ffiArray, arr)

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		pinner.Unpin()
	}

	return LoreMetadataTypeArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreUint32ArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreUint32Array = []uint32

func (arr LoreUint32ArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreUint32ArrayFFI) Get(index int) uint32 {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*uint32)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreUint32ArrayFFI) Clone() []uint32 {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*uint32)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]uint32, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreUint32Array(arr []uint32) (LoreUint32ArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreUint32ArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type is FFI-compatible (primitive, enum, or alias); copy directly.
	ffiArray := make([]uint32, len(arr))
	copy(ffiArray, arr)

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		pinner.Unpin()
	}

	return LoreUint32ArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStoragePutItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStoragePutItemArray = []LoreStoragePutItem

func (arr LoreStoragePutItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStoragePutItemArrayFFI) Get(index int) LoreStoragePutItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStoragePutItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStoragePutItemArrayFFI) Clone() []LoreStoragePutItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStoragePutItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStoragePutItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStoragePutItemArray(arr []LoreStoragePutItem) (LoreStoragePutItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStoragePutItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStoragePutItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStoragePutItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStoragePutItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageGetItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageGetItemArray = []LoreStorageGetItem

func (arr LoreStorageGetItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageGetItemArrayFFI) Get(index int) LoreStorageGetItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageGetItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageGetItemArrayFFI) Clone() []LoreStorageGetItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageGetItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageGetItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageGetItemArray(arr []LoreStorageGetItem) (LoreStorageGetItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageGetItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageGetItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageGetItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageGetItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageGetResolvedItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageGetResolvedItemArray = []LoreStorageGetResolvedItem

func (arr LoreStorageGetResolvedItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageGetResolvedItemArrayFFI) Get(index int) LoreStorageGetResolvedItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageGetResolvedItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageGetResolvedItemArrayFFI) Clone() []LoreStorageGetResolvedItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageGetResolvedItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageGetResolvedItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageGetResolvedItemArray(arr []LoreStorageGetResolvedItem) (LoreStorageGetResolvedItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageGetResolvedItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageGetResolvedItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageGetResolvedItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageGetResolvedItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStoragePutResolvedItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStoragePutResolvedItemArray = []LoreStoragePutResolvedItem

func (arr LoreStoragePutResolvedItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStoragePutResolvedItemArrayFFI) Get(index int) LoreStoragePutResolvedItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStoragePutResolvedItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStoragePutResolvedItemArrayFFI) Clone() []LoreStoragePutResolvedItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStoragePutResolvedItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStoragePutResolvedItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStoragePutResolvedItemArray(arr []LoreStoragePutResolvedItem) (LoreStoragePutResolvedItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStoragePutResolvedItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStoragePutResolvedItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStoragePutResolvedItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStoragePutResolvedItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageGetMetadataItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageGetMetadataItemArray = []LoreStorageGetMetadataItem

func (arr LoreStorageGetMetadataItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageGetMetadataItemArrayFFI) Get(index int) LoreStorageGetMetadataItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageGetMetadataItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageGetMetadataItemArrayFFI) Clone() []LoreStorageGetMetadataItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageGetMetadataItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageGetMetadataItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageGetMetadataItemArray(arr []LoreStorageGetMetadataItem) (LoreStorageGetMetadataItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageGetMetadataItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageGetMetadataItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageGetMetadataItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageGetMetadataItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageObliterateItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageObliterateItemArray = []LoreStorageObliterateItem

func (arr LoreStorageObliterateItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageObliterateItemArrayFFI) Get(index int) LoreStorageObliterateItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageObliterateItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageObliterateItemArrayFFI) Clone() []LoreStorageObliterateItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageObliterateItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageObliterateItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageObliterateItemArray(arr []LoreStorageObliterateItem) (LoreStorageObliterateItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageObliterateItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageObliterateItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageObliterateItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageObliterateItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageMutableLoadItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageMutableLoadItemArray = []LoreStorageMutableLoadItem

func (arr LoreStorageMutableLoadItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageMutableLoadItemArrayFFI) Get(index int) LoreStorageMutableLoadItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageMutableLoadItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageMutableLoadItemArrayFFI) Clone() []LoreStorageMutableLoadItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageMutableLoadItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageMutableLoadItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageMutableLoadItemArray(arr []LoreStorageMutableLoadItem) (LoreStorageMutableLoadItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageMutableLoadItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageMutableLoadItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageMutableLoadItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageMutableLoadItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageMutableStoreItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageMutableStoreItemArray = []LoreStorageMutableStoreItem

func (arr LoreStorageMutableStoreItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageMutableStoreItemArrayFFI) Get(index int) LoreStorageMutableStoreItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageMutableStoreItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageMutableStoreItemArrayFFI) Clone() []LoreStorageMutableStoreItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageMutableStoreItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageMutableStoreItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageMutableStoreItemArray(arr []LoreStorageMutableStoreItem) (LoreStorageMutableStoreItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageMutableStoreItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageMutableStoreItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageMutableStoreItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageMutableStoreItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageMutableCompareAndSwapItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageMutableCompareAndSwapItemArray = []LoreStorageMutableCompareAndSwapItem

func (arr LoreStorageMutableCompareAndSwapItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageMutableCompareAndSwapItemArrayFFI) Get(index int) LoreStorageMutableCompareAndSwapItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageMutableCompareAndSwapItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageMutableCompareAndSwapItemArrayFFI) Clone() []LoreStorageMutableCompareAndSwapItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageMutableCompareAndSwapItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageMutableCompareAndSwapItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageMutableCompareAndSwapItemArray(arr []LoreStorageMutableCompareAndSwapItem) (LoreStorageMutableCompareAndSwapItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageMutableCompareAndSwapItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageMutableCompareAndSwapItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageMutableCompareAndSwapItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageMutableCompareAndSwapItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageMutableListItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageMutableListItemArray = []LoreStorageMutableListItem

func (arr LoreStorageMutableListItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageMutableListItemArrayFFI) Get(index int) LoreStorageMutableListItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageMutableListItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageMutableListItemArrayFFI) Clone() []LoreStorageMutableListItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageMutableListItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageMutableListItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageMutableListItemArray(arr []LoreStorageMutableListItem) (LoreStorageMutableListItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageMutableListItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageMutableListItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageMutableListItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageMutableListItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageCopyItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageCopyItemArray = []LoreStorageCopyItem

func (arr LoreStorageCopyItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageCopyItemArrayFFI) Get(index int) LoreStorageCopyItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageCopyItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageCopyItemArrayFFI) Clone() []LoreStorageCopyItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageCopyItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageCopyItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageCopyItemArray(arr []LoreStorageCopyItem) (LoreStorageCopyItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageCopyItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageCopyItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageCopyItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageCopyItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStoragePutFileItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStoragePutFileItemArray = []LoreStoragePutFileItem

func (arr LoreStoragePutFileItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStoragePutFileItemArrayFFI) Get(index int) LoreStoragePutFileItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStoragePutFileItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStoragePutFileItemArrayFFI) Clone() []LoreStoragePutFileItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStoragePutFileItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStoragePutFileItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStoragePutFileItemArray(arr []LoreStoragePutFileItem) (LoreStoragePutFileItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStoragePutFileItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStoragePutFileItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStoragePutFileItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStoragePutFileItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageGetFileItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageGetFileItemArray = []LoreStorageGetFileItem

func (arr LoreStorageGetFileItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageGetFileItemArrayFFI) Get(index int) LoreStorageGetFileItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageGetFileItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageGetFileItemArrayFFI) Clone() []LoreStorageGetFileItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageGetFileItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageGetFileItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageGetFileItemArray(arr []LoreStorageGetFileItem) (LoreStorageGetFileItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageGetFileItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageGetFileItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageGetFileItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageGetFileItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreStorageUploadItemArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreStorageUploadItemArray = []LoreStorageUploadItem

func (arr LoreStorageUploadItemArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreStorageUploadItemArrayFFI) Get(index int) LoreStorageUploadItem {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreStorageUploadItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreStorageUploadItemArrayFFI) Clone() []LoreStorageUploadItem {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreStorageUploadItem)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreStorageUploadItem, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreStorageUploadItemArray(arr []LoreStorageUploadItem) (LoreStorageUploadItemArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreStorageUploadItemArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreStorageUploadItemFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreStorageUploadItem(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreStorageUploadItemArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreRevisionTreeAddEntryArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreRevisionTreeAddEntryArray = []LoreRevisionTreeAddEntry

func (arr LoreRevisionTreeAddEntryArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreRevisionTreeAddEntryArrayFFI) Get(index int) LoreRevisionTreeAddEntry {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreRevisionTreeAddEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreRevisionTreeAddEntryArrayFFI) Clone() []LoreRevisionTreeAddEntry {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreRevisionTreeAddEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreRevisionTreeAddEntry, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreRevisionTreeAddEntryArray(arr []LoreRevisionTreeAddEntry) (LoreRevisionTreeAddEntryArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreRevisionTreeAddEntryArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreRevisionTreeAddEntryFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreRevisionTreeAddEntry(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreRevisionTreeAddEntryArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreRevisionTreeDeleteEntryArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreRevisionTreeDeleteEntryArray = []LoreRevisionTreeDeleteEntry

func (arr LoreRevisionTreeDeleteEntryArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreRevisionTreeDeleteEntryArrayFFI) Get(index int) LoreRevisionTreeDeleteEntry {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreRevisionTreeDeleteEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreRevisionTreeDeleteEntryArrayFFI) Clone() []LoreRevisionTreeDeleteEntry {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreRevisionTreeDeleteEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreRevisionTreeDeleteEntry, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreRevisionTreeDeleteEntryArray(arr []LoreRevisionTreeDeleteEntry) (LoreRevisionTreeDeleteEntryArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreRevisionTreeDeleteEntryArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreRevisionTreeDeleteEntryFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreRevisionTreeDeleteEntry(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreRevisionTreeDeleteEntryArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreRevisionTreeModifyEntryArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreRevisionTreeModifyEntryArray = []LoreRevisionTreeModifyEntry

func (arr LoreRevisionTreeModifyEntryArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreRevisionTreeModifyEntryArrayFFI) Get(index int) LoreRevisionTreeModifyEntry {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreRevisionTreeModifyEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreRevisionTreeModifyEntryArrayFFI) Clone() []LoreRevisionTreeModifyEntry {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreRevisionTreeModifyEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreRevisionTreeModifyEntry, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreRevisionTreeModifyEntryArray(arr []LoreRevisionTreeModifyEntry) (LoreRevisionTreeModifyEntryArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreRevisionTreeModifyEntryArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreRevisionTreeModifyEntryFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreRevisionTreeModifyEntry(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreRevisionTreeModifyEntryArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreRevisionTreeMoveEntryArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreRevisionTreeMoveEntryArray = []LoreRevisionTreeMoveEntry

func (arr LoreRevisionTreeMoveEntryArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreRevisionTreeMoveEntryArrayFFI) Get(index int) LoreRevisionTreeMoveEntry {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreRevisionTreeMoveEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreRevisionTreeMoveEntryArrayFFI) Clone() []LoreRevisionTreeMoveEntry {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreRevisionTreeMoveEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreRevisionTreeMoveEntry, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreRevisionTreeMoveEntryArray(arr []LoreRevisionTreeMoveEntry) (LoreRevisionTreeMoveEntryArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreRevisionTreeMoveEntryArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreRevisionTreeMoveEntryFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreRevisionTreeMoveEntry(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreRevisionTreeMoveEntryArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreRevisionTreeMetadataSetEntryArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreRevisionTreeMetadataSetEntryArray = []LoreRevisionTreeMetadataSetEntry

func (arr LoreRevisionTreeMetadataSetEntryArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreRevisionTreeMetadataSetEntryArrayFFI) Get(index int) LoreRevisionTreeMetadataSetEntry {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreRevisionTreeMetadataSetEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreRevisionTreeMetadataSetEntryArrayFFI) Clone() []LoreRevisionTreeMetadataSetEntry {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreRevisionTreeMetadataSetEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreRevisionTreeMetadataSetEntry, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreRevisionTreeMetadataSetEntryArray(arr []LoreRevisionTreeMetadataSetEntry) (LoreRevisionTreeMetadataSetEntryArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreRevisionTreeMetadataSetEntryArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreRevisionTreeMetadataSetEntryFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreRevisionTreeMetadataSetEntry(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreRevisionTreeMetadataSetEntryArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreRevisionTreeMetadataGetEntryArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreRevisionTreeMetadataGetEntryArray = []LoreRevisionTreeMetadataGetEntry

func (arr LoreRevisionTreeMetadataGetEntryArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreRevisionTreeMetadataGetEntryArrayFFI) Get(index int) LoreRevisionTreeMetadataGetEntry {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreRevisionTreeMetadataGetEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreRevisionTreeMetadataGetEntryArrayFFI) Clone() []LoreRevisionTreeMetadataGetEntry {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreRevisionTreeMetadataGetEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreRevisionTreeMetadataGetEntry, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreRevisionTreeMetadataGetEntryArray(arr []LoreRevisionTreeMetadataGetEntry) (LoreRevisionTreeMetadataGetEntryArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreRevisionTreeMetadataGetEntryArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreRevisionTreeMetadataGetEntryFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreRevisionTreeMetadataGetEntry(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreRevisionTreeMetadataGetEntryArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreRevisionTreeMetadataClearEntryArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreRevisionTreeMetadataClearEntryArray = []LoreRevisionTreeMetadataClearEntry

func (arr LoreRevisionTreeMetadataClearEntryArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreRevisionTreeMetadataClearEntryArrayFFI) Get(index int) LoreRevisionTreeMetadataClearEntry {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreRevisionTreeMetadataClearEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreRevisionTreeMetadataClearEntryArrayFFI) Clone() []LoreRevisionTreeMetadataClearEntry {
	if arr.Count == 0 {
		return nil
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	cDataSlice := unsafe.Slice((*LoreRevisionTreeMetadataClearEntry)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreRevisionTreeMetadataClearEntry, arr.Count)
	copy(result, cDataSlice)
	return result
}

func NewLoreRevisionTreeMetadataClearEntryArray(arr []LoreRevisionTreeMetadataClearEntry) (LoreRevisionTreeMetadataClearEntryArrayFFI, func()) {
	if len(arr) == 0 {
		return LoreRevisionTreeMetadataClearEntryArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Element type has separate Go and FFI representations; convert each item
	// through its NewXxx() builder so the FFI buffer contains FFI-layout structs.
	ffiArray := make([]LoreRevisionTreeMetadataClearEntryFFI, len(arr))
	cleanups := make([]func(), len(arr))
	for i := range arr {
		ffiArray[i], cleanups[i] = NewLoreRevisionTreeMetadataClearEntry(arr[i])
	}

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&ffiArray[0])
	arrayPtr := uintptr(unsafe.Pointer(&ffiArray[0]))

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
		pinner.Unpin()
	}

	return LoreRevisionTreeMetadataClearEntryArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(arr)),
	}, cleanup
}

type LoreTraceLocationFFI struct {
	/* The source file path. */
	File LoreString
	/* The line number in the source file. */
	Line uint32
	/* The column number in the source file. */
	Column uint32
	/* The context describing the operation at this location, or an empty
	string when the location has none. */
	Context LoreString
}

type LoreTraceLocation struct {
	/* The source file path. */
	File string
	/* The line number in the source file. */
	Line uint32
	/* The column number in the source file. */
	Column uint32
	/* The context describing the operation at this location, or an empty
	string when the location has none. */
	Context string
}

func NewLoreTraceLocation(opts LoreTraceLocation) (LoreTraceLocationFFI, func()) {
	valFile, cleanupFile := NewLoreString(opts.File)
	valContext, cleanupContext := NewLoreString(opts.Context)

	cleanup := func() {
		cleanupFile()
		cleanupContext()
	}

	return LoreTraceLocationFFI{
		File:    valFile,
		Line:    opts.Line,
		Column:  opts.Column,
		Context: valContext,
	}, cleanup
}

func (e *LoreTraceLocationFFI) Clone() LoreTraceLocation {
	return LoreTraceLocation{
		File:    e.File.Clone(),
		Line:    e.Line,
		Column:  e.Column,
		Context: e.Context.Clone(),
	}
}

type LoreErrorDetailFFI struct {
	/* The error's error code. `0` on success; `-1` for an internal error. */
	ErrorCode int32
	/* The error message, taken from the error's `Display` output. Empty on
	success. */
	Message LoreString
	/* The captured trace, one location per trace entry. Empty when
	`track-locations` is off or the error carries no trace. */
	TraceLocations LoreTraceLocationArrayFFI
}

type LoreErrorDetail struct {
	/* The error's error code. `0` on success; `-1` for an internal error. */
	ErrorCode int32
	/* The error message, taken from the error's `Display` output. Empty on
	success. */
	Message string
	/* The captured trace, one location per trace entry. Empty when
	`track-locations` is off or the error carries no trace. */
	TraceLocations LoreTraceLocationArray
}

func NewLoreErrorDetail(opts LoreErrorDetail) (LoreErrorDetailFFI, func()) {
	valMessage, cleanupMessage := NewLoreString(opts.Message)
	valTraceLocations, cleanupTraceLocations := NewLoreTraceLocationArray(opts.TraceLocations)

	cleanup := func() {
		cleanupMessage()
		cleanupTraceLocations()
	}

	return LoreErrorDetailFFI{
		ErrorCode:      opts.ErrorCode,
		Message:        valMessage,
		TraceLocations: valTraceLocations,
	}, cleanup
}

func (e *LoreErrorDetailFFI) Clone() LoreErrorDetail {
	return LoreErrorDetail{
		ErrorCode:      e.ErrorCode,
		Message:        e.Message.Clone(),
		TraceLocations: e.TraceLocations.Clone(),
	}
}

type LoreBranchPointFFI struct {
	/* The branch. */
	Branch LoreBranchId
	/* The revision on the branch. */
	Revision LoreHash
}

type LoreBranchPoint struct {
	/* The branch. */
	Branch LoreBranchId
	/* The revision on the branch. */
	Revision LoreHash
}

func NewLoreBranchPoint(opts LoreBranchPoint) (LoreBranchPointFFI, func()) {

	cleanup := func() {
	}

	return LoreBranchPointFFI{
		Branch:   opts.Branch,
		Revision: opts.Revision,
	}, cleanup
}

func (e *LoreBranchPointFFI) Clone() LoreBranchPoint {
	return LoreBranchPoint{
		Branch:   e.Branch,
		Revision: e.Revision,
	}
}

type LoreBranchDiffNodeDataFFI struct {
	/* File action applied to the node. */
	Action LoreFileAction
	/* Path of the node. */
	Path LoreString
	/* Set when the change was merged automatically. */
	Automerged uint8
	/* Previous path of the node when it was moved or copied. Empty otherwise. */
	FromPath LoreString
}

type LoreBranchDiffNodeData struct {
	/* File action applied to the node. */
	Action LoreFileAction
	/* Path of the node. */
	Path string
	/* Set when the change was merged automatically. */
	Automerged bool
	/* Previous path of the node when it was moved or copied. Empty otherwise. */
	FromPath string
}

func NewLoreBranchDiffNodeData(opts LoreBranchDiffNodeData) (LoreBranchDiffNodeDataFFI, func()) {
	valPath, cleanupPath := NewLoreString(opts.Path)
	valAutomerged, cleanupAutomerged := Newuint8(opts.Automerged)
	valFromPath, cleanupFromPath := NewLoreString(opts.FromPath)

	cleanup := func() {
		cleanupPath()
		cleanupAutomerged()
		cleanupFromPath()
	}

	return LoreBranchDiffNodeDataFFI{
		Action:     opts.Action,
		Path:       valPath,
		Automerged: valAutomerged,
		FromPath:   valFromPath,
	}, cleanup
}

func (e *LoreBranchDiffNodeDataFFI) Clone() LoreBranchDiffNodeData {
	return LoreBranchDiffNodeData{
		Action:     e.Action,
		Path:       e.Path.Clone(),
		Automerged: e.Automerged != 0,
		FromPath:   e.FromPath.Clone(),
	}
}

type LoreBranchSwitchDataFFI struct {
	/* Branch identifier. */
	Id LoreBranchId
	/* Branch name. */
	Name LoreString
	/* Latest revision known locally for the branch. */
	LatestLocal LoreHash
	/* Latest revision known on the remote for the branch. */
	LatestRemote LoreHash
	/* Revision the branch is switched to. */
	Revision LoreHash
	/* Where the branch exists: local, remote, or both. */
	Location LoreBranchLocation
}

type LoreBranchSwitchData struct {
	/* Branch identifier. */
	Id LoreBranchId
	/* Branch name. */
	Name string
	/* Latest revision known locally for the branch. */
	LatestLocal LoreHash
	/* Latest revision known on the remote for the branch. */
	LatestRemote LoreHash
	/* Revision the branch is switched to. */
	Revision LoreHash
	/* Where the branch exists: local, remote, or both. */
	Location LoreBranchLocation
}

func NewLoreBranchSwitchData(opts LoreBranchSwitchData) (LoreBranchSwitchDataFFI, func()) {
	valName, cleanupName := NewLoreString(opts.Name)

	cleanup := func() {
		cleanupName()
	}

	return LoreBranchSwitchDataFFI{
		Id:           opts.Id,
		Name:         valName,
		LatestLocal:  opts.LatestLocal,
		LatestRemote: opts.LatestRemote,
		Revision:     opts.Revision,
		Location:     opts.Location,
	}, cleanup
}

func (e *LoreBranchSwitchDataFFI) Clone() LoreBranchSwitchData {
	return LoreBranchSwitchData{
		Id:           e.Id,
		Name:         e.Name.Clone(),
		LatestLocal:  e.LatestLocal,
		LatestRemote: e.LatestRemote,
		Revision:     e.Revision,
		Location:     e.Location,
	}
}

type LoreFileResetCountDataFFI struct {
	/* Number of directories that were reset. */
	DirectoryResetCount uint64
	/* Number of directories that were deleted. */
	DirectoryDeleteCount uint64
	/* Number of files that were reset. */
	FileResetCount uint64
	/* Number of files that were deleted. */
	FileDeleteCount uint64
}

type LoreFileResetCountData struct {
	/* Number of directories that were reset. */
	DirectoryResetCount uint64
	/* Number of directories that were deleted. */
	DirectoryDeleteCount uint64
	/* Number of files that were reset. */
	FileResetCount uint64
	/* Number of files that were deleted. */
	FileDeleteCount uint64
}

func NewLoreFileResetCountData(opts LoreFileResetCountData) (LoreFileResetCountDataFFI, func()) {

	cleanup := func() {
	}

	return LoreFileResetCountDataFFI{
		DirectoryResetCount:  opts.DirectoryResetCount,
		DirectoryDeleteCount: opts.DirectoryDeleteCount,
		FileResetCount:       opts.FileResetCount,
		FileDeleteCount:      opts.FileDeleteCount,
	}, cleanup
}

func (e *LoreFileResetCountDataFFI) Clone() LoreFileResetCountData {
	return LoreFileResetCountData{
		DirectoryResetCount:  e.DirectoryResetCount,
		DirectoryDeleteCount: e.DirectoryDeleteCount,
		FileResetCount:       e.FileResetCount,
		FileDeleteCount:      e.FileDeleteCount,
	}
}

type LoreFileStageCountDataFFI struct {
	/* Number of directories staged as modified. */
	DirectoryModifyCount uint64
	/* Number of directories staged as added. */
	DirectoryAddCount uint64
	/* Number of directories staged as deleted. */
	DirectoryDeleteCount uint64
	/* Number of directories staged as moved. */
	DirectoryMoveCount uint64
	/* Number of files staged as modified. */
	FileModifyCount uint64
	/* Number of files staged as added. */
	FileAddCount uint64
	/* Number of files staged as deleted. */
	FileDeleteCount uint64
	/* Number of files staged as moved. */
	FileMoveCount uint64
	/* Total number of items processed. */
	TotalCount uint64
}

type LoreFileStageCountData struct {
	/* Number of directories staged as modified. */
	DirectoryModifyCount uint64
	/* Number of directories staged as added. */
	DirectoryAddCount uint64
	/* Number of directories staged as deleted. */
	DirectoryDeleteCount uint64
	/* Number of directories staged as moved. */
	DirectoryMoveCount uint64
	/* Number of files staged as modified. */
	FileModifyCount uint64
	/* Number of files staged as added. */
	FileAddCount uint64
	/* Number of files staged as deleted. */
	FileDeleteCount uint64
	/* Number of files staged as moved. */
	FileMoveCount uint64
	/* Total number of items processed. */
	TotalCount uint64
}

func NewLoreFileStageCountData(opts LoreFileStageCountData) (LoreFileStageCountDataFFI, func()) {

	cleanup := func() {
	}

	return LoreFileStageCountDataFFI{
		DirectoryModifyCount: opts.DirectoryModifyCount,
		DirectoryAddCount:    opts.DirectoryAddCount,
		DirectoryDeleteCount: opts.DirectoryDeleteCount,
		DirectoryMoveCount:   opts.DirectoryMoveCount,
		FileModifyCount:      opts.FileModifyCount,
		FileAddCount:         opts.FileAddCount,
		FileDeleteCount:      opts.FileDeleteCount,
		FileMoveCount:        opts.FileMoveCount,
		TotalCount:           opts.TotalCount,
	}, cleanup
}

func (e *LoreFileStageCountDataFFI) Clone() LoreFileStageCountData {
	return LoreFileStageCountData{
		DirectoryModifyCount: e.DirectoryModifyCount,
		DirectoryAddCount:    e.DirectoryAddCount,
		DirectoryDeleteCount: e.DirectoryDeleteCount,
		DirectoryMoveCount:   e.DirectoryMoveCount,
		FileModifyCount:      e.FileModifyCount,
		FileAddCount:         e.FileAddCount,
		FileDeleteCount:      e.FileDeleteCount,
		FileMoveCount:        e.FileMoveCount,
		TotalCount:           e.TotalCount,
	}
}

type LoreFileUnstageCountDataFFI struct {
	/* Number of directories that were unstaged. */
	DirectoryUnstagedCount uint64
	/* Number of directories that were discarded. */
	DirectoryDiscardedCount uint64
	/* Number of files that were unstaged. */
	FileUnstagedCount uint64
	/* Number of files that were discarded. */
	FileDiscardedCount uint64
	/* Total number of items processed. */
	TotalCount uint64
}

type LoreFileUnstageCountData struct {
	/* Number of directories that were unstaged. */
	DirectoryUnstagedCount uint64
	/* Number of directories that were discarded. */
	DirectoryDiscardedCount uint64
	/* Number of files that were unstaged. */
	FileUnstagedCount uint64
	/* Number of files that were discarded. */
	FileDiscardedCount uint64
	/* Total number of items processed. */
	TotalCount uint64
}

func NewLoreFileUnstageCountData(opts LoreFileUnstageCountData) (LoreFileUnstageCountDataFFI, func()) {

	cleanup := func() {
	}

	return LoreFileUnstageCountDataFFI{
		DirectoryUnstagedCount:  opts.DirectoryUnstagedCount,
		DirectoryDiscardedCount: opts.DirectoryDiscardedCount,
		FileUnstagedCount:       opts.FileUnstagedCount,
		FileDiscardedCount:      opts.FileDiscardedCount,
		TotalCount:              opts.TotalCount,
	}, cleanup
}

func (e *LoreFileUnstageCountDataFFI) Clone() LoreFileUnstageCountData {
	return LoreFileUnstageCountData{
		DirectoryUnstagedCount:  e.DirectoryUnstagedCount,
		DirectoryDiscardedCount: e.DirectoryDiscardedCount,
		FileUnstagedCount:       e.FileUnstagedCount,
		FileDiscardedCount:      e.FileDiscardedCount,
		TotalCount:              e.TotalCount,
	}
}

type LoreFragmentFFI struct {
	/* Flags */
	Flags uint32
	/* Payload size */
	SizePayload uint32
	/* Size of the uncompressed and reassembled content */
	SizeContent uint64
}

type LoreFragment struct {
	/* Flags */
	Flags uint32
	/* Payload size */
	SizePayload uint32
	/* Size of the uncompressed and reassembled content */
	SizeContent uint64
}

func NewLoreFragment(opts LoreFragment) (LoreFragmentFFI, func()) {

	cleanup := func() {
	}

	return LoreFragmentFFI{
		Flags:       opts.Flags,
		SizePayload: opts.SizePayload,
		SizeContent: opts.SizeContent,
	}, cleanup
}

func (e *LoreFragmentFFI) Clone() LoreFragment {
	return LoreFragment{
		Flags:       e.Flags,
		SizePayload: e.SizePayload,
		SizeContent: e.SizeContent,
	}
}

type LoreRepositoryCloneCountDataFFI struct {
	/* Number of files finished. */
	FileComplete uint64
	/* Number of files kept as they already matched. */
	FileRetain uint64
	/* Number of files replaced. */
	FileReplace uint64
	/* Total number of files discovered to process. */
	FileCount uint64
	/* Number of files currently being processed. */
	FileInflight uint64
	/* Number of fragment fetches currently in flight. */
	FragmentInflight uint64
	/* Number of bytes transferred so far. */
	BytesTransferred uint64
	/* Total number of bytes to transfer. */
	BytesTotal uint64
	/* Non-zero once file discovery has finished. */
	DiscoveryComplete uint8
}

type LoreRepositoryCloneCountData struct {
	/* Number of files finished. */
	FileComplete uint64
	/* Number of files kept as they already matched. */
	FileRetain uint64
	/* Number of files replaced. */
	FileReplace uint64
	/* Total number of files discovered to process. */
	FileCount uint64
	/* Number of files currently being processed. */
	FileInflight uint64
	/* Number of fragment fetches currently in flight. */
	FragmentInflight uint64
	/* Number of bytes transferred so far. */
	BytesTransferred uint64
	/* Total number of bytes to transfer. */
	BytesTotal uint64
	/* Non-zero once file discovery has finished. */
	DiscoveryComplete bool
}

func NewLoreRepositoryCloneCountData(opts LoreRepositoryCloneCountData) (LoreRepositoryCloneCountDataFFI, func()) {
	valDiscoveryComplete, cleanupDiscoveryComplete := Newuint8(opts.DiscoveryComplete)

	cleanup := func() {
		cleanupDiscoveryComplete()
	}

	return LoreRepositoryCloneCountDataFFI{
		FileComplete:      opts.FileComplete,
		FileRetain:        opts.FileRetain,
		FileReplace:       opts.FileReplace,
		FileCount:         opts.FileCount,
		FileInflight:      opts.FileInflight,
		FragmentInflight:  opts.FragmentInflight,
		BytesTransferred:  opts.BytesTransferred,
		BytesTotal:        opts.BytesTotal,
		DiscoveryComplete: valDiscoveryComplete,
	}, cleanup
}

func (e *LoreRepositoryCloneCountDataFFI) Clone() LoreRepositoryCloneCountData {
	return LoreRepositoryCloneCountData{
		FileComplete:      e.FileComplete,
		FileRetain:        e.FileRetain,
		FileReplace:       e.FileReplace,
		FileCount:         e.FileCount,
		FileInflight:      e.FileInflight,
		FragmentInflight:  e.FragmentInflight,
		BytesTransferred:  e.BytesTransferred,
		BytesTotal:        e.BytesTotal,
		DiscoveryComplete: e.DiscoveryComplete != 0,
	}
}

type LoreRevisionCommitCountDataFFI struct {
	/* Number of directories processed so far. */
	DirectoryCount uint64
	/* Total number of directories to process. */
	DirectoryTotal uint64
	/* Number of files processed so far. */
	FileCount uint64
	/* Total number of files to process. */
	FileTotal uint64
	/* Number of directories deleted. */
	DirectoryDeleteCount uint64
	/* Number of files modified. */
	FileModifyCount uint64
	/* Number of files deleted. */
	FileDeleteCount uint64
	/* Number of content bytes transferred so far. */
	BytesTransferred uint64
	/* Total number of content bytes to transfer. */
	BytesTotal uint64
	/* Set when file and directory discovery has finished. */
	DiscoveryComplete uint8
}

type LoreRevisionCommitCountData struct {
	/* Number of directories processed so far. */
	DirectoryCount uint64
	/* Total number of directories to process. */
	DirectoryTotal uint64
	/* Number of files processed so far. */
	FileCount uint64
	/* Total number of files to process. */
	FileTotal uint64
	/* Number of directories deleted. */
	DirectoryDeleteCount uint64
	/* Number of files modified. */
	FileModifyCount uint64
	/* Number of files deleted. */
	FileDeleteCount uint64
	/* Number of content bytes transferred so far. */
	BytesTransferred uint64
	/* Total number of content bytes to transfer. */
	BytesTotal uint64
	/* Set when file and directory discovery has finished. */
	DiscoveryComplete bool
}

func NewLoreRevisionCommitCountData(opts LoreRevisionCommitCountData) (LoreRevisionCommitCountDataFFI, func()) {
	valDiscoveryComplete, cleanupDiscoveryComplete := Newuint8(opts.DiscoveryComplete)

	cleanup := func() {
		cleanupDiscoveryComplete()
	}

	return LoreRevisionCommitCountDataFFI{
		DirectoryCount:       opts.DirectoryCount,
		DirectoryTotal:       opts.DirectoryTotal,
		FileCount:            opts.FileCount,
		FileTotal:            opts.FileTotal,
		DirectoryDeleteCount: opts.DirectoryDeleteCount,
		FileModifyCount:      opts.FileModifyCount,
		FileDeleteCount:      opts.FileDeleteCount,
		BytesTransferred:     opts.BytesTransferred,
		BytesTotal:           opts.BytesTotal,
		DiscoveryComplete:    valDiscoveryComplete,
	}, cleanup
}

func (e *LoreRevisionCommitCountDataFFI) Clone() LoreRevisionCommitCountData {
	return LoreRevisionCommitCountData{
		DirectoryCount:       e.DirectoryCount,
		DirectoryTotal:       e.DirectoryTotal,
		FileCount:            e.FileCount,
		FileTotal:            e.FileTotal,
		DirectoryDeleteCount: e.DirectoryDeleteCount,
		FileModifyCount:      e.FileModifyCount,
		FileDeleteCount:      e.FileDeleteCount,
		BytesTransferred:     e.BytesTransferred,
		BytesTotal:           e.BytesTotal,
		DiscoveryComplete:    e.DiscoveryComplete != 0,
	}
}

type LoreStorageRemoteConfigFFI struct {
	/* gRPC endpoint of the peer storage service; authenticated with the open call's `globals.identity` */
	RemoteUrl LoreString
}

type LoreStorageRemoteConfig struct {
	/* gRPC endpoint of the peer storage service; authenticated with the open call's `globals.identity` */
	RemoteUrl string
}

func NewLoreStorageRemoteConfig(opts LoreStorageRemoteConfig) (LoreStorageRemoteConfigFFI, func()) {
	valRemoteUrl, cleanupRemoteUrl := NewLoreString(opts.RemoteUrl)

	cleanup := func() {
		cleanupRemoteUrl()
	}

	return LoreStorageRemoteConfigFFI{
		RemoteUrl: valRemoteUrl,
	}, cleanup
}

func (e *LoreStorageRemoteConfigFFI) Clone() LoreStorageRemoteConfig {
	return LoreStorageRemoteConfig{
		RemoteUrl: e.RemoteUrl.Clone(),
	}
}

type LoreStoreFFI struct {
	/* Registry key; `0` is the reserved invalid/unregistered sentinel (zero-init = null handle) */
	HandleId uint64
}

type LoreStore struct {
	/* Registry key; `0` is the reserved invalid/unregistered sentinel (zero-init = null handle) */
	HandleId uint64
}

func NewLoreStore(opts LoreStore) (LoreStoreFFI, func()) {

	cleanup := func() {
	}

	return LoreStoreFFI{
		HandleId: opts.HandleId,
	}, cleanup
}

func (e *LoreStoreFFI) Clone() LoreStore {
	return LoreStore{
		HandleId: e.HandleId,
	}
}

type LoreStoragePutItemFFI struct {
	/* Caller-chosen id echoed back in `PUT_ITEM_COMPLETE` */
	Id uint64
	/* Target partition; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Dedup tag stored alongside the content hash in the resulting address */
	Context LoreContext
	/* Borrowed view into caller memory; bytes must live until `Complete` fires */
	Data LoreBytesFFI
	/* Opt into remote upload — honored on the remote path, ignored local-only */
	RemoteWrite uint8
	/* Tag the fragment with `PayloadLocalCachePriority` so future remote reads always cache it locally */
	LocalCache uint8
	/* Leaf fragment size cap for large buffers; `0` lets `write_content` choose. Ignored
	for buffers under `FRAGMENT_SIZE_THRESHOLD` */
	FixedSizeChunk uint64
}

type LoreStoragePutItem struct {
	/* Caller-chosen id echoed back in `PUT_ITEM_COMPLETE` */
	Id uint64
	/* Target partition; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Dedup tag stored alongside the content hash in the resulting address */
	Context LoreContext
	/* Borrowed view into caller memory; bytes must live until `Complete` fires */
	Data LoreBytes
	/* Opt into remote upload — honored on the remote path, ignored local-only */
	RemoteWrite bool
	/* Tag the fragment with `PayloadLocalCachePriority` so future remote reads always cache it locally */
	LocalCache bool
	/* Leaf fragment size cap for large buffers; `0` lets `write_content` choose. Ignored
	for buffers under `FRAGMENT_SIZE_THRESHOLD` */
	FixedSizeChunk uint64
}

func NewLoreStoragePutItem(opts LoreStoragePutItem) (LoreStoragePutItemFFI, func()) {
	valData, cleanupData := NewLoreBytes(opts.Data)
	valRemoteWrite, cleanupRemoteWrite := Newuint8(opts.RemoteWrite)
	valLocalCache, cleanupLocalCache := Newuint8(opts.LocalCache)

	cleanup := func() {
		cleanupData()
		cleanupRemoteWrite()
		cleanupLocalCache()
	}

	return LoreStoragePutItemFFI{
		Id:             opts.Id,
		Partition:      opts.Partition,
		Context:        opts.Context,
		Data:           valData,
		RemoteWrite:    valRemoteWrite,
		LocalCache:     valLocalCache,
		FixedSizeChunk: opts.FixedSizeChunk,
	}, cleanup
}

func (e *LoreStoragePutItemFFI) Clone() LoreStoragePutItem {
	return LoreStoragePutItem{
		Id:             e.Id,
		Partition:      e.Partition,
		Context:        e.Context,
		Data:           e.Data.Clone(),
		RemoteWrite:    e.RemoteWrite != 0,
		LocalCache:     e.LocalCache != 0,
		FixedSizeChunk: e.FixedSizeChunk,
	}
}

type LoreStorageGetItemFFI struct {
	/* Caller-chosen id echoed back in every event for this item */
	Id uint64
	/* Partition to read from; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Content address to read; `hash == Hash::default()` short-circuits to an empty buffer */
	Address LoreAddress
	/* First content byte to read, counted from the start of the decompressed content.
	Past the end of the content rejects with `INVALID_ARGUMENTS` */
	Offset uint64
	/* Content bytes to read from `offset`; `0` reads to the end. A range reaching past the
	end is clamped to it, so `GET_DATA` may carry fewer bytes than asked for */
	Length uint64
	/* Stream one `GET_DATA` per leaf fragment instead of a single reassembled buffer. A read
	that fails partway reports the failure on `GET_ITEM_COMPLETE` rather than ending short
	with a success code */
	Streaming uint8
	/* Cache fetched bytes back to the local store even without the producer's
	`PayloadLocalCachePriority` hint */
	LocalCache uint8
}

type LoreStorageGetItem struct {
	/* Caller-chosen id echoed back in every event for this item */
	Id uint64
	/* Partition to read from; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Content address to read; `hash == Hash::default()` short-circuits to an empty buffer */
	Address LoreAddress
	/* First content byte to read, counted from the start of the decompressed content.
	Past the end of the content rejects with `INVALID_ARGUMENTS` */
	Offset uint64
	/* Content bytes to read from `offset`; `0` reads to the end. A range reaching past the
	end is clamped to it, so `GET_DATA` may carry fewer bytes than asked for */
	Length uint64
	/* Stream one `GET_DATA` per leaf fragment instead of a single reassembled buffer. A read
	that fails partway reports the failure on `GET_ITEM_COMPLETE` rather than ending short
	with a success code */
	Streaming bool
	/* Cache fetched bytes back to the local store even without the producer's
	`PayloadLocalCachePriority` hint */
	LocalCache bool
}

func NewLoreStorageGetItem(opts LoreStorageGetItem) (LoreStorageGetItemFFI, func()) {
	valStreaming, cleanupStreaming := Newuint8(opts.Streaming)
	valLocalCache, cleanupLocalCache := Newuint8(opts.LocalCache)

	cleanup := func() {
		cleanupStreaming()
		cleanupLocalCache()
	}

	return LoreStorageGetItemFFI{
		Id:         opts.Id,
		Partition:  opts.Partition,
		Address:    opts.Address,
		Offset:     opts.Offset,
		Length:     opts.Length,
		Streaming:  valStreaming,
		LocalCache: valLocalCache,
	}, cleanup
}

func (e *LoreStorageGetItemFFI) Clone() LoreStorageGetItem {
	return LoreStorageGetItem{
		Id:         e.Id,
		Partition:  e.Partition,
		Address:    e.Address,
		Offset:     e.Offset,
		Length:     e.Length,
		Streaming:  e.Streaming != 0,
		LocalCache: e.LocalCache != 0,
	}
}

type LoreStorageGetResolvedItemFFI struct {
	/* Caller-chosen id echoed back in every event for this item */
	Id uint64
	/* Partition to resolve and read within; the zero/default partition rejects with
	`INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Mutable key to resolve, always read as `KeyType::Resolve` */
	Key LoreHash
	/* Paired with the resolved hash to address the immutable read; the mutable store yields
	only a hash. */
	Context LoreContext
	/* Stream one `GET_DATA` per leaf fragment instead of a single reassembled buffer, as
	`lore_storage_get` does. Peak memory then follows the fragment size rather than the
	content size, which is what makes a key naming something large usable. A read that fails
	partway reports the failure on `GET_ITEM_COMPLETE` rather than ending short with a
	success code */
	Streaming uint8
	/* Cache fetched bytes back to the local store even without the producer's
	`PayloadLocalCachePriority` hint */
	LocalCache uint8
}

type LoreStorageGetResolvedItem struct {
	/* Caller-chosen id echoed back in every event for this item */
	Id uint64
	/* Partition to resolve and read within; the zero/default partition rejects with
	`INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Mutable key to resolve, always read as `KeyType::Resolve` */
	Key LoreHash
	/* Paired with the resolved hash to address the immutable read; the mutable store yields
	only a hash. */
	Context LoreContext
	/* Stream one `GET_DATA` per leaf fragment instead of a single reassembled buffer, as
	`lore_storage_get` does. Peak memory then follows the fragment size rather than the
	content size, which is what makes a key naming something large usable. A read that fails
	partway reports the failure on `GET_ITEM_COMPLETE` rather than ending short with a
	success code */
	Streaming bool
	/* Cache fetched bytes back to the local store even without the producer's
	`PayloadLocalCachePriority` hint */
	LocalCache bool
}

func NewLoreStorageGetResolvedItem(opts LoreStorageGetResolvedItem) (LoreStorageGetResolvedItemFFI, func()) {
	valStreaming, cleanupStreaming := Newuint8(opts.Streaming)
	valLocalCache, cleanupLocalCache := Newuint8(opts.LocalCache)

	cleanup := func() {
		cleanupStreaming()
		cleanupLocalCache()
	}

	return LoreStorageGetResolvedItemFFI{
		Id:         opts.Id,
		Partition:  opts.Partition,
		Key:        opts.Key,
		Context:    opts.Context,
		Streaming:  valStreaming,
		LocalCache: valLocalCache,
	}, cleanup
}

func (e *LoreStorageGetResolvedItemFFI) Clone() LoreStorageGetResolvedItem {
	return LoreStorageGetResolvedItem{
		Id:         e.Id,
		Partition:  e.Partition,
		Key:        e.Key,
		Context:    e.Context,
		Streaming:  e.Streaming != 0,
		LocalCache: e.LocalCache != 0,
	}
}

type LoreStoragePutResolvedItemFFI struct {
	/* Caller-chosen id echoed back in `PUT_ITEM_COMPLETE` */
	Id uint64
	/* Target partition; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Mutable key to publish the stored hash under; a zero key rejects with `INVALID_ARGUMENTS` */
	Key LoreHash
	/* Dedup tag stored alongside the content hash in the resulting address, and the context a
	later `get_resolved` must read the key at */
	Context LoreContext
	/* Borrowed view into caller memory; bytes must live until `Complete` fires. A zero-length
	buffer removes the key's mapping instead of publishing one */
	Data LoreBytesFFI
	/* Also publish the content and the mapping to the remote; ignored when the handle has no
	remote or the call is offline/local */
	RemoteWrite uint8
	/* Tag the fragment with `PayloadLocalCachePriority` so future remote reads always cache it
	locally */
	LocalCache uint8
	/* Leaf fragment size cap for large buffers; `0` lets the writer choose. Ignored for buffers
	under `FRAGMENT_SIZE_THRESHOLD` */
	FixedSizeChunk uint64
}

type LoreStoragePutResolvedItem struct {
	/* Caller-chosen id echoed back in `PUT_ITEM_COMPLETE` */
	Id uint64
	/* Target partition; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Mutable key to publish the stored hash under; a zero key rejects with `INVALID_ARGUMENTS` */
	Key LoreHash
	/* Dedup tag stored alongside the content hash in the resulting address, and the context a
	later `get_resolved` must read the key at */
	Context LoreContext
	/* Borrowed view into caller memory; bytes must live until `Complete` fires. A zero-length
	buffer removes the key's mapping instead of publishing one */
	Data LoreBytes
	/* Also publish the content and the mapping to the remote; ignored when the handle has no
	remote or the call is offline/local */
	RemoteWrite bool
	/* Tag the fragment with `PayloadLocalCachePriority` so future remote reads always cache it
	locally */
	LocalCache bool
	/* Leaf fragment size cap for large buffers; `0` lets the writer choose. Ignored for buffers
	under `FRAGMENT_SIZE_THRESHOLD` */
	FixedSizeChunk uint64
}

func NewLoreStoragePutResolvedItem(opts LoreStoragePutResolvedItem) (LoreStoragePutResolvedItemFFI, func()) {
	valData, cleanupData := NewLoreBytes(opts.Data)
	valRemoteWrite, cleanupRemoteWrite := Newuint8(opts.RemoteWrite)
	valLocalCache, cleanupLocalCache := Newuint8(opts.LocalCache)

	cleanup := func() {
		cleanupData()
		cleanupRemoteWrite()
		cleanupLocalCache()
	}

	return LoreStoragePutResolvedItemFFI{
		Id:             opts.Id,
		Partition:      opts.Partition,
		Key:            opts.Key,
		Context:        opts.Context,
		Data:           valData,
		RemoteWrite:    valRemoteWrite,
		LocalCache:     valLocalCache,
		FixedSizeChunk: opts.FixedSizeChunk,
	}, cleanup
}

func (e *LoreStoragePutResolvedItemFFI) Clone() LoreStoragePutResolvedItem {
	return LoreStoragePutResolvedItem{
		Id:             e.Id,
		Partition:      e.Partition,
		Key:            e.Key,
		Context:        e.Context,
		Data:           e.Data.Clone(),
		RemoteWrite:    e.RemoteWrite != 0,
		LocalCache:     e.LocalCache != 0,
		FixedSizeChunk: e.FixedSizeChunk,
	}
}

type LoreStorageGetMetadataItemFFI struct {
	/* Caller-chosen id echoed back in `GET_METADATA_ITEM_COMPLETE` */
	Id uint64
	/* Partition to look up; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Content address to look up; `hash == Hash::default()` short-circuits to an empty fragment */
	Address LoreAddress
}

type LoreStorageGetMetadataItem struct {
	/* Caller-chosen id echoed back in `GET_METADATA_ITEM_COMPLETE` */
	Id uint64
	/* Partition to look up; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Content address to look up; `hash == Hash::default()` short-circuits to an empty fragment */
	Address LoreAddress
}

func NewLoreStorageGetMetadataItem(opts LoreStorageGetMetadataItem) (LoreStorageGetMetadataItemFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageGetMetadataItemFFI{
		Id:        opts.Id,
		Partition: opts.Partition,
		Address:   opts.Address,
	}, cleanup
}

func (e *LoreStorageGetMetadataItemFFI) Clone() LoreStorageGetMetadataItem {
	return LoreStorageGetMetadataItem{
		Id:        e.Id,
		Partition: e.Partition,
		Address:   e.Address,
	}
}

type LoreStorageObliterateItemFFI struct {
	/* Caller-chosen id echoed back in `OBLITERATE_ITEM_COMPLETE` */
	Id uint64
	/* Partition to delete from; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Content address to delete; absence on a side is idempotent success for that side */
	Address LoreAddress
}

type LoreStorageObliterateItem struct {
	/* Caller-chosen id echoed back in `OBLITERATE_ITEM_COMPLETE` */
	Id uint64
	/* Partition to delete from; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Content address to delete; absence on a side is idempotent success for that side */
	Address LoreAddress
}

func NewLoreStorageObliterateItem(opts LoreStorageObliterateItem) (LoreStorageObliterateItemFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageObliterateItemFFI{
		Id:        opts.Id,
		Partition: opts.Partition,
		Address:   opts.Address,
	}, cleanup
}

func (e *LoreStorageObliterateItemFFI) Clone() LoreStorageObliterateItem {
	return LoreStorageObliterateItem{
		Id:        e.Id,
		Partition: e.Partition,
		Address:   e.Address,
	}
}

type LoreStorageMutableLoadItemFFI struct {
	/* Caller-chosen id echoed back in `MUTABLE_LOAD_ITEM_COMPLETE` */
	Id uint64
	/* Partition (repository) to read from; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Key to read */
	Key LoreHash
	/* Kind of value the key refers to */
	KeyType LoreKeyType
}

type LoreStorageMutableLoadItem struct {
	/* Caller-chosen id echoed back in `MUTABLE_LOAD_ITEM_COMPLETE` */
	Id uint64
	/* Partition (repository) to read from; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Key to read */
	Key LoreHash
	/* Kind of value the key refers to */
	KeyType LoreKeyType
}

func NewLoreStorageMutableLoadItem(opts LoreStorageMutableLoadItem) (LoreStorageMutableLoadItemFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageMutableLoadItemFFI{
		Id:        opts.Id,
		Partition: opts.Partition,
		Key:       opts.Key,
		KeyType:   opts.KeyType,
	}, cleanup
}

func (e *LoreStorageMutableLoadItemFFI) Clone() LoreStorageMutableLoadItem {
	return LoreStorageMutableLoadItem{
		Id:        e.Id,
		Partition: e.Partition,
		Key:       e.Key,
		KeyType:   e.KeyType,
	}
}

type LoreStorageMutableStoreItemFFI struct {
	/* Caller-chosen id echoed back in `MUTABLE_STORE_ITEM_COMPLETE` */
	Id uint64
	/* Partition (repository) to write to; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Key to write */
	Key LoreHash
	/* Value to store; the null value (`Hash::default()`) removes the key */
	Value LoreHash
	/* Kind of value the key refers to */
	KeyType LoreKeyType
}

type LoreStorageMutableStoreItem struct {
	/* Caller-chosen id echoed back in `MUTABLE_STORE_ITEM_COMPLETE` */
	Id uint64
	/* Partition (repository) to write to; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Key to write */
	Key LoreHash
	/* Value to store; the null value (`Hash::default()`) removes the key */
	Value LoreHash
	/* Kind of value the key refers to */
	KeyType LoreKeyType
}

func NewLoreStorageMutableStoreItem(opts LoreStorageMutableStoreItem) (LoreStorageMutableStoreItemFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageMutableStoreItemFFI{
		Id:        opts.Id,
		Partition: opts.Partition,
		Key:       opts.Key,
		Value:     opts.Value,
		KeyType:   opts.KeyType,
	}, cleanup
}

func (e *LoreStorageMutableStoreItemFFI) Clone() LoreStorageMutableStoreItem {
	return LoreStorageMutableStoreItem{
		Id:        e.Id,
		Partition: e.Partition,
		Key:       e.Key,
		Value:     e.Value,
		KeyType:   e.KeyType,
	}
}

type LoreStorageMutableCompareAndSwapItemFFI struct {
	/* Caller-chosen id echoed back in `MUTABLE_COMPARE_AND_SWAP_ITEM_COMPLETE` */
	Id uint64
	/* Partition (repository) to act on; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Key to swap */
	Key LoreHash
	/* Value the key must currently hold for the swap to take effect (null matches an absent key) */
	Expected LoreHash
	/* Value to store when the swap takes effect; the null value removes the key */
	Value LoreHash
	/* Kind of value the key refers to */
	KeyType LoreKeyType
}

type LoreStorageMutableCompareAndSwapItem struct {
	/* Caller-chosen id echoed back in `MUTABLE_COMPARE_AND_SWAP_ITEM_COMPLETE` */
	Id uint64
	/* Partition (repository) to act on; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Key to swap */
	Key LoreHash
	/* Value the key must currently hold for the swap to take effect (null matches an absent key) */
	Expected LoreHash
	/* Value to store when the swap takes effect; the null value removes the key */
	Value LoreHash
	/* Kind of value the key refers to */
	KeyType LoreKeyType
}

func NewLoreStorageMutableCompareAndSwapItem(opts LoreStorageMutableCompareAndSwapItem) (LoreStorageMutableCompareAndSwapItemFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageMutableCompareAndSwapItemFFI{
		Id:        opts.Id,
		Partition: opts.Partition,
		Key:       opts.Key,
		Expected:  opts.Expected,
		Value:     opts.Value,
		KeyType:   opts.KeyType,
	}, cleanup
}

func (e *LoreStorageMutableCompareAndSwapItemFFI) Clone() LoreStorageMutableCompareAndSwapItem {
	return LoreStorageMutableCompareAndSwapItem{
		Id:        e.Id,
		Partition: e.Partition,
		Key:       e.Key,
		Expected:  e.Expected,
		Value:     e.Value,
		KeyType:   e.KeyType,
	}
}

type LoreStorageMutableListItemFFI struct {
	/* Caller-chosen id echoed back on every entry and the terminal event */
	Id uint64
	/* Partition (repository) to list; the zero/default partition lists every accessible partition */
	Partition LorePartition
	/* Kind of value to list */
	KeyType LoreKeyType
}

type LoreStorageMutableListItem struct {
	/* Caller-chosen id echoed back on every entry and the terminal event */
	Id uint64
	/* Partition (repository) to list; the zero/default partition lists every accessible partition */
	Partition LorePartition
	/* Kind of value to list */
	KeyType LoreKeyType
}

func NewLoreStorageMutableListItem(opts LoreStorageMutableListItem) (LoreStorageMutableListItemFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageMutableListItemFFI{
		Id:        opts.Id,
		Partition: opts.Partition,
		KeyType:   opts.KeyType,
	}, cleanup
}

func (e *LoreStorageMutableListItemFFI) Clone() LoreStorageMutableListItem {
	return LoreStorageMutableListItem{
		Id:        e.Id,
		Partition: e.Partition,
		KeyType:   e.KeyType,
	}
}

type LoreStorageCopyItemFFI struct {
	/* Caller-chosen id echoed back in `COPY_ITEM_COMPLETE` */
	Id uint64
	/* Source partition; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	SourcePartition LorePartition
	/* Destination partition; zero/default rejects, as does an exact `(source_partition, source
	context)` match (no-op) — a different `target_context` enables in-partition duplication */
	TargetPartition LorePartition
	/* Source content address; its `hash` carries over to the destination address unchanged */
	SourceAddress LoreAddress
	/* Dedup tag for the destination address `(target_partition, source_address.hash,
	target_context)`; may match the source tag or re-tag the payload */
	TargetContext LoreContext
}

type LoreStorageCopyItem struct {
	/* Caller-chosen id echoed back in `COPY_ITEM_COMPLETE` */
	Id uint64
	/* Source partition; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	SourcePartition LorePartition
	/* Destination partition; zero/default rejects, as does an exact `(source_partition, source
	context)` match (no-op) — a different `target_context` enables in-partition duplication */
	TargetPartition LorePartition
	/* Source content address; its `hash` carries over to the destination address unchanged */
	SourceAddress LoreAddress
	/* Dedup tag for the destination address `(target_partition, source_address.hash,
	target_context)`; may match the source tag or re-tag the payload */
	TargetContext LoreContext
}

func NewLoreStorageCopyItem(opts LoreStorageCopyItem) (LoreStorageCopyItemFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageCopyItemFFI{
		Id:              opts.Id,
		SourcePartition: opts.SourcePartition,
		TargetPartition: opts.TargetPartition,
		SourceAddress:   opts.SourceAddress,
		TargetContext:   opts.TargetContext,
	}, cleanup
}

func (e *LoreStorageCopyItemFFI) Clone() LoreStorageCopyItem {
	return LoreStorageCopyItem{
		Id:              e.Id,
		SourcePartition: e.SourcePartition,
		TargetPartition: e.TargetPartition,
		SourceAddress:   e.SourceAddress,
		TargetContext:   e.TargetContext,
	}
}

type LoreStoragePutFileItemFFI struct {
	/* Caller-chosen id echoed back in `PUT_ITEM_COMPLETE` */
	Id uint64
	/* Target partition; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Dedup tag stored alongside the content hash in the resulting address */
	Context LoreContext
	/* Source path; empty, missing, or non-file rejects with `INVALID_ARGUMENTS`; a zero-length
	file maps to the zero-hash address */
	Path LoreString
	/* Opt into remote upload — honored on the remote path, ignored local-only */
	RemoteWrite uint8
	/* Tag the resulting fragment with `PayloadLocalCachePriority` so future remote reads always cache it locally */
	LocalCache uint8
	/* Leaf fragment size cap for large files; `0` lets `write_content` choose */
	FixedSizeChunk uint64
}

type LoreStoragePutFileItem struct {
	/* Caller-chosen id echoed back in `PUT_ITEM_COMPLETE` */
	Id uint64
	/* Target partition; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Dedup tag stored alongside the content hash in the resulting address */
	Context LoreContext
	/* Source path; empty, missing, or non-file rejects with `INVALID_ARGUMENTS`; a zero-length
	file maps to the zero-hash address */
	Path string
	/* Opt into remote upload — honored on the remote path, ignored local-only */
	RemoteWrite bool
	/* Tag the resulting fragment with `PayloadLocalCachePriority` so future remote reads always cache it locally */
	LocalCache bool
	/* Leaf fragment size cap for large files; `0` lets `write_content` choose */
	FixedSizeChunk uint64
}

func NewLoreStoragePutFileItem(opts LoreStoragePutFileItem) (LoreStoragePutFileItemFFI, func()) {
	valPath, cleanupPath := NewLoreString(opts.Path)
	valRemoteWrite, cleanupRemoteWrite := Newuint8(opts.RemoteWrite)
	valLocalCache, cleanupLocalCache := Newuint8(opts.LocalCache)

	cleanup := func() {
		cleanupPath()
		cleanupRemoteWrite()
		cleanupLocalCache()
	}

	return LoreStoragePutFileItemFFI{
		Id:             opts.Id,
		Partition:      opts.Partition,
		Context:        opts.Context,
		Path:           valPath,
		RemoteWrite:    valRemoteWrite,
		LocalCache:     valLocalCache,
		FixedSizeChunk: opts.FixedSizeChunk,
	}, cleanup
}

func (e *LoreStoragePutFileItemFFI) Clone() LoreStoragePutFileItem {
	return LoreStoragePutFileItem{
		Id:             e.Id,
		Partition:      e.Partition,
		Context:        e.Context,
		Path:           e.Path.Clone(),
		RemoteWrite:    e.RemoteWrite != 0,
		LocalCache:     e.LocalCache != 0,
		FixedSizeChunk: e.FixedSizeChunk,
	}
}

type LoreStorageGetFileItemFFI struct {
	/* Caller-chosen id echoed back in `GET_ITEM_COMPLETE` */
	Id uint64
	/* Partition to read from; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Content address to read; `hash == Hash::default()` truncates `path` to zero bytes */
	Address LoreAddress
	/* Destination path; empty rejects with `INVALID_ARGUMENTS`. Multi-fragment writes
	stage via `<path>.loretmp` then atomically rename */
	Path LoreString
	/* First content byte to write, counted from the start of the decompressed content.
	Past the end of the content rejects with `INVALID_ARGUMENTS` */
	Offset uint64
	/* Content bytes to write from `offset`; `0` writes to the end. The file holds exactly the
	requested range starting at its own first byte, and is sized to it */
	Length uint64
	/* Cache fetched fragments back to the local store, not just write them to `path` */
	LocalCache uint8
}

type LoreStorageGetFileItem struct {
	/* Caller-chosen id echoed back in `GET_ITEM_COMPLETE` */
	Id uint64
	/* Partition to read from; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Content address to read; `hash == Hash::default()` truncates `path` to zero bytes */
	Address LoreAddress
	/* Destination path; empty rejects with `INVALID_ARGUMENTS`. Multi-fragment writes
	stage via `<path>.loretmp` then atomically rename */
	Path string
	/* First content byte to write, counted from the start of the decompressed content.
	Past the end of the content rejects with `INVALID_ARGUMENTS` */
	Offset uint64
	/* Content bytes to write from `offset`; `0` writes to the end. The file holds exactly the
	requested range starting at its own first byte, and is sized to it */
	Length uint64
	/* Cache fetched fragments back to the local store, not just write them to `path` */
	LocalCache bool
}

func NewLoreStorageGetFileItem(opts LoreStorageGetFileItem) (LoreStorageGetFileItemFFI, func()) {
	valPath, cleanupPath := NewLoreString(opts.Path)
	valLocalCache, cleanupLocalCache := Newuint8(opts.LocalCache)

	cleanup := func() {
		cleanupPath()
		cleanupLocalCache()
	}

	return LoreStorageGetFileItemFFI{
		Id:         opts.Id,
		Partition:  opts.Partition,
		Address:    opts.Address,
		Path:       valPath,
		Offset:     opts.Offset,
		Length:     opts.Length,
		LocalCache: valLocalCache,
	}, cleanup
}

func (e *LoreStorageGetFileItemFFI) Clone() LoreStorageGetFileItem {
	return LoreStorageGetFileItem{
		Id:         e.Id,
		Partition:  e.Partition,
		Address:    e.Address,
		Path:       e.Path.Clone(),
		Offset:     e.Offset,
		Length:     e.Length,
		LocalCache: e.LocalCache != 0,
	}
}

type LoreStorageUploadItemFFI struct {
	/* Caller-chosen id echoed back in `UPLOAD_ITEM_COMPLETE` */
	Id uint64
	/* Partition of the local content to push; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Local content address to push; `hash == Hash::default()` is no-op success with `already_durable=1` */
	Address LoreAddress
}

type LoreStorageUploadItem struct {
	/* Caller-chosen id echoed back in `UPLOAD_ITEM_COMPLETE` */
	Id uint64
	/* Partition of the local content to push; the zero/default partition rejects with `INVALID_ARGUMENTS` */
	Partition LorePartition
	/* Local content address to push; `hash == Hash::default()` is no-op success with `already_durable=1` */
	Address LoreAddress
}

func NewLoreStorageUploadItem(opts LoreStorageUploadItem) (LoreStorageUploadItemFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageUploadItemFFI{
		Id:        opts.Id,
		Partition: opts.Partition,
		Address:   opts.Address,
	}, cleanup
}

func (e *LoreStorageUploadItemFFI) Clone() LoreStorageUploadItem {
	return LoreStorageUploadItem{
		Id:        e.Id,
		Partition: e.Partition,
		Address:   e.Address,
	}
}

type LoreLogConfigFFI struct {
	/* Enable logging to a file (disabled by default) */
	File uint8
	/* Enable daily rolling logfile */
	FileRolling uint8
	/* Path to the log file */
	FilePath LoreString
	/* Prefix for log files */
	FilePrefix LoreString
	/* Minimum log level */
	Level LoreLogLevel
	/* Log categories bitflags (local, remote, transport) */
	Categories uint32
	/* Maximum log file size */
	FileMaxSize uint32
	/* Maximum log file count */
	FileMaxCount uint32
}

type LoreLogConfig struct {
	/* Enable logging to a file (disabled by default) */
	File bool
	/* Enable daily rolling logfile */
	FileRolling bool
	/* Path to the log file */
	FilePath string
	/* Prefix for log files */
	FilePrefix string
	/* Minimum log level */
	Level LoreLogLevel
	/* Log categories bitflags (local, remote, transport) */
	Categories uint32
	/* Maximum log file size */
	FileMaxSize uint32
	/* Maximum log file count */
	FileMaxCount uint32
}

func NewLoreLogConfig(opts LoreLogConfig) (LoreLogConfigFFI, func()) {
	valFile, cleanupFile := Newuint8(opts.File)
	valFileRolling, cleanupFileRolling := Newuint8(opts.FileRolling)
	valFilePath, cleanupFilePath := NewLoreString(opts.FilePath)
	valFilePrefix, cleanupFilePrefix := NewLoreString(opts.FilePrefix)

	cleanup := func() {
		cleanupFile()
		cleanupFileRolling()
		cleanupFilePath()
		cleanupFilePrefix()
	}

	return LoreLogConfigFFI{
		File:         valFile,
		FileRolling:  valFileRolling,
		FilePath:     valFilePath,
		FilePrefix:   valFilePrefix,
		Level:        opts.Level,
		Categories:   opts.Categories,
		FileMaxSize:  opts.FileMaxSize,
		FileMaxCount: opts.FileMaxCount,
	}, cleanup
}

func (e *LoreLogConfigFFI) Clone() LoreLogConfig {
	return LoreLogConfig{
		File:         e.File != 0,
		FileRolling:  e.FileRolling != 0,
		FilePath:     e.FilePath.Clone(),
		FilePrefix:   e.FilePrefix.Clone(),
		Level:        e.Level,
		Categories:   e.Categories,
		FileMaxSize:  e.FileMaxSize,
		FileMaxCount: e.FileMaxCount,
	}
}

type LoreRevisionTreeFFI struct {
	/* Registry key; `0` is the reserved invalid/unregistered sentinel (zero-init = null handle) */
	HandleId uint64
}

type LoreRevisionTree struct {
	/* Registry key; `0` is the reserved invalid/unregistered sentinel (zero-init = null handle) */
	HandleId uint64
}

func NewLoreRevisionTree(opts LoreRevisionTree) (LoreRevisionTreeFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeFFI{
		HandleId: opts.HandleId,
	}, cleanup
}

func (e *LoreRevisionTreeFFI) Clone() LoreRevisionTree {
	return LoreRevisionTree{
		HandleId: e.HandleId,
	}
}

type LoreRevisionTreeAddEntryFFI struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `ADD_COMPLETE` */
	EntryId uint64
	/* Parent for the new node; the invalid-node sentinel selects `parent_entry_index` */
	ParentNodeId uint32
	/* Index of an earlier entry in this batch whose new node is the parent;
	read only when `parent_node_id` is the invalid-node sentinel */
	ParentEntryIndex uint32
	/* UTF-8 name of the new child within its parent */
	Name LoreString
	/* `LoreNodeType` encoding: `DIRECTORY = 0`, `FILE = 1`, `LINK = 2` */
	Kind uint32
	/* POSIX permission bits for the new node */
	Mode uint16
	/* Content size in bytes (leaf nodes); `0` for a directory */
	Size uint64
	/* Content address `(hash, file_id context)` of the new node */
	Address LoreAddress
}

type LoreRevisionTreeAddEntry struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `ADD_COMPLETE` */
	EntryId uint64
	/* Parent for the new node; the invalid-node sentinel selects `parent_entry_index` */
	ParentNodeId uint32
	/* Index of an earlier entry in this batch whose new node is the parent;
	read only when `parent_node_id` is the invalid-node sentinel */
	ParentEntryIndex uint32
	/* UTF-8 name of the new child within its parent */
	Name string
	/* `LoreNodeType` encoding: `DIRECTORY = 0`, `FILE = 1`, `LINK = 2` */
	Kind uint32
	/* POSIX permission bits for the new node */
	Mode uint16
	/* Content size in bytes (leaf nodes); `0` for a directory */
	Size uint64
	/* Content address `(hash, file_id context)` of the new node */
	Address LoreAddress
}

func NewLoreRevisionTreeAddEntry(opts LoreRevisionTreeAddEntry) (LoreRevisionTreeAddEntryFFI, func()) {
	valName, cleanupName := NewLoreString(opts.Name)

	cleanup := func() {
		cleanupName()
	}

	return LoreRevisionTreeAddEntryFFI{
		EntryId:          opts.EntryId,
		ParentNodeId:     opts.ParentNodeId,
		ParentEntryIndex: opts.ParentEntryIndex,
		Name:             valName,
		Kind:             opts.Kind,
		Mode:             opts.Mode,
		Size:             opts.Size,
		Address:          opts.Address,
	}, cleanup
}

func (e *LoreRevisionTreeAddEntryFFI) Clone() LoreRevisionTreeAddEntry {
	return LoreRevisionTreeAddEntry{
		EntryId:          e.EntryId,
		ParentNodeId:     e.ParentNodeId,
		ParentEntryIndex: e.ParentEntryIndex,
		Name:             e.Name.Clone(),
		Kind:             e.Kind,
		Mode:             e.Mode,
		Size:             e.Size,
		Address:          e.Address,
	}
}

type LoreRevisionTreeDeleteEntryFFI struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `DELETE_COMPLETE` */
	EntryId uint64
	/* Root of the subtree to remove, including its transitive children */
	NodeId uint32
}

type LoreRevisionTreeDeleteEntry struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `DELETE_COMPLETE` */
	EntryId uint64
	/* Root of the subtree to remove, including its transitive children */
	NodeId uint32
}

func NewLoreRevisionTreeDeleteEntry(opts LoreRevisionTreeDeleteEntry) (LoreRevisionTreeDeleteEntryFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeDeleteEntryFFI{
		EntryId: opts.EntryId,
		NodeId:  opts.NodeId,
	}, cleanup
}

func (e *LoreRevisionTreeDeleteEntryFFI) Clone() LoreRevisionTreeDeleteEntry {
	return LoreRevisionTreeDeleteEntry{
		EntryId: e.EntryId,
		NodeId:  e.NodeId,
	}
}

type LoreRevisionTreeModifyEntryFFI struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `MODIFY_COMPLETE` */
	EntryId uint64
	/* Leaf node to rewrite; non-leaf targets are rejected */
	NodeId uint32
	/* New POSIX permission bits */
	Mode uint16
	/* New content size in bytes */
	Size uint64
	/* New content address; a zero `context` preserves the node's file id */
	Address LoreAddress
}

type LoreRevisionTreeModifyEntry struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `MODIFY_COMPLETE` */
	EntryId uint64
	/* Leaf node to rewrite; non-leaf targets are rejected */
	NodeId uint32
	/* New POSIX permission bits */
	Mode uint16
	/* New content size in bytes */
	Size uint64
	/* New content address; a zero `context` preserves the node's file id */
	Address LoreAddress
}

func NewLoreRevisionTreeModifyEntry(opts LoreRevisionTreeModifyEntry) (LoreRevisionTreeModifyEntryFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeModifyEntryFFI{
		EntryId: opts.EntryId,
		NodeId:  opts.NodeId,
		Mode:    opts.Mode,
		Size:    opts.Size,
		Address: opts.Address,
	}, cleanup
}

func (e *LoreRevisionTreeModifyEntryFFI) Clone() LoreRevisionTreeModifyEntry {
	return LoreRevisionTreeModifyEntry{
		EntryId: e.EntryId,
		NodeId:  e.NodeId,
		Mode:    e.Mode,
		Size:    e.Size,
		Address: e.Address,
	}
}

type LoreRevisionTreeMoveEntryFFI struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `MOVE_COMPLETE` */
	EntryId uint64
	/* Node to move; its `file_id` is preserved across the move */
	NodeId uint32
	/* Parent node the moved node is reparented under; its current parent renames it */
	DestinationParentId uint32
	/* UTF-8 name the moved node takes at the destination */
	DstName LoreString
}

type LoreRevisionTreeMoveEntry struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `MOVE_COMPLETE` */
	EntryId uint64
	/* Node to move; its `file_id` is preserved across the move */
	NodeId uint32
	/* Parent node the moved node is reparented under; its current parent renames it */
	DestinationParentId uint32
	/* UTF-8 name the moved node takes at the destination */
	DstName string
}

func NewLoreRevisionTreeMoveEntry(opts LoreRevisionTreeMoveEntry) (LoreRevisionTreeMoveEntryFFI, func()) {
	valDstName, cleanupDstName := NewLoreString(opts.DstName)

	cleanup := func() {
		cleanupDstName()
	}

	return LoreRevisionTreeMoveEntryFFI{
		EntryId:             opts.EntryId,
		NodeId:              opts.NodeId,
		DestinationParentId: opts.DestinationParentId,
		DstName:             valDstName,
	}, cleanup
}

func (e *LoreRevisionTreeMoveEntryFFI) Clone() LoreRevisionTreeMoveEntry {
	return LoreRevisionTreeMoveEntry{
		EntryId:             e.EntryId,
		NodeId:              e.NodeId,
		DestinationParentId: e.DestinationParentId,
		DstName:             e.DstName.Clone(),
	}
}

type LoreRevisionTreeMetadataSetEntryFFI struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `METADATA_SET_COMPLETE` */
	EntryId uint64
	/* Metadata key; a later entry naming it overwrites this one */
	Key LoreString
	/* Value to store, stored under the kind it carries */
	Value LoreMetadataFFI
}

type LoreRevisionTreeMetadataSetEntry struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `METADATA_SET_COMPLETE` */
	EntryId uint64
	/* Metadata key; a later entry naming it overwrites this one */
	Key string
	/* Value to store, stored under the kind it carries */
	Value LoreMetadata
}

func NewLoreRevisionTreeMetadataSetEntry(opts LoreRevisionTreeMetadataSetEntry) (LoreRevisionTreeMetadataSetEntryFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)
	valValue, cleanupValue := NewLoreMetadata(opts.Value)

	cleanup := func() {
		cleanupKey()
		cleanupValue()
	}

	return LoreRevisionTreeMetadataSetEntryFFI{
		EntryId: opts.EntryId,
		Key:     valKey,
		Value:   valValue,
	}, cleanup
}

func (e *LoreRevisionTreeMetadataSetEntryFFI) Clone() LoreRevisionTreeMetadataSetEntry {
	return LoreRevisionTreeMetadataSetEntry{
		EntryId: e.EntryId,
		Key:     e.Key.Clone(),
		Value:   e.Value.Clone(),
	}
}

type LoreRevisionTreeMetadataGetEntryFFI struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `METADATA_GET_COMPLETE` */
	EntryId uint64
	/* Metadata key to read */
	Key LoreString
}

type LoreRevisionTreeMetadataGetEntry struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `METADATA_GET_COMPLETE` */
	EntryId uint64
	/* Metadata key to read */
	Key string
}

func NewLoreRevisionTreeMetadataGetEntry(opts LoreRevisionTreeMetadataGetEntry) (LoreRevisionTreeMetadataGetEntryFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)

	cleanup := func() {
		cleanupKey()
	}

	return LoreRevisionTreeMetadataGetEntryFFI{
		EntryId: opts.EntryId,
		Key:     valKey,
	}, cleanup
}

func (e *LoreRevisionTreeMetadataGetEntryFFI) Clone() LoreRevisionTreeMetadataGetEntry {
	return LoreRevisionTreeMetadataGetEntry{
		EntryId: e.EntryId,
		Key:     e.Key.Clone(),
	}
}

type LoreRevisionTreeMetadataClearEntryFFI struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `METADATA_CLEAR_COMPLETE` */
	EntryId uint64
	/* Metadata key to remove; a key that is not set is a no-op */
	Key LoreString
}

type LoreRevisionTreeMetadataClearEntry struct {
	/* Caller-chosen id echoed back as `entry_id` on this entry's `METADATA_CLEAR_COMPLETE` */
	EntryId uint64
	/* Metadata key to remove; a key that is not set is a no-op */
	Key string
}

func NewLoreRevisionTreeMetadataClearEntry(opts LoreRevisionTreeMetadataClearEntry) (LoreRevisionTreeMetadataClearEntryFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)

	cleanup := func() {
		cleanupKey()
	}

	return LoreRevisionTreeMetadataClearEntryFFI{
		EntryId: opts.EntryId,
		Key:     valKey,
	}, cleanup
}

func (e *LoreRevisionTreeMetadataClearEntryFFI) Clone() LoreRevisionTreeMetadataClearEntry {
	return LoreRevisionTreeMetadataClearEntry{
		EntryId: e.EntryId,
		Key:     e.Key.Clone(),
	}
}

type LoreRevisionTreeCommitOptionsFFI struct {
	/* Also upload the new revision to remote (local-only by default) */
	RemoteWrite uint8
}

type LoreRevisionTreeCommitOptions struct {
	/* Also upload the new revision to remote (local-only by default) */
	RemoteWrite bool
}

func NewLoreRevisionTreeCommitOptions(opts LoreRevisionTreeCommitOptions) (LoreRevisionTreeCommitOptionsFFI, func()) {
	valRemoteWrite, cleanupRemoteWrite := Newuint8(opts.RemoteWrite)

	cleanup := func() {
		cleanupRemoteWrite()
	}

	return LoreRevisionTreeCommitOptionsFFI{
		RemoteWrite: valRemoteWrite,
	}, cleanup
}

func (e *LoreRevisionTreeCommitOptionsFFI) Clone() LoreRevisionTreeCommitOptions {
	return LoreRevisionTreeCommitOptions{
		RemoteWrite: e.RemoteWrite != 0,
	}
}

// String returns the hash as a lowercase hexadecimal string
// Implements fmt.Stringer interface
func (h LoreHash) String() string {
	return hex.EncodeToString(h.Data[:])
}

// String returns the context as a lowercase hexadecimal string
// Implements fmt.Stringer interface
func (c LoreContext) String() string {
	return hex.EncodeToString(c.Data[:])
}

// String returns the context as a lowercase hexadecimal string
// Implements fmt.Stringer interface
func (c LorePartition) String() string {
	return hex.EncodeToString(c.Data[:])
}

// String returns the context as a lowercase hexadecimal string
// Implements fmt.Stringer interface
func (c LoreBranchId) String() string {
	return hex.EncodeToString(c.Data[:])
}

// String returns the context as a lowercase hexadecimal string
// Implements fmt.Stringer interface
func (c LoreInstanceId) String() string {
	return hex.EncodeToString(c.Data[:])
}

// String returns the context as a lowercase hexadecimal string
// Implements fmt.Stringer interface
func (c LoreRepositoryId) String() string {
	return hex.EncodeToString(c.Data[:])
}

// The user's callback function type
// The event pointer is only valid during the callback
type LoreEventCallback func(*LoreEventFFI, uint64)

// LoreEventCallbackConfig configures the event callback
type LoreEventCallbackConfig struct {
	Callback    LoreEventCallback
	UserContext uint64 // Optional user-provided context value
}

// C-compatible structures for FFI
// These must match the C layout exactly

// LoreString is a C-compatible string representation
type LoreString struct {
	ptr    uintptr // char*
	Length uint64  // uintptr_t
}

type LoreHash struct {
	Data [32]uint8
}

func (h LoreHash) Clone() LoreHash {
	var result LoreHash
	copy(result.Data[:], h.Data[:])
	return result
}

type LoreContext struct {
	Data [16]uint8
}

func (c LoreContext) Clone() LoreContext {
	var result LoreContext
	copy(result.Data[:], c.Data[:])
	return result
}

type LorePartition struct {
	Data [16]uint8
}

func (c LorePartition) Clone() LorePartition {
	var result LorePartition
	copy(result.Data[:], c.Data[:])
	return result
}

type LoreInstanceId struct {
	Data [16]uint8
}

func (c LoreInstanceId) Clone() LoreInstanceId {
	var result LoreInstanceId
	copy(result.Data[:], c.Data[:])
	return result
}

type LoreAddress struct {
	Hash    LoreHash
	Context LoreContext
}

type LoreBranchId LoreContext

func (c LoreBranchId) Clone() LoreBranchId {
	var result LoreBranchId
	copy(result.Data[:], c.Data[:])
	return result
}

type LoreRepositoryId LorePartition

func (c LoreRepositoryId) Clone() LoreRepositoryId {
	var result LoreRepositoryId
	copy(result.Data[:], c.Data[:])
	return result
}

func (a LoreAddress) Clone() LoreAddress {
	return LoreAddress{
		Hash:    a.Hash.Clone(),
		Context: a.Context.Clone(),
	}
}

// String converts LoreString to a Go string
// Returns empty string if the pointer is null or length is 0
// Implements fmt.Stringer interface
func (u LoreString) String() string {
	if u.ptr == 0 || u.Length == 0 {
		return ""
	}
	return readCString(u.ptr, int(u.Length))
}

// Clone converts LoreString from FFI memory to Go string
func (u LoreString) Clone() string {
	return u.String()
}

// LoreStringArrayFFI is a C-compatible string array representation
type LoreStringArrayFFI struct {
	Ptr   uintptr // pointer to array of LoreString
	Count uint64  // uintptr_t
}

// Len returns the number of elements in the array
func (arr LoreStringArrayFFI) Len() int {
	return int(arr.Count)
}

// Get returns the string at the specified index
// Converts LoreString to Go string
// Panics if index is out of bounds
func (arr LoreStringArrayFFI) Get(index int) string {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreString)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index].String()
}

// Clone converts LoreStringArrayFFI from FFI memory to Go []string
// This creates a copy of the data that remains valid after the callback returns
func (arr LoreStringArrayFFI) Clone() []string {
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	if arr.Count == 0 {
		return nil
	}
	cStrings := unsafe.Slice((*LoreString)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]string, arr.Count)
	for i, s := range cStrings {
		result[i] = s.String()
	}
	return result
}

// LoreUint8ArrayFFI is a C-compatible array representation
type LoreUint8ArrayFFI struct {
	Ptr   uintptr // pointer to array of uint8
	Count uint64  // uintptr_t
}

type LoreUint8Array = []bool

// Len returns the number of elements in the array
func (arr LoreUint8ArrayFFI) Len() int {
	return int(arr.Count)
}

// Get returns the element at the specified index from FFI memory
// Panics if index is out of bounds
func (arr LoreUint8ArrayFFI) Get(index int) bool {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*uint8)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index] != 0
}

// Clone converts LoreUint8ArrayFFI from FFI memory to Go []bool
// This creates a copy of the data that remains valid after the callback returns
func (arr LoreUint8ArrayFFI) Clone() []bool {
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	if arr.Count == 0 {
		return nil
	}
	cUint8s := unsafe.Slice((*bool)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]bool, arr.Count)
	copy(result, cUint8s)
	return result
}

type LoreBinaryFFI struct {
	Payload uintptr
	Length  uint64
}

type LoreBinary = []byte

func NewLoreBinary(data LoreBinary) (LoreBinaryFFI, func()) {
	if len(data) == 0 {
		return LoreBinaryFFI{Payload: 0, Length: 0}, func() {}
	}

	// Allocate Go memory with null terminator
	// The slice will be kept alive by the caller until cleanup
	bytes := make([]byte, len(data))
	copy(bytes, data)

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&bytes[0])
	ptr := uintptr(unsafe.Pointer(&bytes[0]))

	cleanup := func() {
		pinner.Unpin()
	}

	return LoreBinaryFFI{
		Payload: ptr,
		Length:  uint64(len(data)),
	}, cleanup
}

func (data *LoreBinaryFFI) Clone() LoreBinary {
	if data.Payload == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	if data.Length == 0 {
		return nil
	}
	cDataSlice := unsafe.Slice((*byte)(unsafe.Pointer(data.Payload)), data.Length)
	result := make([]byte, data.Length)
	copy(result, cDataSlice)
	return result
}

type LoreBytesFFI = LoreBinaryFFI
type LoreBytes = LoreBinary

func NewLoreBytes(data LoreBinary) (LoreBinaryFFI, func()) {
	return NewLoreBinary(data)
}

// LoreMetadataFFI is a C-compatible representation of lore_metadata_t
type LoreMetadataFFI struct {
	Tag     LoreMetadataType
	padding [4]byte // Ensure union starts at 8-byte boundary
	// Union storage, sized and aligned for the largest member: LoreAddress,
	// a 32-byte hash followed by a 16-byte context. Read through the As*()
	// accessors and written by NewLoreMetadata, both via unsafe pointer
	// arithmetic.
	union [6]uint64
}

// LoreMetadata is a Go-compatible representation of lore_metadata_t
type LoreMetadata struct {
	Tag     LoreMetadataType
	Address *LoreAddress
	Boolean *bool
	Binary  *LoreBinary
	Context *LoreContext
	Hash    *LoreHash
	Numeric *uint64
	String  *string
}

// Precomputed offset to the union data in LoreMetadata
const loreMetadataUnionOffset = unsafe.Sizeof(LoreMetadataType(0)) + unsafe.Sizeof([4]byte{})

// AsAddress returns the metadata value as LoreAddress
// Only valid if Tag == LoreMetadataType_ADDRESS
func (m *LoreMetadataFFI) AsLoreAddress() *LoreAddress {
	unionPtr := unsafe.Add(unsafe.Pointer(m), loreMetadataUnionOffset)
	return (*LoreAddress)(unionPtr)
}

// AsBoolean returns the metadata value as bool
// Only valid if Tag == LoreMetadataType_BOOLEAN
func (m *LoreMetadataFFI) AsBoolean() bool {
	unionPtr := unsafe.Add(unsafe.Pointer(m), loreMetadataUnionOffset)
	return *(*uint8)(unionPtr) != 0
}

// AsBinary returns the metadata value as LoreBinaryFFI
// Only valid if Tag == LoreMetadataType_BINARY
func (m *LoreMetadataFFI) AsLoreBinary() *LoreBinaryFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(m), loreMetadataUnionOffset)
	return (*LoreBinaryFFI)(unionPtr)
}

// AsContext returns the metadata value as LoreContext
// Only valid if Tag == LoreMetadataType_CONTEXT
func (m *LoreMetadataFFI) AsLoreContext() *LoreContext {
	unionPtr := unsafe.Add(unsafe.Pointer(m), loreMetadataUnionOffset)
	return (*LoreContext)(unionPtr)
}

// AsHash returns the metadata value as LoreHash
// Only valid if Tag == LoreMetadataType_HASH
func (m *LoreMetadataFFI) AsLoreHash() *LoreHash {
	unionPtr := unsafe.Add(unsafe.Pointer(m), loreMetadataUnionOffset)
	return (*LoreHash)(unionPtr)
}

// AsNumeric returns the metadata value as uint64
// Only valid if Tag == LoreMetadataType_NUMERIC
func (m *LoreMetadataFFI) AsNumeric() uint64 {
	unionPtr := unsafe.Add(unsafe.Pointer(m), loreMetadataUnionOffset)
	return *(*uint64)(unionPtr)
}

// AsString returns the metadata value as LoreString
// Only valid if Tag == LoreMetadataType_STRING
func (m *LoreMetadataFFI) AsLoreString() *LoreString {
	unionPtr := unsafe.Add(unsafe.Pointer(m), loreMetadataUnionOffset)
	return (*LoreString)(unionPtr)
}

// NewLoreMetadata converts a LoreMetadata to LoreMetadataFFI, writing the
// payload into the union storage under the kind named by Tag. A nil payload
// pointer is written as that kind's zero value. Panics on a Tag that names no
// kind, including the zero value, which the C API does not accept.
func NewLoreMetadata(m LoreMetadata) (LoreMetadataFFI, func()) {
	ffi := LoreMetadataFFI{Tag: m.Tag}
	unionPtr := unsafe.Pointer(&ffi.union)

	switch m.Tag {
	case LoreMetadataType_ADDRESS:
		if m.Address != nil {
			*(*LoreAddress)(unionPtr) = *m.Address
		}
	case LoreMetadataType_BOOLEAN:
		if m.Boolean != nil && *m.Boolean {
			*(*uint8)(unionPtr) = 1
		}
	case LoreMetadataType_BINARY:
		var binary LoreBinary
		if m.Binary != nil {
			binary = *m.Binary
		}
		val, cleanup := NewLoreBinary(binary)
		*(*LoreBinaryFFI)(unionPtr) = val
		return ffi, cleanup
	case LoreMetadataType_CONTEXT:
		if m.Context != nil {
			*(*LoreContext)(unionPtr) = *m.Context
		}
	case LoreMetadataType_HASH:
		if m.Hash != nil {
			*(*LoreHash)(unionPtr) = *m.Hash
		}
	case LoreMetadataType_NUMERIC:
		if m.Numeric != nil {
			*(*uint64)(unionPtr) = *m.Numeric
		}
	case LoreMetadataType_STRING:
		var str string
		if m.String != nil {
			str = *m.String
		}
		val, cleanup := NewLoreString(str)
		*(*LoreString)(unionPtr) = val
		return ffi, cleanup
	default:
		panic(fmt.Sprintf("invalid metadata type: %d", m.Tag))
	}

	return ffi, func() {}
}

func (m *LoreMetadataFFI) Clone() LoreMetadata {
	switch m.Tag {
	case LoreMetadataType_ADDRESS:
		addr := m.AsLoreAddress().Clone()
		return LoreMetadata{
			Tag:     m.Tag,
			Address: &addr,
		}
	case LoreMetadataType_BOOLEAN:
		boolVal := m.AsBoolean()
		return LoreMetadata{
			Tag:     m.Tag,
			Boolean: &boolVal,
		}
	case LoreMetadataType_BINARY:
		binary := m.AsLoreBinary()
		// Copy binary data to Go-owned memory
		binaryCopy := binary.Clone()
		return LoreMetadata{
			Tag:    m.Tag,
			Binary: &binaryCopy,
		}
	case LoreMetadataType_CONTEXT:
		ctx := m.AsLoreContext().Clone()
		return LoreMetadata{
			Tag:     m.Tag,
			Context: &ctx,
		}
	case LoreMetadataType_HASH:
		hash := m.AsLoreHash().Clone()
		return LoreMetadata{
			Tag:  m.Tag,
			Hash: &hash,
		}
	case LoreMetadataType_NUMERIC:
		num := m.AsNumeric()
		return LoreMetadata{
			Tag:     m.Tag,
			Numeric: &num,
		}
	case LoreMetadataType_STRING:
		str := m.AsLoreString().String()
		return LoreMetadata{
			Tag:    m.Tag,
			String: &str,
		}
	default:
		return LoreMetadata{
			Tag: m.Tag,
		}
	}
}

// LoreEventCallbackConfigFFI is a C-compatible representation of lore_event_callback_config_t
type LoreEventCallbackConfigFFI struct {
	UserContext uint64
	FuncPtr     uintptr
}

// readCString reads a C string of known length into a Go string
func readCString(ptr uintptr, length int) string {
	if ptr == 0 || length == 0 {
		return ""
	}
	// Use unsafe.Slice to create a zero-copy view of the C memory
	// Then convert to string (which does copy, but avoids the manual loop)
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), length)
	return string(bytes)
}

// allocateCString allocates a C string from a Go string
// Returns the pointer and a cleanup function that releases the pin
func allocateCString(s string) (uintptr, func()) {
	if s == "" {
		return 0, func() {}
	}

	// Allocate Go memory with null terminator
	bytes := make([]byte, len(s)+1)
	copy(bytes, s)
	bytes[len(s)] = 0

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&bytes[0])
	ptr := uintptr(unsafe.Pointer(&bytes[0]))

	cleanup := func() {
		pinner.Unpin()
	}

	return ptr, cleanup
}

// NewLoreString converts a Go string to LoreString
func NewLoreString(s string) (LoreString, func()) {
	if s == "" {
		return LoreString{ptr: 0, Length: 0}, func() {}
	}

	p, cleanup := allocateCString(s)
	return LoreString{
		ptr:    p,
		Length: uint64(len(s)),
	}, cleanup
}

// NewLoreStringArray converts []string to LoreStringArrayFFI
func NewLoreStringArray(strs []string) (LoreStringArrayFFI, func()) {
	if len(strs) == 0 {
		return LoreStringArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// Allocate Go array of LoreString
	cStrings := make([]LoreString, len(strs))
	cleanups := make([]func(), 0, len(strs))

	// Fill the array
	for i, s := range strs {
		cStr, cleanup := NewLoreString(s)
		cleanups = append(cleanups, cleanup)
		cStrings[i] = cStr
	}

	// Pins the LoreString array itself; each element's own cleanup unpins the
	// character buffer it points at.
	var pinner runtime.Pinner
	pinner.Pin(&cStrings[0])
	arrayPtr := uintptr(unsafe.Pointer(&cStrings[0]))

	cleanup := func() {
		pinner.Unpin()
		for _, c := range cleanups {
			c()
		}
	}

	return LoreStringArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(strs)),
	}, cleanup
}

// NewLoreUint8Array converts []bool to LoreUint8ArrayFFI
func NewLoreUint8Array(values []bool) (LoreUint8ArrayFFI, func()) {
	if len(values) == 0 {
		return LoreUint8ArrayFFI{Ptr: 0, Count: 0}, func() {}
	}

	// cast bools as uint8, these are memory compatible
	boolsAsUint8s := unsafe.Slice((*uint8)(unsafe.Pointer(&values[0])), len(values))
	// Allocate Go array of uint8
	uint8s := make([]uint8, len(values))
	copy(uint8s, boolsAsUint8s)

	// Pin the buffer so it can be neither moved nor freed while the call is in flight.
	var pinner runtime.Pinner
	pinner.Pin(&uint8s[0])
	arrayPtr := uintptr(unsafe.Pointer(&uint8s[0]))

	cleanup := func() {
		pinner.Unpin()
	}

	return LoreUint8ArrayFFI{
		Ptr:   arrayPtr,
		Count: uint64(len(values)),
	}, cleanup
}

// Newuint8 converts Go bool to C uint8_t
func Newuint8(b bool) (uint8, func()) {
	if b {
		return 1, func() {}
	}
	return 0, func() {}
}
