package com.github.akakou.scrappy

import android.content.Context
import android.os.Bundle
import android.view.View
import android.widget.EditText
import androidx.appcompat.app.AppCompatActivity
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.preferencesDataStore
import androidx.lifecycle.lifecycleScope
import com.github.akakou.scrappy.stores.Stores
import kotlinx.coroutines.launch

val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "settings")
val oneMinute = 60
val oneSecond = 1000


class MainActivity : AppCompatActivity() {
    lateinit var resultEditText: EditText
    lateinit var issuerURLEditText: EditText
    lateinit var ipkEditText: EditText
    lateinit var scrappySigner: ScrappySigner

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        issuerURLEditText = findViewById<EditText>(R.id.issuer_edit_text)
        ipkEditText = findViewById(R.id.ipk_edit_text)
        resultEditText = findViewById<EditText>(R.id.result_edit_text)

        val stores = Stores(this)
        scrappySigner = ScrappySigner(stores, null)
        scrappySigner.protocol = "http://"
    }

     fun onClick(v: View) {
        val issuerDomain = issuerURLEditText.text.toString()
        val ipk = ipkEditText.text.toString()

        var msg = "ok"

         lifecycleScope.launch {
             try {
                 scrappySigner.join(issuerDomain, ipk)
             } catch (e: java.lang.Exception) {
                 msg = e.toString()
             }

             resultEditText.setText(msg)
         }
    }
}
