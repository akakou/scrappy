package com.github.akakou.scrappy

import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Button
import android.widget.Toast
import android_scrappy.Android_scrappy
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import androidx.room.Room
import com.github.akakou.scrappy.stores.AppDatabase
import com.github.akakou.scrappy.stores.SignerLog
import com.github.akakou.scrappy.stores.Stores
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.launch
import java.net.URL
import kotlin.concurrent.thread


class SignActivity : AppCompatActivity() {
    lateinit var callback: String
    lateinit var parsedCallback: URL
    var period : Long = 0
    lateinit var scrappySigner: ScrappySigner

    lateinit var db: AppDatabase
    var allLogs: List<SignerLog>? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_sign)

        db = Room.databaseBuilder(
            applicationContext,
            AppDatabase::class.java, "database-name"
        ).build()

        val i = intent
        val urlString = i.dataString

        val stores = Stores(this)
        scrappySigner = ScrappySigner(stores, db)

        val uri = Uri.parse(urlString)
        val button = findViewById<Button>(R.id.callback_check_button)

        callback = uri.getQueryParameter("callback")!!
        parsedCallback = URL(callback)

        period = uri.getQueryParameter("period")?.toLong()!!
        button.text = "Do you come from ${parsedCallback.host}:${parsedCallback.port} ?"

        thread {
            allLogs = db.signerLogDao().getAll()

            for(log in allLogs!!) {
                db.signerLogDao().delete(log)
            }

            allLogs = db.signerLogDao().getAll()
        }
    }

    fun onClick(view : View) {
        var msg = ""

        for(log in allLogs!!) {
            msg += "log: $log\n"
        }

        Toast.makeText(this@SignActivity, msg, Toast.LENGTH_SHORT).show()

        GlobalScope.launch {
            var signature = ""
            var errorMsg = ""
//            var beforeTime = System.currentTimeMillis()
//            var afterTime :Long = 0
            val origin = "http://${parsedCallback.host}:${parsedCallback.port}"
            try {
                signature = scrappySigner.sign(origin, period)
                Log.d("SIGNATYRE", "${origin}, ${period}")
            } catch (e: java.lang.Exception) {
                errorMsg = e.toString()
            }

            lifecycleScope.launch {
                if (errorMsg != "") {
                    Toast.makeText(this@SignActivity, "error: $errorMsg", Toast.LENGTH_LONG).show()
                    return@launch
                }

                val url = "${callback}#${signature}"
                Toast.makeText(this@SignActivity,"url: $url", Toast.LENGTH_SHORT).show()

                val browserIntent = Intent(Intent.ACTION_VIEW, Uri.parse(url))
                startActivity(browserIntent)
            }

        }
    }
}