package com.penndev.socks5.ui.home

import android.R
import android.content.Context
import android.content.SharedPreferences
import android.os.Bundle
import android.text.InputType
import android.util.Log
import android.view.MenuItem
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.edit
import com.penndev.socks5.databinding.ActivityHomeNodeBinding


class HomeNodeActivity : AppCompatActivity() {

    private lateinit var binding: ActivityHomeNodeBinding

    private var sharedPreferences: SharedPreferences? = null
        get() {
            return getSharedPreferences(packageName, Context.MODE_PRIVATE)
        }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityHomeNodeBinding.inflate(layoutInflater)
        setContentView(binding.root)
        initView()
    }

    private fun initView() {
        supportActionBar?.setDisplayHomeAsUpEnabled(true)
        binding.inputPassShow.setOnCheckedChangeListener { _, isChecked ->
            binding.inputPass.inputType = if (isChecked) InputType.TYPE_CLASS_TEXT
            else InputType.TYPE_CLASS_TEXT or InputType.TYPE_TEXT_VARIATION_PASSWORD
        }
        binding.inputSubmit.setOnClickListener {
            var host = binding.inputHost.text.toString()
            var port = binding.inputPort.text.toString().toInt()
            var user = binding.inputUser.text.toString()
            var pass = binding.inputPass.text.toString()
            sharedPreferences?.edit {
                putString("host", host)
                putInt("port", port)
                putString("user", user)
                putString("pass", pass)
                apply()
            }
            finish()
        }
        binding.inputHost.setText(sharedPreferences?.getString("host", ""))
        val port = sharedPreferences?.getInt("port", 0)
        if (port != null && port > 0) binding.inputPort.setText(port.toString())
        binding.inputUser.setText(sharedPreferences?.getString("user", ""))
        binding.inputPass.setText(sharedPreferences?.getString("pass", ""))
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        return when (item.itemId) {
            R.id.home -> {
                finish()
                true
            }
            else -> super.onOptionsItemSelected(item)
        }
    }

}