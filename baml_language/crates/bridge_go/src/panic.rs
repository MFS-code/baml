use std::panic::{catch_unwind, AssertUnwindSafe};

pub fn ffi_safe<T, F: FnOnce() -> Result<T, String>>(f: F) -> Option<T> {
    match catch_unwind(AssertUnwindSafe(f)) {
        Ok(Ok(result)) => Some(result),
        Ok(Err(e)) => {
            eprintln!("FFI error: {e}");
            None
        }
        Err(panic_info) => {
            let msg = if let Some(s) = panic_info.downcast_ref::<&str>() {
                s.to_string()
            } else if let Some(s) = panic_info.downcast_ref::<String>() {
                s.clone()
            } else {
                "Unknown panic".to_string()
            };
            eprintln!("FFI panic: {msg}");
            None
        }
    }
}

pub fn ffi_safe_ptr<F: FnOnce() -> Result<*const libc::c_void, String>>(
    f: F,
) -> *const libc::c_void {
    ffi_safe(f).unwrap_or(std::ptr::null())
}
