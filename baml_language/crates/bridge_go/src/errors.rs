use bridge_cffi::error::BridgeError;

pub fn bridge_error_to_string(err: &BridgeError) -> String {
    match err {
        BridgeError::Ctypes(e) => format!("BamlError: BamlInvalidArgumentError: {e}"),
        BridgeError::NotInitialized => {
            "BamlError: BamlInvalidArgumentError: Engine not initialized".to_string()
        }
        BridgeError::NullFunctionName => {
            "BamlError: BamlInvalidArgumentError: Function name is null".to_string()
        }
        BridgeError::InvalidFunctionName(e) => {
            format!("BamlError: BamlInvalidArgumentError: Invalid function name: {e}")
        }
        BridgeError::FunctionNotFound { name } => {
            format!("BamlError: BamlInvalidArgumentError: Function not found: {name}")
        }
        BridgeError::MissingArgument { function, parameter } => {
            format!(
                "BamlError: BamlInvalidArgumentError: Missing argument '{parameter}' for function '{function}'"
            )
        }
        BridgeError::NotImplemented(msg) => {
            format!("BamlError: BamlInvalidArgumentError: Not implemented: {msg}")
        }
        BridgeError::DuplicateCallId(id) => {
            format!("BamlError: BamlInvalidArgumentError: Duplicate call ID: {id}")
        }
        BridgeError::ProjectNotInitialized => {
            "BamlError: BamlClientError: Project not initialized".to_string()
        }
        BridgeError::LockPoisoned => {
            "BamlError: BamlClientError: Lock poisoned".to_string()
        }
        BridgeError::Internal(msg) => format!("BamlError: BamlClientError: {msg}"),
        BridgeError::Runtime(re) => runtime_error_to_string(re),
    }
}

fn runtime_error_to_string(err: &bex_project::RuntimeError) -> String {
    use bex_project::RuntimeError;
    match err {
        RuntimeError::InvalidArgument { .. } => {
            format!("BamlError: BamlInvalidArgumentError: {err}")
        }
        RuntimeError::Engine(engine_err) => {
            use bex_project::EngineError;
            match engine_err {
                EngineError::FunctionNotFound { .. } => {
                    format!("BamlError: BamlInvalidArgumentError: {engine_err}")
                }
                EngineError::Cancelled => {
                    format!("BamlError: BamlCancelledError: {engine_err}")
                }
                _ => format!("BamlError: BamlClientError: {engine_err}"),
            }
        }
        _ => format!("BamlError: BamlClientError: {err}"),
    }
}
