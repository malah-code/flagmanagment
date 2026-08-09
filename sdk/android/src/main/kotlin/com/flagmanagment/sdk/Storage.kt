package com.flagmanagment.sdk

import android.content.Context
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import org.json.JSONObject

class FlagStorage(context: Context) {
    private val masterKey = MasterKey.Builder(context)
        .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
        .build()

    private val sharedPreferences = EncryptedSharedPreferences.create(
        context,
        "flagmanagment_cache",
        masterKey,
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )

    fun save(flags: Map<String, Any>) {
        val json = JSONObject(flags).toString()
        sharedPreferences.edit().putString("flags_cache", json).apply()
    }

    fun load(): Map<String, Any>? {
        val jsonString = sharedPreferences.getString("flags_cache", null) ?: return null
        val json = JSONObject(jsonString)
        val map = mutableMapOf<String, Any>()
        
        json.keys().forEach { key ->
            map[key] = json.get(key)
        }
        return map
    }
}
