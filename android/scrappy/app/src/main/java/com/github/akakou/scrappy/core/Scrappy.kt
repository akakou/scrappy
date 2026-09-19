package com.github.akakou.scrappy

import android.util.Base64
import android.util.Log
import android_scrappy.Android_scrappy
import com.github.akakou.scrappy.stores.AppDatabase
import com.github.akakou.scrappy.stores.SignerLog
import com.github.akakou.scrappy.stores.Stores
import org.json.JSONObject


class ScrappySigner(stores: Stores, db: AppDatabase?){
    var stores = stores
    var protocol = "https://"
    var db = db
    init {
        Android_scrappy.hello()
    }

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
//        val rawSK = stores.secret.load()
//        val cred = stores.cred.load()
//        val ipk = stores.ipk.load()
//        val bytesSK = Base64.encode(rawSK, Base64.DEFAULT)
//        val sk = String(bytesSK)

        val ipk = "Pv+DAwEBEE1pZGRsZUVuY29kZWRJUEsB/4QAAQUBAVgBCgABAVkBCgABAUMBCgABAlNYAQoAAQJTWQEKAAAA/+//hAFBA30MD1Iy22eykuAdyHNr6Tq6FZnuA9rPXt7g5MijeCFznKS/AiDG1JKVJErAH9OcJm6wlBpmOQd11c2qe6YRnCwBQQJFtU4fPvk3T+ztsyggBjqFnW0hX+ruRQ12oE96oZUsk6Hdy0MTMZAB6rXT7Mp/1Gx+XNHZ8MtWN0gLDG1d+IQjASDpNgHm6zoCm2NtdJffTyMpMLfMilELABmVsz8h56zL0AEggN1aTJ0ITuHJcGDEPNc5YGhDTh5vXY66fqkni0yaqVoBIEn/vaanooPDqFtDPmWqry33TfCSB+n7iKp7n2uA3yikAA=="
        val cred = "Pf+FAwEBF01pZGRsZUVuY29kZWRDcmVkZW50aWFsAf+GAAEEAQFBAQoAAQFCAQoAAQFDAQoAAQFEAQoAAAD/j/+GASED1drR5R/6p7ek9+mxtv4kFnIyozi2qkydYkP9lixXZ90BIQJf2LyMysShJ8bSMmlcI9L641LoztRG7WqgY/eZgWlPJgEhAsWoIN+FPObSo5n3Ximjud81oTwY3yIbIJYzynSJCrJUASECN7BTDejfB8LMMGi5vD1Xjd39sO1Y7w6t8jRDWyOyiTwA"
        val sk = "5dsrSqK+PyiQZ2UoU+VWKKGhhCgocZxp1f7tfh/2Cf0="



        val basename = "$origin:$unixTime"

        val beforeTime1 = System.currentTimeMillis()
        val logDao = db!!.signerLogDao()
        val hasExist = logDao.hasExist(basename)

        if (hasExist == 1) {
            error("$basename has exist!")
        }

        val log = SignerLog(0, basename)
        logDao.insertAll(log)
        val afterTime1 = System.currentTimeMillis()

        Log.d("PERFORMANCE DB", "signing cost: ${afterTime1 - beforeTime1}")

        val beforeTime2 = System.currentTimeMillis()
        val resp = Android_scrappy.androidSign(origin, unixTime, sk, cred, ipk)

        val afterTime2 = System.currentTimeMillis()
        Log.d("PERFORMANCE ECDAA", "signing cost: ${afterTime2 - beforeTime2}")

        return parseLibraryResponse(resp).getString("data")
    }
}
