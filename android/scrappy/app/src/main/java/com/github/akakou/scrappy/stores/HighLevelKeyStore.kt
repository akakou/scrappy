package com.github.akakou.scrappy.stores

import android.annotation.SuppressLint
import android.content.Context
import android.security.KeyPairGeneratorSpec
import android.util.Base64
import java.math.BigInteger
import java.security.KeyPairGenerator
import java.security.KeyStore
import java.util.*
import javax.crypto.Cipher
import javax.security.auth.x500.X500Principal

class HighLevelKeyStore (keyAlias: String, context: Context){
    val KEY_PROVIDER = "AndroidKeyStore"
    val KEY_ALGORITHM = "RSA"
    val ALGORITHM = "RSA/ECB/PKCS1Padding"
    val keyAlias = keyAlias
    val context = context

    val keyStore = KeyStore.getInstance("AndroidKeyStore")!!

    init {
        keyStore.load(null)
        createNewKey()
    }

    @SuppressLint("SuspiciousIndentation")
    fun createNewKey() {
        if (keyStore.containsAlias(keyAlias)) {
            return
        }

        val start = Calendar.getInstance()
        val end = Calendar.getInstance()
        end.add(Calendar.YEAR, 100)

        val keys = KeyPairGeneratorSpec.Builder(context)
            .setAlias(keyAlias)
            .setSubject(X500Principal(String.format("CN=%s", keyAlias)))
            .setSerialNumber(BigInteger.valueOf(1000000))
            .setStartDate(start.time)
            .setEndDate(end.time)
            .build()

        val keyPairGenerator = KeyPairGenerator.getInstance(KEY_ALGORITHM, KEY_PROVIDER)
            keyPairGenerator.initialize(keys)
            keyPairGenerator.generateKeyPair()
    }

    fun encrypt(plainText: ByteArray): String {
        val publicKey = keyStore.getCertificate(keyAlias).publicKey

        val cipher = Cipher.getInstance(ALGORITHM)
        cipher.init(Cipher.ENCRYPT_MODE, publicKey)
        val bytes = cipher.doFinal(plainText)

        return Base64.encodeToString(bytes, Base64.DEFAULT)
    }

    fun decrypt(cipherText: String): ByteArray {
        val keyStore = KeyStore.getInstance(KEY_PROVIDER)
        keyStore.load(null)
        if (!keyStore.containsAlias(keyAlias)) {
            error("there are no ekey")
        }

        val privateKey = keyStore.getKey(keyAlias, null)
        val cipher = Cipher.getInstance(ALGORITHM)
        cipher.init(Cipher.DECRYPT_MODE, privateKey)
        val bytes = Base64.decode(cipherText, Base64.DEFAULT)

        val b = cipher.doFinal(bytes)
        return b
    }
}
