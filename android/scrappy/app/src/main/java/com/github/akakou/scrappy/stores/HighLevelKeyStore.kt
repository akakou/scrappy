package com.github.akakou.scrappy.stores

import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import android.util.Base64
import java.io.ByteArrayInputStream
import java.io.ByteArrayOutputStream
import java.security.KeyPairGenerator
import java.security.KeyStore
import java.security.PrivateKey
import javax.crypto.Cipher
import javax.crypto.CipherInputStream
import javax.crypto.CipherOutputStream

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
