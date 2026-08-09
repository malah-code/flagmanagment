using System;
using System.Collections.Concurrent;
using System.IO;
using System.Net.Http;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;

namespace FlagManagment.Sdk
{
    /// <summary>
    /// SSE streaming client that connects to the FlagManagment backend,
    /// maintains a thread-safe in-memory cache of flag definitions,
    /// and supports exponential backoff reconnection.
    /// </summary>
    public class Client
    {
        private readonly string _apiKey;
        private readonly string _streamUrl;
        private ConcurrentDictionary<string, JsonElement> _flags = new();
        private CancellationTokenSource _cts;
        private Task _streamTask;

        public Client(string apiKey, string streamUrl)
        {
            _apiKey = apiKey;
            _streamUrl = streamUrl;
        }

        /// <summary>Start background SSE streaming.</summary>
        public void Connect()
        {
            if (_streamTask != null) return;
            _cts = new CancellationTokenSource();
            _streamTask = Task.Run(() => StreamAsync(_cts.Token));
        }

        /// <summary>Gracefully stop the SSE stream.</summary>
        public void Shutdown()
        {
            _cts?.Cancel();
            _streamTask?.Wait(TimeSpan.FromSeconds(5));
        }

        private async Task StreamAsync(CancellationToken ct)
        {
            int attempt = 0;
            using var httpClient = new HttpClient { Timeout = Timeout.InfiniteTimeSpan };

            while (!ct.IsCancellationRequested)
            {
                try
                {
                    var request = new HttpRequestMessage(HttpMethod.Get, _streamUrl);
                    request.Headers.Add("Authorization", $"Bearer {_apiKey}");
                    request.Headers.Add("Accept", "text/event-stream");

                    using var response = await httpClient.SendAsync(
                        request, HttpCompletionOption.ResponseHeadersRead, ct);

                    if (!response.IsSuccessStatusCode)
                    {
                        attempt++;
                        await Backoff(attempt, ct);
                        continue;
                    }

                    // Connected — reset backoff
                    attempt = 0;
                    Console.WriteLine($"[flagmanagment-sdk] connected to {_streamUrl}");

                    using var stream = await response.Content.ReadAsStreamAsync(ct);
                    using var reader = new StreamReader(stream);

                    string currentEvent = "";
                    while (!reader.EndOfStream && !ct.IsCancellationRequested)
                    {
                        string line = await reader.ReadLineAsync();
                        if (line == null) break;

                        if (line.StartsWith("event:"))
                        {
                            currentEvent = line.Substring(6).Trim();
                        }
                        else if (line.StartsWith("data:"))
                        {
                            string data = line.Substring(5).Trim();
                            HandleEvent(currentEvent, data);
                        }
                        else if (line.Length == 0)
                        {
                            currentEvent = "";
                        }
                    }
                }
                catch (OperationCanceledException)
                {
                    return;
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"[flagmanagment-sdk] connection error: {ex.Message}");
                }

                if (!ct.IsCancellationRequested)
                {
                    attempt++;
                    Console.WriteLine($"[flagmanagment-sdk] reconnecting (attempt {attempt})...");
                    await Backoff(attempt, ct);
                }
            }
        }

        private void HandleEvent(string eventType, string data)
        {
            try
            {
                using var doc = JsonDocument.Parse(data);
                var root = doc.RootElement;

                switch (eventType)
                {
                    case "bootstrap":
                        if (root.TryGetProperty("flags", out var flagsElement))
                        {
                            var newFlags = new ConcurrentDictionary<string, JsonElement>();
                            foreach (var prop in flagsElement.EnumerateObject())
                            {
                                newFlags[prop.Name] = prop.Value.Clone();
                            }
                            _flags = newFlags;
                            Console.WriteLine($"[flagmanagment-sdk] bootstrapped {newFlags.Count} flags");
                        }
                        break;

                    case "flag_updated":
                        if (root.TryGetProperty("flagKey", out var keyEl) &&
                            root.TryGetProperty("flag", out var flagEl))
                        {
                            _flags[keyEl.GetString()!] = flagEl.Clone();
                            Console.WriteLine($"[flagmanagment-sdk] updated flag: {keyEl.GetString()}");
                        }
                        break;

                    case "ping":
                        break;
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine($"[flagmanagment-sdk] failed to parse event: {ex.Message}");
            }
        }

        private async Task Backoff(int attempt, CancellationToken ct)
        {
            int delaySeconds = Math.Min(attempt * attempt, 60);
            if (delaySeconds < 1) delaySeconds = 1;
            await Task.Delay(TimeSpan.FromSeconds(delaySeconds), ct);
        }

        /// <summary>Return the raw flag JSON for the given key, or null.</summary>
        public JsonElement? GetFlag(string key)
        {
            return _flags.TryGetValue(key, out var flag) ? flag : null;
        }
    }
}
