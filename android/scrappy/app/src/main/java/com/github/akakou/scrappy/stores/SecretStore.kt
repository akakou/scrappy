package com.github.akakou.scrappy.stores

import android.content.Context


class SecretStore (keyAlias: String, dataStoreAlias: String, context: Context){
    var dataStore : HighLevelDataStore
    var keyStore : HighLevelKeyStore

    init {
        dataStore = HighLevelDataStore(dataStoreAlias, context)
        keyStore = HighLevelKeyStore(keyAlias)
    }

    suspend fun store(secret: ByteArray) {
        val cipher = keyStore.encrypt(secret)
        dataStore.store(cipher)
    }

    fun load(): ByteArray {
        val cipher = dataStore.load()
        return keyStore.decrypt(cipher)
    }
}