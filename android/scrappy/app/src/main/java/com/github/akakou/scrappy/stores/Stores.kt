package com.github.akakou.scrappy.stores

import android.content.Context

class Stores(context: Context){
    var secret = SecretStore("scrappy-secret-2", "scrappy-secret", context)
    var cred = HighLevelDataStore("scrappy-credential", context)
    var ipk = HighLevelDataStore("scrappy-ipk", context)
}