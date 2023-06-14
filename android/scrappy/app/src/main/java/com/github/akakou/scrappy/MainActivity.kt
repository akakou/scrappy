package com.github.akakou.scrappy

import android.content.Context
import android.os.Bundle
import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import android.util.Base64
import android.view.View
import android.widget.EditText
import android_scrappy.Android_scrappy.androidJoin
import android_scrappy.Android_scrappy.androidSign
import androidx.appcompat.app.AppCompatActivity
import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import androidx.lifecycle.lifecycleScope
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import kotlinx.coroutines.launch
import org.json.JSONObject
import java.io.ByteArrayInputStream
import java.io.ByteArrayOutputStream
import java.security.KeyPairGenerator
import java.security.KeyStore
import java.security.PrivateKey
import javax.crypto.Cipher
import javax.crypto.CipherInputStream
import javax.crypto.CipherOutputStream

val Context.dataStore: DataStore<Preferences> by preferencesDataStore(name = "settings")

class HighLevelKeyStore (keyAlias: String){
    val KEY_PROVIDER = "AndroidKeyStore"
    val ALGORITHM = "RSA/ECB/OAEPPadding"
    val keyAlias = keyAlias

    val keyStore = KeyStore.getInstance("AndroidKeyStore")!!

    init {
        keyStore.load(null)

        createNewKey()
    }

    fun createNewKey() {
        if (keyStore.containsAlias(keyAlias)) {
            return
        }

        val keyPairGenerator = KeyPairGenerator.getInstance(
            KeyProperties.KEY_ALGORITHM_RSA, KEY_PROVIDER
        )
        keyPairGenerator.initialize(
            KeyGenParameterSpec.Builder(
                keyAlias,
                KeyProperties.PURPOSE_DECRYPT
            )
                .setDigests(KeyProperties.DIGEST_SHA256, KeyProperties.DIGEST_SHA512)
                .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_RSA_OAEP)
                .build()
        )

        keyPairGenerator.generateKeyPair()
    }

    fun encrypt(plainText: ByteArray): String {
        val publicKey = keyStore.getCertificate(keyAlias).publicKey

        val cipher = Cipher.getInstance(ALGORITHM)
        cipher.init(Cipher.ENCRYPT_MODE, publicKey)

        val outputStream = ByteArrayOutputStream()
        val cipherOutputStream = CipherOutputStream(
            outputStream, cipher
        )
        cipherOutputStream.write(plainText)
        cipherOutputStream.close()

        val bytes = outputStream.toByteArray()
        return Base64.encodeToString(bytes, Base64.DEFAULT)
    }

    fun decrypt(cipherText: String): ByteArray {
        val privateKey = keyStore.getKey(keyAlias, null) as PrivateKey

        val cipher = Cipher.getInstance(ALGORITHM)
        cipher.init(Cipher.DECRYPT_MODE, privateKey)

        val cipherInputStream = CipherInputStream(
            ByteArrayInputStream(Base64.decode(cipherText, Base64.DEFAULT)), cipher
        )
        val outputStream = ByteArrayOutputStream()

        var b: Int
        while (cipherInputStream.read().also { b = it } != -1) {
            outputStream.write(b)
        }
        outputStream.close()


        return outputStream.toByteArray()
    }
}

class HighLevelDataStore(dataStoreAlias: String, context: Context) {
    val preference = stringPreferencesKey(dataStoreAlias)
    var context: Context

    init {
        this.context = context
    }

    suspend fun store(data: String) {
        context.dataStore.edit { settings ->
            settings[preference] = data
        }
    }

    fun load():String {
        val loadFlow: Flow<String> = context.dataStore.data
            .map { preferences ->
                preferences[preference] ?: ""
            }

        return loadFlow.toString()
    }
}

class SecretStore (keyAlias: String, dataStoreAlias: String, context: Context){
    var dataStore : HighLevelDataStore
    var keyStore : HighLevelKeyStore

    init {
        dataStore = HighLevelDataStore(dataStoreAlias, context)
        keyStore = HighLevelKeyStore(keyAlias)
    }

    suspend fun store(secret: ByteArray) {
        val cipher = keyStore.encrypt(secret)
        dataStore.store(cipher)
    }

    fun load(): ByteArray {
        val cipher = dataStore.load()
        return keyStore.decrypt(cipher)
    }
}

class Stores(context: Context){
    var secret = SecretStore("scrappy-secret", "scrappy-secret", context)
    var cred = HighLevelDataStore("scrappy-credential", context)
    var ipk = HighLevelDataStore("scrappy-ipk", context)
}

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

        return jsonObject.getJSONObject("data")
    }

    suspend fun join(issuerDomain: String, ipk: String) {
        val response = androidJoin("$protocol$issuerDomain", ipk)
        val config = parseLibraryResponse(response)

        val cred = config["Cred"].toString()
        val base64Sk = config["SK"].toString()
        val bytesSK = Base64.decode(base64Sk, Base64.DEFAULT)

        stores.secret.store(bytesSK)
        stores.cred.store(cred)
        stores.ipk.store(ipk)
    }

    fun sign(origin: String, unixTime: Long): String {
        val rawSK = stores.secret.load()
        val cred = stores.cred.load()
        val ipk = stores.ipk.load()

        val bytesSK = Base64.encode(rawSK, Base64.DEFAULT)

        return androidSign(origin, unixTime, bytesSK.toString(), cred, ipk)
    }
}


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
        scrappySigner = ScrappySigner(stores)
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
