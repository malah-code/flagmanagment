package com.flagmanagment.sdk;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.logging.Level;
import java.util.logging.Logger;

import org.json.JSONObject;

/**
 * SSE streaming client that connects to the FlagManagment backend,
 * maintains a thread-safe in-memory cache of flag definitions,
 * and supports exponential backoff reconnection.
 */
public class Client {
    private static final Logger LOG = Logger.getLogger(Client.class.getName());

    private final String apiKey;
    private final String streamUrl;
    private final Map<String, Object> flags = new ConcurrentHashMap<>();
    private final AtomicBoolean running = new AtomicBoolean(false);
    private volatile Thread streamThread;

    public Client(String apiKey, String streamUrl) {
        this.apiKey = apiKey;
        this.streamUrl = streamUrl;
    }

    /** Start background SSE streaming. */
    public void connect() {
        if (running.getAndSet(true)) return;

        streamThread = new Thread(() -> {
            int attempt = 0;
            while (running.get()) {
                try {
                    URL url = new URL(streamUrl);
                    HttpURLConnection conn = (HttpURLConnection) url.openConnection();
                    conn.setRequestMethod("GET");
                    conn.setRequestProperty("Authorization", "Bearer " + apiKey);
                    conn.setRequestProperty("Accept", "text/event-stream");
                    conn.setConnectTimeout(10_000);
                    conn.setReadTimeout(0); // no read timeout for SSE

                    if (conn.getResponseCode() != 200) {
                        conn.disconnect();
                        attempt++;
                        backoff(attempt);
                        continue;
                    }

                    // Connection established — reset backoff
                    attempt = 0;
                    LOG.info("[flagmanagment-sdk] connected to " + streamUrl);

                    try (BufferedReader reader = new BufferedReader(
                            new InputStreamReader(conn.getInputStream()))) {
                        String line;
                        String currentEvent = "";
                        while ((line = reader.readLine()) != null && running.get()) {
                            if (line.startsWith("event:")) {
                                currentEvent = line.substring(6).trim();
                            } else if (line.startsWith("data:")) {
                                String data = line.substring(5).trim();
                                handleEvent(currentEvent, data);
                            } else if (line.isEmpty()) {
                                currentEvent = "";
                            }
                        }
                    }
                    conn.disconnect();
                } catch (Exception e) {
                    LOG.log(Level.WARNING, "[flagmanagment-sdk] connection error", e);
                }

                if (running.get()) {
                    attempt++;
                    LOG.info("[flagmanagment-sdk] reconnecting (attempt " + attempt + ")...");
                    backoff(attempt);
                }
            }
        }, "flagmanagment-sse");
        streamThread.setDaemon(true);
        streamThread.start();
    }

    /** Gracefully stop the SSE stream. */
    public void shutdown() {
        running.set(false);
        if (streamThread != null) {
            streamThread.interrupt();
        }
    }

    private void handleEvent(String eventType, String data) {
        try {
            JSONObject json = new JSONObject(data);
            switch (eventType) {
                case "bootstrap":
                    JSONObject flagsObj = json.optJSONObject("flags");
                    if (flagsObj != null) {
                        flags.clear();
                        for (String key : flagsObj.keySet()) {
                            flags.put(key, flagsObj.getJSONObject(key));
                        }
                        LOG.info("[flagmanagment-sdk] bootstrapped " + flags.size() + " flags");
                    }
                    break;

                case "flag_updated":
                    String flagKey = json.optString("flagKey");
                    JSONObject flag = json.optJSONObject("flag");
                    if (flagKey != null && flag != null) {
                        flags.put(flagKey, flag);
                        LOG.info("[flagmanagment-sdk] updated flag: " + flagKey);
                    }
                    break;

                case "ping":
                    break;

                default:
                    break;
            }
        } catch (Exception e) {
            LOG.log(Level.WARNING, "[flagmanagment-sdk] failed to parse event", e);
        }
    }

    private void backoff(int attempt) {
        long delayMs = Math.min((long) attempt * attempt * 1000L, 60_000L);
        if (delayMs < 1000L) delayMs = 1000L;
        try {
            Thread.sleep(delayMs);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    /** Returns the raw flag object for the given key, or null. */
    public Object getFlag(String key) {
        return flags.get(key);
    }

    /** Returns the full flags map (for evaluation). */
    public Map<String, Object> getFlags() {
        return flags;
    }
}
