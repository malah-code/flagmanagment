use std::collections::HashMap;
use std::sync::RwLock;
use serde_json::Value;
use crate::evaluator::{EvaluationResult, Evaluator};

pub struct Client {
    flags: RwLock<HashMap<String, Value>>,
}

impl Client {
    pub fn new() -> Self {
        Self {
            flags: RwLock::new(HashMap::new()),
        }
    }

    pub fn set_flags(&self, flags: HashMap<String, Value>) {
        let mut lock = self.flags.write().unwrap();
        *lock = flags;
    }

    pub fn update_flag(&self, key: String, flag: Value) {
        let mut lock = self.flags.write().unwrap();
        lock.insert(key, flag);
    }

    pub fn get_boolean_value(
        &self,
        flag_key: &str,
        default_value: bool,
        targeting_key: Option<&str>,
    ) -> bool {
        let lock = self.flags.read().unwrap();
        if let Some(flag) = lock.get(flag_key) {
            let res = Evaluator::evaluate(flag, targeting_key, None);
            res.value.as_bool().unwrap_or(default_value)
        } else {
            default_value
        }
    }

    pub fn get_string_value(
        &self,
        flag_key: &str,
        default_value: &str,
        targeting_key: Option<&str>,
    ) -> String {
        let lock = self.flags.read().unwrap();
        if let Some(flag) = lock.get(flag_key) {
            let res = Evaluator::evaluate(flag, targeting_key, None);
            res.value
                .as_str()
                .map(|s| s.to_string())
                .unwrap_or_else(|| default_value.to_string())
        } else {
            default_value.to_string()
        }
    }

    pub fn evaluate_details(
        &self,
        flag_key: &str,
        targeting_key: Option<&str>,
    ) -> Option<EvaluationResult> {
        let lock = self.flags.read().unwrap();
        lock.get(flag_key).map(|flag| Evaluator::evaluate(flag, targeting_key, None))
    }
}
