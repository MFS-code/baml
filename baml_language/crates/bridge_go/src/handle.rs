use bridge_ctypes::HANDLE_TABLE;

#[unsafe(no_mangle)]
pub extern "C" fn clone_handle(key: u64, _handle_type: i32) -> u64 {
    HANDLE_TABLE.clone_handle(key).unwrap_or(0)
}

#[unsafe(no_mangle)]
pub extern "C" fn release_handle(key: u64, _handle_type: i32) {
    HANDLE_TABLE.release(key);
}
