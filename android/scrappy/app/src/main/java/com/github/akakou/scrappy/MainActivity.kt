package com.github.akakou.scrappy

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.view.View
import android.widget.EditText
import android_scrappy.Android_scrappy.androidJoin
import kotlin.concurrent.thread
import android.util.Base64

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
