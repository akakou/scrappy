package com.github.akakou.scrappy

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.Toast

class MainActivity : AppCompatActivity() {
    external fun stringFromJNI(): String

    companion object {
        init {
            System.loadLibrary("ecdaa_android")
        }
    }
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        var hello = stringFromJNI()
        Toast.makeText(this, hello, Toast.LENGTH_LONG).show()
    }
}