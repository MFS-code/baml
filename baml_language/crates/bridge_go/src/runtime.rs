use bridge_cffi::Buffer;
use std::ffi::CStr;

use crate::panic::ffi_safe_ptr;

#[unsafe(no_mangle)]
pub extern "C" fn version() -> Buffer {
    Buffer::from(env!("CARGO_PKG_VERSION").as_bytes().to_vec())
}

#[unsafe(no_mangle)]
pub extern "C" fn create_baml_runtime(
    root_path: *const libc::c_char,
    src_files_json: *const libc::c_char,
) -> *const libc::c_void {
    ffi_safe_ptr(|| {
        let root_path_str = unsafe { CStr::from_ptr(root_path) }
            .to_str()
            .map_err(|e| format!("Invalid root_path UTF-8: {e}"))?;
        let src_files_str = unsafe { CStr::from_ptr(src_files_json) }
            .to_str()
            .map_err(|e| format!("Invalid src_files_json UTF-8: {e}"))?;
        let src_files: std::collections::HashMap<String, String> =
            serde_json::from_str(src_files_str)
                .map_err(|e| format!("Invalid src_files JSON: {e}"))?;

        bridge_cffi::engine::initialize_runtime(root_path_str, src_files)
            .map_err(|e| format!("{e}"))?;

        Ok(std::ptr::dangling::<libc::c_void>())
    })
}

#[unsafe(no_mangle)]
pub extern "C" fn destroy_baml_runtime(_runtime: *const libc::c_void) {
    // No-op — runtime is global and persists for the process lifetime.
}
