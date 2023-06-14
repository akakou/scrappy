package com.github.akakou.scrappy

import android.util.Base64
import android_scrappy.Android_scrappy
import com.github.akakou.scrappy.stores.Stores
import org.json.JSONObject


class ScrappySigner(stores: Stores){
    var stores = stores
    var protocol = "https://"

    fun parseLibraryResponse(response: String) : JSONObject {
        val jsonObject = JSONObject(response)
        val status = jsonObject.getString("status").toString()

        if(status != "ok") {
            val error = jsonObject.getString("error").toString()
            error("Scrappy not work: $error")
        }

        return jsonObject
    }

    suspend fun join(issuerDomain: String, ipk: String) {
        val response = Android_scrappy.androidJoin("$protocol$issuerDomain", ipk)
        val config = parseLibraryResponse(response).getJSONObject("data")

        val cred = config["Cred"].toString()
        val base64Sk = config["SK"].toString()
        val bytesSK = Base64.decode(base64Sk, Base64.DEFAULT)

        stores.secret.store(bytesSK)
        stores.cred.store(cred)
        stores.ipk.store(ipk)
    }

    suspend fun sign(origin: String, unixTime: Long): String {
        val rawSK = stores.secret.load()
        val cred = stores.cred.load()
        val ipk = stores.ipk.load()

        val bytesSK = Base64.encode(rawSK, Base64.DEFAULT)
        val sk = String(bytesSK)

        val resp = Android_scrappy.androidSign(origin, unixTime, sk, cred, ipk)
        return parseLibraryResponse(resp).getString("data")
    }
}
