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
    lateinit var resultEditText: EditText
    lateinit var issuerURLEditText: EditText
    lateinit var ipkEditText: EditText

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        issuerURLEditText = findViewById<EditText>(R.id.issuer_edit_text)
        ipkEditText = findViewById(R.id.ipk_edit_text)
        resultEditText = findViewById<EditText>(R.id.result_edit_text)

        Scrappy.protocol = "http://"
    }

    fun onClick(view: View) {
        val issuerDomain = issuerURLEditText.text.toString()
        val base64ipk = ipkEditText.text.toString()

        val ipk = Base64.decode(base64ipk, Base64.DEFAULT)

        var msg = "ok"
        try {
            Scrappy.join(issuerDomain, ipk)
        } catch (e: java.lang.Exception) {
            msg = e.toString()
        }

        resultEditText.setText(msg)
    }
}
