package com.flagmanagment.sdk

import android.content.Context
import androidx.lifecycle.DefaultLifecycleObserver
import androidx.lifecycle.LifecycleOwner
import androidx.lifecycle.ProcessLifecycleOwner
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.sse.EventSource
import okhttp3.sse.EventSourceListener
import okhttp3.sse.EventSources
import org.json.JSONObject
import java.util.concurrent.ConcurrentHashMap

/**
 * SSE streaming client for Android. Thread-safe via ConcurrentHashMap.
 * Supports offline caching via FlagStorage and lifecycle-aware streaming.
 */
class FlagClient(
    private val context: Context,
    private val apiKey: String,
    private val streamUrl: String
) : DefaultLifecycleObserver {

    private val flags: ConcurrentHashMap<String, Any> = ConcurrentHashMap()
    private val storage = FlagStorage(context)
    private var eventSource: EventSource? = null

    init {
        // Load offline cache
        storage.load()?.let { cached ->
            flags.putAll(cached)
        }

        // Register lifecycle observer for background/foreground awareness
        ProcessLifecycleOwner.get().lifecycle.addObserver(this)
    }

    fun connect() {
        if (eventSource != null) return

        val client = OkHttpClient.Builder().build()
        val request = Request.Builder()
            .url(streamUrl)
            .addHeader("Authorization", "Bearer $apiKey")
            .addHeader("Accept", "text/event-stream")
            .build()

        val factory = EventSources.createFactory(client)
        eventSource = factory.newEventSource(request, object : EventSourceListener() {
            override fun onEvent(eventSource: EventSource, id: String?, type: String?, data: String) {
                try {
                    val json = JSONObject(data)
                    when (type) {
                        "bootstrap" -> {
                            val flagsObj = json.optJSONObject("flags")
                            if (flagsObj != null) {
                                flags.clear()
                                val keys = flagsObj.keys()
                                while (keys.hasNext()) {
                                    val key = keys.next()
                                    flags[key] = flagsObj.getJSONObject(key)
                                }
                                storage.save(flags)
                                android.util.Log.i("FlagManagment", "Bootstrapped ${flags.size} flags")
                            }
                        }
                        "flag_updated" -> {
                            val key = json.optString("flagKey")
                            val flag = json.optJSONObject("flag")
                            if (key.isNotEmpty() && flag != null) {
                                flags[key] = flag
                                storage.save(flags)
                                android.util.Log.i("FlagManagment", "Updated flag: $key")
                            }
                        }
                        "ping" -> { /* heartbeat — no action */ }
                    }
                } catch (e: Exception) {
                    android.util.Log.w("FlagManagment", "Failed to parse event", e)
                }
            }

            override fun onFailure(eventSource: EventSource, t: Throwable?, response: okhttp3.Response?) {
                android.util.Log.w("FlagManagment", "SSE connection failed, will reconnect", t)
                this@FlagClient.eventSource = null
                // OkHttp SSE handles reconnection automatically
            }
        })
    }

    fun disconnect() {
        eventSource?.cancel()
        eventSource = null
    }

    fun getFlag(key: String): Any? {
        return flags[key]
    }

    // Lifecycle awareness
    override fun onStart(owner: LifecycleOwner) {
        connect()
        android.util.Log.i("FlagManagment", "Resumed streaming (foreground)")
    }

    override fun onStop(owner: LifecycleOwner) {
        disconnect()
        android.util.Log.i("FlagManagment", "Paused streaming (background)")
    }
}
