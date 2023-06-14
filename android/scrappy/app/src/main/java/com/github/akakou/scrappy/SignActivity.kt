package com.github.akakou.scrappy

import android.net.Uri
import android.os.Bundle
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity


class SignActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_sign)

        val i = intent
        val urlString = i.dataString

        val uri = Uri.parse(urlString)

        Toast.makeText(this, "hello, ${uri.getQueryParameter("test")}", Toast.LENGTH_LONG).show()
    }
}