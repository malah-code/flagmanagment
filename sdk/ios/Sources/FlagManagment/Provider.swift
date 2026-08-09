import Foundation
import OpenFeature

public class FlagManagmentProvider: FeatureProvider {
    public var metadata: ProviderMetadata = FlagManagmentMetadata()
    public var hooks: [any Hook] = []
    
    private let client: FlagClient
    
    public init(client: FlagClient) {
        self.client = client
    }
    
    public func getBooleanEvaluation(key: String, defaultValue: Bool, context: EvaluationContext?) throws -> ProviderResolutionDetail<Bool> {
        guard let flag = client.getFlag(key: key) as? [String: Any],
              let variants = flag["variants"] as? [String: Any],
              let defaultVariant = flag["defaultVariant"] as? String,
              let value = variants[defaultVariant] as? Bool else {
            return ProviderResolutionDetail(value: defaultValue, reason: .error)
        }
        return ProviderResolutionDetail(value: value, reason: .static)
    }
    
    public func getStringEvaluation(key: String, defaultValue: String, context: EvaluationContext?) throws -> ProviderResolutionDetail<String> {
        guard let flag = client.getFlag(key: key) as? [String: Any],
              let variants = flag["variants"] as? [String: Any],
              let defaultVariant = flag["defaultVariant"] as? String,
              let value = variants[defaultVariant] as? String else {
            return ProviderResolutionDetail(value: defaultValue, reason: .error)
        }
        return ProviderResolutionDetail(value: value, reason: .static)
    }
    
    public func getIntegerEvaluation(key: String, defaultValue: Int64, context: EvaluationContext?) throws -> ProviderResolutionDetail<Int64> {
        guard let flag = client.getFlag(key: key) as? [String: Any],
              let variants = flag["variants"] as? [String: Any],
              let defaultVariant = flag["defaultVariant"] as? String,
              let value = variants[defaultVariant] as? Int64 else {
            return ProviderResolutionDetail(value: defaultValue, reason: .error)
        }
        return ProviderResolutionDetail(value: value, reason: .static)
    }
    
    public func getDoubleEvaluation(key: String, defaultValue: Double, context: EvaluationContext?) throws -> ProviderResolutionDetail<Double> {
        guard let flag = client.getFlag(key: key) as? [String: Any],
              let variants = flag["variants"] as? [String: Any],
              let defaultVariant = flag["defaultVariant"] as? String,
              let value = variants[defaultVariant] as? Double else {
            return ProviderResolutionDetail(value: defaultValue, reason: .error)
        }
        return ProviderResolutionDetail(value: value, reason: .static)
    }
    
    public func getObjectEvaluation(key: String, defaultValue: Value, context: EvaluationContext?) throws -> ProviderResolutionDetail<Value> {
        return ProviderResolutionDetail(value: defaultValue, reason: .error)
    }
}

class FlagManagmentMetadata: ProviderMetadata {
    var name: String = "FlagManagment-iOS-Provider"
}
