package com.github.akakou.scrappy.stores

import android.content.Context
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import com.github.akakou.scrappy.dataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

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

    fun load():String {
        val loadFlow: Flow<String> = context.dataStore.data
            .map { preferences ->
                preferences[preference] ?: ""
            }

        return loadFlow.toString()
    }
}