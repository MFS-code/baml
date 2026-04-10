package cffi

/*
#cgo LDFLAGS: -ldl
#include "bridge.h"
*/
import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

var libHandle unsafe.Pointer

// Init loads the bridge_go shared library and resolves all symbols.
// Call this from TestMain or package init.
func Init(libraryPath string) error {
	cPath := C.CString(libraryPath)
	defer C.free(unsafe.Pointer(cPath))

	handle := C.dlopen(cPath, C.RTLD_LAZY|C.RTLD_LOCAL)
	if handle == nil {
		errStr := C.GoString(C.dlerror())
		return fmt.Errorf("dlopen(%s): %s", libraryPath, errStr)
	}
	libHandle = handle

	symbols := map[string]func(unsafe.Pointer){
		"version":              func(p unsafe.Pointer) { C.setVersionFn(p) },
		"create_baml_runtime":  func(p unsafe.Pointer) { C.setCreateBamlRuntimeFn(p) },
		"destroy_baml_runtime": func(p unsafe.Pointer) { C.setDestroyBamlRuntimeFn(p) },
		"free_buffer":          func(p unsafe.Pointer) { C.setFreeBufferFn(p) },
		"call_function":        func(p unsafe.Pointer) { C.setCallFunctionFn(p) },
		"register_callback":   func(p unsafe.Pointer) { C.setRegisterCallbackFn(p) },
		"cancel_function_call": func(p unsafe.Pointer) { C.setCancelFunctionCallFn(p) },
		"clone_handle":         func(p unsafe.Pointer) { C.setCloneHandleFn(p) },
		"release_handle":       func(p unsafe.Pointer) { C.setReleaseHandleFn(p) },
		"flush_events":         func(p unsafe.Pointer) { C.setFlushEventsFn(p) },
	}
	for name, setter := range symbols {
		sym, err := resolveSymbol(handle, name)
		if err != nil {
			return fmt.Errorf("resolving %s: %w", name, err)
		}
		setter(sym)
	}
	return nil
}

func resolveSymbol(handle unsafe.Pointer, name string) (unsafe.Pointer, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.dlerror() // clear
	sym := C.dlsym(handle, cName)
	if sym == nil {
		errStr := C.GoString(C.dlerror())
		return nil, fmt.Errorf("dlsym: %s", errStr)
	}
	return sym, nil
}

// FindLibrary locates the bridge_go shared library relative to the source tree.
func FindLibrary() (string, error) {
	// Walk up from this source file to find the cargo target directory
	_, thisFile, _, _ := runtime.Caller(0)
	crateDir := filepath.Dir(filepath.Dir(thisFile)) // cffi/ -> bridge_go/
	// Try common cargo output locations
	var ext string
	switch runtime.GOOS {
	case "darwin":
		ext = "dylib"
	case "windows":
		ext = "dll"
	default:
		ext = "so"
	}
	candidates := []string{
		filepath.Join(crateDir, "..", "..", "target", "debug", fmt.Sprintf("libbridge_go.%s", ext)),
		filepath.Join(crateDir, "..", "..", "target", "release", fmt.Sprintf("libbridge_go.%s", ext)),
	}
	for _, p := range candidates {
		abs, _ := filepath.Abs(p)
		if _, err := os.Stat(abs); err == nil {
			return abs, nil
		}
	}
	return "", fmt.Errorf("bridge_go shared library not found; tried: %v", candidates)
}

// Version returns the BAML engine version string.
func Version() string {
	buf := C.wrapVersion()
	if buf.ptr == nil || buf.len == 0 {
		return ""
	}
	s := C.GoStringN((*C.char)(unsafe.Pointer(buf.ptr)), C.int(buf.len))
	C.wrapFreeBuffer(buf.ptr, buf.len)
	return s
}

// CreateBamlRuntime initializes the global BAML runtime from source files.
// srcFilesJSON is a JSON-encoded map[string]string of filename -> content.
// Returns an opaque runtime pointer (sentinel, not a real pointer).
func CreateBamlRuntime(rootPath, srcFilesJSON string) (unsafe.Pointer, error) {
	cRootPath := C.CString(rootPath)
	defer C.free(unsafe.Pointer(cRootPath))
	cSrcFiles := C.CString(srcFilesJSON)
	defer C.free(unsafe.Pointer(cSrcFiles))

	ptr := C.wrapCreateBamlRuntime(cRootPath, cSrcFiles)
	if ptr == nil {
		return nil, fmt.Errorf("create_baml_runtime failed (returned null)")
	}
	return ptr, nil
}

// DestroyBamlRuntime is a no-op -- the runtime is global.
func DestroyBamlRuntime(ptr unsafe.Pointer) {
	C.wrapDestroyBamlRuntime(ptr)
}

// RegisterCallback registers a single Go callback function pointer with Rust.
// cb must be a C-callable function pointer (//export function cast to unsafe.Pointer).
func RegisterCallback(cb unsafe.Pointer) {
	C.wrapRegisterCallback((C.CallbackFn)(cb))
}

// CallFunction dispatches an async function call to Rust.
// Results and errors are delivered via the registered callback.
func CallFunction(functionName string, encodedArgs []byte, id uint32) {
	cName := C.CString(functionName)
	defer C.free(unsafe.Pointer(cName))

	var cArgs *C.uint8_t
	if len(encodedArgs) > 0 {
		cArgs = (*C.uint8_t)(unsafe.Pointer(&encodedArgs[0]))
	}

	C.wrapCallFunction((*C.char)(unsafe.Pointer(cName)), cArgs, C.size_t(len(encodedArgs)), C.uint32_t(id))
}

// CancelFunctionCall signals the Rust engine to cancel an in-flight function call.
func CancelFunctionCall(id uint32) {
	C.wrapCancelFunctionCall(C.uint32_t(id))
}

// CloneHandle clones an opaque handle, returning a new key.
func CloneHandle(key uint64, handleType int32) uint64 {
	return uint64(C.wrapCloneHandle(C.uint64_t(key), C.int32_t(handleType)))
}

// ReleaseHandle releases an opaque handle.
func ReleaseHandle(key uint64, handleType int32) {
	C.wrapReleaseHandle(C.uint64_t(key), C.int32_t(handleType))
}

// FlushEvents flushes the event sink.
func FlushEvents() {
	C.wrapFlushEvents()
}
