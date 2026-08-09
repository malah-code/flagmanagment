import Foundation
#if canImport(UIKit)
import UIKit
#endif

/// SSE streaming client for iOS. Thread-safe via a serial dispatch queue.
/// Supports offline caching via FlagStorage and lifecycle-aware streaming.
public class FlagClient {
    private let apiKey: String
    private let streamUrl: URL
    private var flags: [String: Any] = [:]
    private let storage = FlagStorage()
    private let queue = DispatchQueue(label: "com.flagmanagment.sdk.client", attributes: .concurrent)
    private var task: URLSessionDataTask?
    private var isConnected = false
    
    public init(apiKey: String, streamUrl: URL) {
        self.apiKey = apiKey
        self.streamUrl = streamUrl
        
        // Load offline cache on init
        if let cached = storage.load() {
            self.flags = cached
        }
        
        #if canImport(UIKit)
        // Lifecycle awareness: pause in background, resume in foreground
        NotificationCenter.default.addObserver(
            self,
            selector: #selector(appDidEnterBackground),
            name: UIApplication.didEnterBackgroundNotification,
            object: nil
        )
        NotificationCenter.default.addObserver(
            self,
            selector: #selector(appWillEnterForeground),
            name: UIApplication.willEnterForegroundNotification,
            object: nil
        )
        #endif
    }
    
    deinit {
        NotificationCenter.default.removeObserver(self)
        disconnect()
    }
    
    public func connect() {
        guard !isConnected else { return }
        isConnected = true
        
        var request = URLRequest(url: streamUrl)
        request.setValue("Bearer \(apiKey)", forHTTPHeaderField: "Authorization")
        request.setValue("text/event-stream", forHTTPHeaderField: "Accept")
        
        let session = URLSession(configuration: .default)
        task = session.dataTask(with: request) { [weak self] data, response, error in
            guard let self = self, let data = data else {
                // Reconnect with backoff
                DispatchQueue.global().asyncAfter(deadline: .now() + 5) { [weak self] in
                    self?.isConnected = false
                    self?.connect()
                }
                return
            }
            
            if let text = String(data: data, encoding: .utf8) {
                self.parseSSE(text)
            }
        }
        task?.resume()
        print("[flagmanagment-sdk] connected to \(streamUrl)")
    }
    
    public func disconnect() {
        task?.cancel()
        task = nil
        isConnected = false
    }
    
    /// Thread-safe flag retrieval.
    public func getFlag(key: String) -> Any? {
        var result: Any?
        queue.sync {
            result = flags[key]
        }
        return result
    }
    
    internal func updateFlags(_ newFlags: [String: Any]) {
        queue.async(flags: .barrier) { [weak self] in
            self?.flags = newFlags
            self?.storage.save(flags: newFlags)
        }
    }
    
    private func parseSSE(_ text: String) {
        let lines = text.components(separatedBy: "\n")
        var currentEvent = ""
        
        for line in lines {
            if line.hasPrefix("event:") {
                currentEvent = String(line.dropFirst(6)).trimmingCharacters(in: .whitespaces)
            } else if line.hasPrefix("data:") {
                let dataStr = String(line.dropFirst(5)).trimmingCharacters(in: .whitespaces)
                handleEvent(eventType: currentEvent, data: dataStr)
            } else if line.isEmpty {
                currentEvent = ""
            }
        }
    }
    
    private func handleEvent(eventType: String, data: String) {
        guard let jsonData = data.data(using: .utf8),
              let json = try? JSONSerialization.jsonObject(with: jsonData) as? [String: Any] else {
            return
        }
        
        switch eventType {
        case "bootstrap":
            if let flags = json["flags"] as? [String: Any] {
                updateFlags(flags)
                print("[flagmanagment-sdk] bootstrapped \(flags.count) flags")
            }
        case "flag_updated":
            if let flagKey = json["flagKey"] as? String, let flag = json["flag"] {
                queue.async(flags: .barrier) { [weak self] in
                    self?.flags[flagKey] = flag
                    if let allFlags = self?.flags {
                        self?.storage.save(flags: allFlags)
                    }
                }
                print("[flagmanagment-sdk] updated flag: \(flagKey)")
            }
        case "ping":
            break
        default:
            break
        }
    }
    
    #if canImport(UIKit)
    @objc private func appDidEnterBackground() {
        disconnect()
        print("[flagmanagment-sdk] paused streaming (background)")
    }
    
    @objc private func appWillEnterForeground() {
        connect()
        print("[flagmanagment-sdk] resumed streaming (foreground)")
    }
    #endif
}
