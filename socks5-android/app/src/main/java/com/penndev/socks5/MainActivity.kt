package com.penndev.socks5

import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.content.Intent
import android.content.SharedPreferences
import android.net.VpnService
import android.os.Build
import android.os.Bundle
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.edit
import com.penndev.socks5.databinding.ActivityMainBinding
import com.penndev.socks5.service.Socks5Service

class MainActivity : AppCompatActivity() {
    //获取vpn启动权限标志 result 返回码
    private val allowCreateService = 1

    // xml UI 实例
    private lateinit var binding: ActivityMainBinding

    private lateinit var sharedPreferences: SharedPreferences

    // 绑定UI和Socks5Service状态
    private fun onStatus(): Socks5Service.OnStatus {
        // 初始化UI状态
        fun startUI() {
            binding.handleActionIcon.setImageResource(R.drawable.activity_main_start_close)
            binding.handleActionText.setText(R.string.activity_main_start_close_text)
        }

        fun closeUI() {
            binding.handleActionIcon.setImageResource(R.drawable.activity_main_start)
            binding.handleActionText.setText(R.string.activity_main_start_text)
        }
        if (Socks5Service.status) startUI() else closeUI()
        return object : Socks5Service.OnStatus {
            override fun start() {
                super.start()
                startUI()
                // 保存表单数据
                sharedPreferences.edit {
                    putString("inputHost", binding.inputHost.text.toString())
                    putString("inputPortStr", binding.inputPort.text.toString())
                    putString("inputUser", binding.inputUser.text.toString())
                    putString("inputPass", binding.inputPass.text.toString())
                    apply()
                }
            }

            override fun close() {
                super.close()
                closeUI()
            }
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)
        //设置各种初始化
        sharedPreferences = getSharedPreferences(packageName, Context.MODE_PRIVATE)
        binding.inputHost.setText(sharedPreferences.getString("inputHost", ""))
        binding.inputPort.setText(sharedPreferences.getString("inputPortStr", ""))
        binding.inputUser.setText(sharedPreferences.getString("inputUser", ""))
        binding.inputPass.setText(sharedPreferences.getString("inputPass", ""))
        binding.handleActionIcon.setOnClickListener {
            if (Socks5Service.status) onStopSocks5Service() else onStartSocks5Service()
        }
        Socks5Service.onStatus = onStatus()
        createNotificationChannel()
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val name = Socks5Service.notifyChannelName
            val descriptionText = getString(R.string.notify_channel_description)
            val importance = NotificationManager.IMPORTANCE_DEFAULT
            val channel = NotificationChannel(Socks5Service.notifyChannelID, name, importance).apply {
                description = descriptionText
            }
            // Register the channel with the system.
            val notificationManager: NotificationManager =
                getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
            notificationManager.createNotificationChannel(channel)
        }
    }

    private fun onStopSocks5Service() {
        startService(Intent(this, Socks5Service::class.java).putExtra("close", true))
    }

    // 启动 Socks5Service的详细步骤和流程。
    private fun onStartSocks5Service() {
        val intentPrepare = VpnService.prepare(this)
        if (intentPrepare != null) {
            startActivityForResult(intentPrepare, allowCreateService)
            return
        }

        val inputHost = binding.inputHost.text.toString()
        val inputPortStr = binding.inputPort.text.toString()
        val inputUser = binding.inputUser.text.toString()
        val inputPass = binding.inputPass.text.toString()
        if (inputHost == "" || inputPortStr == "") {
            Toast.makeText(this, getString(R.string.toast_main_param_error), Toast.LENGTH_LONG).show()
            return
        }
        val inputPort = inputPortStr.toInt()

        val bundle = Bundle().apply {
            putString("host", inputHost)
            putInt("port", inputPort)
            putString("user", inputUser)
            putString("pass", inputPass)
        }
        startService(Intent(this, Socks5Service::class.java).putExtras(bundle))
    }

    // 请求用户授权开启vpn的结果
    override fun onActivityResult(requestCode: Int, resultCode: Int, data: Intent?) {
        super.onActivityResult(requestCode, resultCode, data)
        when (requestCode) {
            allowCreateService -> {
                if (resultCode == RESULT_OK) {
                    onStartSocks5Service()
                } else {
                    Toast.makeText(this, "您拒绝了VPN请求[$resultCode]", Toast.LENGTH_LONG).show()
                }
            }
        }
    }
}
