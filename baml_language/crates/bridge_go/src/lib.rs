mod callbacks;
mod errors;
mod functions;
mod handle;
mod panic;
mod runtime;

// Re-export Buffer type from bridge_cffi.
pub use bridge_cffi::Buffer;

/// Free a buffer returned by FFI functions.
/// Implemented directly here to avoid calling bridge_cffi::free_buffer
/// through the dynamic symbol, which would be infinitely recursive since
/// both functions export the same symbol name "_free_buffer".
#[unsafe(no_mangle)]
pub extern "C" fn free_buffer(buf: Buffer) {
    if !buf.ptr.is_null() {
        unsafe {
            // Buffer was created from a boxed slice (len == cap)
            let _ = Vec::from_raw_parts(buf.ptr as *mut u8, buf.len, buf.len);
        }
    }
}

#[unsafe(no_mangle)]
pub extern "C" fn flush_events() {
    bridge_cffi::flush_event_sink();
}
