package com.github.akakou.scrappy.stores

import android.content.Context
import android.util.Log


class SecretStore (keyAlias: String, dataStoreAlias: String, context: Context){
    var dataStore : HighLevelDataStore
    var keyStore : HighLevelKeyStore
    var context = context

    init {
        dataStore = HighLevelDataStore(dataStoreAlias, context)
        keyStore = HighLevelKeyStore(keyAlias, context)
    }

    suspend fun store(secret: ByteArray) {
        val cipher = keyStore.encrypt(secret)
        dataStore.store(cipher)
    }

    suspend fun load(): ByteArray {
        val cipher = dataStore.load()
        return keyStore.decrypt(cipher)
    }
}