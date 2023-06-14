package com.github.akakou.scrappy

import android.os.Bundle
import android.util.Base64
import android.view.View
import android.widget.EditText
import android_scrappy.Android_scrappy.androidJoin
import android_scrappy.Android_scrappy.androidSign
import androidx.appcompat.app.AppCompatActivity
import org.json.JSONObject


object Scrappy {
    var config : ByteArray? = null
    var protocol = "https://"

    fun parseLibraryResponse(response: String) : ByteArray {
        val jsonObject = JSONObject(response)
        val status = jsonObject.getString("status").toString()
        val base64Buffer = jsonObject.getString("buffer").toString()

        val buffer = Base64.decode(base64Buffer, Base64.DEFAULT)

        if(status != "ok") {
            val msg = String(buffer, Charsets.UTF_8)
            error("Scrappy not work: $msg")
        }

        return buffer
    }

    fun join(issuerDomain: String, ipk: ByteArray): ByteArray? {
        val response = androidJoin("http://$issuerDomain", ipk)
        config = parseLibraryResponse(response)

        return config
    }

    fun sign(origin: String, unixTime: Long): String {
        return androidSign(origin, unixTime, config)
    }
}

class MainActivity : AppCompatActivity() {
    var msg = ""
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        val base64Ipk = "Pv+DAwEBEE1pZGRsZUVuY29kZWRJUEsB/4QAAQUBAVgBCgABAVkBCgABAUMBCgABAlNYAQoAAQJTWQEKAAAA/+//hAFBA9PozgtFt37RUMaNvDfZxpBieWBd5YjwtJ21x95CjKgTSW9nQG1wewd0iUhXkez+5kq29FjRD6cDoCpZ7gz/nyEBQQISsnYUKfeTQfkdFtwx4qAoWm/qIEV2nkOAWXiFdsCsDnhpPh+R3mq9KxkC6fxLYPYQcDvnbznqBdoG0Ugob5fQASD7AKgvD/eIy/m+Fd10pLiA0k84SrTBHkUekvttOvXc8wEgBM+O1JEIY4li9+NSopfslm06hFUWzSxHE9FnEWm7UG0BIDnachXqfuCeUFaKdzNmiJ1N6kXGJwhXb/eEAoYMP4r/AA=="
        val issuerAddress = "http://192.168.10.105:8080"

        val ipk = Base64.decode(base64Ipk, Base64.DEFAULT)

        thread{
            val config = androidJoin(issuerAddress, ipk)
            msg = config
        }
    }

    fun onClick(view: View) {
        var editText = findViewById<EditText>(R.id.result_edit_text)
        editText.setText(msg)
    }
}
