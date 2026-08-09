import Foundation
import OpenFeature

public class FlagManagmentProvider: FeatureProvider {
    public var metadata: ProviderMetadata = FlagManagmentMetadata()
    public var hooks: [any Hook] = []
    
    private let client: FlagClient
    
    public init(client: FlagClient) {
        self.client = client
    }
    
    private func evaluateFlag<T>(key: String, defaultValue: T, context: EvaluationContext?) -> ProviderResolutionDetail<T> {
        guard let flag = client.getFlag(key: key) as? [String: Any] else {
            return ProviderResolutionDetail(value: defaultValue, reason: .error)
        }
        
        let enabled = flag["enabled"] as? Bool ?? false
        let defaultVariant = flag["defaultVariant"] as? String ?? ""
        let variants = flag["variants"] as? [String: Any] ?? [:]
        
        if !enabled {
            if let rawVal = variants[defaultVariant] {
                if let val = rawVal as? T {
                    return ProviderResolutionDetail(value: val, reason: .disabled)
                }
                return ProviderResolutionDetail(value: defaultValue, reason: .error)
            }
            return ProviderResolutionDetail(value: defaultValue, reason: .disabled)
        }
        
        if let targetingKey = context?.targetingKey,
           let rules = flag["rules"] as? [[String: Any]] {
            for rule in rules {
                if let rollout = rule["rollout"] as? [String: Int] {
                    let bucket = MurmurHash3.bucketUser(targetingKey)
                    var currentThreshold = 0
                    
                    let sortedVariants = rollout.keys.sorted()
                    for variant in sortedVariants {
                        if let weight = rollout[variant] {
                            currentThreshold += weight
                            if bucket < currentThreshold {
                                if let rawVal = variants[variant] {
                                    if let value = rawVal as? T {
                                        return ProviderResolutionDetail(value: value, reason: .targetingMatch)
                                    }
                                    return ProviderResolutionDetail(value: defaultValue, reason: .error)
                                }
                                return ProviderResolutionDetail(value: defaultValue, reason: .error)
                            }
                        }
                    }
                } else if let rolloutRaw = rule["rollout"] as? [String: Any] {
                    let bucket = MurmurHash3.bucketUser(targetingKey)
                    var currentThreshold = 0
                    
                    let sortedVariants = rolloutRaw.keys.sorted()
                    for variant in sortedVariants {
                        if let weight = rolloutRaw[variant] as? Int {
                            currentThreshold += weight
                            if bucket < currentThreshold {
                                if let rawVal = variants[variant] {
                                    if let value = rawVal as? T {
                                        return ProviderResolutionDetail(value: value, reason: .targetingMatch)
                                    }
                                    return ProviderResolutionDetail(value: defaultValue, reason: .error)
                                }
                                return ProviderResolutionDetail(value: defaultValue, reason: .error)
                            }
                        }
                    }
                }
            }
        }
        
        if let rawVal = variants[defaultVariant] {
            if let value = rawVal as? T {
                return ProviderResolutionDetail(value: value, reason: .static)
            }
            return ProviderResolutionDetail(value: defaultValue, reason: .error)
        }
        return ProviderResolutionDetail(value: defaultValue, reason: .static)
    }
    
    public func getBooleanEvaluation(key: String, defaultValue: Bool, context: EvaluationContext?) throws -> ProviderResolutionDetail<Bool> {
        return evaluateFlag(key: key, defaultValue: defaultValue, context: context)
    }
    
    public func getStringEvaluation(key: String, defaultValue: String, context: EvaluationContext?) throws -> ProviderResolutionDetail<String> {
        return evaluateFlag(key: key, defaultValue: defaultValue, context: context)
    }
    
    public func getIntegerEvaluation(key: String, defaultValue: Int64, context: EvaluationContext?) throws -> ProviderResolutionDetail<Int64> {
        return evaluateFlag(key: key, defaultValue: defaultValue, context: context)
    }
    
    public func getDoubleEvaluation(key: String, defaultValue: Double, context: EvaluationContext?) throws -> ProviderResolutionDetail<Double> {
        return evaluateFlag(key: key, defaultValue: defaultValue, context: context)
    }
    
    public func getObjectEvaluation(key: String, defaultValue: Value, context: EvaluationContext?) throws -> ProviderResolutionDetail<Value> {
        return ProviderResolutionDetail(value: defaultValue, reason: .error)
    }
}

class FlagManagmentMetadata: ProviderMetadata {
    var name: String = "FlagManagment-iOS-Provider"
}
