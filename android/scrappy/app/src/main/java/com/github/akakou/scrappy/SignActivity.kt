package com.github.akakou.scrappy

import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.view.View
import android.widget.Button
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import com.github.akakou.scrappy.stores.Stores
import kotlinx.coroutines.launch


class SignActivity : AppCompatActivity() {
    lateinit var callback: String
    lateinit var parsedCallback: Uri
    var timestamp : Long = 0
    lateinit var scrappySigner: ScrappySigner

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_sign)

        val i = intent
        val urlString = i.dataString

        val stores = Stores(this)
        scrappySigner = ScrappySigner(stores)
        scrappySigner.protocol = "http://"

        val uri = Uri.parse(urlString)
        val button = findViewById<Button>(R.id.callback_check_button)

        callback = uri.getQueryParameter("callback")!!
        parsedCallback = Uri.parse(callback)

        timestamp = uri.getQueryParameter("timestamp")?.toLong()!!

        button.text = "Do you come from ${callback} ?"
    }

    fun onClick(view : View) {
        lifecycleScope.launch {
            var signature = ""
            try {
                signature = scrappySigner.sign(callback, timestamp)
            } catch (e: java.lang.Exception) {
                Toast.makeText(this@SignActivity, "error: ${e.toString()}", Toast.LENGTH_LONG).show()
            }

            val url = "${scrappySigner.protocol}${callback}#${signature}"

            val browserIntent = Intent(Intent.ACTION_VIEW, Uri.parse(url))
            startActivity(browserIntent)
        }
    }
}