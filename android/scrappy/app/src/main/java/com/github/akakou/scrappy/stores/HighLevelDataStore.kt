package com.github.akakou.scrappy.stores

import android.content.Context
import android.util.Log
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import com.github.akakou.scrappy.dataStore
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.runBlocking
import kotlin.coroutines.resume
import kotlin.coroutines.suspendCoroutine

class HighLevelDataStore(dataStoreAlias: String, context: Context) {
    val preference = stringPreferencesKey(dataStoreAlias)
    var context: Context

    init {
        this.context = context
    }

    suspend fun store(data: String) {
        context.dataStore.edit { settings ->
            settings[preference] = data
        }
    }

    suspend fun load():String {
        val loadFlow: Flow<String> = context.dataStore.data
            .map { preferences ->
                preferences[preference] ?: ""
            }

        val result: String = loadFlow.first()

        return result

    }
}